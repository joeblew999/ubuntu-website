// Package googleauth provides shared Google OAuth token management.
// Used by gmail, calendar, and other Google API clients.
package googleauth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DefaultTokenPath is the standard location for google-mcp-server tokens
const DefaultTokenPath = "~/.google-mcp-accounts"

// Account represents a Google account from google-mcp-server
type Account struct {
	Email string `json:"email"`
	Token Token  `json:"token"`
}

// Token represents OAuth token data
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Expiry       string `json:"expiry"`
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	if t.Expiry == "" {
		return false // No expiry means we don't know
	}
	expiry, err := time.Parse(time.RFC3339, t.Expiry)
	if err != nil {
		return false // Can't parse, assume not expired
	}
	return time.Now().After(expiry)
}

// LoadAccount loads a Google account from the token path
// tokenPath can be a directory (uses first .json file) or a specific file
func LoadAccount(tokenPath string) (*Account, error) {
	path, err := resolveTokenPath(tokenPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file %s: %w", path, err)
	}

	var account Account
	if err := json.Unmarshal(data, &account); err != nil {
		return nil, fmt.Errorf("failed to parse token file: %w", err)
	}

	if account.Token.AccessToken == "" {
		return nil, fmt.Errorf("no access token found in token file")
	}

	return &account, nil
}

// LoadAccessToken loads just the access token string (convenience function)
func LoadAccessToken(tokenPath string) (string, error) {
	account, err := LoadAccount(tokenPath)
	if err != nil {
		return "", err
	}
	return account.Token.AccessToken, nil
}

// ListAccounts returns all accounts in the token directory
func ListAccounts(tokenPath string) ([]*Account, error) {
	path := expandPath(tokenPath)

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat token path %s: %w", path, err)
	}

	if !info.IsDir() {
		// Single file, load it
		account, err := LoadAccount(tokenPath)
		if err != nil {
			return nil, err
		}
		return []*Account{account}, nil
	}

	// Directory - load all .json files
	files, err := filepath.Glob(filepath.Join(path, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list token files: %w", err)
	}

	var accounts []*Account
	for _, file := range files {
		account, err := LoadAccount(file)
		if err != nil {
			continue // Skip invalid files
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// resolveTokenPath expands ~ and finds the token file
func resolveTokenPath(tokenPath string) (string, error) {
	path := expandPath(tokenPath)

	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to stat token path %s: %w", path, err)
	}

	if info.IsDir() {
		// Find first .json file in directory
		files, err := filepath.Glob(filepath.Join(path, "*.json"))
		if err != nil {
			return "", fmt.Errorf("failed to list token files: %w", err)
		}
		if len(files) == 0 {
			return "", fmt.Errorf("no token files found in %s", path)
		}
		return files[0], nil
	}

	return path, nil
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
