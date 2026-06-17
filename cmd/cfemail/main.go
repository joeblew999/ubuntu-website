// cfemail - send transactional email via Cloudflare Email Service.
//
// A dependency-free replacement for the SMTP2GO/Gmail sending relay. It does
// not manage subscriber lists or campaigns — that remains MailerLite's job.
//
// Usage:
//
//	cfemail send --to=EMAIL --subject="..." --text="..." [--html="..."]
//	cfemail check
//	cfemail version
//
// Environment:
//
//	CLOUDFLARE_API_TOKEN   API token with 'Email Sending: Edit'
//	CF_ACCOUNT_ID          Cloudflare account ID
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/cfemail"
)

var version = "dev"

func main() {
	os.Exit(cfemail.Run(os.Args, version, os.Stdout, os.Stderr))
}
