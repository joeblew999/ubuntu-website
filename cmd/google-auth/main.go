// Google MCP configuration tool for Claude Code
//
// Manages MCP server configuration in various locations.
//
// Usage:
//
//	google-auth add [-location=vscode|project|claude]
//	google-auth remove [-location=vscode|project|claude]
//	google-auth status [-location=vscode|project|claude]
//
// Locations:
//
//	vscode  - .vscode/mcp.json (default, VSCode extension)
//	project - .mcp.json (CLI project-level)
//	claude  - .claude/mcp.json (Claude folder)
package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// MCPConfig represents the mcp.json file structure
type MCPConfig struct {
	Schema     string               `json:"$schema,omitempty"`
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

const (
	schemaURL   = "https://modelcontextprotocol.io/schema/config.json"
	serverName  = "google"
	serverCmd   = "google-mcp-server"
	envClientID = "${GOOGLE_CLIENT_ID}"
	envSecretID = "${GOOGLE_CLIENT_SECRET}"
	// settings.json is used for project-level Claude permissions
	// VSCode extension reads this file (settings.local.json didn't work)
	settingsFile = "settings.json"
)

// ClaudeSettings represents the .claude/settings.json file structure
type ClaudeSettings struct {
	Permissions                Permissions       `json:"permissions"`
	EnableAllProjectMcpServers bool              `json:"enableAllProjectMcpServers"`
	Env                        map[string]string `json:"env,omitempty"`
}

// Permissions represents the permissions block in settings.json
type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

// Google MCP tool permissions - all tools must be listed individually
// The blanket "mcp__google" permission doesn't work consistently in VSCode
var googleMCPPermissions = []string{
	// Account management
	"mcp__google__accounts_list",
	"mcp__google__accounts_details",
	"mcp__google__accounts_add",
	"mcp__google__accounts_remove",
	"mcp__google__accounts_refresh",
	// Calendar
	"mcp__google__calendar_list",
	"mcp__google__calendar_events_list",
	"mcp__google__calendar_event_create",
	"mcp__google__calendar_events_list_all_accounts",
	// Drive
	"mcp__google__drive_files_list",
	"mcp__google__drive_files_search",
	"mcp__google__drive_file_download",
	"mcp__google__drive_file_upload",
	"mcp__google__drive_markdown_upload",
	"mcp__google__drive_markdown_replace",
	"mcp__google__drive_file_get_metadata",
	"mcp__google__drive_file_update_metadata",
	"mcp__google__drive_folder_create",
	"mcp__google__drive_file_move",
	"mcp__google__drive_file_copy",
	"mcp__google__drive_file_delete",
	"mcp__google__drive_file_trash",
	"mcp__google__drive_file_restore",
	"mcp__google__drive_shared_link_create",
	"mcp__google__drive_permissions_list",
	"mcp__google__drive_permissions_create",
	"mcp__google__drive_permissions_delete",
	"mcp__google__drive_files_list_all_accounts",
	// Gmail
	"mcp__google__gmail_messages_list",
	"mcp__google__gmail_message_get",
	"mcp__google__gmail_messages_list_all_accounts",
	// Sheets
	"mcp__google__sheets_spreadsheet_get",
	"mcp__google__sheets_values_get",
	"mcp__google__sheets_values_update",
	// Docs
	"mcp__google__docs_document_get",
	"mcp__google__docs_document_create",
	"mcp__google__docs_document_update",
	// Slides
	"mcp__google__slides_presentation_create",
	"mcp__google__slides_presentation_get",
	"mcp__google__slides_slide_create",
	"mcp__google__slides_slide_delete",
	"mcp__google__slides_slide_duplicate",
	"mcp__google__slides_markdown_create",
	"mcp__google__slides_markdown_update",
	"mcp__google__slides_markdown_append",
	"mcp__google__slides_add_text",
	"mcp__google__slides_add_image",
	"mcp__google__slides_add_table",
	"mcp__google__slides_add_shape",
	"mcp__google__slides_set_layout",
	"mcp__google__slides_export_pdf",
	"mcp__google__slides_share",
	"mcp__google__slides_presentations_list_all_accounts",
}

// Location configs
var locations = map[string]struct {
	dir  string // subdirectory (empty = project root)
	file string // filename
}{
	"vscode":  {dir: ".vscode", file: "mcp.json"},
	"project": {dir: "", file: ".mcp.json"},
	"claude":  {dir: ".claude", file: "mcp.json"},
}

func main() {
	// Define flags
	location := flag.String("location", "vscode", "Config location: vscode, project, or claude")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	cmd := args[0]

	// Validate location
	loc, ok := locations[*location]
	if !ok {
		fmt.Printf("Unknown location: %s\n", *location)
		fmt.Println("Valid locations: vscode, project, claude")
		os.Exit(1)
	}

	// Find project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Printf("Error finding project root: %v\n", err)
		os.Exit(1)
	}

	// Build paths
	var targetDir, mcpFile string
	if loc.dir != "" {
		targetDir = filepath.Join(projectRoot, loc.dir)
		mcpFile = filepath.Join(targetDir, loc.file)
	} else {
		targetDir = projectRoot
		mcpFile = filepath.Join(projectRoot, loc.file)
	}

	switch cmd {
	case "add":
		if err := addServer(targetDir, mcpFile, loc.dir != ""); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "remove":
		if err := removeServer(mcpFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "status":
		if err := showStatus(mcpFile, *location); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "guide":
		showSetupGuide()
	case "check":
		checkSetup(mcpFile)
	case "open":
		if len(args) < 2 {
			printOpenUsage()
			os.Exit(1)
		}
		openURL(args[1])
	case "gcloud-auth":
		// Parse gcloud-auth specific flags
		gcloudFlags := flag.NewFlagSet("gcloud-auth", flag.ExitOnError)
		serverOnly := gcloudFlags.Bool("server-only", false, "Only start the callback server (for Playwright automation)")
		autoMode := gcloudFlags.Bool("auto", false, "Fully automated mode using embedded Playwright")
		assistedMode := gcloudFlags.Bool("assisted", false, "Assisted mode: Claude navigates UI, user handles passkey")
		headless := gcloudFlags.Bool("headless", false, "Run browser in headless mode (with -auto)")
		accountHint := gcloudFlags.String("account", "", "Google account email to pre-select (with -assisted)")
		timeout := gcloudFlags.Int("timeout", 120, "Timeout in seconds for auth flow")
		gcloudFlags.Parse(args[1:])

		if *autoMode {
			if err := runGcloudAuthAuto(*headless, *timeout); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else if *assistedMode {
			if err := runGcloudAuthAssisted(*accountHint, *timeout); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		} else if err := runGcloudAuth(*serverOnly, *timeout); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// Google Cloud Console URLs - single source of truth for all Google URLs
// The Taskfile open:* tasks call this binary instead of hardcoding URLs
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
	urlGitHubRepo   = "https://github.com/ngs/google-mcp-server"
	accountsDir     = ".google-mcp-accounts"
)

// OAuth constants for gcloud application-default credentials
const (
	// These are gcloud's official OAuth client credentials (public, not secret)
	gcloudClientID     = "764086051850-6qr4p6gpi6hn506pt8ejuq83di341hur.apps.googleusercontent.com"
	gcloudClientSecret = "d-FL95Q19q7MQmFpd7hHD0Ty"
	oauthRedirectURI   = "http://localhost:8085/"
	oauthTokenURL      = "https://oauth2.googleapis.com/token"
	oauthAuthURL       = "https://accounts.google.com/o/oauth2/v2/auth"

	// Default scopes for application-default credentials
	defaultScopes = "openid https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/cloud-platform https://www.googleapis.com/auth/sqlservice.login https://www.googleapis.com/auth/accounts.reauth"

	// OAuth callback HTML templates
	oauthSuccessHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Authentication Successful - Ubuntu Software</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
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
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      max-width: 480px;
    }
    .logo { width: 180px; height: auto; margin-bottom: 2rem; }
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
    h1 { font-size: 1.75rem; font-weight: 600; margin-bottom: 0.75rem; }
    p { color: rgba(255, 255, 255, 0.7); font-size: 1rem; line-height: 1.5; }
    .hint { margin-top: 2rem; font-size: 0.875rem; color: rgba(255, 255, 255, 0.5); }
  </style>
</head>
<body>
  <div class="container">
    <img src="https://www.ubuntusoftware.net/images/logo.svg" alt="Ubuntu Software" class="logo">
    <div class="check">
      <svg fill="none" stroke="white" stroke-width="3" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7"/>
      </svg>
    </div>
    <h1>Authentication Successful</h1>
    <p>Your Google Cloud credentials have been configured. You can now use Terraform and other Google Cloud tools.</p>
    <p class="hint">You may close this window.</p>
  </div>
</body>
</html>`

	oauthErrorHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Authentication Failed - Ubuntu Software</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
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
      backdrop-filter: blur(10px);
      border: 1px solid rgba(255, 255, 255, 0.1);
      max-width: 480px;
    }
    .logo { width: 180px; height: auto; margin-bottom: 2rem; }
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
    h1 { font-size: 1.75rem; font-weight: 600; margin-bottom: 0.75rem; }
    p { color: rgba(255, 255, 255, 0.7); font-size: 1rem; line-height: 1.5; }
    .error-msg {
      margin-top: 1rem;
      padding: 1rem;
      background: rgba(239, 68, 68, 0.2);
      border-radius: 8px;
      font-family: monospace;
      font-size: 0.875rem;
    }
    .hint { margin-top: 2rem; font-size: 0.875rem; color: rgba(255, 255, 255, 0.5); }
  </style>
</head>
<body>
  <div class="container">
    <img src="https://www.ubuntusoftware.net/images/logo.svg" alt="Ubuntu Software" class="logo">
    <div class="error-icon">
      <svg fill="none" stroke="white" stroke-width="3" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12"/>
      </svg>
    </div>
    <h1>Authentication Failed</h1>
    <p>There was a problem authenticating with Google Cloud.</p>
    <div class="error-msg">%s</div>
    <p class="hint">Please try again or check your Google Cloud configuration.</p>
  </div>
</body>
</html>`
)

// openTargets maps target names to URLs and descriptions
var openTargets = map[string]struct {
	url  string
	desc string
}{
	"project":     {urlProject, "Create or select a Google Cloud project"},
	"consent":     {urlOAuthConsent, "Configure OAuth consent screen"},
	"credentials": {urlCredentials, "Create OAuth credentials"},
	"apis":        {urlCredentials, "Enable required Google APIs"}, // special handling
	"gmail":       {urlGmailAPI, "Gmail API"},
	"calendar":    {urlCalendarAPI, "Calendar API"},
	"drive":       {urlDriveAPI, "Drive API"},
	"sheets":      {urlSheetsAPI, "Sheets API"},
	"docs":        {urlDocsAPI, "Docs API"},
	"slides":      {urlSlidesAPI, "Slides API"},
	"repo":        {urlGitHubRepo, "GitHub repository"},
}

// showSetupGuide prints the full setup guide
func showSetupGuide() {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║         Google MCP Server - Setup Guide                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Println("This MCP server gives Claude access to:")
	fmt.Println("  • Gmail      - read messages")
	fmt.Println("  • Calendar   - list/create/update events")
	fmt.Println("  • Drive      - list/upload/download files")
	fmt.Println("  • Sheets     - read/update spreadsheets")
	fmt.Println("  • Docs       - create/update documents")
	fmt.Println("  • Slides     - create presentations")
	fmt.Println("")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("STEP 1: Create Google Cloud Project")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("  " + urlProject)
	fmt.Println("")
	fmt.Println("  • Give it any name (e.g., 'Claude MCP')")
	fmt.Println("  • Click CREATE")
	fmt.Println("")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("STEP 2: Configure OAuth Consent Screen")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("  " + urlOAuthConsent)
	fmt.Println("")
	fmt.Println("  • User Type: External → CREATE")
	fmt.Println("  • App name: anything (e.g., 'Claude MCP')")
	fmt.Println("  • User support email: your email")
	fmt.Println("  • Developer contact: your email")
	fmt.Println("  • Click SAVE AND CONTINUE through all screens")
	fmt.Println("  • On 'Test users' screen: ADD USERS → add your email")
	fmt.Println("")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("STEP 3: Enable APIs (click ENABLE on each page)")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("  Gmail:    " + urlGmailAPI)
	fmt.Println("  Calendar: " + urlCalendarAPI)
	fmt.Println("  Drive:    " + urlDriveAPI)
	fmt.Println("  Sheets:   " + urlSheetsAPI)
	fmt.Println("  Docs:     " + urlDocsAPI)
	fmt.Println("  Slides:   " + urlSlidesAPI)
	fmt.Println("")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("STEP 4: Create OAuth Credentials")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("  " + urlCredentials)
	fmt.Println("")
	fmt.Println("  • Click CREATE CREDENTIALS → OAuth client ID")
	fmt.Println("  • Application type: Desktop app")
	fmt.Println("  • Name: anything")
	fmt.Println("  • Click CREATE")
	fmt.Println("  • COPY the Client ID and Client Secret!")
	fmt.Println("")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("STEP 5: Save Credentials & Authenticate")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("  Run these commands:")
	fmt.Println("")
	fmt.Println("    # Add to .env file:")
	fmt.Println("    echo \"GOOGLE_CLIENT_ID='your-client-id'\" >> .env")
	fmt.Println("    echo \"GOOGLE_CLIENT_SECRET='your-secret'\" >> .env")
	fmt.Println("")
	fmt.Println("    # Load and authenticate:")
	fmt.Println("    source .env && google-mcp-server")
	fmt.Println("")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("STEP 6: Add to Claude Code")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("    google-auth add")
	fmt.Println("")
	fmt.Println("    Then FULLY RESTART VSCode (quit and reopen)")
	fmt.Println("")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("Run 'google-auth check' to see what's already configured.")
	fmt.Println("════════════════════════════════════════════════════════════════")
}

// checkSetup checks what's configured and shows next step
func checkSetup(mcpFile string) {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║         Google MCP Server - Setup Check                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")

	allGood := true
	nextStep := ""

	// Check 1: google-mcp-server binary
	_, err := exec.LookPath("google-mcp-server")
	if err != nil {
		fmt.Println("❌ google-mcp-server binary: NOT INSTALLED")
		fmt.Println("   Run: go install go.ngs.io/google-mcp-server@latest")
		allGood = false
		if nextStep == "" {
			nextStep = "go install go.ngs.io/google-mcp-server@latest"
		}
	} else {
		fmt.Println("✅ google-mcp-server binary: installed")
	}

	// Check 2: Environment variables
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		fmt.Println("❌ Credentials: NOT SET in environment")
		fmt.Println("   Add to .env and run: source .env")
		allGood = false
		if nextStep == "" {
			nextStep = "google-auth guide  # then follow Steps 1-4"
		}
	} else {
		fmt.Println("✅ Credentials: set in environment")
	}

	// Check 3: Authenticated accounts
	homeDir, _ := os.UserHomeDir()
	accountsPath := filepath.Join(homeDir, accountsDir)
	accounts, _ := filepath.Glob(filepath.Join(accountsPath, "*.json"))
	if len(accounts) == 0 {
		fmt.Println("❌ Google account: NOT AUTHENTICATED")
		fmt.Println("   Run: google-mcp-server")
		allGood = false
		if nextStep == "" {
			nextStep = "source .env && google-mcp-server"
		}
	} else {
		fmt.Printf("✅ Google account: %d authenticated\n", len(accounts))
		for _, acc := range accounts {
			name := filepath.Base(acc)
			name = name[:len(name)-5] // remove .json
			fmt.Printf("   • %s\n", name)
		}
	}

	// Check 4: MCP config
	if _, err := os.Stat(mcpFile); os.IsNotExist(err) {
		fmt.Println("❌ MCP config: NOT CONFIGURED")
		fmt.Println("   Run: google-auth add")
		allGood = false
		if nextStep == "" {
			nextStep = "google-auth add"
		}
	} else {
		config, err := loadConfig(mcpFile)
		if err == nil {
			if _, exists := config.MCPServers[serverName]; exists {
				fmt.Printf("✅ MCP config: %s\n", mcpFile)
			} else {
				fmt.Printf("⚠️  MCP config exists but no google server: %s\n", mcpFile)
				fmt.Println("   Run: google-auth add")
				allGood = false
				if nextStep == "" {
					nextStep = "google-auth add"
				}
			}
		}
	}

	// Check 5: Permissions
	projectRoot, _ := findProjectRoot()
	settingsPath := filepath.Join(projectRoot, ".claude", settingsFile)
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		fmt.Println("❌ Permissions: NOT CONFIGURED")
		fmt.Println("   Run: google-auth add")
		allGood = false
		if nextStep == "" {
			nextStep = "google-auth add"
		}
	} else {
		settings, err := loadSettings(settingsPath)
		if err == nil {
			googlePermCount := 0
			for _, p := range settings.Permissions.Allow {
				if len(p) > 12 && p[:12] == "mcp__google_" {
					googlePermCount++
				}
			}
			if googlePermCount > 0 {
				fmt.Printf("✅ Permissions: %d google tools allowed\n", googlePermCount)
			} else {
				fmt.Println("⚠️  Permissions file exists but no google permissions")
				fmt.Println("   Run: google-auth add")
				allGood = false
				if nextStep == "" {
					nextStep = "google-auth add"
				}
			}
		}
	}

	fmt.Println("")
	if allGood {
		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("✅ ALL CONFIGURED! Ready to use.")
		fmt.Println("")
		fmt.Println("If Claude still prompts for permissions, FULLY RESTART VSCode.")
		fmt.Println("════════════════════════════════════════════════════════════════")
	} else {
		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("NEXT STEP:")
		fmt.Printf("  %s\n", nextStep)
		fmt.Println("")
		fmt.Println("For full guide: google-auth guide")
		fmt.Println("════════════════════════════════════════════════════════════════")
	}
}

