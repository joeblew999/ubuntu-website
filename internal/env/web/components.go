package web

import (
	"strings"

	"github.com/go-via/via/h"
	"github.com/joeblew999/ubuntu-website/internal/env"
)

// renderNavigation renders the shared navigation menu
// currentPage: "home", "cloudflare", or "claude"
func renderNavigation(currentPage string) h.H {
	// Helper to render a nav item (link or bold text)
	navItem := func(page, label, href string) h.H {
		if currentPage == page {
			return h.Li(h.Strong(h.Text(label)))
		}
		return h.Li(h.A(h.Href(href), h.Text(label)))
	}

	return h.Nav(
		h.Ul(
			navItem("home", "Overview", "/"),
			navItem("cloudflare", "Cloudflare", "/cloudflare"),
			navItem("claude", "Claude AI", "/claude"),
		),
	)
}

// BuildCloudflareURL builds a Cloudflare dashboard URL, replacing :account placeholder with actual account ID when available
func BuildCloudflareURL(baseURL, accountID string) string {
	if accountID != "" && !env.IsPlaceholder(accountID) {
		return strings.Replace(baseURL, ":account", accountID, 1)
	}
	return baseURL
}

// ConfigTableRow represents a row in the configuration overview table
type ConfigTableRow struct {
	Display   string // Human-readable display name
	Key       string // Environment variable key name
	Value     string // The actual value (formatted/masked)
	Required  string // "Yes" or "-"
	Validated string // "✓", "✗", "-"
	Error     string // Error message or "-"
}

// BuildConfigTableRows builds the configuration overview table data
func BuildConfigTableRows(mockMode bool) ([]ConfigTableRow, string, error) {
	svc := env.NewService(mockMode)

	// Get current config
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		return nil, "", err
	}

	// Get env file path
	envPath, err := env.GetEnvPath()
	if err != nil {
		envPath = ".env" // fallback
	}

	// Validate all fields to get current status
	validationResults := env.ValidateAllWithMode(cfg, mockMode)

	// Build table rows from validation results
	webRows := make([]ConfigTableRow, 0, len(validationResults))
	for _, result := range validationResults {
		// Get the key name from the display name (result.Name is DisplayName)
		keyName := env.GetKeyFromDisplayName(result.Name)
		fieldInfo := env.GetFieldInfo(keyName)

		// Determine "Required" status
		required := "-"
		if fieldInfo != nil && fieldInfo.SyncToGitHub {
			required = "Yes"
		}

		// Determine "Validated" status
		validated := "-"
		if !result.Skipped {
			if result.Valid {
				validated = "✓"
			} else {
				validated = "✗"
			}
		}

		// Get error message
		errorMsg := "-"
		if result.Error != nil {
			errorMsg = result.Error.Error()
		}

		// Get display value (formatted/masked)
		value := cfg.Get(keyName)
		displayValue := formatValueForDisplay(value)

		webRows = append(webRows, ConfigTableRow{
			Display:   result.Name,
			Key:       keyName,
			Value:     displayValue,
			Required:  required,
			Validated: validated,
			Error:     errorMsg,
		})
	}

	return webRows, envPath, nil
}

// formatValueForDisplay formats a value for display (masks sensitive data)
func formatValueForDisplay(value string) string {
	if env.IsPlaceholder(value) {
		return "(not set)"
	}
	// Show preview for secrets
	if len(value) > 24 {
		preview := value[:10] + "..." + value[len(value)-10:]
		return preview
	}
	return value
}
