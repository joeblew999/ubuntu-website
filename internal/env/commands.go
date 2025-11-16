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

// printWizardStep prints a wizard step header
func printWizardStep(step, title, envKey string) {
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("%s: %s\n", step, title)
	fmt.Printf("Setting: %s\n", envKey)
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()
}

// printClaudeValidationSuccess prints standardized Claude validation success messages
func printClaudeValidationSuccess() {
	fmt.Println(Success("Claude API key is valid"))
	fmt.Println(Success("API key has active credits"))
	fmt.Println()
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

	// Find longest key name for alignment
	maxKeyLen := 0
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		envKey := getEnvKey(field)
		if envKey != "" && len(envKey) > maxKeyLen {
			maxKeyLen = len(envKey)
		}
	}

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
			fmt.Println(Colorize(comment, ColorCyan))
			fmt.Println()
			lastComment = comment
		}

		// Get value
		value := v.Field(i).String()

		// Format value for display
		var displayValue string
		var statusIcon string
		if isPlaceholder(value) {
			displayValue = Colorize("<not set>", ColorGray)
			statusIcon = "○"
		} else {
			// Show preview for secrets
			preview := value
			if len(preview) > 40 {
				preview = preview[:20] + "..." + preview[len(preview)-17:]
			}
			displayValue = preview
			statusIcon = Colorize("✓", ColorGreen)
		}

		// Format key with padding
		keyDisplay := Colorize(fmt.Sprintf("%-*s", maxKeyLen, envKey), ColorBlue)

		required := ""
		if isRequired(field) {
			required = Colorize("*", ColorRed)
		} else {
			required = " "
		}

		fmt.Printf("  %s %s  %-*s  %s\n", statusIcon, required, maxKeyLen, keyDisplay, displayValue)
	}

	fmt.Println()
	printFooter("→ To update: task env:local:setup  (* = required)")

	return nil
}

// ValidateAll validates all configured credentials
func ValidateAll() error {
	printHeader("Validating Credentials", "")

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
		tokenName, err := ValidateCloudflareToken(cfg.CloudflareToken)
		if err != nil {
			fmt.Println(Error(fmt.Sprintf("Token invalid: %v", err)))
			hasErrors = true
		} else {
			if tokenName != "" {
				fmt.Println(Success(fmt.Sprintf("Token is valid: %s", tokenName)))
			} else {
				fmt.Println(Success("Token is valid"))
			}

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
	printHeader("Validating Cloudflare Credentials", "")

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	if isPlaceholder(cfg.CloudflareToken) {
		return fmt.Errorf("cloudflare token not configured in .env")
	}

	fmt.Println("Validating token...")
	tokenName, err := ValidateCloudflareToken(cfg.CloudflareToken)
	if err != nil {
		return fmt.Errorf("token is invalid: %w", err)
	}

	fmt.Println(Success(fmt.Sprintf("Cloudflare API token is valid: %s", tokenName)))

	accountName := ""
	if !isPlaceholder(cfg.CloudflareAccount) {
		fmt.Println()
		fmt.Println("Validating account ID...")
		accountName, err = ValidateCloudflareAccount(cfg.CloudflareToken, cfg.CloudflareAccount)
		if err != nil {
			return fmt.Errorf("account ID is invalid: %w", err)
		}
		fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
	}

	fmt.Println()
	return nil
}

// ValidateClaudeCredentials validates only Claude API credentials
func ValidateClaudeCredentials() error {
	printHeader("Validating Claude API Key", "")

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

	printClaudeValidationSuccess()
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
		fmt.Println(Colorize("No secrets configured", ColorGray))
		fmt.Println()
	} else {
		// Find longest secret name for alignment
		maxNameLen := 0
		for _, secret := range secrets {
			if len(secret.Name) > maxNameLen {
				maxNameLen = len(secret.Name)
			}
		}

		fmt.Println(Colorize("GitHub Secrets", ColorCyan))
		fmt.Println()

		for _, secret := range secrets {
			nameDisplay := Colorize(fmt.Sprintf("%-*s", maxNameLen, secret.Name), ColorBlue)
			fmt.Printf("  %s  %s  %s\n",
				Colorize("✓", ColorGreen),
				nameDisplay,
				Colorize(secret.UpdatedAt, ColorGray))
		}
		fmt.Println()
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
	created := 0
	updated := 0
	skipped := 0
	failed := 0

	// Find longest secret name for alignment
	maxNameLen := 0
	for _, result := range results {
		if len(result.Name) > maxNameLen {
			maxNameLen = len(result.Name)
		}
	}

	fmt.Println(Colorize("Secrets Status", ColorCyan))
	fmt.Println()

	for _, result := range results {
		var icon, status string
		var nameColor string
		switch result.Status {
		case "synced":
			icon = Colorize("✓", ColorGreen)
			nameColor = ColorGreen
			if result.Reason == "created" {
				status = "Created new"
				created++
			} else {
				status = "Updated"
				updated++
			}
		case "would-sync":
			icon = "→"
			nameColor = ColorBlue
			if result.Reason == "would create new" {
				status = "Would create"
				created++
			} else {
				status = "Would update"
				updated++
			}
		case "skipped":
			icon = Colorize("○", ColorGray)
			nameColor = ColorGray
			status = result.Reason
			skipped++
		case "failed":
			icon = Colorize("✗", ColorRed)
			nameColor = ColorRed
			status = fmt.Sprintf("%s: %v", result.Reason, result.Error)
			failed++
		}

		nameDisplay := Colorize(fmt.Sprintf("%-*s", maxNameLen, result.Name), nameColor)
		fmt.Printf("  %s  %s  %s\n", icon, nameDisplay, status)
	}

	fmt.Println()

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
