// OAuth utilities for browser-based authentication flows.
// Supports PKCE, callback servers, and token exchange.
package browser

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// OAuthConfig holds configuration for an OAuth flow.
type OAuthConfig struct {
	// ClientID and ClientSecret for the OAuth application
	ClientID     string
	ClientSecret string

	// AuthURL is the authorization endpoint
	AuthURL string

	// TokenURL is the token exchange endpoint
	TokenURL string

	// RedirectURI for the callback (default: http://localhost:8085/)
	RedirectURI string

	// Scopes to request
	Scopes []string

	// CallbackPort for the local server (default: 8085)
	CallbackPort int

	// Timeout for the entire OAuth flow
	Timeout time.Duration
}

// OAuthResult represents the result of an OAuth flow.
type OAuthResult struct {
	// AuthCode is the authorization code received
	AuthCode string

	// AccessToken after exchange
	AccessToken string

	// RefreshToken after exchange
	RefreshToken string

	// ExpiresIn seconds
	ExpiresIn int

	// TokenType (usually "Bearer")
	TokenType string

	// Error if the flow failed
	Error string
}

// OAuthTokenResponse represents the token endpoint response.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// PKCE holds the code verifier and challenge for PKCE flow.
type PKCE struct {
	Verifier  string
	Challenge string
}

// GeneratePKCE generates code_verifier and code_challenge for PKCE.
func GeneratePKCE() (*PKCE, error) {
	// Generate 32 random bytes for verifier
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Base64 URL encode without padding
	verifier := base64.RawURLEncoding.EncodeToString(b)

	// SHA256 hash the verifier
	h := sha256.Sum256([]byte(verifier))

	// Base64 URL encode the hash without padding
	challenge := base64.RawURLEncoding.EncodeToString(h[:])

	return &PKCE{
		Verifier:  verifier,
		Challenge: challenge,
	}, nil
}

// BuildAuthURL constructs the OAuth authorization URL.
func BuildAuthURL(config *OAuthConfig, pkce *PKCE, state string) string {
	params := url.Values{}
	params.Set("client_id", config.ClientID)
	params.Set("redirect_uri", config.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(config.Scopes, " "))
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")

	if pkce != nil {
		params.Set("code_challenge", pkce.Challenge)
		params.Set("code_challenge_method", "S256")
	}

	if state != "" {
		params.Set("state", state)
	}

	return config.AuthURL + "?" + params.Encode()
}

// ExchangeCodeForTokens exchanges an authorization code for tokens.
func ExchangeCodeForTokens(config *OAuthConfig, code string, pkce *PKCE) (*OAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", config.RedirectURI)

	if pkce != nil {
		data.Set("code_verifier", pkce.Verifier)
	}

	resp, err := http.Post(config.TokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("token error: %s - %s", tokenResp.Error, tokenResp.ErrorDesc)
	}

	return &tokenResp, nil
}

// CallbackServer handles OAuth callbacks.
type CallbackServer struct {
	Port     int
	server   *http.Server
	codeChan chan string
	errChan  chan error

	// HTML templates for success/error pages
	SuccessHTML string
	ErrorHTML   string
}

// DefaultSuccessHTML is the default success page HTML.
const DefaultSuccessHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Authentication Successful</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
    }
    .container {
      text-align: center;
      padding: 3rem;
      background: rgba(255, 255, 255, 0.05);
      border-radius: 16px;
      max-width: 480px;
    }
    .check {
      width: 80px; height: 80px;
      background: #22c55e;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 1.5rem;
    }
    .check svg { width: 40px; height: 40px; }
    h1 { font-size: 1.75rem; margin-bottom: 0.75rem; }
    p { color: rgba(255, 255, 255, 0.7); }
    .hint { margin-top: 2rem; font-size: 0.875rem; color: rgba(255, 255, 255, 0.5); }
  </style>
</head>
<body>
  <div class="container">
    <div class="check">
      <svg fill="none" stroke="white" stroke-width="3" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
      </svg>
    </div>
    <h1>Authentication Successful</h1>
    <p>You can now close this window.</p>
    <p class="hint">Credentials have been saved.</p>
  </div>
