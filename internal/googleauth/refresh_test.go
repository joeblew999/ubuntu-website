package googleauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// fakeGoogle points tokenURL at a server answering as Google does, and PATH at a fnox that prints a dummy client id
// and secret, so no test reads the keychain or calls Google. It counts the refreshes.
func fakeGoogle(t *testing.T, status int, body string) *int {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("fnox stub is a shell script")
	}
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		r.ParseForm()
		if r.Form.Get("grant_type") != "refresh_token" || r.Form.Get("refresh_token") != "old-refresh" ||
			r.Form.Get("client_id") != "dummy-GOOGLE_CLIENT_ID" || r.Form.Get("client_secret") != "dummy-GOOGLE_CLIENT_SECRET" {
			t.Errorf("unexpected form %v", r.Form)
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	old := tokenURL
	tokenURL = srv.URL
	t.Cleanup(func() { tokenURL = old })

	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "fnox"), []byte("#!/bin/sh\necho \"dummy-$2\"\n"), 0o755)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return &calls
}

// account writes an account file as google-mcp-server does, its token expiring at expiry.
func account(t *testing.T, expiry time.Time) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "me_at_example_com.json")
	data := `{
  "email": "me@example.com",
  "name": "Me",
  "token": {"access_token": "old-access", "token_type": "Bearer", "refresh_token": "old-refresh", "expiry": "` +
		expiry.Format(time.RFC3339Nano) + `", "expires_in": 3599},
  "last_used": "2026-01-01T00:00:00Z"
}`
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadAccessTokenRefreshes(t *testing.T) {
	calls := fakeGoogle(t, http.StatusOK, `{"access_token":"new-access","expires_in":3599,"token_type":"Bearer"}`)
	path := account(t, time.Now().Add(-time.Hour))

	token, err := LoadAccessToken(path)
	if err != nil || token != "new-access" || *calls != 1 {
		t.Fatalf("got %q, %v, %d refreshes", token, err, *calls)
	}
	data, _ := os.ReadFile(path)
	var file map[string]any
	if err := json.Unmarshal(data, &file); err != nil {
		t.Fatal(err)
	}
	tok := file["token"].(map[string]any)
	if tok["access_token"] != "new-access" || tok["refresh_token"] != "old-refresh" || file["name"] != "Me" || file["last_used"] == nil {
		t.Errorf("saved file lost something:\n%s", data)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Errorf("saved file is %v, want 0600", info.Mode().Perm())
	}
	// The saved token is good for an hour: the next load uses it.
	if token, err := LoadAccessToken(path); err != nil || token != "new-access" || *calls != 1 {
		t.Errorf("second load: %q, %v, %d refreshes", token, err, *calls)
	}
}

func TestLoadAccessTokenFresh(t *testing.T) {
	calls := fakeGoogle(t, http.StatusOK, `{}`)
	path := account(t, time.Now().Add(time.Hour))
	if token, err := LoadAccessToken(path); err != nil || token != "old-access" || *calls != 0 {
		t.Errorf("got %q, %v, %d refreshes", token, err, *calls)
	}
}

func TestLoadAccessTokenRevoked(t *testing.T) {
	fakeGoogle(t, http.StatusBadRequest, `{"error":"invalid_grant"}`)
	path := account(t, time.Now().Add(30*time.Second)) // within the minute: refreshed
	before, _ := os.ReadFile(path)
	_, err := LoadAccessToken(path)
	if err == nil || !strings.Contains(err.Error(), "mise run google-mcp:auth") {
		t.Errorf("got %v, want a hint to sign in again", err)
	}
	if after, _ := os.ReadFile(path); string(after) != string(before) {
		t.Errorf("a refused refresh changed the file")
	}
}
