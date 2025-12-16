package main

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
	"github.com/joeblew999/ubuntu-website/internal/mcp"
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

func handleAuth(args []string) {
	if len(args) < 1 {
		printAuthUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "add":
		authAdd(cmdArgs)
	case "remove":
		authRemove(cmdArgs)
	case "status":
		authStatus(cmdArgs)
	case "check":
		authCheck(cmdArgs)
	case "guide":
		authGuide()
	case "open":
		authOpen(cmdArgs)
	case "login":
		authLogin(cmdArgs)
	default:
		fmt.Fprintf(os.Stderr, "Unknown auth command: %s\n", cmd)
		printAuthUsage()
		os.Exit(1)
	}
}

func authAdd(args []string) {
	location := "vscode"
	for _, a := range args {
		if strings.HasPrefix(a, "--location=") {
			location = strings.TrimPrefix(a, "--location=")
		}
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		exitError(fmt.Sprintf("Error finding project root: %v", err))
	}

	var mcpFile string
	switch location {
	case "vscode":
		targetDir := filepath.Join(projectRoot, ".vscode")
		os.MkdirAll(targetDir, 0755)
		mcpFile = filepath.Join(targetDir, "mcp.json")
	case "project":
		mcpFile = filepath.Join(projectRoot, ".mcp.json")
	case "claude":
		targetDir := filepath.Join(projectRoot, ".claude")
		os.MkdirAll(targetDir, 0755)
		mcpFile = filepath.Join(targetDir, "mcp.json")
	default:
		exitError(fmt.Sprintf("Unknown location: %s", location))
	}

	config, err := mcp.LoadConfig(mcpFile)
	if err != nil {
		exitError(err.Error())
	}

	if mcp.AddGoogleServer(config) {
		if err := mcp.SaveConfig(mcpFile, config); err != nil {
			exitError(err.Error())
		}
		fmt.Println("Added Google MCP server")
		fmt.Printf("   File: %s\n", mcpFile)
	} else {
		fmt.Printf("Google MCP server already configured in %s\n", mcpFile)
	}

	// Add permissions
	if err := mcp.EnsureClaudeDir(projectRoot); err != nil {
		exitError(err.Error())
	}
	settingsPath := mcp.GetSettingsPath(projectRoot)
	settings, err := mcp.LoadSettings(settingsPath)
	if err != nil {
		exitError(err.Error())
	}

	if mcp.AddGooglePermissions(settings) {
		if err := mcp.SaveSettings(settingsPath, settings); err != nil {
			exitError(err.Error())
		}
		fmt.Println("Added Google MCP permissions")
		fmt.Printf("   File: %s\n", settingsPath)
	}

	fmt.Println("")
	fmt.Println("Restart VSCode for changes to take effect.")
}

func authRemove(args []string) {
	location := "vscode"
	for _, a := range args {
		if strings.HasPrefix(a, "--location=") {
			location = strings.TrimPrefix(a, "--location=")
		}
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		exitError(fmt.Sprintf("Error finding project root: %v", err))
	}

	var mcpFile string
	switch location {
	case "vscode":
		mcpFile = filepath.Join(projectRoot, ".vscode", "mcp.json")
	case "project":
		mcpFile = filepath.Join(projectRoot, ".mcp.json")
	case "claude":
		mcpFile = filepath.Join(projectRoot, ".claude", "mcp.json")
	default:
		exitError(fmt.Sprintf("Unknown location: %s", location))
	}

	config, err := mcp.LoadConfig(mcpFile)
	if err != nil {
		exitError(err.Error())
	}

	if mcp.RemoveGoogleServer(config) {
		if len(config.MCPServers) == 0 {
			os.Remove(mcpFile)
			fmt.Println("Removed Google MCP server (deleted empty config)")
		} else {
			mcp.SaveConfig(mcpFile, config)
			fmt.Println("Removed Google MCP server")
		}
	} else {
		fmt.Println("Google MCP server not found")
	}

	// Remove permissions
	settingsPath := mcp.GetSettingsPath(projectRoot)
	settings, _ := mcp.LoadSettings(settingsPath)
	if settings != nil && mcp.RemoveGooglePermissions(settings) {
		mcp.SaveSettings(settingsPath, settings)
		fmt.Println("Removed Google MCP permissions")
	}
}

func authStatus(args []string) {
	location := "vscode"
	for _, a := range args {
		if strings.HasPrefix(a, "--location=") {
			location = strings.TrimPrefix(a, "--location=")
		}
	}

	projectRoot, _ := findProjectRoot()
	var mcpFile string
	switch location {
	case "vscode":
		mcpFile = filepath.Join(projectRoot, ".vscode", "mcp.json")
	case "project":
		mcpFile = filepath.Join(projectRoot, ".mcp.json")
	case "claude":
		mcpFile = filepath.Join(projectRoot, ".claude", "mcp.json")
	}

	fmt.Printf("=== Google MCP Status (%s) ===\n\n", location)

	config, err := mcp.LoadConfig(mcpFile)
	if err != nil {
		fmt.Printf("No config at %s\n", mcpFile)
		return
	}

	if config.HasServer(mcp.GoogleServerName) {
		fmt.Printf("MCP Server: configured in %s\n", mcpFile)
	} else {
		fmt.Println("MCP Server: not configured")
	}

	settingsPath := mcp.GetSettingsPath(projectRoot)
	settings, _ := mcp.LoadSettings(settingsPath)
	if settings != nil {
		count := mcp.CountGooglePermissions(settings)
		fmt.Printf("Permissions: %d google tools allowed\n", count)
	}
}

