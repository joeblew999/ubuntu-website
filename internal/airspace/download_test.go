package airspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// arcgis is a FeatureServer query endpoint holding count features, answering pages of what page asks for; errAt, when
// not negative, is the offset it answers ArcGIS's HTTP 200 error at (failing times times, every time if negative),
// and serve caps the features it serves.
func arcgis(t *testing.T, count, serve, errAt, failing int) *httptest.Server {
	t.Helper()
	old := retryDelay
	retryDelay = 0
	t.Cleanup(func() { retryDelay = old })
	failed := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("returnCountOnly") == "true" {
			fmt.Fprintf(w, `{"count":%d}`, count)
			return
		}
		offset, _ := strconv.Atoi(q.Get("resultOffset"))
		size, _ := strconv.Atoi(q.Get("resultRecordCount"))
		if offset == errAt && (failing < 0 || failed < failing) {
			failed++
			w.Write([]byte(`{"error":{"code":499,"message":"Token Required"}}`))
			return
		}
		var features []string
		for i := offset; i < offset+size && i < serve; i++ {
			features = append(features, fmt.Sprintf(`{"type":"Feature","id":%d,"geometry":{"type":"Point","coordinates":[0,0]},"properties":{}}`, i))
		}
		fmt.Fprintf(w, `{"type":"FeatureCollection","features":[%s]}`, strings.Join(features, ","))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestDownloadPaginated(t *testing.T) {
	for _, tc := range []struct {
		name                         string
		count, serve, errAt, failing int
		want                         string // "" = succeeds
	}{
		{"every page", 5, 5, -1, 0, ""},
		{"a page that errs twice: retried", 5, 5, 2, 2, ""},
		{"a page that keeps erring", 5, 5, 2, -1, "page at offset 2, after 3 attempts: ArcGIS error 499 Token Required"},
		{"the first page keeps erring", 5, 5, 0, -1, "page at offset 0"},
		{"fewer than counted", 5, 3, -1, 0, "got 3 of the 5 features"},
		{"counts none", 0, 0, -1, 0, "reports none"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := arcgis(t, tc.count, tc.serve, tc.errAt, tc.failing)
			out := filepath.Join(t.TempDir(), "out.geojson")
			err := DownloadPaginated(srv.Client(), Dataset{BaseURL: srv.URL, PageSize: 2}, out)
			if tc.want != "" {
				if err == nil || !strings.Contains(err.Error(), tc.want) {
					t.Fatalf("got %v, want %q", err, tc.want)
				}
				if _, err := os.Stat(out); err == nil {
					t.Error("a failed download wrote its file")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			var fc struct{ Features []any }
			data, _ := os.ReadFile(out)
			if json.Unmarshal(data, &fc) != nil || len(fc.Features) != tc.count {
				t.Errorf("wrote %d features, want %d", len(fc.Features), tc.count)
			}
		})
	}
}
