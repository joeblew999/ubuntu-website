// Command google is one CLI for the Google services the site's work uses, on the accounts signed in
// (~/.google-mcp-accounts, shared with the Google MCP server): Gmail, Calendar, Drive, Sheets, Docs, Slides, and
// signing in. Each service is a file of its own and a task of .mise/tasks/google.toml; its subcommands and flags are
// internal/google's, passed through:
//
//	mise run google:calendar today
//	mise run google:gmail send --to x@y.com --subject Hello --body Message
//	mise run google:auth login
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/google"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "auth":
		auth()
	case "calendar":
		calendar()
	case "gmail":
		gmail()
	case "drive":
		drive()
	case "sheets":
		sheets()
	case "docs":
		docs()
	case "slides":
		slides()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/google's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return google.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
