package env

import (
	"fmt"
)

// ValidationResult holds the result of validating a credential
type ValidationResult struct {
	Name   string
	Valid  bool
	Error  error
	Skipped bool
}

// ValidateField validates a single field by env key
func ValidateField(envKey, value string, cfg *EnvConfig, mockMode bool) ValidationResult {
	// Get display name and validation requirement from envFieldsInOrder
	displayName := envKey
	requiresValidation := false
	for _, field := range envFieldsInOrder {
		if field.Key == envKey {
			displayName = field.DisplayName
			requiresValidation = field.Validate
			break
		}
	}

	// Skip placeholder values only for fields that don't require validation
	if IsPlaceholder(value) {
		if !requiresValidation {
			return ValidationResult{
				Name:    displayName,
				Skipped: true,
			}
		}
		// For fields that require validation, treat empty/placeholder as invalid
		return ValidationResult{
			Name:  displayName,
			Valid: false,
			Error: fmt.Errorf("%s is required", displayName),
		}
	}

	// Mock mode - simple length check
	if mockMode {
		valid := len(value) > 5
		var err error
		if !valid {
			err = fmt.Errorf("%s must be longer than 5 characters (mock validation)", displayName)
		}
		return ValidationResult{
			Name:  displayName,
			Valid: valid,
			Error: err,
		}
	}

	// Real validation using env key
	var err error
	switch envKey {
	case KeyCloudflareAPIToken:
		_, err = ValidateCloudflareToken(value)
	case KeyCloudflareAPITokenName:
		// Token name is just metadata - just check it exists
		if len(value) == 0 {
			err = fmt.Errorf("token name is required")
		}
	case KeyCloudflareAccountID:
		token := cfg.Get(KeyCloudflareAPIToken)
		_, err = ValidateCloudflareAccount(token, value)
	case KeyCloudflarePageProject:
		err = ValidateCloudflareProjectName(value)
	case KeyClaudeAPIKey:
		err = ValidateClaudeAPIKey(value)
	default:
		// Unknown field - skip validation
		return ValidationResult{
			Name:    displayName,
			Skipped: true,
		}
	}

	return ValidationResult{
		Name:  displayName,
		Valid: err == nil,
		Error: err,
	}
}

// ValidateAll validates all credentials in the config
func ValidateAll(cfg *EnvConfig) []ValidationResult {
	return ValidateAllWithMode(cfg, false)
}

// ValidateAllWithMode validates all credentials with optional mock mode
func ValidateAllWithMode(cfg *EnvConfig, mockMode bool) []ValidationResult {
	results := []ValidationResult{}

	// Iterate over all fields (show all, not just validated ones)
	for _, field := range envFieldsInOrder {
		value := cfg.Get(field.Key)
		// Always include the field in results, even if validation is skipped
		if field.Validate {
			results = append(results, ValidateField(field.Key, value, cfg, mockMode))
		} else {
			// Non-validated fields show as info only
			results = append(results, ValidationResult{
				Name:    field.DisplayName,
				Skipped: true,
			})
		}
	}

	return results
}

// PrintValidationResults prints validation results in a table format
func PrintValidationResults(results []ValidationResult, cfg *EnvConfig) {
	printHeader("Credential Validation", "")

	// Build table rows using helper
	rows := buildRowsFromValidation(results, cfg)
	renderCredentialTable(rows)

	fmt.Println()

	// Count results for summary
	valid, invalid, skipped := countValidationResults(results)

	// Summary
	parts := []string{}
	if valid > 0 {
		parts = append(parts, fmt.Sprintf("Valid: %d", valid))
	}
	if invalid > 0 {
		parts = append(parts, fmt.Sprintf("Invalid: %d", invalid))
	}
	if skipped > 0 {
		parts = append(parts, fmt.Sprintf("Not set: %d", skipped))
	}

	fmt.Printf("  %s\n", joinParts(parts))
	printFooter("")
}

// HasInvalidCredentials returns true if any credentials are invalid
func HasInvalidCredentials(results []ValidationResult) bool {
	for _, result := range results {
		if !result.Skipped && !result.Valid {
			return true
		}
	}
	return false
}

// GetInvalidFields returns the names of invalid fields
func GetInvalidFields(results []ValidationResult) []string {
	fields := []string{}
	for _, result := range results {
		if !result.Skipped && !result.Valid {
			fields = append(fields, result.Name)
		}
	}
	return fields
}

// HasInvalidCredentialsMap returns true if any credentials are invalid in a map
func HasInvalidCredentialsMap(results map[string]ValidationResult) bool {
	for _, result := range results {
		if !result.Skipped && !result.Valid {
			return true
		}
	}
	return false
}