func printUsage() {
	fmt.Println("Google MCP Configuration Tool")
	fmt.Println("")
	fmt.Println("Manages MCP server configuration for Claude Code.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  google-auth guide                    Show full setup guide")
	fmt.Println("  google-auth check                    Check setup status & next step")
	fmt.Println("  google-auth open <target>            Open Google Cloud Console page")
	fmt.Println("  google-auth [-location=LOC] add      Add google server + permissions")
	fmt.Println("  google-auth [-location=LOC] remove   Remove google server + permissions")
	fmt.Println("  google-auth [-location=LOC] status   Show MCP config status")
	fmt.Println("  google-auth gcloud-auth              Authenticate gcloud (OAuth flow)")
	fmt.Println("")
	fmt.Println("gcloud-auth modes:")
	fmt.Println("  (default)       Opens browser, you complete auth manually")
	fmt.Println("  -assisted       Passkey mode: opens default browser (recommended for passkey users)")
	fmt.Println("  -auto           Fully automated with embedded Playwright (password accounts only)")
	fmt.Println("")
	fmt.Println("gcloud-auth options:")
	fmt.Println("  -account=EMAIL  Pre-select Google account (skips account selection screen)")
	fmt.Println("  -headless       Run browser headless (with -auto)")
	fmt.Println("  -server-only    Only start callback server (for external automation)")
	fmt.Println("  -timeout=120    Timeout in seconds (default: 120)")
	fmt.Println("")
	fmt.Println("Locations (-location flag):")
	fmt.Println("  vscode  - .vscode/mcp.json (default, for VSCode extension)")
	fmt.Println("  project - .mcp.json (CLI project-level)")
	fmt.Println("  claude  - .claude/mcp.json (Claude folder)")
	fmt.Println("")
	fmt.Println("Open targets:")
	fmt.Println("  project, consent, credentials, apis, repo")
	fmt.Println("  gmail, calendar, drive, sheets, docs, slides")
	fmt.Println("")
	fmt.Println("Quick start:")
	fmt.Println("  google-auth guide       # Follow the setup steps")
	fmt.Println("  google-auth check       # See what's configured")
	fmt.Println("  google-auth gcloud-auth # Authenticate for Terraform")
	fmt.Println("  google-auth add         # Add to Claude Code")
}

