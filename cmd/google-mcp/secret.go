package main

import (
	"fmt"
	"os"
	"os/exec"
)

// secret is google-mcp secret <key>, the task google-mcp:secret: GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET into the
// keychain, from a hidden prompt or stdin.
func secret(args []string) {
	if len(args) != 1 || (args[0] != "GOOGLE_CLIENT_ID" && args[0] != "GOOGLE_CLIENT_SECRET") {
		fail("the key is GOOGLE_CLIENT_ID or GOOGLE_CLIENT_SECRET")
	}
	key := args[0]
	cmd := exec.Command("fnox", "set", "-p", "keychain", "-k", "UBUNTU_WEBSITE_"+key, key)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		fail("fnox set " + key + ": " + err.Error())
	}
	fmt.Printf("✓ %s stored in the keychain\n", key)
}
