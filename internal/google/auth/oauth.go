// Package auth provides Google OAuth authentication flows and credential management.
//
// This package supports two authentication patterns:
//
// 1. Application Default Credentials (ADC) - for gcloud/Terraform
//    Uses gcloud's official OAuth client to create credentials compatible
//    with Google Cloud SDK tools.
//
// 2. Google MCP Server tokens - for Claude Code integrations
//    Loads tokens from ~/.google-mcp-accounts for Gmail, Calendar, etc.
//
// Example - gcloud auth flow:
//
//	result, err := auth.RunOAuthFlow(auth.FlowOptions{
//	    Mode:       auth.ModeAssisted,
//	    AccountHint: "user@example.com",
//	    Timeout:    120 * time.Second,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Authenticated: %s\n", result.Email)
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// OAuth constants for gcloud application-default credentials
// These are gcloud's official OAuth client credentials (public, not secret)
const (
	GcloudClientID     = "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com"
	GcloudClientSecret = "d-FL95Q19q7MQmFpd7hHD0Ty"

	OAuthRedirectURI = "http://localhost:8085/"
	OAuthTokenURL    = "https://oauth2.googleapis.com/token"
	OAuthAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"

	// Default scopes for application-default credentials
	DefaultScopes = "openid https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/sqlservice.login https://www.googleapis.com/auth/accounts.reauth"
)

// FlowMode represents the OAuth flow mode
type FlowMode string

const (
	// ModeManual - Opens browser, user completes auth manually
	ModeManual FlowMode = "manual"

	// ModeAssisted - Opens default browser (ideal for passkey auth)
	ModeAssisted FlowMode = "assisted"

	// ModeAuto - Fully automated with Playwright (password accounts only)
	ModeAuto FlowMode = "auto"

	// ModeServerOnly - Only start callback server (for external automation)
	ModeServerOnly FlowMode = "server-only"
)

// FlowOptions configures the OAuth flow
type FlowOptions struct {
	Mode        FlowMode      // Which flow mode to use
	AccountHint string        // Email to pre-select (for assisted/auto modes)
	Timeout     time.Duration // Timeout for the flow (default: 120s)
	Headless    bool          // Run browser headless (for auto mode)
	Port        int           // Callback server port (default: 8085)
}

// FlowResult contains the result of an OAuth flow
type FlowResult struct {
	RefreshToken string // The refresh token for credential storage
	AccessToken  string // The access token (short-lived)
	Email        string // The authenticated user's email (if available)
	ExpiresIn    int    // Token expiration time in seconds
}

// GcloudCredentials represents the application_default_credentials.json format
type GcloudCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

// TokenResponse represents the OAuth token endpoint response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// BuildAuthURL constructs the OAuth authorization URL with PKCE
func BuildAuthURL(codeChallenge string) string {
	return BuildAuthURLWithOptions(codeChallenge, "", "")
}

// BuildAuthURLWithOptions constructs the OAuth authorization URL with PKCE and optional parameters
func BuildAuthURLWithOptions(codeChallenge, accountHint, customScopes string) string {
	params := url.Values{}
	params.Set("client_id", GcloudClientID)
	params.Set("redirect_uri", OAuthRedirectURI)
	params.Set("response_type", "code")

	if customScopes != "" {
		params.Set("scope", customScopes)
	} else {
		params.Set("scope", DefaultScopes)
	}

	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")

	if accountHint != "" {
		params.Set("login_hint", accountHint)
	}

	return OAuthAuthURL + "?" + params.Encode()
}