</body>
</html>`

// DefaultErrorHTML is the default error page HTML (%s is replaced with error message).
const DefaultErrorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Authentication Failed</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
    }
    .container {
      text-align: center;
      padding: 3rem;
      background: rgba(255, 255, 255, 0.05);
      border-radius: 16px;
      max-width: 480px;
    }
    .error-icon {
      width: 80px; height: 80px;
      background: #ef4444;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin: 0 auto 1.5rem;
    }
    .error-icon svg { width: 40px; height: 40px; }
    h1 { font-size: 1.75rem; margin-bottom: 0.75rem; }
    p { color: rgba(255, 255, 255, 0.7); }
    .error-msg {
      margin-top: 1rem;
      padding: 1rem;
      background: rgba(239, 68, 68, 0.2);
      border-radius: 8px;
      font-family: monospace;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="error-icon">
      <svg fill="none" stroke="white" stroke-width="3" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </div>
    <h1>Authentication Failed</h1>
    <p>There was a problem authenticating.</p>
    <div class="error-msg">%s</div>
  </div>
</body>
</html>`

// NewCallbackServer creates a new OAuth callback server.
func NewCallbackServer(port int) *CallbackServer {
	if port == 0 {
		port = 8085
	}
	return &CallbackServer{
		Port:        port,
		codeChan:    make(chan string, 1),
		errChan:     make(chan error, 1),
		SuccessHTML: DefaultSuccessHTML,
		ErrorHTML:   DefaultErrorHTML,
	}
}

// Start starts the callback server.
func (s *CallbackServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, s.SuccessHTML)
			select {
			case s.codeChan <- code:
			default:
			}
		} else if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, s.ErrorHTML, errMsg)
			select {
			case s.errChan <- fmt.Errorf("OAuth error: %s", errMsg):
			default:
			}
		}
	})

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			select {
			case s.errChan <- fmt.Errorf("server error: %w", err):
			default:
			}
		}
	}()

	return nil
}

// WaitForCode waits for the authorization code or an error.
func (s *CallbackServer) WaitForCode(ctx context.Context) (string, error) {
	select {
	case code := <-s.codeChan:
		return code, nil
	case err := <-s.errChan:
		return "", err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Stop stops the callback server.
func (s *CallbackServer) Stop(ctx context.Context) {
	if s.server != nil {
		s.server.Shutdown(ctx)
	}
}

// RunOAuthFlow runs a complete OAuth flow with PKCE.
func RunOAuthFlow(config *OAuthConfig) (*OAuthResult, error) {
	// Set defaults
	if config.RedirectURI == "" {
		config.RedirectURI = fmt.Sprintf("http://localhost:%d/", config.CallbackPort)
	}
	if config.CallbackPort == 0 {
		config.CallbackPort = 8085
	}
	if config.Timeout == 0 {
		config.Timeout = 120 * time.Second
	}

	// Generate PKCE
	pkce, err := GeneratePKCE()
	if err != nil {
		return nil, err
	}

	// Build auth URL
	authURL := BuildAuthURL(config, pkce, "")

	// Start callback server
	server := NewCallbackServer(config.CallbackPort)
	if err := server.Start(); err != nil {
		return nil, err
	}
	defer server.Stop(context.Background())

	// Open browser
	if err := OpenURL(authURL); err != nil {
		fmt.Printf("Warning: Could not open browser: %v\n", err)
		fmt.Println("Please open this URL manually:")
		fmt.Println(authURL)
	}

	// Wait for code
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	code, err := server.WaitForCode(ctx)
	if err != nil {
		return nil, err
	}

	// Exchange for tokens
	tokenResp, err := ExchangeCodeForTokens(config, code, pkce)
	if err != nil {
		return nil, err
	}

	return &OAuthResult{
		AuthCode:     code,
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresIn:    tokenResp.ExpiresIn,
		TokenType:    tokenResp.TokenType,
	}, nil
}
