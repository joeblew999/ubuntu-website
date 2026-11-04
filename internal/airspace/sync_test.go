package airspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// faa serves what Sync talks to: each dataset's FeatureServer layer (its version, doc.go rule 1), its download, and
// the manifest published in R2 (rule 2), with the versions published records, or a 404 when published is nil.
// broken names a dataset whose layer answers an ArcGIS error. It counts the downloads.
func faa(t *testing.T, versions, published map[string]int64, broken string) (*httptest.Server, map[string]int) {
	t.Helper()
	downloads := map[string]int{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/") // layer/<key>, data/<key>, manifest
		switch {
		case parts[0] == "manifest":
			if published == nil {
				http.NotFound(w, r)
				return
			}
			m := ManifestUSA{Layers: map[string]ManifestLayer{}}
			for key, v := range published {
				m.Layers[key] = ManifestLayer{SourceVersion: time.UnixMilli(v).UTC().Format(time.RFC3339Nano), Features: 7}
			}
			json.NewEncoder(w).Encode(m)
		case parts[0] == "layer" && parts[1] == broken:
			w.Write([]byte(`{"error":{"code":500,"message":"Internal error"}}`))
		case parts[0] == "layer":
			fmt.Fprintf(w, `{"editingInfo":{"dataLastEditDate":%d}}`, versions[parts[1]])
		case parts[0] == "data":
			downloads[parts[1]]++
			w.Write([]byte(`{"type":"FeatureCollection","features":[{"type":"Feature","geometry":{"type":"Point","coordinates":[0,0]},"properties":{}}]}`))
		}
	}))
	t.Cleanup(srv.Close)
	return srv, downloads
}

// testDatasets puts datasets a and b, served by srv, in the registry for the test.
func testDatasets(t *testing.T, srv *httptest.Server) []string {
	keys := []string{"test-a", "test-b"}
	for _, key := range keys {
		Datasets[key] = Dataset{Key: key, GeoJSON: key + ".geojson", BaseURL: srv.URL + "/data/" + key, LayerURL: srv.URL + "/layer/" + key}
	}
	t.Cleanup(func() {
		for _, key := range keys {
			delete(Datasets, key)
		}
	})
	return keys
}

// TestSyncRules: the sync rules of doc.go, rules 1 to 4.
func TestSyncRules(t *testing.T) {
	now := map[string]int64{"test-a": 1793761023497, "test-b": 1793750000000}
	for _, tc := range []struct {
		name      string
		published map[string]int64
		broken    string
		force     bool
		want      map[string]string // status per dataset
	}{
		{"nothing changed: nothing downloaded, though no file is here", now, "", false, map[string]string{"test-a": "unchanged", "test-b": "unchanged"}},
		{"one changed: only it", map[string]int64{"test-a": now["test-a"], "test-b": 1}, "", false, map[string]string{"test-a": "unchanged", "test-b": "updated"}},
		{"never published: every one", nil, "", false, map[string]string{"test-a": "updated", "test-b": "updated"}},
		{"a layer that can't say its version: an error", now, "test-a", false, map[string]string{"test-a": "error", "test-b": "unchanged"}},
		{"forced: every one", now, "", true, map[string]string{"test-a": "updated", "test-b": "updated"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv, downloads := faa(t, now, tc.published, tc.broken)
			dir := t.TempDir()
			opts := SyncOptions{OutputDir: filepath.Join(dir, "geojson"), DataDir: filepath.Join(dir, "data"), Force: tc.force,
				Timeout: 10 * time.Second, Datasets: testDatasets(t, srv), PublishedURL: srv.URL + "/manifest"}
			result, err := Sync(opts)
			if err != nil {
				t.Fatal(err)
			}
			for key, status := range tc.want {
				got := result.Datasets[key]
				if got.Status != status {
					t.Errorf("%s: %s (%s), want %s", key, got.Status, got.Error, status)
				}
				_, statErr := os.Stat(filepath.Join(opts.OutputDir, key+".geojson"))
				if downloaded := downloads[key] > 0; downloaded != (status == "updated") || (statErr == nil) != downloaded {
					t.Errorf("%s: %d downloads, file there: %v, for %s", key, downloads[key], statErr == nil, status)
				}
				if status != "error" && got.Version != time.UnixMilli(now[key]).UTC().Format(time.RFC3339Nano) {
					t.Errorf("%s: version %q", key, got.Version)
				}
			}
			if result.HasChanges != strings.Contains(fmt.Sprint(tc.want), "updated") {
				t.Errorf("has_changes %v", result.HasChanges)
			}
		})
	}
}

// TestManifestMerge: the manifest describes every layer in R2 (doc.go, rule 5): one this run downloaded from this
// machine's files and the run's version; any other from what's published, even with an older file here.
func TestManifestMerge(t *testing.T) {
	dir := t.TempDir()
	tiles, geojson := filepath.Join(dir, "tiles"), filepath.Join(dir, "geojson")
	os.MkdirAll(tiles, 0o755)
	os.MkdirAll(geojson, 0o755)
	for _, key := range []string{"uas", "sua"} {
		os.WriteFile(filepath.Join(geojson, Datasets[key].GeoJSON), []byte(`{"type":"Feature"},{"type":"Feature"}`), 0o644)
		os.WriteFile(filepath.Join(tiles, Datasets[key].PMTiles), make([]byte, 1<<20), 0o644)
	}
	run := SyncResult{Datasets: map[string]DatasetSync{"uas": {Status: "updated", Version: "v-new"}, "sua": {Status: "unchanged", Version: "v-sua"}}}
	published := &ManifestUSA{Layers: map[string]ManifestLayer{
		"laanc": {SourceVersion: "v-old", Features: 100, SizeMB: 9}, // the uas dataset's layer
		"sua":   {SourceVersion: "v-sua", Features: 50, SizeMB: 3},
	}}
	if err := GenerateManifests(dir, tiles, geojson, run, published); err != nil {
		t.Fatal(err)
	}
	var m ManifestUSA
	data, _ := os.ReadFile(filepath.Join(dir, FileUSAManifest))
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatal(err)
	}
	if l := m.Layers["laanc"]; l.SourceVersion != "v-new" || l.Features != 2 || l.SizeMB != 1 {
		t.Errorf("uas, downloaded: %+v", l)
	}
	if l := m.Layers["sua"]; l.SourceVersion != "v-sua" || l.Features != 50 || l.SizeMB != 3 {
		t.Errorf("sua, unchanged: %+v (want what's published, not this machine's older file)", l)
	}
	if l := m.Layers["airports"]; l.SourceVersion != "" {
		t.Errorf("airports, in neither: version %q, want none (so the next sync downloads it)", l.SourceVersion)
	}
}
