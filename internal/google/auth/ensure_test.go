package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRefresh: a refresh token Google accepts works; one it answers invalid_grant to (expired, revoked) means
// signing in again; anything else is an error to show, not a sign-in.
func TestRefresh(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   error
	}{
		{"accepted", http.StatusOK, `{"access_token":"a","expires_in":3599}`, nil},
		{"expired or revoked", http.StatusBadRequest, `{"error":"invalid_grant"}`, ErrSignInAgain},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
			}))
			defer srv.Close()
			old := tokenURL
			tokenURL = srv.URL
			defer func() { tokenURL = old }()
			if err := refresh("id", "secret", "refresh"); !errors.Is(err, tc.want) {
				t.Errorf("got %v, want %v", err, tc.want)
			}
		})
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"invalid_client"}`))
	}))
	defer srv.Close()
	old := tokenURL
	tokenURL = srv.URL
	defer func() { tokenURL = old }()
	if err := refresh("id", "secret", "refresh"); err == nil || errors.Is(err, ErrSignInAgain) {
		t.Errorf("invalid_client: got %v, want an error that isn't a sign-in", err)
	}
}
