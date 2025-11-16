package env

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"strings"
)

// GitHubSecret represents a GitHub repository secret
type GitHubSecret struct {
	Name      string `json:"name"`
	UpdatedAt string `json:"updatedAt"`
}

// SecretConfig defines how a secret should be handled
type SecretConfig struct {
	Name        string
	Description string
	Required    bool // Required for CI/CD
	Validate    bool // Should validate before syncing
}

// SecretsToSync defines which secrets should be synced to GitHub
var SecretsToSync = []SecretConfig{
	{
		Name:        EnvCloudflareToken,
		Description: "Cloudflare API Token (required for deployment)",
		Required:    true,
		Validate:    true,
	},
	{
		Name:        EnvCloudflareAccount,
		Description: "Cloudflare Account ID",
		Required:    true,
		Validate:    false,
	},
	{
		Name:        EnvClaudeAPIKey,
		Description: "Claude API Key (optional, for CI translation)",
		Required:    false,
		Validate:    true,
	},
}

// CheckGitHubCLI checks if gh CLI is installed and authenticated
func CheckGitHubCLI() error {
	// Check if gh is installed
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed. Install from: https://cli.github.com/")
	}

	// Check if authenticated
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not authenticated with GitHub CLI. Run: gh auth login")
	}

	return nil
}

// GetRepositoryInfo returns the current repository owner and name
func GetRepositoryInfo() (string, string, error) {
	cmd := exec.Command("gh", "repo", "view", "--json", "owner,name")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get repository info: %w", err)
	}

	var repo struct {
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(output, &repo); err != nil {
		return "", "", fmt.Errorf("failed to parse repository info: %w", err)
	}

	return repo.Owner.Login, repo.Name, nil
}

// ListGitHubSecrets lists all secrets in the current repository
func ListGitHubSecrets() ([]GitHubSecret, error) {
	cmd := exec.Command("gh", "secret", "list", "--json", "name,updatedAt")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	var secrets []GitHubSecret
	if err := json.Unmarshal(output, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets: %w", err)
	}

	return secrets, nil
}

// SecretExists checks if a secret already exists
func SecretExists(name string, secrets []GitHubSecret) bool {
	for _, s := range secrets {
		if s.Name == name {
			return true
		}
	}
	return false
}

// SetGitHubSecret sets a single secret in the repository
func SetGitHubSecret(name, value string) error {
	cmd := exec.Command("gh", "secret", "set", name, "--body", value)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set secret %s: %w", name, err)
	}
	return nil
}

// SyncOptions contains options for syncing secrets
type SyncOptions struct {
	DryRun   bool
	Force    bool
	Validate bool
}

// SyncResult represents the result of syncing a secret
type SyncResult struct {
	Name   string
	Status string // "synced", "skipped", "failed"
	Reason string
	Error  error
}

// validateSecret validates a secret value by calling the appropriate validation function
func validateSecret(cfg *EnvConfig, envKey, value string) error {
	// Use reflection to get the validate tag for this env key
	v := reflect.ValueOf(&EnvConfig{}).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if getEnvKey(field) == envKey {
			validateName := getValidateName(field)
			if validateName == "" {
				return nil // No validation configured
			}

			// Call the appropriate validation function
			switch validateName {
			case "cloudflare_token":
				return ValidateCloudflareToken(value)
			case "cloudflare_account":
				// Account validation needs the token
				token, _ := getFieldByEnvKey(cfg, EnvCloudflareToken)
				_, err := ValidateCloudflareAccount(token, value)
				return err
			case "claude_api_key":
				return ValidateClaudeAPIKey(value)
			default:
				return nil
			}
		}
	}
	return nil
}

// SyncSecretsToGitHub syncs environment variables to GitHub secrets
func SyncSecretsToGitHub(opts SyncOptions) ([]SyncResult, error) {
	var results []SyncResult

	// Load .env configuration
	cfg, err := LoadEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	// Get existing secrets
	existingSecrets, err := ListGitHubSecrets()
	if err != nil {
		return nil, err
	}

	// Process each secret
	for _, secretCfg := range SecretsToSync {
		result := SyncResult{Name: secretCfg.Name}

		// Get value from config using reflection
		value, found := getFieldByEnvKey(cfg, secretCfg.Name)
		if !found {
			result.Status = "failed"
			result.Reason = "unknown field"
			result.Error = fmt.Errorf("no field found for env key: %s", secretCfg.Name)
			results = append(results, result)
			continue
		}

		// Skip placeholders
		if isPlaceholder(value) {
			result.Status = "skipped"
			result.Reason = "placeholder value"
			results = append(results, result)
			continue
		}

		// Validate if requested
		if opts.Validate && secretCfg.Validate {
			if validationErr := validateSecret(cfg, secretCfg.Name, value); validationErr != nil {
				result.Status = "failed"
				result.Reason = "validation failed"
				result.Error = validationErr
				results = append(results, result)
				continue
			}
		}

		// Check if exists
		exists := SecretExists(secretCfg.Name, existingSecrets)
		if exists && !opts.Force {
			result.Status = "skipped"
			result.Reason = "already exists (use --force to overwrite)"
			results = append(results, result)
			continue
		}

		// Dry run - don't actually set
		if opts.DryRun {
			result.Status = "would-sync"
			if exists {
				result.Reason = "would overwrite existing"
			} else {
				result.Reason = "would create new"
			}
			results = append(results, result)
			continue
		}

		// Actually set the secret
		if err := SetGitHubSecret(secretCfg.Name, value); err != nil {
			result.Status = "failed"
			result.Reason = "failed to set"
			result.Error = err
			results = append(results, result)
			continue
		}

		result.Status = "synced"
		if exists {
			result.Reason = "updated"
		} else {
			result.Reason = "created"
		}
		results = append(results, result)
	}

	return results, nil
}

// GetRepositoryURL returns the GitHub repository URL
func GetRepositoryURL() (string, error) {
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://github.com/%s/%s", owner, name), nil
}

// ValidateGitHubSetup checks if everything is ready for syncing
func ValidateGitHubSetup() error {
	// Check gh CLI
	if err := CheckGitHubCLI(); err != nil {
		return err
	}

	// Check repository
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return err
	}

	if owner == "" || name == "" {
		return fmt.Errorf("not in a GitHub repository")
	}

	return nil
}

// FormatSecretValue returns a preview of the secret value
func FormatSecretValue(value string, maxLen int) string {
	if isPlaceholder(value) {
		return "<not set>"
	}

	if len(value) <= maxLen {
		return strings.Repeat("*", len(value))
	}

	preview := value[:min(20, len(value))]
	return preview + "..." + strings.Repeat("*", maxLen-len(preview)-3)
}
