package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

const site = "https://www.ubuntusoftware.net"

// console is google-mcp console, the task google-mcp:console: the Google Cloud OAuth screens (Branding, Audience,
// client) opened, and the exact values to enter on them printed. Google has no API for these screens.
func console() {
	for _, path := range []string{"/", "/privacy-policy/", "/terms-of-service/"} {
		code := 0
		if resp, err := http.Get(site + path); err == nil {
			code = resp.StatusCode
			resp.Body.Close()
		}
		if code != http.StatusOK {
			fail(fmt.Sprintf("%s%s answers %d: fix the page before using it on the Branding screen", site, path, code))
		}
	}
	client, err := googleauth.Keychain("GOOGLE_CLIENT_ID")
	if err != nil {
		fail(err.Error())
	}
	p := os.Getenv("UBUNTU_WEBSITE_GCP_PROJECT")
	fmt.Printf(`1. Branding  https://console.cloud.google.com/auth/branding?project=%[1]s
     App name                  anything WITHOUT the word "Google" (Google rejects it), e.g. ged-mcp
     User support email        your Gmail (pick it from the dropdown)
     App logo                  leave EMPTY (a logo forces brand verification)
     Application home page     %[2]s/
     Privacy policy link       %[2]s/privacy-policy/
     Terms of service link     %[2]s/terms-of-service/
     Authorized domains        ubuntusoftware.net
     Developer contact         your Gmail
     → Save
2. Audience  https://console.cloud.google.com/auth/audience?project=%[1]s
     Publishing status must be "In production" → Publish app → Confirm
     ("Testing" refresh tokens expire after 7 days). After publishing: mise run google-mcp:auth
3. Client    https://console.cloud.google.com/auth/clients/%[3]s?project=%[1]s
     Type: Desktop app. Secrets: mise run google-mcp:rotate
`, p, site, client)
	browse("https://console.cloud.google.com/auth/branding?project=" + p)
	browse("https://console.cloud.google.com/auth/audience?project=" + p)
}
