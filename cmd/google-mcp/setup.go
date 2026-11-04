package main

import "fmt"

// setup is google-mcp setup, the task google-mcp:setup: the whole setup on a new machine, once the Google Cloud side
// exists (mise run google-mcp:console): the OAuth client into the keychain (only what's missing), then a Google
// account signed in, which ends with google-mcp:check proving it all.
func setup() {
	credentials()
	auth()
	fmt.Println("\n✅ Set up. Restart Claude Code to pick up the server (.mcp.json is checked in).")
}