// ExchangeCodeForTokens exchanges the authorization code for tokens
func ExchangeCodeForTokens(code, codeVerifier string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", GcloudClientID)
	data.Set("client_secret", GcloudClientSecret)
	data.Set("code", code)
	data.Set("code_verifier", codeVerifier)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", OAuthRedirectURI)

	resp, err := http.Post(OAuthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("token error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	return &tokenResp, nil
}

// WriteGcloudCredentials writes the credentials to the gcloud config directory
func WriteGcloudCredentials(refreshToken string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create gcloud config directory if needed
	configDir := filepath.Join(homeDir, ".config", "gcloud")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	creds := GcloudCredentials{
		ClientID:     GcloudClientID,
		ClientSecret: GcloudClientSecret,
		RefreshToken: refreshToken,
		Type:         "authorized_user",
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	credFile := filepath.Join(configDir, "application_default_credentials.json")
	if err := os.WriteFile(credFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials: %w", err)
	}

	return nil
}

// GetGcloudCredentialsPath returns the path to the gcloud credentials file
func GetGcloudCredentialsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".config", "gcloud", "application_default_credentials.json"), nil
}

// CheckGcloudCredentials checks if application-default credentials exist and are valid
func CheckGcloudCredentials() (bool, string, error) {
	credFile, err := GetGcloudCredentialsPath()
	if err != nil {
		return false, "", err
	}

	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		return false, "", nil
	}

	// Try to read and validate
	data, err := os.ReadFile(credFile)
	if err != nil {
		return false, credFile, fmt.Errorf("could not read credentials: %w", err)
	}

	var creds GcloudCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return false, credFile, fmt.Errorf("invalid credentials format: %w", err)
	}

	if creds.RefreshToken == "" {
		return false, credFile, fmt.Errorf("credentials missing refresh_token")
	}

	return true, credFile, nil
}

// GeneratePKCE generates code_verifier and code_challenge for PKCE
func GeneratePKCE() (verifier, challenge string, err error) {
	pkce, err := browser.GeneratePKCE()
	if err != nil {
		return "", "", err
	}
	return pkce.Verifier, pkce.Challenge, nil
}

// RunOAuthFlowAssisted runs OAuth in assisted mode - opens default browser
// This mode is ideal for users who authenticate with passkeys
func RunOAuthFlowAssisted(opts FlowOptions) (*FlowResult, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 120 * time.Second
	}
	if opts.Port == 0 {
		opts.Port = 8085
	}

	// Generate PKCE values
	codeVerifier, codeChallenge, err := GeneratePKCE()
	if err != nil {
		return nil, err
	}

	authURL := BuildAuthURLWithOptions(codeChallenge, opts.AccountHint, "")

	// Channel to receive the auth code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Start the callback server
	server := &http.Server{Addr: fmt.Sprintf(":%d", opts.Port)}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, OAuthSuccessHTML)
			select {
			case codeChan <- code:
			default:
			}
		} else if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, OAuthErrorHTML, errMsg)
			select {
			case errChan <- fmt.Errorf("OAuth error: %s", errMsg):
			default:
			}
		}
	})
	server.Handler = mux

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			select {
			case errChan <- fmt.Errorf("server error: %w", err):
			default:
			}
		}
	}()

	// Open in default browser
	if err := browser.OpenURL(authURL); err != nil {
		// Non-fatal - user can still open URL manually
		fmt.Printf("Warning: Could not open browser: %v\n", err)
		fmt.Println("Please open this URL manually:")
		fmt.Println(authURL)
	}

	// Wait for code, error, or timeout
	select {
	case code := <-codeChan:
		// Shutdown server
		server.Shutdown(ctx)

		// Exchange code for tokens
		tokenResp, err := ExchangeCodeForTokens(code, codeVerifier)
		if err != nil {
			return nil, fmt.Errorf("token exchange failed: %w", err)
		}

		return &FlowResult{
			RefreshToken: tokenResp.RefreshToken,
			AccessToken:  tokenResp.AccessToken,
			ExpiresIn:    tokenResp.ExpiresIn,
		}, nil

	case err := <-errChan:
		server.Shutdown(ctx)
		return nil, err

	case <-ctx.Done():
		server.Shutdown(context.Background())
		return nil, fmt.Errorf("authentication timed out after %v", opts.Timeout)
	}
}
