package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// RunWizard runs the interactive setup wizard
func RunWizard() error {
	printHeader("Environment Setup Wizard", "")

	// Check/create .env file
	if !EnvExists() {
		fmt.Println("Creating .env file...")
		if err := CreateEnv(); err != nil {
			return err
		}
		fmt.Println("✓ Created .env file")
	} else {
		fmt.Println("✓ .env file exists")
	}
	fmt.Println()

	// Setup Cloudflare
	if err := setupCloudflare(); err != nil {
		return err
	}

	// Setup Claude
	if err := setupClaude(); err != nil {
		return err
	}

	// Show next steps
	showNextSteps()

	return nil
}

func setupCloudflare() error {
	printWizardStep("Step 1 of 2", "Cloudflare API Token (Optional)", EnvCloudflareToken)

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	// Try to validate existing credentials
	if cfg.CloudflareToken != "" && cfg.CloudflareToken != PlaceholderToken {
		if validated, err := validateExistingCloudflareSetup(cfg); validated {
			return err // nil on success, error on failure
		}
		// If not validated, continue to prompt for new token
	}

	// Prompt for new token
	if err := promptForCloudflareToken(); err != nil {
		return err
	}

	// Ensure project name is set
	return ensureProjectName()
}

// validateExistingCloudflareSetup validates existing Cloudflare credentials
// Returns (true, nil) if validation succeeded and setup is complete
// Returns (true, err) if validation failed with an error
// Returns (false, nil) if validation failed but should continue to prompt
func validateExistingCloudflareSetup(cfg *EnvConfig) (bool, error) {
	fmt.Println("Validating existing Cloudflare credentials...")
	tokenName, err := ValidateCloudflareToken(cfg.CloudflareToken)
	if err != nil {
		fmt.Println(Error(fmt.Sprintf("Token validation failed: %v", err)))
		fmt.Println()
		fmt.Println(Colorize("Will prompt for new token...", ColorYellow))
		fmt.Println()
		return false, nil
	}

	// Token is valid - handle account ID
	if cfg.CloudflareAccount != "" && !isPlaceholder(cfg.CloudflareAccount) {
		return validateOrFixAccountID(cfg, tokenName)
	}
	
	return fetchAndSaveAccountID(cfg.CloudflareToken, tokenName)
}

// validateOrFixAccountID validates existing account ID or auto-fixes if mismatched
func validateOrFixAccountID(cfg *EnvConfig, tokenName string) (bool, error) {
	accountName, err := ValidateCloudflareAccount(cfg.CloudflareToken, cfg.CloudflareAccount)
	if err == nil {
		// Account ID is valid - setup complete!
		fmt.Println(Success(fmt.Sprintf("Cloudflare API token is valid: %s", tokenName)))
		fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
		fmt.Println()
		return true, ensureProjectName()
	}

	// Account ID doesn't match - try to auto-fix
	fmt.Println(Error(fmt.Sprintf("Account ID validation failed: %v", err)))
	fmt.Println()
	fmt.Println(Colorize("Token is valid but account ID doesn't match", ColorYellow))
	fmt.Println(Colorize("Fetching correct account ID for this token...", ColorYellow))
	fmt.Println()

	accountID, accountName, err := GetCloudflareAccounts(cfg.CloudflareToken)
	if err != nil {
		fmt.Println(Error(fmt.Sprintf("Could not fetch account ID: %v", err)))
		fmt.Println()
		if promptYesNo("Keep token without account ID?", true) {
			fmt.Println("✓ Keeping existing token")
			fmt.Println()
			return true, nil
		}
		fmt.Println("Will prompt for new token...")
		fmt.Println()
		return false, nil
	}

	// Save the correct account ID
	if err := UpdateEnv(EnvCloudflareAccount, accountID); err != nil {
		return true, fmt.Errorf("failed to save account ID: %w", err)
	}
	
	fmt.Println(Success(fmt.Sprintf("Account ID automatically configured: %s", accountName)))
	fmt.Println(Colorize(fmt.Sprintf("  Old ID: %s (was incorrect)", cfg.CloudflareAccount), ColorGray))
	fmt.Println(Colorize(fmt.Sprintf("  New ID: %s", accountID), ColorGray))
	fmt.Println()
	return true, ensureProjectName()
}

