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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	schemaURL    = "https://modelcontextprotocol.io/schema/config.json"
	serverName   = "google"
	serverCmd    = "google-mcp-server"
	envClientID  = "${GOOGLE_CLIENT_ID}"
	envSecretID  = "${GOOGLE_CLIENT_SECRET}"
	// settings.json is used for project-level Claude permissions
	// VSCode extension reads this file (settings.local.json didn't work)
	settingsFile = "settings.json"
)

// ClaudeSettings represents the .claude/settings.json file structure
type ClaudeSettings struct {
	Permissions                 Permissions       `json:"permissions"`
	EnableAllProjectMcpServers  bool              `json:"enableAllProjectMcpServers"`
	Env                         map[string]string `json:"env,omitempty"`
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
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

// Google Cloud Console URLs
const (
	urlProject        = "https://console.cloud.google.com/projectcreate"
	urlOAuthConsent   = "https://console.cloud.google.com/apis/credentials/consent"
	urlCredentials    = "https://console.cloud.google.com/apis/credentials"
	urlGmailAPI       = "https://console.cloud.google.com/apis/library/gmail.googleapis.com"
	urlCalendarAPI    = "https://console.cloud.google.com/apis/library/calendar-json.googleapis.com"
	urlDriveAPI       = "https://console.cloud.google.com/apis/library/drive.googleapis.com"
	urlSheetsAPI      = "https://console.cloud.google.com/apis/library/sheets.googleapis.com"
	urlDocsAPI        = "https://console.cloud.google.com/apis/library/docs.googleapis.com"
	urlSlidesAPI      = "https://console.cloud.google.com/apis/library/slides.googleapis.com"
	accountsDir       = ".google-mcp-accounts"
)

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
	fmt.Println("  google-auth [-location=LOC] add      Add google server + permissions")
	fmt.Println("  google-auth [-location=LOC] remove   Remove google server + permissions")
	fmt.Println("  google-auth [-location=LOC] status   Show MCP config status")
	fmt.Println("")
	fmt.Println("Locations (-location flag):")
	fmt.Println("  vscode  - .vscode/mcp.json (default, for VSCode extension)")
	fmt.Println("  project - .mcp.json (CLI project-level)")
	fmt.Println("  claude  - .claude/mcp.json (Claude folder)")
	fmt.Println("")
	fmt.Println("Quick start:")
	fmt.Println("  google-auth guide   # Follow the setup steps")
	fmt.Println("  google-auth check   # See what's configured")
	fmt.Println("  google-auth add     # Add to Claude Code")
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
