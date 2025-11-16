package env

import (
	"fmt"
	"reflect"
)

// printHeader prints a consistent header for all env commands
func printHeader(title, subtitle string) {
	fmt.Println()
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Printf("  %s\n", title)
	if subtitle != "" {
		fmt.Printf("  %s\n", subtitle)
	}
	fmt.Println("════════════════════════════════════════════════════════════")
	fmt.Println()
}

// printFooter prints a consistent footer with optional action hint
func printFooter(hint string) {
	fmt.Println("────────────────────────────────────────────────────────────")
	if hint != "" {
		fmt.Printf("%s\n", hint)
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

// ShowConfig displays the current environment configuration
func ShowConfig() error {
	// Get absolute path to .env file
	envPath, err := GetEnvPath()
	if err != nil {
		envPath = ".env"
	}

	printHeader("Local Environment Configuration", envPath)

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

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
	printFooter("→ To update: task env:local:setup")

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
	// Validate GitHub setup
	if err := ValidateGitHubSetup(); err != nil {
		return err
	}

	// Get repository info
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return err
	}

	repoName := fmt.Sprintf("%s/%s", owner, name)
	printHeader("GitHub Secrets", repoName)

	// List secrets
	secrets, err := ListGitHubSecrets()
	if err != nil {
		return err
	}

	if len(secrets) == 0 {
		fmt.Println("No secrets configured")
		fmt.Println()
	} else {
		for _, secret := range secrets {
			fmt.Printf("  %s\n", secret.Name)
			fmt.Printf("    Updated: %s\n", secret.UpdatedAt)
			fmt.Println()
		}
	}

	// Show management URL
	repoURL, _ := GetRepositoryURL()
	printFooter(fmt.Sprintf("→ Manage at: %s/settings/secrets/actions", repoURL))

	return nil
}

// SyncSecrets syncs environment variables to GitHub secrets
func SyncSecrets(dryRun, force, validate bool) error {
	// Validate GitHub setup
	if err := ValidateGitHubSetup(); err != nil {
		return err
	}

	// Get repository info
	owner, name, err := GetRepositoryInfo()
	if err != nil {
		return err
	}

	repoName := fmt.Sprintf("%s/%s", owner, name)
	mode := "Push to GitHub"
	if dryRun {
		mode = "Push to GitHub (Dry Run)"
	}
	printHeader(mode, repoName)

	if dryRun {
		fmt.Println("DRY RUN MODE - No secrets will be modified")
		fmt.Println()
	}

	if validate {
		fmt.Println("Validating credentials before push...")
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

	created := 0
	updated := 0
	skipped := 0
	failed := 0

	for _, result := range results {
		var icon, badge, status string
		switch result.Status {
		case "synced":
			icon = "✓"
			if result.Reason == "created" {
				badge = Colorize("[new]", ColorGreen)
				status = "Created new secret"
				created++
			} else {
				badge = Colorize("[upd]", ColorBlue)
				status = "Updated existing secret"
				updated++
			}
		case "would-sync":
			icon = "→"
			if result.Reason == "would create new" {
				badge = Colorize("[new]", ColorGreen)
				status = "Would create new secret"
				created++
			} else {
				badge = Colorize("[upd]", ColorBlue)
				status = "Would update existing secret"
				updated++
			}
		case "skipped":
			icon = "⊘"
			badge = Colorize("[skip]", ColorGray)
			status = result.Reason
			skipped++
		case "failed":
			icon = "✗"
			badge = Colorize("[fail]", ColorRed)
			status = fmt.Sprintf("%s: %v", result.Reason, result.Error)
			failed++
		}

		fmt.Printf("  [%s] %s %s\n", icon, badge, result.Name)
		fmt.Printf("      %s\n", status)
		fmt.Println()
	}

	synced := created + updated

	// Summary
	summary := ""
	if dryRun {
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
		summary = fmt.Sprintf("Would sync: %s", joinParts(parts))
		summary += "\n→ Run without --check to actually sync"
	} else {
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
		summary = fmt.Sprintf("Synced: %s", joinParts(parts))
		if synced > 0 {
			repoURL, _ := GetRepositoryURL()
			summary += fmt.Sprintf("\n→ Verify at: %s/settings/secrets/actions", repoURL)
		}
	}

	printFooter(summary)

	if failed > 0 {
		return fmt.Errorf("failed to sync %d secrets", failed)
	}

	return nil
}