// fetchAndSaveAccountID fetches and saves account ID for a valid token
func fetchAndSaveAccountID(token, tokenName string) (bool, error) {
	fmt.Println(Success(fmt.Sprintf("Cloudflare API token is valid: %s", tokenName)))
	fmt.Println()
	fmt.Println("Fetching account information...")
	
	accountID, accountName, err := GetCloudflareAccounts(token)
	if err != nil {
		fmt.Println(Colorize(fmt.Sprintf("Could not fetch account ID: %v", err), ColorYellow))
		fmt.Println()
		if promptYesNo("Keep token without account ID?", true) {
			fmt.Println("✓ Keeping existing token")
			fmt.Println()
			return true, nil
		}
		fmt.Println("Will prompt for new token...")
		fmt.Println()
		return false, nil
	}

	// Save the account ID
	if err := UpdateEnv(EnvCloudflareAccount, accountID); err != nil {
		return true, fmt.Errorf("failed to save account ID: %w", err)
	}
	
	fmt.Println(Success(fmt.Sprintf("Account ID automatically configured: %s", accountName)))
	fmt.Println(Colorize(fmt.Sprintf("  ID: %s", accountID), ColorGray))
	fmt.Println()
	return true, ensureProjectName()
}

// promptForCloudflareToken prompts user for new Cloudflare token and validates it
func promptForCloudflareToken() error {
	_, repoName, _ := GetRepositoryInfo()
	
	for {
		showCloudflareInstructions(repoName)
		
		token := promptString("Paste your Cloudflare API token (or press Enter to skip)")
		if token == "" {
			fmt.Println()
			fmt.Println("⊘ Skipped - you can add it later in .env")
			fmt.Println("   Without this token, you cannot deploy to Cloudflare Pages.")
			fmt.Println()
			return nil
		}

		// Save and validate token
		if err := UpdateEnv(EnvCloudflareToken, token); err != nil {
			return err
		}
		fmt.Println()
		fmt.Printf("✓ Token saved to %s\n", EnvCloudflareToken)
		fmt.Println()

		fmt.Println("Validating Cloudflare credentials...")
		tokenName, err := ValidateCloudflareToken(token)
		if err != nil {
			fmt.Println(Error(err.Error()))
			fmt.Println()
			fmt.Println(Colorize("Please try again or press Enter to skip...", ColorYellow))
			fmt.Println()
			continue
		}

		fmt.Println(Success(fmt.Sprintf("Cloudflare API token is valid: %s", tokenName)))

		// Try to fetch and save account ID
		if err := handleAccountIDForNewToken(token); err != nil {
			return err
		}
		return nil
	}
}

// handleAccountIDForNewToken handles account ID validation/fetching for newly entered token
func handleAccountIDForNewToken(token string) error {
	cfg, _ := LoadEnv()
	
	// Check if there's an existing account ID to validate
	if cfg.CloudflareAccount != "" && !isPlaceholder(cfg.CloudflareAccount) {
		accountName, err := ValidateCloudflareAccount(token, cfg.CloudflareAccount)
		if err == nil {
			fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
			fmt.Println()
			return nil
		}
		
		// Account ID mismatch - fetch correct one
		fmt.Println(Error(fmt.Sprintf("Account ID validation failed: %v", err)))
		fmt.Println()
		fmt.Println(Colorize("Token is valid but account ID doesn't match", ColorYellow))
		fmt.Println(Colorize("Fetching correct account ID for this token...", ColorYellow))
		fmt.Println()
	}
	
	// Fetch account ID
	fmt.Println()
	fmt.Println("Fetching account information...")
	accountID, accountName, err := GetCloudflareAccounts(token)
	if err != nil {
		fmt.Println(Colorize(fmt.Sprintf("Could not fetch account ID: %v", err), ColorYellow))
		fmt.Println(Colorize("You can add it manually to .env later if needed", ColorYellow))
		fmt.Println()
		return nil
	}

	// Save the account ID
	if err := UpdateEnv(EnvCloudflareAccount, accountID); err != nil {
		return fmt.Errorf("failed to save account ID: %w", err)
	}

	if cfg.CloudflareAccount != "" && !isPlaceholder(cfg.CloudflareAccount) {
		fmt.Println(Success(fmt.Sprintf("Account ID automatically configured: %s", accountName)))
		fmt.Println(Colorize(fmt.Sprintf("  Old ID: %s (was incorrect)", cfg.CloudflareAccount), ColorGray))
		fmt.Println(Colorize(fmt.Sprintf("  New ID: %s", accountID), ColorGray))
	} else {
		fmt.Println(Success(fmt.Sprintf("Account ID automatically configured: %s", accountName)))
		fmt.Println(Colorize(fmt.Sprintf("  ID: %s", accountID), ColorGray))
	}
	fmt.Println()
	return nil
}

