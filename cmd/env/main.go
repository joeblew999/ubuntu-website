package main

import (
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/env"
	"github.com/joeblew999/ubuntu-website/internal/env/web"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	var err error

	switch command {
	case "setup":
		err = setupLocal()
	case "list":
		err = env.List()
	case "web-gui":
		err = web.ServeSetupGUI()
	case "web-gui-mock":
		err = web.ServeSetupGUIMock()
	case "gh-push":
		dryRun, force, validate := parseSyncSecretsFlags()
		err = env.PushGithub(dryRun, force, validate)
	// Deprecated aliases
	case "local-setup":
		err = setupLocal()
	case "local-list", "gh-list":
		err = env.List()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setupLocal() error {
	// Check/create .env file
	if !env.EnvExists() {
		fmt.Println("Creating .env file...")
		if err := env.CreateEnv(); err != nil {
			return err
		}
		fmt.Println("✓ Created .env file")
		fmt.Println()
	}

	// Create service for config operations
	svc := env.NewService(false) // false = real validation, not mock

	// Load and validate credentials using service
	cfg, err := svc.GetCurrentConfig()
	if err != nil {
		return fmt.Errorf("failed to load .env: %w", err)
	}

	// Validate all credentials using service
	resultsMap := svc.ValidateConfig(cfg)
	results := env.ResultsToSlice(resultsMap)
	env.PrintValidationResults(results, cfg)

	// If all valid, we're done
	if !env.HasInvalidCredentials(results) {
		return nil
	}

	// Otherwise, show web GUI links for invalid credentials
	invalidFields := env.GetInvalidFields(results)

	fmt.Println()
	fmt.Println("Fix invalid credentials in web GUI:")
	fmt.Println()

	// Show specific links based on which credentials are invalid
	hasCloudflareIssues := false
	hasClaudeIssues := false

	for _, field := range invalidFields {
		if field == "Cloudflare API Token" || field == "Cloudflare Account ID" {
			hasCloudflareIssues = true
		}
		if field == "Claude API Key" {
			hasClaudeIssues = true
		}
	}

	if hasCloudflareIssues {
		fmt.Printf("  Cloudflare → http://localhost:3000/cloudflare\n")
	}
	if hasClaudeIssues {
		fmt.Printf("  Claude     → http://localhost:3000/claude\n")
	}

	fmt.Println()
	fmt.Println("Start web GUI: go run cmd/env/main.go web-gui")
	fmt.Println()

	return fmt.Errorf("validation failed")
}

func parseSyncSecretsFlags() (dryRun, force, validate bool) {
	force = true    // Default to true - overwrite existing secrets
	validate = true // Default to true - validate before pushing

	for _, arg := range os.Args[2:] {
		switch arg {
		case "--check", "--dry-run":
			dryRun = true
		case "--no-force":
			force = false
		case "--no-validate":
			validate = false
		}
	}

	return dryRun, force, validate
}

func printUsage() {
	fmt.Println("Usage: go run cmd/env/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  setup               Validate local .env credentials")
	fmt.Println("  list                List local and GitHub configuration")
	fmt.Println("  web-gui             Open web GUI for environment setup")
	fmt.Println("  web-gui-mock        Open web GUI with mock validation (for testing)")
	fmt.Println("  gh-push             Push .env to GitHub secrets")
	fmt.Println()
	fmt.Println("Options for gh-push:")
	fmt.Println("  --check, --dry-run  Show what would be synced without syncing")
	fmt.Println("  --no-force          Prompt before overwriting existing secrets")
	fmt.Println("  --no-validate       Skip credential validation before syncing")
}
