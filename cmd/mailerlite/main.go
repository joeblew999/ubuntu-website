// Command mailerlite manages the MailerLite account: subscribers, groups, forms, webhooks, automations, stats, the
// GitHub releases its emails link to, and the webhook server that takes Web3Forms submissions in as subscribers. Each
// command is a file of its own and a task of .mise/tasks/mailerlite.toml; its subcommands and flags are
// internal/mailerlite's, passed through:
//
//	mise run mailerlite:subscribers list
//	mise run mailerlite:subscribers add x@y.com "Name"
//	mise run mailerlite:server PORT=8086 GROUP_ID=12345
//
// Flags, anywhere: -github-issue, -v, -version. MAILERLITE_API_KEY from .env.
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/mailerlite"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "subscribers":
		subscribers()
	case "groups":
		groups()
	case "forms":
		forms()
	case "webhooks":
		webhooks()
	case "automations":
		automations()
	case "stats":
		stats()
	case "open":
		open()
	case "releases":
		releases()
	case "server":
		server()
	case "docs":
		docs()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/mailerlite's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return mailerlite.Main(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
