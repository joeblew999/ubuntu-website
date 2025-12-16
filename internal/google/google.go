// Package google provides a unified import point for all Google service packages.
//
// This package re-exports subpackages for convenient access:
//
//	import "github.com/joeblew999/ubuntu-website/internal/google"
//
// Available services:
//   - gcal: Google Calendar API client
//   - gmail: Gmail API client (TODO)
//   - auth: Google OAuth token management
//
// Future services (planned):
//   - sheets: Google Sheets API
//   - docs: Google Docs API
//   - drive: Google Drive API
//   - slides: Google Slides API
//
// Import individual packages:
//
//	import "github.com/joeblew999/ubuntu-website/internal/google/gcal"
//	import "github.com/joeblew999/ubuntu-website/internal/google/gmail"
//	import "github.com/joeblew999/ubuntu-website/internal/googleauth"
package google

// Google Cloud Console URLs for API enablement
// These are used by google-auth CLI and can be used by other tools
const (
	// Project management
	URLProject = "https://console.cloud.google.com/projectcreate"

	// OAuth configuration
	URLOAuthConsent = "https://console.cloud.google.com/apis/credentials/consent"
	URLCredentials  = "https://console.cloud.google.com/apis/credentials"

	// API enablement pages
	URLGmailAPI    = "https://console.cloud.google.com/apis/library/gmail.googleapis.com"
	URLCalendarAPI = "https://console.cloud.google.com/apis/library/calendar-json.googleapis.com"
	URLDriveAPI    = "https://console.cloud.google.com/apis/library/drive.googleapis.com"
	URLSheetsAPI   = "https://console.cloud.google.com/apis/library/sheets.googleapis.com"
	URLDocsAPI     = "https://console.cloud.google.com/apis/library/docs.googleapis.com"
	URLSlidesAPI   = "https://console.cloud.google.com/apis/library/slides.googleapis.com"

	// External tools
	URLGitHubRepo = "https://github.com/ngs/google-mcp-server"
)

// APIInfo holds information about a Google API
type APIInfo struct {
	Name string // Human-readable name (e.g., "Gmail")
	ID   string // API ID (e.g., "gmail.googleapis.com")
	URL  string // Console URL for enablement
}

// AllAPIs returns information about all supported Google APIs
func AllAPIs() []APIInfo {
	return []APIInfo{
		{Name: "Gmail", ID: "gmail.googleapis.com", URL: URLGmailAPI},
		{Name: "Calendar", ID: "calendar-json.googleapis.com", URL: URLCalendarAPI},
		{Name: "Drive", ID: "drive.googleapis.com", URL: URLDriveAPI},
		{Name: "Sheets", ID: "sheets.googleapis.com", URL: URLSheetsAPI},
		{Name: "Docs", ID: "docs.googleapis.com", URL: URLDocsAPI},
		{Name: "Slides", ID: "slides.googleapis.com", URL: URLSlidesAPI},
	}
}