func authCheck(args []string) {
	fmt.Println("=== Google MCP Setup Check ===\n")

	// Check binary
	if _, err := exec.LookPath(mcp.GoogleServerCmd); err != nil {
		fmt.Println("google-mcp-server: NOT INSTALLED")
		fmt.Println("   Run: go install go.ngs.io/google-mcp-server@latest")
	} else {
		fmt.Println("google-mcp-server: installed")
	}

	// Check env vars
	if os.Getenv("GOOGLE_CLIENT_ID") == "" || os.Getenv("GOOGLE_CLIENT_SECRET") == "" {
		fmt.Println("Credentials: NOT SET")
		fmt.Println("   Add GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET to .env")
	} else {
		fmt.Println("Credentials: set in environment")
	}

	// Check accounts
	homeDir, _ := os.UserHomeDir()
	accounts, _ := filepath.Glob(filepath.Join(homeDir, mcp.GoogleAccountsDir, "*.json"))
	if len(accounts) == 0 {
		fmt.Println("Accounts: NOT AUTHENTICATED")
		fmt.Println("   Run: google-mcp-server")
	} else {
		fmt.Printf("Accounts: %d authenticated\n", len(accounts))
	}
}

func authGuide() {
	fmt.Println(`=== Google MCP Server Setup Guide ===

STEP 1: Create Google Cloud Project
  ` + urlProject + `

STEP 2: Configure OAuth Consent Screen
  ` + urlOAuthConsent + `
  - User Type: External
  - Add your email as test user

STEP 3: Enable APIs
  Gmail:    ` + urlGmailAPI + `
  Calendar: ` + urlCalendarAPI + `
  Drive:    ` + urlDriveAPI + `
  Sheets:   ` + urlSheetsAPI + `
  Docs:     ` + urlDocsAPI + `
  Slides:   ` + urlSlidesAPI + `

STEP 4: Create OAuth Credentials
  ` + urlCredentials + `
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

func authOpen(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: google auth open <target>")
		fmt.Println("\nTargets: project, consent, credentials, gmail, calendar, drive, sheets, docs, slides")
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
	url, ok := targets[target]
	if !ok {
		exitError(fmt.Sprintf("Unknown target: %s", target))
	}

	fmt.Printf("Opening: %s\n", url)
	browser.OpenURL(url)
}

func authLogin(args []string) {
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

	fmt.Println("=== Google OAuth Login ===\n")

	pkce, err := browser.GeneratePKCE()
	if err != nil {
		exitError(err.Error())
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

	fmt.Println("Opening browser for authentication...")
	fmt.Printf("Timeout: %ds\n\n", timeout)
	browser.OpenURL(authURL)

	select {
	case code := <-codeChan:
		server.Shutdown(ctx)
		fmt.Println("Received code, exchanging for tokens...")

		data := url.Values{}
		data.Set("client_id", gcloudClientID)
		data.Set("client_secret", gcloudClientSecret)
		data.Set("code", code)
		data.Set("code_verifier", pkce.Verifier)
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", oauthRedirectURI)

		resp, err := http.Post(oauthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
		if err != nil {
			exitError(err.Error())
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var tokenResp struct {
			RefreshToken string `json:"refresh_token"`
			Error        string `json:"error"`
		}
		json.Unmarshal(body, &tokenResp)

		if tokenResp.Error != "" {
			exitError(tokenResp.Error)
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

		fmt.Println("\nCredentials saved!")
		fmt.Printf("   File: %s\n", credFile)

	case err := <-errChan:
		exitError(err.Error())

	case <-ctx.Done():
		server.Shutdown(context.Background())
		exitError(fmt.Sprintf("Timeout after %d seconds", timeout))
	}
}

func printAuthUsage() {
	fmt.Println(`Usage: google auth <command> [arguments]

Commands:
  add [--location=LOC]      Add Google MCP server to Claude Code
  remove [--location=LOC]   Remove Google MCP server
  status [--location=LOC]   Show configuration status
  check                     Check setup status
  guide                     Show setup guide
  open <target>             Open Google Cloud Console page
  login [--account=EMAIL]   Authenticate with Google (OAuth flow)

Locations:
  vscode  - .vscode/mcp.json (default)
  project - .mcp.json
  claude  - .claude/mcp.json

Open targets:
  project, consent, credentials
  gmail, calendar, drive, sheets, docs, slides

Examples:
  google auth add
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
