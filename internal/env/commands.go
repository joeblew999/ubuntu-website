package env

import (
	"fmt"
	"reflect"
)

// ShowConfig displays the current environment configuration
func ShowConfig() error {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  Current Configuration")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	// Get absolute path to .env file
	envPath, err := GetEnvPath()
	if err != nil {
		envPath = ".env"
	}
	fmt.Printf("File: %s\n", envPath)
	fmt.Println()

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	var lastComment string
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envKey := getEnvKey(field)
		if envKey == "" {
			continue
		}

		// Print comment header if new section
		comment := getComment(field)
		if comment != "" && comment != lastComment {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("# %s\n", comment)
			lastComment = comment
		}

		// Get value
		value := v.Field(i).String()

		// Format value for display
		var displayValue string
		if isPlaceholder(value) {
			displayValue = Colorize("<not set>", ColorGray)
		} else {
			// Show preview for secrets
			preview := value
			if len(preview) > 40 {
				preview = preview[:20] + "..." + preview[len(preview)-17:]
			}
			displayValue = preview
		}

		required := ""
		if isRequired(field) {
			required = Colorize(" (required)", ColorRed)
		}

		fmt.Printf("  %s%s = %s\n", envKey, required, displayValue)
	}

	fmt.Println()
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println("To update: edit .env or run the wizard")
	fmt.Println()

	return nil
}

// ValidateAll validates all configured credentials
func ValidateAll() error {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  Validating Credentials")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	hasErrors := false

	// Validate Cloudflare
	fmt.Println("Cloudflare:")
	if isPlaceholder(cfg.CloudflareToken) {
		fmt.Println(Skipped("Token not configured"))
	} else {
		if err := ValidateCloudflareToken(cfg.CloudflareToken); err != nil {
			fmt.Println(Error(fmt.Sprintf("Token invalid: %v", err)))
			hasErrors = true
		} else {
			fmt.Println(Success("Token is valid"))

			if !isPlaceholder(cfg.CloudflareAccount) {
				if accountName, err := ValidateCloudflareAccount(cfg.CloudflareToken, cfg.CloudflareAccount); err != nil {
					fmt.Println(Error(fmt.Sprintf("Account ID invalid: %v", err)))
					hasErrors = true
				} else {
					fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
				}
			}
		}
	}
	fmt.Println()

	// Validate Claude
	fmt.Println("Claude:")
	if isPlaceholder(cfg.ClaudeAPIKey) {
		fmt.Println(Skipped("API key not configured"))
	} else {
		if err := ValidateClaudeAPIKey(cfg.ClaudeAPIKey); err != nil {
			fmt.Println(Error(fmt.Sprintf("API key invalid: %v", err)))
			hasErrors = true
		} else {
			fmt.Println(Success("API key is valid"))
			fmt.Println(Success("API key has active credits"))
		}
	}
	fmt.Println()

	if hasErrors {
		return fmt.Errorf("validation failed")
	}

	return nil
}

// ValidateCloudflareCredentials validates only Cloudflare credentials
func ValidateCloudflareCredentials() error {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  Validating Cloudflare Credentials")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	if isPlaceholder(cfg.CloudflareToken) {
		return fmt.Errorf("cloudflare token not configured in .env")
	}

	fmt.Println("Validating token...")
	if err := ValidateCloudflareToken(cfg.CloudflareToken); err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}

	fmt.Println(Success("Cloudflare API token is valid"))

	if !isPlaceholder(cfg.CloudflareAccount) {
		fmt.Println()
		fmt.Println("Validating account ID...")
		if accountName, err := ValidateCloudflareAccount(cfg.CloudflareToken, cfg.CloudflareAccount); err != nil {
			return fmt.Errorf("account ID is invalid: %w", err)
		} else {
			fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
		}
	}

	fmt.Println()
	return nil
}

// ValidateClaudeCredentials validates only Claude API credentials
func ValidateClaudeCredentials() error {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  Validating Claude API Key")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	if isPlaceholder(cfg.ClaudeAPIKey) {
		return fmt.Errorf("claude API key not configured in .env")
	}

	fmt.Println("Validating API key...")
	if err := ValidateClaudeAPIKey(cfg.ClaudeAPIKey); err != nil {
		return fmt.Errorf("API key is invalid: %w", err)
	}

	fmt.Println(Success("Claude API key is valid"))
	fmt.Println(Success("API key has active credits"))
	fmt.Println()

	return nil
}

// ShowRemoteSecrets displays GitHub secrets status
func ShowRemoteSecrets() error {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  Remote GitHub Secrets")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()

	// Validate GitHub setup
	if err := ValidateGitHubSetup(); err != nil {
		return err
	}

	// Get repository info
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return err
	}

	fmt.Printf("Repository: %s/%s\n", owner, name)
	fmt.Println()

	// List secrets
	secrets, err := ListGitHubSecrets()
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets configured")
	} else {
		for _, secret := range secrets {
			fmt.Printf("  %s\n", secret.Name)
			fmt.Printf("    Updated: %s\n", secret.UpdatedAt)
			fmt.Println()
		}
	}

	// Show management URL
	repoURL, _ := GetRepositoryURL()
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("→ Manage at: %s/settings/secrets/actions\n", repoURL)
	fmt.Println()

	return nil
}

// SyncSecrets syncs environment variables to GitHub secrets
func SyncSecrets(dryRun, force, validate bool) error {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println("  Sync Environment Variables to GitHub Secrets")
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()

	// Validate GitHub setup
	if err := ValidateGitHubSetup(); err != nil {
		return err
	}

	// Get repository info
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return err
	}

	fmt.Printf("✓ GitHub CLI authenticated\n")
	fmt.Printf("✓ Repository: %s/%s\n", owner, name)
	fmt.Println()

	if dryRun {
		fmt.Println("DRY RUN MODE - No secrets will be modified")
		fmt.Println()
	}

	// Sync secrets
	opts := SyncOptions{
		DryRun:   dryRun,
		Force:    force,
		Validate: validate,
	}

	results, err := SyncSecretsToGitHub(opts)
	if err != nil {
		return err
	}

	// Display results
	fmt.Println("Secrets status:")
	fmt.Println()

	synced := 0
	skipped := 0
	failed := 0

	for _, result := range results {
		var icon, status string
		switch result.Status {
		case "synced":
			icon = "✓"
			status = result.Reason
			synced++
		case "would-sync":
			icon = "→"
			status = result.Reason
			synced++ // Count for dry-run summary
		case "skipped":
			icon = "⊘"
			status = result.Reason
			skipped++
		case "failed":
			icon = "✗"
			status = fmt.Sprintf("%s: %v", result.Reason, result.Error)
			failed++
		}

		fmt.Printf("  [%s] %s\n", icon, result.Name)
		fmt.Printf("      %s\n", status)
		fmt.Println()
	}

	// Summary
	fmt.Println("────────────────────────────────────────────────────────────")
	if dryRun {
		fmt.Printf("Would sync: %d\n", synced)
	} else {
		fmt.Printf("Synced: %d\n", synced)
	}
	fmt.Printf("Skipped: %d\n", skipped)
	if failed > 0 {
		fmt.Printf("Failed: %d\n", failed)
	}
	fmt.Println()

	if !dryRun && synced > 0 {
		repoURL, _ := GetRepositoryURL()
		fmt.Println("Next steps:")
		fmt.Printf("  Verify secrets at: %s/settings/secrets/actions\n", repoURL)
		fmt.Println()
	}

	if dryRun {
		fmt.Println("Run without --check to actually sync secrets")
		fmt.Println()
	}

	return nil
}
