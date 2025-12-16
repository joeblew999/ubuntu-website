// mailerlite - CLI tool for MailerLite subscriber management
//
// Usage:
//
//	mailerlite subscribers list         List all subscribers
//	mailerlite subscribers count        Count total subscribers
//	mailerlite subscribers get EMAIL    Get subscriber by email
//	mailerlite subscribers add EMAIL    Add/update subscriber
//	mailerlite subscribers delete EMAIL Delete subscriber
//	mailerlite groups list              List all groups
//	mailerlite groups create NAME       Create a group
//	mailerlite groups subscribers ID    List subscribers in a group
//	mailerlite groups assign ID EMAIL   Assign subscriber to group
//	mailerlite groups unassign ID EMAIL Remove subscriber from group
//	mailerlite forms list               List all forms
//	mailerlite forms subscribers ID     List subscribers for a form
//	mailerlite webhooks list            List all webhooks
//	mailerlite webhooks create URL      Create a webhook
//	mailerlite webhooks delete ID       Delete a webhook
//	mailerlite automations list         List all automations
//	mailerlite automations get ID       Get automation details
//	mailerlite stats                    Show account statistics
//	mailerlite open [TARGET]            Open dashboard in browser
//	mailerlite releases latest          Show latest GitHub release
//	mailerlite releases list            List all GitHub releases
//	mailerlite releases urls            Show release URLs for emails
//	mailerlite server [PORT=8086]       Start webhook server for Web3Forms
//
// Flags:
//
//	-github-issue    Output markdown for GitHub issue
//	-v               Verbose output
//	-version         Show version
//
// Environment:
//
//	MAILERLITE_API_KEY    API key (required)
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/pkg/mailerlite"
)

var version = "dev"

func main() {
	var (
		githubIssue = flag.Bool("github-issue", false, "Output markdown for GitHub issue")
		verbose     = flag.Bool("v", false, "Verbose output")
		showVersion = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("mailerlite %s\n", version)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		mailerlite.PrintUsage()
		os.Exit(1)
	}

	cli, cancel, err := mailerlite.NewCLI(&mailerlite.CLIConfig{
		GitHubIssue: *githubIssue,
		Verbose:     *verbose,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	defer cancel()

	if err := cli.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