// printOpenUsage shows help for the open command
func printOpenUsage() {
	fmt.Println("Usage: google-auth open <target>")
	fmt.Println("")
	fmt.Println("Opens Google Cloud Console pages in your browser.")
	fmt.Println("")
	fmt.Println("Targets:")
	fmt.Println("  project      Create or select a Google Cloud project")
	fmt.Println("  consent      Configure OAuth consent screen")
	fmt.Println("  credentials  Create OAuth credentials")
	fmt.Println("  apis         Enable all required Google APIs")
	fmt.Println("  repo         GitHub repository for google-mcp-server")
	fmt.Println("")
	fmt.Println("Individual APIs:")
	fmt.Println("  gmail, calendar, drive, sheets, docs, slides")
}

// openURL opens a URL target in the browser
func openURL(target string) {
	// Special case: "apis" opens all API pages
	if target == "apis" {
		openAllAPIs()
		return
	}

	t, ok := openTargets[target]
	if !ok {
		fmt.Printf("Unknown target: %s\n", target)
		fmt.Println("")
		printOpenUsage()
		os.Exit(1)
	}

	fmt.Printf("Opening: %s\n", t.desc)
	fmt.Printf("URL: %s\n", t.url)

	// Print contextual instructions
	switch target {
	case "project":
		fmt.Println("")
		fmt.Println("Instructions:")
		fmt.Println("  • Give your project any name (e.g., 'Claude MCP')")
		fmt.Println("  • Click CREATE")
		fmt.Println("  • Wait for project to be created")
		fmt.Println("")
		fmt.Println("Next: google-auth open consent")
	case "consent":
		fmt.Println("")
		fmt.Println("Instructions:")
		fmt.Println("  1. User Type: External → Create")
		fmt.Println("  2. App name: Google MCP (or any name)")
		fmt.Println("  3. User support email: your email")
		fmt.Println("  4. Developer contact: your email")
		fmt.Println("  5. Save and Continue through remaining screens")
		fmt.Println("  6. On 'Test users': ADD USERS → add your email")
		fmt.Println("")
		fmt.Println("Next: google-auth open apis")
	case "credentials":
		fmt.Println("")
		fmt.Println("Instructions:")
		fmt.Println("  1. Click 'CREATE CREDENTIALS' → 'OAuth client ID'")
		fmt.Println("  2. Application type: Desktop app")
		fmt.Println("  3. Name: Google MCP Client (or any name)")
		fmt.Println("  4. Click CREATE")
		fmt.Println("  5. Copy the Client ID and Client Secret")
		fmt.Println("")
		fmt.Println("Then save credentials:")
		fmt.Println("  task google-mcp:credentials CLIENT_ID='xxx' CLIENT_SECRET='xxx'")
	}

	// Open in browser
	if err := exec.Command("open", t.url).Start(); err != nil {
		fmt.Printf("Failed to open browser: %v\n", err)
		fmt.Printf("Please open manually: %s\n", t.url)
	}
}

