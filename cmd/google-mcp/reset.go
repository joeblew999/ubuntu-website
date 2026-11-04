package main

import (
	"flag"
	"fmt"
	"os"
)

// reset is google-mcp reset [--yes], the task google-mcp:reset: every Google account signed out (~/.google-mcp-accounts
// removed); without --yes it only says what it would remove. The keychain's client and the checked-in MCP configs stay.
func reset(args []string) {
	fs := flag.NewFlagSet("reset", flag.ExitOnError)
	yes := fs.Bool("yes", false, "really remove the accounts")
	fs.Parse(args)

	fmt.Println("This removes " + accountsDir() + "/ (the signed-in accounts' tokens)")
	if !*yes {
		fmt.Println("Nothing removed. To remove it: mise run google-mcp:reset --yes")
		return
	}
	if err := os.RemoveAll(accountsDir()); err != nil {
		fail(err.Error())
	}
	fmt.Println("✅ Removed. Sign in again: mise run google-mcp:auth")
}