func setupClaude() error {
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println("Step 2 of 2: Claude API Key (Optional)")
	fmt.Printf("Setting: %s\n", EnvClaudeAPIKey)
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()

	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	// Check if key exists and is not placeholder
	if cfg.ClaudeAPIKey != "" && cfg.ClaudeAPIKey != PlaceholderKey {
		// Validate existing key first
		fmt.Println("Validating existing Claude API key...")
		if err := ValidateClaudeAPIKey(cfg.ClaudeAPIKey); err != nil {
			fmt.Println(Error(fmt.Sprintf("API key validation failed: %v", err)))
			fmt.Println()
			fmt.Println(Colorize("Will prompt for new key...", ColorYellow))
			fmt.Println()
			cfg.ClaudeAPIKey = ""
		} else {
			// Validation passed - key is valid
			fmt.Println("✓ Claude API key is valid")
			fmt.Println("✓ API key has active credits")
			fmt.Println()
			// Continue to workspace setup
			cfg.ClaudeAPIKey = "valid" // Mark as valid to skip the loop
		}
	}

	// Loop until valid key or skip
	for cfg.ClaudeAPIKey == "" || cfg.ClaudeAPIKey == PlaceholderKey {
		// Get repository name for workspace suggestion
		_, repoName, _ := GetRepositoryInfo()
		showClaudeInstructions(repoName)

		key := promptString("Paste your Claude API key (or press Enter to skip)")
		if key == "" {
			fmt.Println()
			fmt.Println("⊘ Skipped - you can add it later in .env")
			fmt.Println("   Without this key, you cannot use automated translation.")
			fmt.Println()
			return nil
		}

		// Save key
		if err := UpdateEnv(EnvClaudeAPIKey, key); err != nil {
			return err
		}
		fmt.Println()
		fmt.Printf("✓ Key saved to %s\n", EnvClaudeAPIKey)
		fmt.Println()

		// Validate key
		fmt.Println("Validating Claude API key...")
		if err := ValidateClaudeAPIKey(key); err != nil {
			fmt.Println(Error(err.Error()))
			fmt.Println()
			fmt.Println(Colorize("Please try again or press Enter to skip...", ColorYellow))
			fmt.Println()
			continue
		}

		fmt.Println("✓ Claude API key is valid")
		fmt.Println("✓ API key has active credits")
		fmt.Println()
		break
	}

	// Ask for workspace name - use repo name as default
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("Setting: %s (recommended)\n", EnvClaudeWorkspace)
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()

	// Get repository name for default workspace
	_, repoName, _ := GetRepositoryInfo()
	if repoName == "" {
		repoName = "my-project"
	}

	fmt.Println("Claude Workspaces help keep your project's API usage isolated")
	fmt.Println("and organized. Use your project name as the workspace name.")
	fmt.Println()
	fmt.Printf("Workspace name [default: %s]: ", repoName)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	workspace := strings.TrimSpace(input)

	// Use default if empty
	if workspace == "" {
		workspace = repoName
	}

	if err := UpdateEnv(EnvClaudeWorkspace, workspace); err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("✓ Workspace saved to %s: %s\n", EnvClaudeWorkspace, workspace)
	fmt.Println()

	return nil
}

// ensureProjectName ensures CLOUDFLARE_PROJECT_NAME is set with repository name as default
func ensureProjectName() error {
	cfg, err := LoadEnv()
	if err != nil {
		return err
	}

	// Check if project name is already set (not a placeholder)
	if cfg.CloudflareProject != "" && !isPlaceholder(cfg.CloudflareProject) {
		fmt.Println(Success(fmt.Sprintf("Cloudflare project name: %s", cfg.CloudflareProject)))
		fmt.Println()
		return nil
	}

	// Prompt for project name
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("Setting: %s\n", EnvCloudflareProject)
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()

	// Get repository name for default project name
	_, repoName, _ := GetRepositoryInfo()
	if repoName == "" {
		repoName = "my-project"
	}

	fmt.Println("This is the name of your Cloudflare Pages project.")
	fmt.Println("Your site will be deployed to: <project-name>.pages.dev")
	fmt.Println()
	fmt.Printf("Project name [default: %s]: ", repoName)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	projectName := strings.TrimSpace(input)

	// Use default if empty
	if projectName == "" {
		projectName = repoName
	}

	if err := UpdateEnv(EnvCloudflareProject, projectName); err != nil {
		return err
	}
	fmt.Println()
	fmt.Printf("✓ Project name saved to %s: %s\n", EnvCloudflareProject, projectName)
	fmt.Printf("  Your site will deploy to: %s.pages.dev\n", projectName)
	fmt.Println()

	return nil
}

