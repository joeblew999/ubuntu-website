package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
	"github.com/joeblew999/ubuntu-website/internal/claude"
)

// OAuth constants for gcloud application-default credentials
const (
	gcloudClientID     = "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com"
	gcloudClientSecret = "d-FL95Q19q7MQmFpd7hHD0Ty"
	oauthRedirectURI   = "http://localhost:8085/"
	oauthTokenURL      = "https://oauth2.googleapis.com/token"
	oauthAuthURL       = "https://accounts.google.com/o/oauth2/v2/auth"
	defaultScopes      = "openid https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/sqlservice.login https://www.googleapis.com/auth/accounts.reauth"
)

// Google Cloud Console URLs
const (
	urlProject      = "https://console.cloud.google.com/projectcreate"
	urlOAuthConsent = "https://console.cloud.google.com/apis/credentials/consent"
	urlCredentials  = "https://console.cloud.google.com/apis/credentials"
	urlGmailAPI     = "https://console.cloud.google.com/apis/library/gmail.googleapis.com"
	urlCalendarAPI  = "https://console.cloud.google.com/apis/library/calendar-json.googleapis.com"
	urlDriveAPI     = "https://console.cloud.google.com/apis/library/drive.googleapis.com"
	urlSheetsAPI    = "https://console.cloud.google.com/apis/library/sheets.googleapis.com"
	urlDocsAPI      = "https://console.cloud.google.com/apis/library/docs.googleapis.com"
	urlSlidesAPI    = "https://console.cloud.google.com/apis/library/slides.googleapis.com"
)

func (c *cliContext) handleAuth(args []string) {
	if len(args) < 1 {
		c.printAuthUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "add":
		c.authAdd(cmdArgs)
	case "remove":
		c.authRemove(cmdArgs)
	case "status":
		c.authStatus(cmdArgs)
	case "check":
		c.authCheck(cmdArgs)
	case "guide":
		c.authGuide()
	case "open":
		c.authOpen(cmdArgs)
	case "login":
		c.authLogin(cmdArgs)
	default:
		fmt.Fprintf(c.stderr, "Unknown auth command: %s\n", cmd)
		c.printAuthUsage()
		os.Exit(1)
	}
}

// locationToTarget maps CLI --location flag values to Claude targets
func locationToTarget(location string) (claude.Target, error) {
	switch location {
	case "vscode":
		return claude.TargetVSCode, nil
	case "project":
		return claude.TargetProject, nil
	case "claude":
		return claude.TargetClaude, nil
	case "desktop":
		return claude.TargetDesktop, nil
	default:
		return "", fmt.Errorf("unknown location: %s (valid: vscode, project, claude, desktop)", location)
	}
}

func (c *cliContext) authAdd(args []string) {
	location := "vscode"
	for _, a := range args {
		if strings.HasPrefix(a, "--location=") {
			location = strings.TrimPrefix(a, "--location=")
		}
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		c.exitError(fmt.Sprintf("Error finding project root: %v", err))
	}

	target, err := locationToTarget(location)
	if err != nil {
		c.exitError(err.Error())
	}

	result, err := claude.AddMCPServer(claude.GoogleServerName, target, projectRoot)
	if err != nil {
		c.exitError(fmt.Sprintf("Error adding MCP server: %v", err))
	}

	// Show backup info for VSCode targets (safety feature)
	if result.BackupPath != "" {
		fmt.Fprintf(c.stdout, "Backup: %s\n", result.BackupPath)
	}

	if result.ServerAdded {
		fmt.Fprintln(c.stdout, "Added Google MCP server")
		fmt.Fprintf(c.stdout, "   File: %s\n", result.ConfigPath)
	} else {
		fmt.Fprintf(c.stdout, "Google MCP server already configured in %s\n", result.ConfigPath)
	}

	if result.PermissionsSet {
		fmt.Fprintln(c.stdout, "Added Google MCP permissions")
		fmt.Fprintf(c.stdout, "   File: %s\n", result.SettingsPath)
	}

	fmt.Fprintln(c.stdout, "")
	fmt.Fprintln(c.stdout, "Restart VSCode for changes to take effect.")
}