// openAllAPIs opens required Google API pages one at a time
func openAllAPIs() {
	fmt.Println("Enable each Google API one at a time.")
	fmt.Println("Click ENABLE on each page, then press Enter to continue.")
	fmt.Println("")

	apis := []string{"gmail", "calendar", "drive", "sheets", "docs", "slides"}
	for i, api := range apis {
		t := openTargets[api]
		fmt.Printf("[%d/%d] Opening %s API...\n", i+1, len(apis), t.desc)
		fmt.Printf("       %s\n", t.url)

		if err := exec.Command("open", t.url).Start(); err != nil {
			fmt.Printf("  Failed to open: %v\n", err)
			fmt.Printf("  Open manually: %s\n", t.url)
		}

		if i < len(apis)-1 {
			fmt.Print("\nPress Enter after enabling...")
			fmt.Scanln()
		}
	}

	fmt.Println("")
	fmt.Println("✅ All APIs enabled!")
	fmt.Println("")
	fmt.Println("Next: google-auth open credentials")
}

// findProjectRoot looks for the project root by finding go.mod or .git
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check for go.mod or .git
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, use current working directory
			return os.Getwd()
		}
		dir = parent
	}
}

// loadConfig loads the mcp.json file, returning empty config if not exists
func loadConfig(mcpFile string) (*MCPConfig, error) {
	config := &MCPConfig{
		Schema:     schemaURL,
		MCPServers: make(map[string]MCPServer),
	}

	data, err := os.ReadFile(mcpFile)
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", mcpFile, err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", mcpFile, err)
	}

	if config.MCPServers == nil {
		config.MCPServers = make(map[string]MCPServer)
	}

	// Ensure schema is set
	if config.Schema == "" {
		config.Schema = schemaURL
	}

	return config, nil
}

