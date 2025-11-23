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
	case "admin":
		err = web.ServeSetupGUI()
	case "admin-mock":
		err = web.ServeSetupGUIMock()
	case "validate":
		exitCode := env.RunValidateFast()
		os.Exit(exitCode)
	case "validate-deep":
		exitCode := env.RunValidateDeep()
		os.Exit(exitCode)
	case "build":
		err = env.RunBuild()
	case "deploy-preview":
		err = env.RunDeployPreview()
	case "deploy-production":
		err = env.RunDeployProduction()
	case "domain-status":
		err = env.RunDomainStatus()
	case "caddy-start":
		err = env.EnsureCaddyRunning()
	case "caddy-stop":
		err = env.StopCaddy()
	case "caddy-status":
		err = env.PrintCaddyStatus()
	case "kill-all":
		err = env.KillAll()
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

func printUsage() {
	fmt.Println("Usage: go run cmd/env/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  admin               Open admin GUI for environment setup (starts Caddy + Via GUI)")
	fmt.Println("  admin-mock          Open admin GUI with mock validation (for testing)")
	fmt.Println()
	fmt.Println("  validate            Validate .env file (fast - format checks only)")
	fmt.Println("  validate-deep       Validate .env file (deep - includes API verification)")
	fmt.Println()
	fmt.Println("  build               Build Hugo site (starts Caddy + Hugo server)")
	fmt.Println("  deploy-preview      Build + deploy to Cloudflare Pages preview")
	fmt.Println("  deploy-production   Build + deploy to Cloudflare Pages production (main branch)")
	fmt.Println("  domain-status       Check custom domain status and troubleshoot Error 1014")
	fmt.Println()
	fmt.Println("  caddy-start         Start Caddy HTTPS server (port 443)")
	fmt.Println("  caddy-stop          Stop Caddy HTTPS server")
	fmt.Println("  caddy-status        Check if Caddy is running")
	fmt.Println()
	fmt.Println("  kill-all            Stop all services (Caddy, Hugo, Via GUI) and clean up ports")
}