func (c *cliContext) authRemove(args []string) {
	location := "vscode"
	for _, a := range args {
		if strings.HasPrefix(a, "--location=") {
			location = strings.TrimPrefix(a, "--location=")
		}
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		c.exitError(fmt.Sprintf("Error finding project root: %v", err))
	}

	target, err := locationToTarget(location)
	if err != nil {
		c.exitError(err.Error())
	}

	result, err := claude.RemoveMCPServer(claude.GoogleServerName, target, projectRoot)
	if err != nil {
		c.exitError(fmt.Sprintf("Error removing MCP server: %v", err))
	}

	// Show backup info for VSCode targets (safety feature)
	if result.BackupPath != "" {
		fmt.Fprintf(c.stdout, "Backup: %s\n", result.BackupPath)
	}

	if result.ServerRemoved {
		if result.ConfigDeleted {
			fmt.Fprintln(c.stdout, "Removed Google MCP server (deleted empty config)")
		} else {
			fmt.Fprintln(c.stdout, "Removed Google MCP server")
		}
	} else {
		fmt.Fprintln(c.stdout, "Google MCP server not found")
	}

	if result.PermissionsRemoved {
		fmt.Fprintln(c.stdout, "Removed Google MCP permissions")
	}
}

func (c *cliContext) authStatus(args []string) {
	location := "vscode"
	for _, a := range args {
		if strings.HasPrefix(a, "--location=") {
			location = strings.TrimPrefix(a, "--location=")
		}
	}

	projectRoot, _ := findProjectRoot()

	target, err := locationToTarget(location)
	if err != nil {
		c.exitError(err.Error())
	}

	fmt.Fprintf(c.stdout, "=== Google MCP Status (%s) ===\n\n", location)

	status, err := claude.GetMCPServerStatus(claude.GoogleServerName, target, projectRoot)
	if err != nil {
		c.exitError(fmt.Sprintf("Error getting status: %v", err))
	}

	if status.Configured {
		fmt.Fprintf(c.stdout, "MCP Server: configured in %s\n", status.ConfigPath)
	} else {
		fmt.Fprintln(c.stdout, "MCP Server: not configured")
	}

	fmt.Fprintf(c.stdout, "Permissions: %d google tools allowed\n", status.PermissionCount)
}

func (c *cliContext) authCheck(args []string) {
	fmt.Fprintln(c.stdout, "=== Google MCP Setup Check ===")
	fmt.Fprintln(c.stdout)

	// Check binary
	if _, err := exec.LookPath(claude.GoogleServerCmd); err != nil {
		fmt.Fprintln(c.stdout, "google-mcp-server: NOT INSTALLED")
		fmt.Fprintln(c.stdout, "   Run: go install go.ngs.io/google-mcp-server@latest")
	} else {
		fmt.Fprintln(c.stdout, "google-mcp-server: installed")
	}

	// Check env vars
	if os.Getenv("GOOGLE_CLIENT_ID") == "" || os.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		fmt.Fprintln(c.stdout, "Credentials: NOT SET")
		fmt.Fprintln(c.stdout, "   Add GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET to .env")
	} else {
		fmt.Fprintln(c.stdout, "Credentials: set in environment")
	}

	// Check accounts
	homeDir, _ := os.UserHomeDir()
	accounts, _ := filepath.Glob(filepath.Join(homeDir, claude.GoogleAccountsDir, "*.json"))
	if len(accounts) == 0 {
		fmt.Fprintln(c.stdout, "Accounts: NOT AUTHENTICATED")
		fmt.Fprintln(c.stdout, "   Run: google-mcp-server")
	} else {
		fmt.Fprintf(c.stdout, "Accounts: %d authenticated\n", len(accounts))
	}
}

func (c *cliContext) authGuide() {
	fmt.Fprintln(c.stdout, `=== Google MCP Server Setup Guide ===

STEP 1: Create Google Cloud Project
  `+urlProject+`

STEP 2: Configure OAuth Consent Screen
  `+urlOAuthConsent+`
  - User Type: External
  - Add your email as test user

STEP 3: Enable APIs
  Gmail:    `+urlGmailAPI+`
  Calendar: `+urlCalendarAPI+`
  Drive:    `+urlDriveAPI+`
  Sheets:   `+urlSheetsAPI+`
  Docs:     `+urlDocsAPI+`
  Slides:   `+urlSlidesAPI+`

STEP 4: Create OAuth Credentials
  `+urlCredentials+`
  - CREATE CREDENTIALS > OAuth client ID
  - Application type: Desktop app
  - Copy Client ID and Secret

STEP 5: Save Credentials
  echo "GOOGLE_CLIENT_ID='your-id'" >> .env
  echo "GOOGLE_CLIENT_SECRET='your-secret'" >> .env
  source .env && google-mcp-server

STEP 6: Add to Claude Code
  google auth add

Run 'google auth check' to see current status.`)
}