// saveConfig saves the config to mcp.json with pretty formatting
func saveConfig(mcpFile string, config *MCPConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	if err := os.WriteFile(mcpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", mcpFile, err)
	}

	return nil
}

// addServer adds the google MCP server to mcp.json and permissions to settings.json
func addServer(targetDir, mcpFile string, createDir bool) error {
	// Ensure target directory exists if needed
	if createDir {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	config, err := loadConfig(mcpFile)
	if err != nil {
		return err
	}

	// Check if already exists
	mcpExists := false
	if _, exists := config.MCPServers[serverName]; exists {
		mcpExists = true
	}

	if !mcpExists {
		// Add the server
		config.MCPServers[serverName] = MCPServer{
			Command: serverCmd,
			Args:    []string{},
			Env: map[string]string{
				"GOOGLE_CLIENT_ID":     envClientID,
				"GOOGLE_CLIENT_SECRET": envSecretID,
			},
		}

		if err := saveConfig(mcpFile, config); err != nil {
			return err
		}
		fmt.Println("✅ Google MCP server added")
		fmt.Printf("   File: %s\n", mcpFile)
	} else {
		fmt.Printf("✅ Google MCP server already configured in %s\n", mcpFile)
	}

	// Now handle settings.json for permissions (always in .claude/)
	projectRoot, err := findProjectRoot()
	if err != nil {
		return err
	}
	claudeDir := filepath.Join(projectRoot, ".claude")
	settingsPath := filepath.Join(claudeDir, settingsFile)

	// Ensure .claude directory exists
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	// Add permissions
	permissionsAdded, err := addPermissions(settingsPath)
	if err != nil {
		return err
	}

	if permissionsAdded {
		fmt.Println("✅ Google MCP permissions added")
		fmt.Printf("   File: %s\n", settingsPath)
	} else {
		fmt.Printf("✅ Google MCP permissions already configured in %s\n", settingsPath)
	}

	fmt.Println("")
	fmt.Println("The config uses ${GOOGLE_CLIENT_ID} and ${GOOGLE_CLIENT_SECRET}")
	fmt.Println("environment variables - no secrets stored in the file.")
	fmt.Println("")
	fmt.Println("Make sure your .env file has these variables set.")
	fmt.Println("")
	fmt.Println("⚠️  IMPORTANT: You must FULLY RESTART VSCode for permissions to take effect.")
	fmt.Println("   (Window reload is not enough - quit and reopen VSCode)")

	return nil
}

// loadSettings loads the settings.json file, returning empty settings if not exists
func loadSettings(settingsPath string) (*ClaudeSettings, error) {
	settings := &ClaudeSettings{
		Permissions: Permissions{
			Allow: []string{},
			Deny:  []string{},
		},
		EnableAllProjectMcpServers: true,
	}

	data, err := os.ReadFile(settingsPath)
	if os.IsNotExist(err) {
		return settings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", settingsPath, err)
	}

	if err := json.Unmarshal(data, settings); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", settingsPath, err)
	}

	if settings.Permissions.Allow == nil {
		settings.Permissions.Allow = []string{}
	}
	if settings.Permissions.Deny == nil {
		settings.Permissions.Deny = []string{}
	}

	return settings, nil
}

// saveSettings saves the settings to settings.json with pretty formatting
func saveSettings(settingsPath string, settings *ClaudeSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", settingsPath, err)
	}

	return nil
}

// addPermissions adds Google MCP permissions to settings.json
// Returns true if permissions were added, false if already present
func addPermissions(settingsPath string) (bool, error) {
	settings, err := loadSettings(settingsPath)
	if err != nil {
		return false, err
	}

	// Check which permissions are missing
	existingPerms := make(map[string]bool)
	for _, p := range settings.Permissions.Allow {
		existingPerms[p] = true
	}

	added := false
	for _, perm := range googleMCPPermissions {
		if !existingPerms[perm] {
			settings.Permissions.Allow = append(settings.Permissions.Allow, perm)
			added = true
		}
	}

	if !added {
		return false, nil
	}

	// Ensure enableAllProjectMcpServers is true
	settings.EnableAllProjectMcpServers = true

	if err := saveSettings(settingsPath, settings); err != nil {
		return false, err
	}

	return true, nil
}

