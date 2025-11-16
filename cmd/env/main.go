package main

import (
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/env"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "local-setup":
		if err := env.RunWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "local-list":
		if err := env.ShowConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "gh-list":
		if err := env.ShowRemoteSecrets(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "gh-push":
		dryRun, force, validate := parseSyncSecretsFlags()
		if err := env.SyncSecrets(dryRun, force, validate); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "validate":
		if err := env.ValidateAll(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "validate-cloudflare":
		if err := env.ValidateCloudflareCredentials(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "validate-claude":
		if err := env.ValidateClaudeCredentials(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
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
	fmt.Println("Local Commands:")
	fmt.Println("  local-setup         Setup local .env - interactive wizard")
	fmt.Println("  local-list          List local .env configuration")
	fmt.Println()
	fmt.Println("GitHub Commands:")
	fmt.Println("  gh-list             List GitHub secrets")
	fmt.Println("  gh-push             Push .env to GitHub secrets")
	fmt.Println()
	fmt.Println("Validation Commands:")
	fmt.Println("  validate            Validate all credentials")
	fmt.Println("  validate-cloudflare Validate Cloudflare token only")
	fmt.Println("  validate-claude     Validate Claude API key only")
	fmt.Println()
	fmt.Println("Options for gh-push:")
	fmt.Println("  --check, --dry-run  Show what would be synced without syncing")
	fmt.Println("  --no-force          Prompt before overwriting existing secrets")
	fmt.Println("  --no-validate       Skip credential validation before syncing")
}
