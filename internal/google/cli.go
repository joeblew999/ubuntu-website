// Package googlecli provides unified CLI for all Google services.
//
// Services:
//   - auth:     MCP configuration and OAuth authentication
//   - calendar: Google Calendar management
//   - gmail:    Email sending via Gmail
//   - drive:    Google Drive file management
//   - sheets:   Google Sheets operations
//   - docs:     Google Docs operations
//   - slides:   Google Slides operations
package google

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Run executes the google CLI with the given arguments.
// Returns exit code: 0 for success, non-zero for errors.
func Run(args []string, version string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stdout)
		return 1
	}

	ctx := &cliContext{
		stdout: stdout,
		stderr: stderr,
	}

	service := args[1]

	switch service {
	case "auth":
		ctx.handleAuth(args[2:])
	case "calendar", "cal":
		ctx.handleCalendar(args[2:])
	case "gmail", "mail":
		ctx.handleGmail(args[2:])
	case "drive":
		ctx.handleDrive(args[2:])
	case "sheets":
		ctx.handleSheets(args[2:])
	case "docs":
		ctx.handleDocs(args[2:])
	case "slides":
		ctx.handleSlides(args[2:])
	case "-h", "--help", "help":
		printUsage(stdout)
	default:
		fmt.Fprintf(stderr, "Unknown service: %s\n", service)
		printUsage(stderr)
		return 1
	}
	return 0
}

// cliContext holds the CLI state
type cliContext struct {
	stdout io.Writer
	stderr io.Writer
}

func (c *cliContext) exitError(msg string) {
	fmt.Fprintf(c.stderr, "Error: %s\n", msg)
	os.Exit(1)
}

func (c *cliContext) outputJSON(v interface{}) {
	enc := json.NewEncoder(c.stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// hasFlag checks if a flag is present in args
func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

// getFlagValue gets a flag value like --parent=ID
func getFlagValue(args []string, prefix string) string {
	for _, a := range args {
		if len(a) > len(prefix) && a[:len(prefix)] == prefix {
			return a[len(prefix):]
		}
	}
	return ""
}

// filterFlags removes flags from args, returning only positional args
func filterFlags(args []string) []string {
	var result []string
	for _, a := range args {
		if len(a) == 0 || a[0] != '-' {
			result = append(result, a)
		}
	}
	return result
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, `google - Unified CLI for Google services

Usage:
  google <service> <command> [arguments]

Services:
  auth      - MCP configuration and OAuth authentication
  calendar  - Google Calendar management (alias: cal)
  gmail     - Email via Gmail API (alias: mail)
  drive     - Google Drive file management
  sheets    - Google Sheets operations
  docs      - Google Docs operations
  slides    - Google Slides presentations

Quick Start:
  google auth guide              # Setup instructions
  google auth check              # Check configuration
  google auth login              # Authenticate with Google

Auth Commands:
  google auth add                Add MCP server to Claude Code
  google auth remove             Remove MCP server
  google auth status             Show configuration
  google auth guide              Show setup guide
  google auth open <target>      Open Cloud Console page
  google auth login              OAuth login flow

Calendar Commands:
  google calendar list           List upcoming events
  google calendar today          List today's events
  google calendar create ...     Create calendar event
  google calendar check          Verify API connection

Gmail Commands:
  google gmail list             List recent messages
  google gmail send ...          Send email via API
  google gmail compose ...       Open Gmail compose
  google gmail check             Verify API connection

Drive Commands:
  google drive list [FOLDER_ID]  List files
  google drive upload FILE       Upload a file
  google drive download FILE_ID  Download a file

Sheets Commands:
  google sheets get ID RANGE     Get cell values
  google sheets set ID RANGE ... Update cells

Docs Commands:
  google docs get ID             Get document content
  google docs create TITLE       Create new document

Slides Commands:
  google slides get ID           Get presentation
  google slides create TITLE     Create presentation

Global Options:
  --json    Output as JSON

Examples:
  google auth login
  google calendar today
  google gmail send --to=user@example.com --subject="Hi" --body="Hello"
  google drive list
  google sheets get 1abc123 "Sheet1!A1:D10"`)
}