// removePermissions removes Google MCP permissions from settings.json
// Returns true if permissions were removed, false if not present
func removePermissions(settingsPath string) (bool, error) {
	settings, err := loadSettings(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	// Build set of Google permissions to remove
	googlePerms := make(map[string]bool)
	for _, p := range googleMCPPermissions {
		googlePerms[p] = true
	}

	// Filter out Google permissions
	newAllow := []string{}
	removed := false
	for _, p := range settings.Permissions.Allow {
		if googlePerms[p] {
			removed = true
		} else {
			newAllow = append(newAllow, p)
		}
	}

	if !removed {
		return false, nil
	}

	settings.Permissions.Allow = newAllow

	if err := saveSettings(settingsPath, settings); err != nil {
		return false, err
	}

	return true, nil
}

// removeServer removes the google MCP server from mcp.json and permissions from settings.json
func removeServer(mcpFile string) error {
	config, err := loadConfig(mcpFile)
	if err != nil {
		return err
	}

	mcpRemoved := false
	// Check if exists
	if _, exists := config.MCPServers[serverName]; exists {
		// Remove the server
		delete(config.MCPServers, serverName)
		mcpRemoved = true

		// If no servers left, remove the file
		if len(config.MCPServers) == 0 {
			if err := os.Remove(mcpFile); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove %s: %w", mcpFile, err)
			}
			fmt.Println("✅ Google MCP server removed")
			fmt.Printf("   Deleted %s (no servers remaining)\n", mcpFile)
		} else {
			if err := saveConfig(mcpFile, config); err != nil {
				return err
			}
			fmt.Println("✅ Google MCP server removed")
			fmt.Printf("   File: %s\n", mcpFile)
		}
	} else {
		fmt.Printf("ℹ️  Google MCP server not found in %s\n", mcpFile)
	}

	// Also remove permissions from .claude/settings.json
	projectRoot, err := findProjectRoot()
	if err != nil {
		return err
	}
	settingsPath := filepath.Join(projectRoot, ".claude", settingsFile)

	permissionsRemoved, err := removePermissions(settingsPath)
	if err != nil {
		return err
	}

	if permissionsRemoved {
		fmt.Println("✅ Google MCP permissions removed")
		fmt.Printf("   File: %s\n", settingsPath)
	} else if mcpRemoved {
		fmt.Println("ℹ️  No Google MCP permissions found to remove")
	}

	if mcpRemoved || permissionsRemoved {
		fmt.Println("")
		fmt.Println("⚠️  IMPORTANT: You must FULLY RESTART VSCode for changes to take effect.")
		fmt.Println("   (Window reload is not enough - quit and reopen VSCode)")
	}

	return nil
}

// showStatus shows the current mcp.json configuration
func showStatus(mcpFile, location string) error {
	fmt.Println("=== Google MCP Configuration Status ===")
	fmt.Printf("Location: %s\n", location)
	fmt.Println("")

	// Check if file exists
	if _, err := os.Stat(mcpFile); os.IsNotExist(err) {
		fmt.Printf("❌ No config found at %s\n", mcpFile)
		fmt.Println("")
		fmt.Printf("Run: google-auth -location=%s add\n", location)
		return nil
	}

	config, err := loadConfig(mcpFile)
	if err != nil {
		return err
	}

	fmt.Printf("File: %s\n", mcpFile)
	fmt.Println("")

	// Check for google server
	if server, exists := config.MCPServers[serverName]; exists {
		fmt.Println("✅ Google MCP server configured:")
		fmt.Printf("   Command: %s\n", server.Command)
		if len(server.Env) > 0 {
			fmt.Println("   Environment:")
			for k, v := range server.Env {
				fmt.Printf("     %s: %s\n", k, v)
			}
		}
	} else {
		fmt.Println("❌ Google MCP server not configured")
		fmt.Println("")
		fmt.Printf("Run: google-auth -location=%s add\n", location)
	}

	// Show other servers
	otherCount := 0
	for name := range config.MCPServers {
		if name != serverName {
			otherCount++
		}
	}
	if otherCount > 0 {
		fmt.Println("")
		fmt.Printf("Other MCP servers in config: %d\n", otherCount)
		for name := range config.MCPServers {
			if name != serverName {
				fmt.Printf("   - %s\n", name)
			}
		}
	}

	return nil
}

// ============================================================================
// gcloud-auth: Automated OAuth flow for gcloud application-default credentials
// ============================================================================

// GcloudCredentials represents the application_default_credentials.json format
type GcloudCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token"`
	Type         string `json:"type"`
}

// OAuthTokenResponse represents the token endpoint response
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	Error        string `json:"error,omitempty"`
	ErrorDesc    string `json:"error_description,omitempty"`
}

// generatePKCE generates code_verifier and code_challenge for PKCE
func generatePKCE() (verifier, challenge string, err error) {
	// Generate 32 random bytes for verifier
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Base64 URL encode without padding
	verifier = base64.RawURLEncoding.EncodeToString(b)

	// SHA256 hash the verifier
	h := sha256.Sum256([]byte(verifier))

	// Base64 URL encode the hash without padding
	challenge = base64.RawURLEncoding.EncodeToString(h[:])

	return verifier, challenge, nil
}