func showCloudflareInstructions(repoName string) {
	fmt.Println("You need a Cloudflare API token to deploy to Cloudflare Pages.")
	fmt.Println()
	fmt.Println("Follow these steps:")
	fmt.Println()
	fmt.Println("  1. Open: https://dash.cloudflare.com/login")
	fmt.Println("     → Log in (or sign up for free)")
	fmt.Println()
	fmt.Println("  2. Open: https://dash.cloudflare.com/profile/api-tokens")
	fmt.Println("     → Click the blue 'Create Token' button")
	fmt.Println()
	fmt.Println("  3. Click 'Create Custom Token' (there's no Cloudflare Pages template)")

	// Build token name - use repo name if available
	tokenName := repoName
	if tokenName == "" {
		tokenName = "my-project"
	}
	fmt.Printf("     → Token name: Use your project name: '%s'\n", tokenName)
	fmt.Println("     → Permissions: Add these 3 permissions:")
	fmt.Println("       • Account | Cloudflare Pages | Edit")
	fmt.Println("       • Account | Account Settings | Read")
	fmt.Println("       • User | API Tokens | Read")
	fmt.Println("     → Account Resources:")
	fmt.Println("       • Select your specific account")
	fmt.Println("     → Click 'Continue to summary'")
	fmt.Println("     → Click 'Create Token'")
	fmt.Println()
	fmt.Println("  4. Copy the token")
	fmt.Println("     → It starts with letters/numbers (40+ chars)")
	fmt.Println("     → COPY IT NOW - you can't see it again!")
	fmt.Println()
}

func showClaudeInstructions(repoName string) {
	fmt.Println("You need a Claude API key for automated translation.")
	fmt.Println()
	fmt.Println("Follow these steps:")
	fmt.Println()
	fmt.Println("  1. Open: https://console.anthropic.com/")
	fmt.Println("     → Click 'Sign Up' (or 'Log In' if you have account)")
	fmt.Println("     → Complete registration with email")
	fmt.Println()
	fmt.Println("  2. Open: https://console.anthropic.com/settings/billing")
	fmt.Println("     → Click 'Purchase credits'")
	fmt.Println("     → Add at least $5 USD (minimum)")
	fmt.Println("     → REQUIRED: You CANNOT create API keys without credits")
	fmt.Println("     → The $5 will last for many translations")
	fmt.Println()
	fmt.Println("  3. Create a Workspace for this project:")
	fmt.Println("     → Open: https://console.anthropic.com/settings/workspaces")
	fmt.Println("     → Click 'Create Workspace'")

	// Use actual project name
	workspaceName := repoName
	if workspaceName == "" {
		workspaceName = "my-project"
	}
	fmt.Printf("     → Name it: '%s'\n", workspaceName)
	fmt.Println("     → This keeps this project's usage separate")
	fmt.Println()
	fmt.Println("  4. Create API key in your workspace:")
	fmt.Println("     → In your workspace, go to 'API Keys' tab")
	fmt.Println("     → Click 'Create Key'")
	fmt.Println("     → Give it a name (e.g., 'Translation')")
	fmt.Println("     → Click 'Create Key'")
	fmt.Println()
	fmt.Println("  5. Copy the key")
	fmt.Println("     → It starts with 'sk-ant-api03-...'")
	fmt.Println("     → COPY IT NOW - you can't see it again!")
	fmt.Println()
}

func showNextSteps() {
	// Get .env file path
	envPath, err := GetEnvPath()
	if err != nil {
		envPath = ".env"
	}

	printHeader("Setup Complete!", fmt.Sprintf("Configuration saved to: %s", envPath))

	fmt.Println("Next steps:")
	fmt.Println()
	fmt.Println("  task setup      # Install Hugo, Bun, and dependencies")
	fmt.Println("  task dev        # Start development server")
	fmt.Println("  task build      # Build the site")
	fmt.Println()
	fmt.Println("Optional commands:")
	fmt.Println("  task env:local:list    # Show current configuration")
	fmt.Println("  task env:gh:push       # Push secrets to GitHub for CI/CD")
	fmt.Println()
}

func promptString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func promptYesNo(prompt string, defaultYes bool) bool {
	suffix := "[Y/n]"
	if !defaultYes {
		suffix = "[y/N]"
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s %s: ", prompt, suffix)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(strings.ToLower(text))

	if text == "" {
		return defaultYes
	}

	return text == "y" || text == "yes"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
