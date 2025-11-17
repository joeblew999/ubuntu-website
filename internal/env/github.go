package env

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// GitHubSecret represents a GitHub repository secret
type GitHubSecret struct {
	Name      string `json:"name"`
	UpdatedAt string `json:"updatedAt"`
}

// GetSecretsToSync returns fields that should be synced to GitHub
// Returns only fields marked for GitHub CI/CD deployment
func GetSecretsToSync() []FieldInfo {
	secrets := []FieldInfo{}
	for _, field := range envFieldsInOrder {
		// Only sync fields needed by GitHub Actions CI/CD
		if field.SyncToGitHub {
			secrets = append(secrets, field)
		}
	}
	return secrets
}

// CheckGitHubCLI checks if gh CLI is installed and authenticated
func CheckGitHubCLI() error {
	// Check if gh is installed
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("GitHub CLI (gh) is not installed. Install from: %s", GitHubCLIInstallURL)
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

// SyncSecretsToGitHub syncs environment variables to GitHub secrets
func SyncSecretsToGitHub(opts SyncOptions) ([]SyncResult, error) {
	var results []SyncResult

	// Use service to load config
	svc := NewService(false)
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env: %w", err)
	}

	// Get existing secrets
	existingSecrets, err := ListGitHubSecrets()
	if err != nil {
		return nil, err
	}

	// Process each secret
	for _, secretCfg := range GetSecretsToSync() {
		result := SyncResult{Name: secretCfg.Key}

		// Get value from config
		value := cfg.Get(secretCfg.Key)

		// Skip placeholders
		if IsPlaceholder(value) {
			result.Status = SyncStatusSkipped
			result.Reason = SyncReasonPlaceholder
			results = append(results, result)
			continue
		}

		// Validate if requested
		if opts.Validate && secretCfg.Validate {
			validationResult := ValidateField(secretCfg.Key, value, cfg, false)
			if !validationResult.Valid {
				result.Status = "failed"
				result.Reason = "validation failed"
				result.Error = validationResult.Error
				results = append(results, result)
				continue
			}
		}

		// Check if exists
		exists := SecretExists(secretCfg.Key, existingSecrets)
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
		if err := SetGitHubSecret(secretCfg.Key, value); err != nil {
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
	return fmt.Sprintf(GitHubRepoURLTemplate, owner, name), nil
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