// buildAuthURL constructs the OAuth authorization URL with PKCE
func buildAuthURL(codeChallenge string) string {
	params := url.Values{}
	params.Set("client_id", gcloudClientID)
	params.Set("redirect_uri", oauthRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", defaultScopes)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	params.Set("access_type", "offline")
	params.Set("prompt", "consent")

	return oauthAuthURL + "?" + params.Encode()
}

// exchangeCodeForTokens exchanges the authorization code for tokens
func exchangeCodeForTokens(code, codeVerifier string) (*OAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", gcloudClientID)
	data.Set("client_secret", gcloudClientSecret)
	data.Set("code", code)
	data.Set("code_verifier", codeVerifier)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", oauthRedirectURI)

	resp, err := http.Post(oauthTokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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

// writeGcloudCredentials writes the credentials to the gcloud config directory
func writeGcloudCredentials(refreshToken string) error {
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
		ClientID:     gcloudClientID,
		ClientSecret: gcloudClientSecret,
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

// runGcloudAuth runs the gcloud application-default login flow
func runGcloudAuth(serverOnly bool, timeoutSecs int) error {
	// Generate PKCE values
	codeVerifier, codeChallenge, err := generatePKCE()
	if err != nil {
		return err
	}

	authURL := buildAuthURL(codeChallenge)

	// Channel to receive the auth code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create server with context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	server := &http.Server{
		Addr: ":8085",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, oauthSuccessHTML)
			// Send code to channel
			select {
			case codeChan <- code:
			default:
			}
		} else {
			// Handle error response
			errMsg := r.URL.Query().Get("error")
			if errMsg != "" {
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprintf(w, oauthErrorHTML, errMsg)
				select {
				case errChan <- fmt.Errorf("OAuth error: %s", errMsg):
				default:
				}
			}
		}
	})

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			select {
			case errChan <- fmt.Errorf("server error: %w", err):
			default:
			}
		}
	}()

	// Print information
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║         gcloud Application Default Credentials               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")
	fmt.Println("OAuth callback server started on http://localhost:8085")
	fmt.Println("")

	if serverOnly {
		// Server-only mode: just print the URL and wait
		fmt.Println("SERVER-ONLY MODE: Waiting for OAuth callback...")
		fmt.Println("")
		fmt.Println("Authorization URL:")
		fmt.Println(authURL)
		fmt.Println("")
		fmt.Printf("PKCE code_verifier: %s\n", codeVerifier)
		fmt.Println("")
		fmt.Printf("Timeout: %d seconds\n", timeoutSecs)
	} else {
		// Normal mode: open browser
		fmt.Println("Opening browser for authentication...")
		fmt.Println("")
		fmt.Println("If the browser doesn't open, visit this URL:")
		fmt.Println(authURL)
		fmt.Println("")

		// Open browser
		if err := exec.Command("open", authURL).Start(); err != nil {
			fmt.Printf("Warning: Could not open browser: %v\n", err)
			fmt.Println("Please open the URL above manually.")
		}
	}

	// Wait for code, error, or timeout
	select {
	case code := <-codeChan:
		fmt.Println("")
		fmt.Println("Received authorization code, exchanging for tokens...")

		// Shutdown server
		server.Shutdown(ctx)

		// Exchange code for tokens
		tokenResp, err := exchangeCodeForTokens(code, codeVerifier)
		if err != nil {
			return fmt.Errorf("token exchange failed: %w", err)
		}

		// Write credentials
		if err := writeGcloudCredentials(tokenResp.RefreshToken); err != nil {
			return err
		}

		fmt.Println("")
		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("✅ Credentials saved successfully!")
		fmt.Println("")
		homeDir, _ := os.UserHomeDir()
		fmt.Printf("   File: %s/.config/gcloud/application_default_credentials.json\n", homeDir)
		fmt.Println("")
		fmt.Println("You can now use:")
		fmt.Println("   gcloud auth application-default print-access-token")
		fmt.Println("   terraform apply (with google provider)")
		fmt.Println("════════════════════════════════════════════════════════════════")

		return nil

	case err := <-errChan:
		server.Shutdown(ctx)
		return err

	case <-ctx.Done():
		server.Shutdown(context.Background())
		return fmt.Errorf("authentication timed out after %d seconds", timeoutSecs)
	}
}

// checkGcloudAuth checks if application-default credentials exist and are valid
func checkGcloudAuth() (bool, string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, "Could not determine home directory"
	}

	credFile := filepath.Join(homeDir, ".config", "gcloud", "application_default_credentials.json")
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		return false, "No credentials file found"
	}

	// Try to read and validate
	data, err := os.ReadFile(credFile)
	if err != nil {
		return false, fmt.Sprintf("Could not read credentials: %v", err)
	}

	var creds GcloudCredentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return false, fmt.Sprintf("Invalid credentials format: %v", err)
	}

	if creds.RefreshToken == "" {
		return false, "Credentials missing refresh_token"
	}

	return true, credFile
}

// ============================================================================
// Playwright-based automated OAuth flow
// ============================================================================

