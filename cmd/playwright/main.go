// Command playwright is browser automation for flows with no API: OAuth sign-ins (the code from its callback),
// screenshots, and Gmail's and Cloudflare's settings screens. Each command is a file of its own and a task of
// .mise/tasks/playwright.toml; its arguments and flags are internal/playwright's, passed through
// (mise run playwright:screenshot https://... shot.png).
//
// Flags, anywhere: --headless, --timeout, --port, --profile (reuses logins), --chrome, --browser.
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/playwright"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "install":
		install()
	case "oauth":
		oauth()
	case "open":
		open()
	case "screenshot":
		screenshot()
	case "gmail-sendas":
		gmailSendas()
	case "gmail-settings":
		gmailSettings()
	case "gmail-logout":
		gmailLogout()
	case "gmail-switch":
		gmailSwitch()
	case "cloudflare":
		cloudflare()
	case "cloudflare-email":
		cloudflareEmail()
	case "cloudflare-email-add":
		cloudflareEmailAdd()
	case "cloudflare-token":
		cloudflareToken()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/playwright's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return playwright.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