func (c *cliContext) authOpen(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(c.stdout, "Usage: google auth open <target>")
		fmt.Fprintln(c.stdout, "\nTargets: project, consent, credentials, gmail, calendar, drive, sheets, docs, slides")
		return
	}

	targets := map[string]string{
		"project":     urlProject,
		"consent":     urlOAuthConsent,
		"credentials": urlCredentials,
		"gmail":       urlGmailAPI,
		"calendar":    urlCalendarAPI,
		"drive":       urlDriveAPI,
		"sheets":      urlSheetsAPI,
		"docs":        urlDocsAPI,
		"slides":      urlSlidesAPI,
	}

	target := args[0]
	targetURL, ok := targets[target]
	if !ok {
		c.exitError(fmt.Sprintf("Unknown target: %s", target))
	}

	fmt.Fprintf(c.stdout, "Opening: %s\n", targetURL)
	browser.OpenURL(targetURL)
}

func (c *cliContext) authLogin(args []string) {
	timeout := 120
	accountHint := ""
	for _, a := range args {
		if strings.HasPrefix(a, "--timeout=") {
			fmt.Sscanf(strings.TrimPrefix(a, "--timeout="), "%d", &timeout)
		}
		if strings.HasPrefix(a, "--account=") {
			accountHint = strings.TrimPrefix(a, "--account=")
		}
	}

	fmt.Fprintln(c.stdout, "=== Google OAuth Login ===")
	fmt.Fprintln(c.stdout)

	pkce, err := browser.GeneratePKCE()
	if err != nil {
		c.exitError(err.Error())
	}

	params := url.Values{}
	params.Set("client_id", gcloudClientID)
	params.Set("redirect_uri", oauthRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", defaultScopes)
	params.Set("code_challenge", pkce.Challenge)
	params.Set("code_challenge_method", "S256")
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")
	if accountHint != "" {
		params.Set("login_hint", accountHint)
	}

	authURL := oauthAuthURL + "?" + params.Encode()

	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	server := &http.Server{Addr: ":8085"}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, "<html><body><h1>Success!</h1><p>You can close this window.</p></body></html>")
			select {
			case codeChan <- code:
			default:
			}
		}
	})

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	fmt.Fprintln(c.stdout, "Opening browser for authentication...")
	fmt.Fprintf(c.stdout, "Timeout: %ds\n\n", timeout)
	browser.OpenURL(authURL)

	select {
	case code := <-codeChan:
		server.Shutdown(ctx)
		fmt.Fprintln(c.stdout, "Received code, exchanging for tokens...")

		data := url.Values{}
		data.Set("client_id", gcloudClientID)
		data.Set("client_secret", gcloudClientSecret)
		data.Set("code", code)
		data.Set("code_verifier", pkce.Verifier)
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", oauthRedirectURI)

		resp, err := http.Post(oauthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		if err != nil {
			c.exitError(err.Error())
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var tokenResp struct {
			RefreshToken string `json:"refresh_token"`
			Error        string `json:"error"`
		}
		json.Unmarshal(body, &tokenResp)

		if tokenResp.Error != "" {
			c.exitError(tokenResp.Error)
		}

		// Write credentials
		homeDir, _ := os.UserHomeDir()
		configDir := filepath.Join(homeDir, ".config", "gcloud")
		os.MkdirAll(configDir, 0755)

		creds := map[string]string{
			"client_id":     gcloudClientID,
			"client_secret": gcloudClientSecret,
			"refresh_token": tokenResp.RefreshToken,
			"type":          "authorized_user",
		}
		credsJSON, _ := json.MarshalIndent(creds, "", "  ")
		credFile := filepath.Join(configDir, "application_default_credentials.json")
		os.WriteFile(credFile, credsJSON, 0600)

		fmt.Fprintln(c.stdout, "\nCredentials saved!")
		fmt.Fprintf(c.stdout, "   File: %s\n", credFile)

	case err := <-errChan:
		c.exitError(err.Error())

	case <-ctx.Done():
		server.Shutdown(context.Background())
		c.exitError(fmt.Sprintf("Timeout after %d seconds", timeout))
	}
}

func (c *cliContext) printAuthUsage() {
	fmt.Fprintln(c.stdout, `Usage: google auth <command> [arguments]

Commands:
  add [--location=LOC]      Add Google MCP server to Claude
  remove [--location=LOC]   Remove Google MCP server
  status [--location=LOC]   Show configuration status
  check                     Check setup status
  guide                     Show setup guide
  open <target>             Open Google Cloud Console page
  login [--account=EMAIL]   Authenticate with Google (OAuth flow)

Locations:
  vscode  - .vscode/mcp.json (default, Claude Code)
  project - .mcp.json (project root)
  claude  - .claude/mcp.json (claude folder)
  desktop - Claude Desktop app config

Open targets:
  project, consent, credentials
  gmail, calendar, drive, sheets, docs, slides

Examples:
  google auth add
  google auth add --location=desktop
  google auth check
  google auth guide
  google auth open credentials
  google auth login --account=user@gmail.com`)
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return os.Getwd()
		}
		dir = parent
	}
}