// runGcloudAuthAssisted runs OAuth in assisted mode - opens default browser for passkey auth
// This mode is ideal for users who authenticate with passkeys, as it uses the system browser
// which has access to the OS keychain and passkey credentials
func runGcloudAuthAssisted(accountHint string, timeoutSecs int) error {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     gcloud Auth - Passkey Mode                               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")

	// Generate PKCE values
	codeVerifier, codeChallenge, err := generatePKCE()
	if err != nil {
		return err
	}

	authURL := buildAuthURL(codeChallenge)

	// If account hint provided, add login_hint to URL
	if accountHint != "" {
		authURL += "&login_hint=" + url.QueryEscape(accountHint)
		fmt.Printf("Account: %s\n", accountHint)
	} else {
		fmt.Println("Account: (will prompt for selection)")
		fmt.Println("")
		fmt.Println("TIP: Specify account with -account=EMAIL to skip selection")
	}
	fmt.Println("")

	// Channel to receive the auth code
	codeChan := make(chan string, 1)
	errChan := make(chan error, 1)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecs)*time.Second)
	defer cancel()

	// Start the callback server
	server := &http.Server{Addr: ":8085"}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, oauthSuccessHTML)
			select {
			case codeChan <- code:
			default:
			}
		} else if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, oauthErrorHTML, errMsg)
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

	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("Opening your default browser...")
	fmt.Println("")
	fmt.Println("Steps:")
	fmt.Println("  1. Select your Google account (or it will be pre-selected)")
	fmt.Println("  2. Authenticate with your passkey when prompted")
	fmt.Println("  3. Click 'Allow' to grant permissions")
	fmt.Println("  4. Return here - credentials will be saved automatically")
	fmt.Println("────────────────────────────────────────────────────────────────")
	fmt.Println("")
	fmt.Printf("Waiting for authentication (timeout: %ds)...\n", timeoutSecs)

	// Open in default browser
	if err := exec.Command("open", authURL).Start(); err != nil {
		fmt.Printf("\nWarning: Could not open browser: %v\n", err)
		fmt.Println("Please open this URL manually:")
		fmt.Println(authURL)
	}

	// Wait for code, error, or timeout
	select {
	case code := <-codeChan:
		fmt.Println("")
		fmt.Println("Received authorization code, exchanging for tokens...")

		// Shutdown server
		server.Shutdown(ctx)

		// Exchange code for tokens
		tokenResp, err := exchangeCodeForTokens(code, codeVerifier)
		if err != nil {
			return fmt.Errorf("token exchange failed: %w", err)
		}

		// Write credentials
		if err := writeGcloudCredentials(tokenResp.RefreshToken); err != nil {
			return err
		}

		fmt.Println("")
		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Println("✅ Credentials saved successfully!")
		fmt.Println("")
		homeDir, _ := os.UserHomeDir()
		fmt.Printf("   File: %s/.config/gcloud/application_default_credentials.json\n", homeDir)
		fmt.Println("")
		fmt.Println("You can now use:")
		fmt.Println("   gcloud auth application-default print-access-token")
		fmt.Println("   terraform apply (with google provider)")
		fmt.Println("════════════════════════════════════════════════════════════════")

		return nil

	case err := <-errChan:
		server.Shutdown(ctx)
		return err

	case <-ctx.Done():
		server.Shutdown(context.Background())
		return fmt.Errorf("authentication timed out after %d seconds", timeoutSecs)
	}
}

// PlaywrightOAuthResult represents the JSON output from the playwright oauth command
type PlaywrightOAuthResult struct {
	Code  string            `json:"code,omitempty"`
	Token string            `json:"token,omitempty"`
	Error string            `json:"error,omitempty"`
	Query map[string]string `json:"query,omitempty"`
}

// findPlaywrightBinary finds the playwright binary in common locations
func findPlaywrightBinary() (string, error) {
	// Check if playwright is in PATH
	if path, err := exec.LookPath("playwright"); err == nil {
		return path, nil
	}

	// Check common build locations relative to this binary
	execPath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(execPath)
		candidates := []string{
			filepath.Join(dir, "playwright"),
			filepath.Join(dir, "..", "bin", "playwright"),
		}
		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return candidate, nil
			}
		}
	}

	// Check if we can build it
	return "", fmt.Errorf("playwright binary not found - run: go build -o $(go env GOPATH)/bin/playwright ./cmd/playwright")
}

// runGcloudAuthAuto runs fully automated OAuth using the playwright CLI tool
func runGcloudAuthAuto(headless bool, timeoutSecs int) error {
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║     gcloud Auth - Automated Mode (Playwright)                ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Println("")

	// Find playwright binary
	playwrightBin, err := findPlaywrightBinary()
	if err != nil {
		return err
	}
	fmt.Printf("Using playwright: %s\n", playwrightBin)

	// Generate PKCE values
	codeVerifier, codeChallenge, err := generatePKCE()
	if err != nil {
		return err
	}

	authURL := buildAuthURL(codeChallenge)

	// Build playwright command args
	args := []string{
		fmt.Sprintf("-timeout=%d", timeoutSecs),
		"-port=8085",
	}
	if headless {
		args = append(args, "-headless")
	}
	args = append(args, "oauth", authURL)

	fmt.Println("")
	fmt.Println("Starting Playwright OAuth flow...")
	fmt.Println("Please complete authentication in the browser window.")
	fmt.Println("")

	// Run playwright oauth command
	cmd := exec.Command(playwrightBin, args...)
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("playwright oauth failed: %w", err)
	}

	// Parse the JSON result
	var result PlaywrightOAuthResult
	if err := json.Unmarshal(output, &result); err != nil {
		return fmt.Errorf("failed to parse playwright output: %w\nOutput: %s", err, string(output))
	}

	if result.Error != "" {
		return fmt.Errorf("OAuth error: %s", result.Error)
	}

	if result.Code == "" {
		return fmt.Errorf("no authorization code received")
	}

	fmt.Println("Received authorization code, exchanging for tokens...")

	// Exchange code for tokens
	tokenResp, err := exchangeCodeForTokens(result.Code, codeVerifier)
	if err != nil {
		return fmt.Errorf("token exchange failed: %w", err)
	}

	// Write credentials
	if err := writeGcloudCredentials(tokenResp.RefreshToken); err != nil {
		return err
	}

	fmt.Println("")
	fmt.Println("════════════════════════════════════════════════════════════════")
	fmt.Println("✅ Credentials saved successfully!")
	fmt.Println("")
	homeDir, _ := os.UserHomeDir()
	fmt.Printf("   File: %s/.config/gcloud/application_default_credentials.json\n", homeDir)
	fmt.Println("")
	fmt.Println("You can now use:")
	fmt.Println("   gcloud auth application-default print-access-token")
	fmt.Println("   terraform apply (with google provider)")
	fmt.Println("════════════════════════════════════════════════════════════════")

	return nil
}
