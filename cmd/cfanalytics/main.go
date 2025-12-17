// cfanalytics fetches Cloudflare Web Analytics and reports changes.
//
// Compares current metrics to the previous run (stored in .analytics-state.json)
// and reports significant changes (>20% threshold). Can run locally or in GitHub
// Actions to create issues when traffic changes significantly.
//
// Usage:
//
//	go run cmd/cfanalytics/main.go                   # Print report to terminal
//	go run cmd/cfanalytics/main.go -webhook URL     # Post to webhook if changed
//	go run cmd/cfanalytics/main.go -days 14         # Compare last 14 days
//	go run cmd/cfanalytics/main.go -github-issue   # Output markdown for GitHub Issue
//	task seo:report                                  # Via Taskfile
//
// GitHub Actions:
//
//	Runs weekly via .github/workflows/analytics-report.yml
//	Creates a GitHub Issue when visits or pageviews change >20%
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/cfanalytics"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	exitCode := cfanalytics.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
