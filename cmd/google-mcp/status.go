package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// status is google-mcp status, the task google-mcp:status: each part of the setup, there or not: the server binary
// (mise pins it), the OAuth client in the keychain, the accounts signed in, and .mcp.json's entry for Claude Code. The
// proof that it all works is google-mcp:check.
func status() {
	if p, err := exec.LookPath("google-mcp-server"); err == nil {
		fmt.Println("✅ google-mcp-server: " + p)
	} else {
		fmt.Println("❌ google-mcp-server: not installed (mise install)")
	}
	for _, key := range []string{"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET"} {
		if _, err := secretValue(key); err == nil {
			fmt.Printf("✅ %s: in the keychain\n", key)
		} else {
			fmt.Printf("❌ %s: missing (mise run google-mcp:credentials)\n", key)
		}
	}
	if names := signedIn(); len(names) > 0 {
		fmt.Printf("✅ Accounts (%d):\n", len(names))
		for _, n := range names {
			fmt.Println("   - " + n)
		}
	} else {
		fmt.Println("❌ Accounts: none signed in (mise run google-mcp:auth)")
	}
	if b, err := os.ReadFile(".mcp.json"); err == nil && strings.Contains(string(b), "google-mcp-server") {
		fmt.Println("✅ .mcp.json: Claude Code runs google-mcp-server (through fnox)")
	} else {
		fmt.Println("❌ .mcp.json: no google-mcp-server entry (it's checked in: restore it from git)")
	}
	if b, err := os.ReadFile(".vscode/mcp.json"); err == nil && !strings.Contains(string(b), `"fnox"`) {
		fmt.Println("⚠️  .vscode/mcp.json: doesn't go through fnox, so it has no client secret")
	}
}
