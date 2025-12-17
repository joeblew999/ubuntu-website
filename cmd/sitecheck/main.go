// sitecheck checks site reachability from multiple global locations.
//
// Uses the check-host.net API to verify a URL is accessible from
// different geographic regions (US, EU, Asia, etc.).
//
// Usage:
//
//	go run cmd/sitecheck/main.go                           # HTTP check (default)
//	go run cmd/sitecheck/main.go -type dns                 # DNS resolution check
//	go run cmd/sitecheck/main.go -type tcp                 # TCP port 443 check
//	go run cmd/sitecheck/main.go -type redirect            # Apex->www redirect check
//	go run cmd/sitecheck/main.go -type all                 # Run all checks
//	go run cmd/sitecheck/main.go -url https://example.com  # Check custom URL
//	go run cmd/sitecheck/main.go -github-issue             # Output markdown for GitHub Issue
//	task site:check                                        # Via Taskfile
//
// GitHub Actions:
//
//	Runs every 6 hours via .github/workflows/site-monitor.yml
//	Creates a GitHub Issue when failures detected or significant changes occur
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/sitecheck"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	exitCode := sitecheck.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
