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
	var err error

	switch command {
	case "local-setup":
		err = env.RunWizard()
	case "local-list":
		err = env.ShowConfig()
	case "gh-list":
		err = env.ShowRemoteSecrets()
	case "gh-push":
		dryRun, force, validate := parseSyncSecretsFlags()
		err = env.SyncSecrets(dryRun, force, validate)
	case "validate":
		err = env.ValidateAll()
	case "validate-cloudflare":
		err = env.ValidateCloudflareCredentials()
	case "validate-claude":
		err = env.ValidateClaudeCredentials()
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
