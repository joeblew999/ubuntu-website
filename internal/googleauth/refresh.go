package googleauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// tokenURL is Google's OAuth token endpoint (a var so tests can point it at their own server).
var tokenURL = "https://oauth2.googleapis.com/token"

// Keychain is a secret of this repo's fnox.toml (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET), from the keychain.
func Keychain(key string) (string, error) {
	out, err := exec.Command("fnox", "get", key).Output()
	if err != nil {
		return "", fmt.Errorf("fnox get %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// RefreshError is Google refusing a refresh: Code is its error, "invalid_grant" (the refresh token expired or was
// revoked) or "invalid_client" (the client id or secret is wrong).
type RefreshError struct {
	Status int
	Code   string
}

func (e *RefreshError) Error() string {
	switch e.Code {
	case "invalid_grant":
		return "Google refused the refresh token (expired or revoked): run mise run google-mcp:auth"
	case "invalid_client":
		return "Google rejected the client id/secret: run mise run google-mcp:rotate (or google-mcp:secret)"
	}
	return fmt.Sprintf("Google's token endpoint answered HTTP %d (%s)", e.Status, e.Code)
}

// Refresh trades a refresh token for a new access token at Google, with the OAuth client's id and secret.
func Refresh(refreshToken, clientID, clientSecret string) (*Token, error) {
	resp, err := http.PostForm(tokenURL, url.Values{
		"grant_type": {"refresh_token"}, "refresh_token": {refreshToken}, "client_id": {clientID}, "client_secret": {clientSecret},
	})
	if err != nil {
		return nil, fmt.Errorf("Google's token endpoint: %w", err)
	}
	defer resp.Body.Close()
	var body struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		Error        string `json:"error"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if resp.StatusCode != http.StatusOK || body.AccessToken == "" {
		if body.Error == "" {
			body.Error = "unknown"
		}
		return nil, &RefreshError{Status: resp.StatusCode, Code: body.Error}
	}
	return &Token{
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
		TokenType:    body.TokenType,
		Expiry:       time.Now().Add(time.Duration(body.ExpiresIn) * time.Second).Format(time.RFC3339Nano),
	}, nil
}

// refreshAccount refreshes the token of the account file at path with the keychain's client id and secret, and saves
// the new one there.
func refreshAccount(path string, account *Account) (string, error) {
	if account.Token.RefreshToken == "" {
		return "", fmt.Errorf("%s: the access token expired and there's no refresh token: run mise run google-mcp:auth", account.Email)
	}
	id, err := Keychain("GOOGLE_CLIENT_ID")
	if err != nil {
		return "", fmt.Errorf("%s: the access token expired, and no client id to refresh it with: %w (mise run google-mcp:credentials)", account.Email, err)
	}
	secret, err := Keychain("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return "", fmt.Errorf("%s: the access token expired, and no client secret to refresh it with: %w (mise run google-mcp:credentials)", account.Email, err)
	}
	t, err := Refresh(account.Token.RefreshToken, id, secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", account.Email, err)
	}
	if err := saveToken(path, t); err != nil {
		return "", fmt.Errorf("%s: refreshed, but not saved: %w", account.Email, err)
	}
	return t.AccessToken, nil
}

// saveToken writes a refreshed token into the account file at path, keeping every other field google-mcp-server
// wrote there, with the file's own permissions; through a temporary file renamed over it, so no reader sees half.
func saveToken(path string, t *Token) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var file map[string]any
	if err := json.Unmarshal(data, &file); err != nil {
		return err
	}
	token, _ := file["token"].(map[string]any)
	if token == nil {
		token = map[string]any{}
		file["token"] = token
	}
	token["access_token"] = t.AccessToken
	token["expiry"] = t.Expiry
	if t.TokenType != "" {
		token["token_type"] = t.TokenType
	}
	if t.RefreshToken != "" {
		token["refresh_token"] = t.RefreshToken // Google sends one only when it rotates it
	}
	out, err := json.MarshalIndent(file, "", "  ")
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".token-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(append(out, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(info.Mode().Perm()); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
