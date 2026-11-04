package main

import "os"

// mcp is claude-cli mcp, the task claude-cli:mcp: the MCP servers Claude Code runs: list, status, add, enable, restart, refresh-tokens, kill.
func mcp() { os.Exit(run(os.Args[1:])) }
