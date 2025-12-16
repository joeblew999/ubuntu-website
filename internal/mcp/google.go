// Google MCP server configuration helpers
//
// IMPORTANT: The permissions defined here are specifically for CLAUDE CODE
// (VSCode extension). Different Claude products have different permission models:
//
// - Claude Code (VSCode): Uses settings.json with allow/deny permission lists
// - Claude Desktop: Uses OAuth-based permissions (no file-based permissions)
// - Claude CLI: Uses internal permission management
// - Claude Cloud: Uses API keys and organization policies
//
// When we add KRONK API support, it will have its own permission model that
// may be API-based rather than file-based like Claude Code.
//
// See targets.go for the full list of Claude targets and their config locations.
package mcp

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

// NewGoogleServer creates a Google MCP server configuration
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
func AddGoogleServer(config *Config) bool {
	if config.HasServer(GoogleServerName) {
		return false // Already exists
	}
	config.AddServer(GoogleServerName, NewGoogleServer())
	return true
}

// RemoveGoogleServer removes the Google MCP server from a config
func RemoveGoogleServer(config *Config) bool {
	return config.RemoveServer(GoogleServerName)
}

// AddGooglePermissions adds all Google MCP permissions to settings
func AddGooglePermissions(settings *ClaudeSettings) bool {
	return settings.AddPermissions(GoogleMCPPermissions)
}

// RemoveGooglePermissions removes all Google MCP permissions from settings
func RemoveGooglePermissions(settings *ClaudeSettings) bool {
	return settings.RemovePermissions(GoogleMCPPermissions)
}

// CountGooglePermissions returns how many Google permissions are in settings
func CountGooglePermissions(settings *ClaudeSettings) int {
	return settings.CountPermissionsWithPrefix("mcp__google_")
}
