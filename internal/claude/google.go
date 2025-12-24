// Google MCP server definition and registration.
//
// This file registers the Google MCP server with the MCP registry.
// The server provides access to Google services (Calendar, Drive, Gmail, etc.)
// via the Model Context Protocol.
//
// Different Claude products handle permissions differently:
// - Claude Code (VSCode): Uses settings.json with allow/deny permission lists
// - Claude Desktop: Uses OAuth-based permissions (no file-based permissions)
// - Claude CLI: Uses internal permission management
//
// See targets.go for the full list of Claude targets and their config locations.
package claude

func init() {
	// Register the Google MCP server
	RegisterMCPServer(MCPServerDef{
		Name:        GoogleServerName,
		Command:     GoogleServerCmd,
		Args:        []string{},
		Env: map[string]string{
			"GOOGLE_CLIENT_ID":     EnvClientID,
			"GOOGLE_CLIENT_SECRET": EnvClientSecret,
		},
		Permissions: GoogleMCPPermissions,
		AccountsDir: GoogleAccountsDir,
	})
}

// GoogleServerName is the standard name for the Google MCP server
const GoogleServerName = "google"

// GoogleServerCmd is the command to run the Google MCP server
const GoogleServerCmd = "google-mcp-server"

// GoogleAccountsDir is the directory where google-mcp-server stores tokens
const GoogleAccountsDir = ".google-mcp-accounts"

// EnvClientID is the environment variable reference for client ID
const EnvClientID = "${GOOGLE_CLIENT_ID}"

// EnvClientSecret is the environment variable reference for client secret
const EnvClientSecret = "${GOOGLE_CLIENT_SECRET}"

// GoogleMCPPermissions lists all Google MCP tool permissions for Claude Code (VSCode)
// These must be listed individually for VSCode extension compatibility.
// Note: Claude Desktop and Claude CLI do NOT use this permission format.
var GoogleMCPPermissions = []string{
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

// Legacy helper functions - these are kept for backwards compatibility
// but new code should use the MCP registry functions instead.

// NewGoogleServer creates a Google MCP server configuration
// Deprecated: Use AddMCPServer("google", target, projectRoot) instead
func NewGoogleServer() Server {
	return Server{
		Command: GoogleServerCmd,
		Args:    []string{},
		Env: map[string]string{
			"GOOGLE_CLIENT_ID":     EnvClientID,
			"GOOGLE_CLIENT_SECRET": EnvClientSecret,
		},
	}
}

// AddGoogleServer adds the Google MCP server to a config
// Deprecated: Use AddMCPServer("google", target, projectRoot) instead
func AddGoogleServer(config *Config) bool {
	if config.HasServer(GoogleServerName) {
		return false // Already exists
	}
	config.AddServer(GoogleServerName, NewGoogleServer())
	return true
}

// RemoveGoogleServer removes the Google MCP server from a config
// Deprecated: Use RemoveMCPServer("google", target, projectRoot) instead
func RemoveGoogleServer(config *Config) bool {
	return config.RemoveServer(GoogleServerName)
}

// AddGooglePermissions adds all Google MCP permissions to settings
// Deprecated: Use AddMCPServer which handles permissions automatically
func AddGooglePermissions(settings *ClaudeSettings) bool {
	return settings.AddPermissions(GoogleMCPPermissions)
}

// RemoveGooglePermissions removes all Google MCP permissions from settings
// Deprecated: Use RemoveMCPServer which handles permissions automatically
func RemoveGooglePermissions(settings *ClaudeSettings) bool {
	return settings.RemovePermissions(GoogleMCPPermissions)
}

// CountGooglePermissions returns how many Google permissions are in settings
func CountGooglePermissions(settings *ClaudeSettings) int {
	return settings.CountPermissionsWithPrefix("mcp__google_")
}
