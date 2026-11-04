package main

import (
	"fmt"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// credentials is google-mcp credentials, the task google-mcp:credentials: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET
// into the keychain, prompting (hidden) only for whichever is missing. To replace one: mise run google-mcp:secret
// GOOGLE_CLIENT_SECRET; to rotate the secret: mise run google-mcp:rotate.
func credentials() {
	for _, key := range []string{"GOOGLE_CLIENT_ID", "GOOGLE_CLIENT_SECRET"} {
		if _, err := googleauth.Keychain(key); err == nil {
			fmt.Printf("✅ %s: already in the keychain\n", key)
			continue
		}
		secret([]string{key})
	}
}
