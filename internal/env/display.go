package env

import (
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

// printHeader prints a consistent header for all env commands
func printHeader(title, subtitle string) {
	fmt.Println()
	fmt.Println(title)
	if subtitle != "" {
		fmt.Println(subtitle)
	}
	fmt.Println()
}

// printFooter prints a consistent footer with optional action hint
func printFooter(hint string) {
	if hint != "" {
		fmt.Println()
		fmt.Println(hint)
	}
	fmt.Println()
}

// joinParts joins string parts with bullet separator
func joinParts(parts []string) string {
	if len(parts) == 0 {
		return "0"
	}
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += " • "
		}
		result += part
	}
	return result
}

// printSection prints a section header
func printSection(title string) {
	fmt.Println(title + ":")
}

// formatValueForDisplay formats a value for display (handles placeholders and truncation)
func formatValueForDisplay(value string) string {
	if IsPlaceholder(value) {
		return "(not set)"
	}
	// Show preview for secrets
	if len(value) > 24 {
		preview := value[:10] + "..." + value[len(value)-10:]
		return preview
	}
	return value
}

// buildSummaryParts builds summary parts from counts
func buildSummaryParts(created, updated, skipped, failed int) []string {
	parts := []string{}
	if created > 0 {
		parts = append(parts, fmt.Sprintf("New: %d", created))
	}
	if updated > 0 {
		parts = append(parts, fmt.Sprintf("Update: %d", updated))
	}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("Skip: %d", skipped))
	}
	if failed > 0 {
		parts = append(parts, fmt.Sprintf("Fail: %d", failed))
	}
	return parts
}

// formatGitHubTimestamp formats GitHub API timestamp (ISO 8601) to simple format
func formatGitHubTimestamp(githubTime string) string {
	// GitHub returns ISO 8601: "2024-01-15T10:30:45Z"
	// Parse it and format to simple: "2006-01-02 15:04"
	t, err := time.Parse(time.RFC3339, githubTime)
	if err != nil {
		return githubTime // Return as-is if parse fails
	}
	return t.Local().Format("2006-01-02 15:04")
}

// credentialTableRow represents a row in a credential table
// All tables show: Display | Key | Value | Required | Validated | Error
type credentialTableRow struct {
	Display   string // Human-readable display name (e.g., "Cloudflare API Token")
	Key       string // Environment variable key name (e.g., "CLOUDFLARE_API_TOKEN")
	Value     string // The actual value (or formatted value)
	Required  string // Required for deployment ("Yes" or "-")
	Validated string // Validation status ("✓", "✗", "-")
	Error     string // Error message if validation failed (or "-")
}

// renderCredentialTable renders a unified table with consistent UX
// Always shows: Display | Key | Value | Required | Validated | Error (6 columns for all tables)
func renderCredentialTable(rows []credentialTableRow) {
	// Always show all 6 columns for consistent spacing and clarity
	header := []string{"Display", "Key", "Value", "Required", "Validated", "Error"}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(header),
	)

	// Build rows - always include all 6 columns
	tableRows := [][]string{}
	for _, row := range rows {
		tableRow := []string{row.Display, row.Key, row.Value, row.Required, row.Validated, row.Error}
		tableRows = append(tableRows, tableRow)
	}

	table.Bulk(tableRows)
	table.Render()
}

// formatRequired returns "Yes" if should sync to GitHub, "-" otherwise
func formatRequired(syncToGitHub bool) string {
	if syncToGitHub {
		return "Yes"
	}
	return "-"
}


// buildRowFromValidation creates a table row from a ValidationResult
// Handles the DisplayName->Key lookup and error formatting automatically
func buildRowFromValidation(result ValidationResult, cfg *EnvConfig) credentialTableRow {
	var validated, errorMsg string

	if result.Skipped {
		validated = "-"
		errorMsg = "-"
	} else if result.Valid {
		validated = "✓"
		errorMsg = "-"
	} else {
		validated = "✗"
		if result.Error != nil {
			errorMsg = result.Error.Error()
			// Truncate very long error messages
			if len(errorMsg) > 80 {
				errorMsg = errorMsg[:77] + "..."
			}
		} else {
			errorMsg = "Invalid"
		}
	}

	// Get the key name and field info from result.Name (which is DisplayName)
	keyName := ""
	var fieldInfo *FieldInfo
	for _, field := range envFieldsInOrder {
		if field.DisplayName == result.Name {
			keyName = field.Key
			fieldInfo = &field
			break
		}
	}
	if keyName == "" {
		keyName = result.Name // Fallback to DisplayName if not found
	}

	// Get required status from field info
	required := "-"
	if fieldInfo != nil && fieldInfo.SyncToGitHub {
		required = "Yes"
	}

	// Get the value from config
	value := cfg.Get(keyName)

	return credentialTableRow{
		Display:   result.Name,
		Key:       keyName,
		Value:     formatValueForDisplay(value),
		Required:  required,
		Validated: validated,
		Error:     errorMsg,
	}
}

// buildRowsFromValidation creates table rows from all validation results
// Used by validation command
func buildRowsFromValidation(results []ValidationResult, cfg *EnvConfig) []credentialTableRow {
	rows := make([]credentialTableRow, 0, len(results))
	for _, result := range results {
		rows = append(rows, buildRowFromValidation(result, cfg))
	}
	return rows
}


// countValidationResults counts results by status for summaries
func countValidationResults(results []ValidationResult) (valid, invalid, skipped int) {
	for _, result := range results {
		if result.Skipped {
			skipped++
		} else if result.Valid {
			valid++
		} else {
			invalid++
		}
	}
	return valid, invalid, skipped
}

// buildUnifiedRows creates a single table combining local .env and GitHub secrets status
// Shows: Display | Key | Local Value | GitHub Status | Required | Validated
func buildUnifiedRows(cfg *EnvConfig, fields []FieldInfo, githubSecrets []GitHubSecret) []credentialTableRow {
	// Create map of GitHub secrets for quick lookup
	ghMap := make(map[string]GitHubSecret)
	for _, secret := range githubSecrets {
		ghMap[secret.Name] = secret
	}

	rows := make([]credentialTableRow, 0, len(fields))
	for _, field := range fields {
		localValue := cfg.Get(field.Key)

		// Determine local validation status
		validated := "✓"
		errorMsg := "-"
		if IsPlaceholder(localValue) {
			validated = "✗"
			if field.SyncToGitHub {
				errorMsg = "Not set (required for GitHub deployment)"
			}
		}

		// Build value display: local value + GitHub sync status
		valueDisplay := formatValueForDisplay(localValue)

		// If field syncs to GitHub, show GitHub status
		if field.SyncToGitHub {
			if ghSecret, exists := ghMap[field.Key]; exists {
				// Show that it's synced to GitHub with timestamp
				valueDisplay += " (GH: ✓ " + formatGitHubTimestamp(ghSecret.UpdatedAt) + ")"
			} else {
				// Should be synced but isn't
				valueDisplay += " (GH: ✗ not synced)"
				if validated == "✓" {
					// Local value is set but not on GitHub
					errorMsg = "Not synced to GitHub (run 'gh-push')"
				}
			}
		}

		rows = append(rows, credentialTableRow{
			Display:   field.DisplayName,
			Key:       field.Key,
			Value:     valueDisplay,
			Required:  formatRequired(field.SyncToGitHub),
			Validated: validated,
			Error:     errorMsg,
		})
	}
	return rows
}
