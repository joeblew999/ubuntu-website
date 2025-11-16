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

	// Check if token exists and is not placeholder
	if cfg.CloudflareToken != "" && cfg.CloudflareToken != PlaceholderToken {
		// Validate existing token first
		fmt.Println("Validating existing Cloudflare credentials...")
		tokenName, err := ValidateCloudflareToken(cfg.CloudflareToken)
		if err != nil {
			fmt.Println(Error(fmt.Sprintf("Token validation failed: %v", err)))
			fmt.Println()
			fmt.Println(Colorize("Will prompt for new token...", ColorYellow))
			fmt.Println()
			cfg.CloudflareToken = ""
		} else {
			// Validate or fetch account ID
			if cfg.CloudflareAccount != "" && !isPlaceholder(cfg.CloudflareAccount) {
				// Account ID exists - validate it
				if accountName, err := ValidateCloudflareAccount(cfg.CloudflareToken, cfg.CloudflareAccount); err == nil {
					fmt.Println(Success(fmt.Sprintf("Cloudflare API token is valid: %s", tokenName)))
					fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
					fmt.Println()
					return nil
				} else {
					fmt.Println(Error(fmt.Sprintf("Account ID validation failed: %v", err)))
					fmt.Println()
					fmt.Println(Colorize("Will prompt for new token...", ColorYellow))
					fmt.Println()
					cfg.CloudflareToken = ""
				}
			} else {
				// Token valid but no account ID - try to fetch it
				fmt.Println(Success(fmt.Sprintf("Cloudflare API token is valid: %s", tokenName)))
				fmt.Println()
				fmt.Println("Fetching account information...")
				accountID, accountName, err := GetCloudflareAccounts(cfg.CloudflareToken)
				if err != nil {
					fmt.Println(Colorize(fmt.Sprintf("Could not fetch account ID: %v", err), ColorYellow))
					fmt.Println()
					keep := promptYesNo("Keep token without account ID?", true)
					if keep {
						fmt.Println("✓ Keeping existing token")
						fmt.Println()
						return nil
					} else {
						cfg.CloudflareToken = ""
						fmt.Println("Will prompt for new token...")
						fmt.Println()
					}
				} else {
					// Save the account ID
					if err := UpdateEnv(EnvCloudflareAccount, accountID); err != nil {
						fmt.Println(Error(fmt.Sprintf("Failed to save account ID: %v", err)))
						fmt.Println()
						return err
					}
					fmt.Println(Success(fmt.Sprintf("Account ID automatically configured: %s", accountName)))
					fmt.Println(Colorize(fmt.Sprintf("  ID: %s", accountID), ColorGray))
					fmt.Println()
					return nil
				}
			}
		}
	}

	// Loop until valid token or skip
	for cfg.CloudflareToken == "" || cfg.CloudflareToken == PlaceholderToken {
		// Get repo name for token suggestion (ignore errors - it's just a suggestion)
		_, repoName, _ := GetRepositoryInfo()
		showCloudflareInstructions(repoName)

		token := promptString("Paste your Cloudflare API token (or press Enter to skip)")
		if token == "" {
			fmt.Println()
			fmt.Println("⊘ Skipped - you can add it later in .env")
			fmt.Println("   Without this token, you cannot deploy to Cloudflare Pages.")
			fmt.Println()
			return nil
		}

		// Save token
		if err := UpdateEnv(EnvCloudflareToken, token); err != nil {
			return err
		}
		fmt.Println()
		fmt.Printf("✓ Token saved to %s\n", EnvCloudflareToken)
		fmt.Println()

		// Validate token
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

		// Fetch or validate account ID
		cfg, _ := LoadEnv()
		if cfg.CloudflareAccount != "" && !isPlaceholder(cfg.CloudflareAccount) {
			// Account ID exists - validate it
			if accountName, err := ValidateCloudflareAccount(token, cfg.CloudflareAccount); err == nil {
				fmt.Println(Success(fmt.Sprintf("Account ID is valid: %s", accountName)))
				fmt.Println()
				break
			} else {
				fmt.Println(Error(fmt.Sprintf("Account ID validation failed: %v", err)))
				fmt.Println()
				fmt.Println(Colorize("Please try again or press Enter to skip...", ColorYellow))
				fmt.Println()
				continue
			}
		} else {
			// No account ID - try to fetch it automatically
			fmt.Println()
			fmt.Println("Fetching account information...")
			accountID, accountName, err := GetCloudflareAccounts(token)
			if err != nil {
				fmt.Println(Colorize(fmt.Sprintf("Could not fetch account ID automatically: %v", err), ColorYellow))
				fmt.Println(Colorize("You can add it manually to .env later if needed", ColorYellow))
				fmt.Println()
				break
			}

			// Save the account ID
			if err := UpdateEnv(EnvCloudflareAccount, accountID); err != nil {
				fmt.Println(Error(fmt.Sprintf("Failed to save account ID: %v", err)))
				fmt.Println()
				continue
			}

			fmt.Println(Success(fmt.Sprintf("Account ID automatically configured: %s", accountName)))
			fmt.Println(Colorize(fmt.Sprintf("  ID: %s", accountID), ColorGray))
			fmt.Println()
			break
		}
	}

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
			return nil
		}
	}

	// Loop until valid key or skip
	for cfg.ClaudeAPIKey == "" || cfg.ClaudeAPIKey == PlaceholderKey {
		showClaudeInstructions()

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

	// Ask for workspace name
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Printf("Setting: %s (recommended)\n", EnvClaudeWorkspace)
	fmt.Println("────────────────────────────────────────────────────────────")
	fmt.Println()
	fmt.Println("Enter your Claude Workspace name to keep this project's")
	fmt.Println("usage isolated and organized.")
	fmt.Println()

	workspace := promptString("Workspace name (or press Enter to skip)")
	if workspace != "" {
		if err := UpdateEnv(EnvClaudeWorkspace, workspace); err != nil {
			return err
		}
		fmt.Println()
		fmt.Printf("✓ Workspace saved to %s\n", EnvClaudeWorkspace)
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("⊘ Skipped - using default workspace")
		fmt.Println()
	}

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

	// Build token name suggestion
	tokenSuggestion := "My Pages Deploy"
	if repoName != "" {
		tokenSuggestion = repoName
	}
	fmt.Printf("     → Token name: Give it a descriptive name (e.g., '%s')\n", tokenSuggestion)
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

func showClaudeInstructions() {
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
	fmt.Println("     → Name it after your project (e.g., 'ubuntu-website')")
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
