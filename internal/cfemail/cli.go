package cfemail

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// Run is the CLI entry point. cmd/cfemail/main.go calls this.
// Returns a process exit code (0 = success).
func Run(args []string, version string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		printUsage(stdout)
		return 1
	}

	cmd := args[1]
	rest := args[2:]

	switch cmd {
	case "send":
		return runSend(rest, stdout, stderr)
	case "check":
		return runCheck(stdout, stderr)
	case "version", "-version", "--version":
		fmt.Fprintf(stdout, "cfemail %s\n", version)
		return 0
	case "help", "-h", "--help":
		printUsage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", cmd)
		printUsage(stderr)
		return 1
	}
}

// runSend handles `cfemail send`.
func runSend(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("send", flag.ContinueOnError)
	fs.SetOutput(stderr)

	var (
		to      = fs.String("to", "", "Recipient(s), comma-separated. Supports \"Name <email>\".")
		cc      = fs.String("cc", "", "CC recipient(s), comma-separated")
		bcc     = fs.String("bcc", "", "BCC recipient(s), comma-separated")
		from    = fs.String("from", DefaultFrom, "From address")
		replyTo = fs.String("reply-to", "", "Reply-To address")
		subject = fs.String("subject", "", "Subject line")
		text    = fs.String("text", "", "Plain-text body")
		html    = fs.String("html", "", "HTML body")
	)
	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *to == "" {
		fmt.Fprintln(stderr, "Error: --to is required")
		return 1
	}
	if *text == "" && *html == "" {
		fmt.Fprintln(stderr, "Error: provide --text and/or --html")
		return 1
	}

	client, err := NewClientFromEnv()
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	msg := Message{
		From:    ParseAddress(*from),
		To:      ParseAddressList(*to),
		CC:      ParseAddressList(*cc),
		BCC:     ParseAddressList(*bcc),
		ReplyTo: *replyTo,
		Subject: *subject,
		Text:    *text,
		HTML:    *html,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := client.Send(ctx, msg)
	if err != nil {
		fmt.Fprintf(stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "Sent from %s\n", msg.From.Email)
	if len(result.Delivered) > 0 {
		fmt.Fprintf(stdout, "  delivered: %v\n", result.Delivered)
	}
	if len(result.Queued) > 0 {
		fmt.Fprintf(stdout, "  queued:    %v\n", result.Queued)
	}
	if len(result.PermanentBounces) > 0 {
		fmt.Fprintf(stdout, "  bounced:   %v\n", result.PermanentBounces)
	}
	return 0
}

// runCheck verifies that the required environment is present.
func runCheck(stdout, stderr io.Writer) int {
	ok := true
	if os.Getenv("CLOUDFLARE_API_TOKEN") == "" {
		fmt.Fprintln(stderr, "✗ CLOUDFLARE_API_TOKEN not set")
		ok = false
	} else {
		fmt.Fprintln(stdout, "✓ CLOUDFLARE_API_TOKEN set")
	}
	if os.Getenv("CF_ACCOUNT_ID") == "" {
		fmt.Fprintln(stderr, "✗ CF_ACCOUNT_ID not set")
		ok = false
	} else {
		fmt.Fprintln(stdout, "✓ CF_ACCOUNT_ID set")
	}
	if !ok {
		fmt.Fprintf(stderr, "\nCreate a token with 'Email Sending: Edit' at %s\n", APITokenURL)
		return 1
	}
	fmt.Fprintf(stdout, "Default sender: %s\n", DefaultFrom)
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `cfemail - send transactional email via Cloudflare Email Service

Usage:
  cfemail send --to=EMAIL --subject="..." --text="..." [--html="..."]
  cfemail check                 Verify env (CLOUDFLARE_API_TOKEN, CF_ACCOUNT_ID)
  cfemail version

Flags for send:
  --to        Recipient(s), comma-separated (supports "Name <email>")
  --cc        CC recipient(s)
  --bcc       BCC recipient(s)
  --from      From address (default: `+DefaultFrom+`)
  --reply-to  Reply-To address
  --subject   Subject line
  --text      Plain-text body
  --html      HTML body

Environment:
  CLOUDFLARE_API_TOKEN   API token with 'Email Sending: Edit'
  CF_ACCOUNT_ID          Cloudflare account ID
`)
}
