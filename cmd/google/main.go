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
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	service := os.Args[1]

	switch service {
	case "auth":
		handleAuth(os.Args[2:])
	case "calendar", "cal":
		handleCalendar(os.Args[2:])
	case "gmail", "mail":
		handleGmail(os.Args[2:])
	case "drive":
		handleDrive(os.Args[2:])
	case "sheets":
		handleSheets(os.Args[2:])
	case "docs":
		handleDocs(os.Args[2:])
	case "slides":
		handleSlides(os.Args[2:])
	case "-h", "--help", "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown service: %s\n", service)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`google - Unified CLI for Google services

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

func exitError(msg string) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	os.Exit(1)
}

func outputJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
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
