package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// check is google-mcp check, the task google-mcp:check: the setup proved end to end, a token refresh at Google (so the
// keychain's client id and secret are right) and a real Calendar call through the server.
func check() {
	if err := checkAll(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// checkAll is check, its failure as the message to print.
func checkAll() error {
	home, _ := os.UserHomeDir()
	files, _ := filepath.Glob(filepath.Join(home, ".google-mcp-accounts", "*.json"))
	if len(files) == 0 {
		return errors.New("✗ no signed-in Google account: run mise run google-mcp:auth")
	}
	id, err := secretValue("GOOGLE_CLIENT_ID")
	if err != nil {
		return errors.New("✗ " + err.Error())
	}
	secret, err := secretValue("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return errors.New("✗ " + err.Error())
	}
	for _, f := range files {
		name := strings.TrimSuffix(filepath.Base(f), ".json")
		refresh := refreshToken(f)
		if refresh == "" {
			return fmt.Errorf("✗ %s: no refresh token: run mise run google-mcp:auth", filepath.Base(f))
		}
		resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
			"grant_type": {"refresh_token"}, "refresh_token": {refresh}, "client_id": {id}, "client_secret": {secret},
		})
		if err != nil {
			return fmt.Errorf("✗ Google token endpoint: %v", err)
		}
		var body struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			switch body.Error {
			case "invalid_grant":
				return fmt.Errorf("✗ %s: token expired or revoked: run mise run google-mcp:auth", name)
			case "invalid_client":
				return errors.New("✗ client id/secret rejected: run mise run google-mcp:rotate (or google-mcp:secret)")
			case "":
				body.Error = "unknown"
			}
			return fmt.Errorf("✗ Google token endpoint answered HTTP %d (%s)", resp.StatusCode, body.Error)
		}
		fmt.Printf("✓ %s: Google refreshed the token with the keychain client id + secret\n", name)
	}

	srv, err := startServer()
	if err != nil {
		return fmt.Errorf("✗ google-mcp-server didn't start: %v", err)
	}
	defer srv.stop()
	srv.initialize()
	srv.send(`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"calendar_list","arguments":{}}}`)
	line, ok := srv.wait(srv.outLines, regexp.MustCompile(`"id":2[,}]`), 60*time.Second)
	if !ok {
		return fmt.Errorf("✗ google-mcp-server did not answer:\n%s", srv.tail(srv.errLines, 5))
	}
	if strings.Contains(line, `"error":`) || strings.Contains(line, `"isError":true`) {
		return fmt.Errorf("✗ Calendar call failed: %s", regexp.MustCompile(`"message":"[^"]*"`).FindString(line))
	}
	version := "?"
	if out, err := exec.Command("google-mcp-server", "--version").Output(); err == nil {
		if f := strings.Fields(string(out)); len(f) > 1 {
			version = f[1]
		}
	}
	fmt.Printf("✓ google-mcp-server %s: Calendar API answered (%d calendars)\n", version, strings.Count(line, `\"id\"`))
	return nil
}

// refreshToken is the first refresh_token anywhere in an account file, or "".
func refreshToken(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	var v any
	if json.Unmarshal(b, &v) != nil {
		return ""
	}
	var find func(any) string
	find = func(v any) string {
		switch v := v.(type) {
		case map[string]any:
			if s, ok := v["refresh_token"].(string); ok && s != "" {
				return s
			}
			for _, x := range v {
				if s := find(x); s != "" {
					return s
				}
			}
		case []any:
			for _, x := range v {
				if s := find(x); s != "" {
					return s
				}
			}
		}
		return ""
	}
	return find(v)
}
