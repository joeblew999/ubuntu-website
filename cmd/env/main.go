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
	case "wizard":
		if err := env.RunWizard(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "show":
		if err := env.ShowConfig(); err != nil {
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

	case "sync-secrets":
		dryRun, force, validate := parseSyncSecretsFlags()
		if err := env.SyncSecrets(dryRun, force, validate); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "show-remote":
		if err := env.ShowRemoteSecrets(); err != nil {
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
	validate = true // Default to true

	for _, arg := range os.Args[2:] {
		switch arg {
		case "--check", "--dry-run":
			dryRun = true
		case "--force":
			force = true
		case "--no-validate":
			validate = false
		}
	}

	return dryRun, force, validate
}

func printUsage() {
	fmt.Println("Usage: setup <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  wizard              Run interactive setup wizard")
	fmt.Println("  show                Show current configuration")
	fmt.Println("  validate            Validate all credentials")
	fmt.Println("  validate-cloudflare Validate Cloudflare token only")
	fmt.Println("  validate-claude     Validate Claude API key only")
	fmt.Println("  sync-secrets        Sync .env secrets to GitHub repository")
	fmt.Println("  show-remote         Show GitHub secrets status")
	fmt.Println()
	fmt.Println("Options for sync-secrets:")
	fmt.Println("  --check, --dry-run  Show what would be synced without syncing")
	fmt.Println("  --force             Overwrite existing secrets without prompting")
	fmt.Println("  --no-validate       Skip credential validation before syncing")
}
