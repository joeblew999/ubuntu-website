// playwright - Reusable browser automation CLI
//
// A thin CLI wrapper around browser automation packages.
// Domain logic lives in internal/browser and internal/google/gmail.
//
// Usage:
//
//	playwright oauth <url>             Start OAuth flow with callback server
//	playwright open <url>              Open URL in Playwright browser
//	playwright screenshot <url> <file> Take screenshot of URL
//	playwright install                 Install Playwright browsers
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/playwright"
)

var version = "dev"

func main() {
	exitCode := playwright.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
