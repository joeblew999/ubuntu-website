// Command google provides unified CLI for all Google services.
//
// Services:
//   - auth:     MCP configuration and OAuth authentication
//   - calendar: Google Calendar management
//   - gmail:    Email sending via Gmail
//   - drive:    Google Drive file management
//   - sheets:   Google Sheets operations
//   - docs:     Google Docs operations
//   - slides:   Google Slides operations
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/googlecli"
)

var version = "dev"

func main() {
	exitCode := googlecli.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
