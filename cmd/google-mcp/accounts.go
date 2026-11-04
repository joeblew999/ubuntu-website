package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// accounts is google-mcp accounts, the task google-mcp:accounts: the Google accounts signed in to the MCP server, one
// file each in ~/.google-mcp-accounts.
func accounts() {
	names := signedIn()
	if len(names) == 0 {
		fmt.Println("No accounts signed in: mise run google-mcp:auth")
		return
	}
	for _, n := range names {
		fmt.Println(n)
	}
}

// accountsDir is where the MCP server keeps each account's tokens, outside the repo: they open the account.
func accountsDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".google-mcp-accounts")
}

// signedIn are the accounts in accountsDir, by file name.
func signedIn() []string {
	files, _ := filepath.Glob(filepath.Join(accountsDir(), "*.json"))
	var names []string
	for _, f := range files {
		names = append(names, strings.TrimSuffix(filepath.Base(f), ".json"))
	}
	return names
}
