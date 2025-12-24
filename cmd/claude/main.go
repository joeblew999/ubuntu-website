// Command claude provides unified CLI for Claude Code management.
//
// This includes:
// - Claude CLI binary detection and version management
// - MCP server management (list, restart, status)
// - Migration from npm/bun to native binary
//
// Commands:
//   - status: Show Claude CLI installation status
//   - check: Check if native binary is installed
//   - detect: List all detected installations
//   - cleanup: Show what needs cleanup
//   - version-check: Verify installed version meets minimum
//   - ensure: Ensure native CLI is installed and meets version
//   - mcp: Manage MCP servers
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/claude"
)

var version = "dev"

func main() {
	cli := claude.NewCLI(version, os.Stdout, os.Stderr)
	exitCode := cli.Run(os.Args[1:])
	os.Exit(exitCode)
}
