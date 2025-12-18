package google

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/google/gmail"
)

func (c *cliContext) handleGmail(args []string) {
	if len(args) < 1 {
		c.printGmailUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "list":
		c.gmailList(cmdArgs)
	case "send":
		c.gmailSend(cmdArgs)
	case "compose":
		c.gmailCompose(cmdArgs)
	case "check":
		c.gmailCheck(cmdArgs)
	case "server":
		c.gmailServer(cmdArgs)
	default:
		fmt.Fprintf(c.stderr, "Unknown gmail command: %s\n", cmd)
		c.printGmailUsage()
		os.Exit(1)
	}
}

func (c *cliContext) gmailList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	maxResults := fs.Int("max", 10, "Maximum messages to list")
	query := fs.String("query", "", "Search query (optional)")
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	config := gmail.DefaultConfig()
	client, err := gmail.NewAPISender(config)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create API client: %v", err))
	}

	result, err := client.List(*maxResults, *query)
	if err != nil {
		c.exitError(fmt.Sprintf("List failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
		return
	}

	if len(result.Messages) == 0 {
		fmt.Fprintln(c.stdout, "No messages found.")
		return
	}

	fmt.Fprintf(c.stdout, "Recent messages (%d):\n", len(result.Messages))
	for _, msg := range result.Messages {
		flagStr := ""
		if msg.Unread {
			flagStr = "[UNREAD] "
		}
		fmt.Fprintf(c.stdout, "\n  %s%s\n", flagStr, msg.Subject)
		if msg.From != "" {
			fmt.Fprintf(c.stdout, "    From: %s\n", msg.From)
		}
		if !msg.Date.IsZero() {
			fmt.Fprintf(c.stdout, "    Date: %s\n", msg.Date.Format("2006-01-02 15:04"))
		}
		if msg.Snippet != "" {
			fmt.Fprintf(c.stdout, "    Snippet: %s\n", msg.Snippet)
		}
	}
}

func (c *cliContext) gmailSend(args []string) {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	to := fs.String("to", "", "Recipient email")
	subject := fs.String("subject", "", "Email subject")
	body := fs.String("body", "", "Email body")
	bodyFile := fs.String("body-file", "", "Read body from file")
	mode := fs.String("mode", "api", "Send mode: api or browser")
	signature := fs.String("signature", "", "Override signature")
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	bodyText := *body
	if *bodyFile != "" {
		data, err := os.ReadFile(*bodyFile)
		if err != nil {
			c.exitError(fmt.Sprintf("Failed to read body file: %v", err))
		}
		bodyText = string(data)
	}

	if *to == "" {
		c.exitError("--to is required")
	}
	if *subject == "" {
		c.exitError("--subject is required")
	}
	if bodyText == "" {
		c.exitError("--body or --body-file is required")
	}

	config := gmail.DefaultConfig()
	if *signature != "" {
		config.Signature = *signature
	}

	email := &gmail.Email{
		To:      *to,
		Subject: *subject,
		Body:    bodyText,
	}

	var sender gmail.Sender
	var err error

	switch strings.ToLower(*mode) {
	case "api":
		sender, err = gmail.NewAPISender(config)
		if err != nil {
			c.exitError(fmt.Sprintf("Failed to create API sender: %v", err))
		}
	case "browser":
		sender = gmail.NewBrowserSender(config, false)
	default:
		c.exitError(fmt.Sprintf("Invalid mode: %s (use 'api' or 'browser')", *mode))
	}

	result, err := sender.Send(email)
	if err != nil {
		c.exitError(fmt.Sprintf("Send failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
	} else {
		fmt.Fprintln(c.stdout, "Email sent successfully!")
		fmt.Fprintf(c.stdout, "  To: %s\n", *to)
		fmt.Fprintf(c.stdout, "  Subject: %s\n", *subject)
		fmt.Fprintf(c.stdout, "  Mode: %s\n", result.Mode)
		if result.MessageID != "" {
			fmt.Fprintf(c.stdout, "  Message ID: %s\n", result.MessageID)
		}
	}
}

func (c *cliContext) gmailCompose(args []string) {
	fs := flag.NewFlagSet("compose", flag.ExitOnError)
	to := fs.String("to", "", "Recipient email")
	subject := fs.String("subject", "", "Email subject")
	body := fs.String("body", "", "Email body")
	bodyFile := fs.String("body-file", "", "Read body from file")
	signature := fs.String("signature", "", "Override signature")
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	bodyText := *body
	if *bodyFile != "" {
		data, err := os.ReadFile(*bodyFile)
		if err != nil {
			c.exitError(fmt.Sprintf("Failed to read body file: %v", err))
		}
		bodyText = string(data)
	}

	if *to == "" {
		c.exitError("--to is required")
	}
	if *subject == "" {
		c.exitError("--subject is required")
	}
	if bodyText == "" {
		c.exitError("--body or --body-file is required")
	}

	config := gmail.DefaultConfig()
	if *signature != "" {
		config.Signature = *signature
	}

	email := &gmail.Email{
		To:      *to,
		Subject: *subject,
		Body:    bodyText,
	}

	sender := gmail.NewBrowserSender(config, true)
	result, err := sender.Send(email)
	if err != nil {
		c.exitError(fmt.Sprintf("Compose failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
	} else {
		fmt.Fprintln(c.stdout, "Gmail compose opened!")
		fmt.Fprintf(c.stdout, "  To: %s\n", *to)
		fmt.Fprintf(c.stdout, "  Subject: %s\n", *subject)
		fmt.Fprintf(c.stdout, "  From: %s (change this before sending!)\n", config.FromAddress)
		fmt.Fprintln(c.stdout, "\nReview and click Send in the browser.")
	}
}

func (c *cliContext) gmailCheck(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	config := gmail.DefaultConfig()

	sender, err := gmail.NewAPISender(config)
	if err != nil {
		if *jsonOutput {
			c.outputJSON(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		} else {
			c.exitError(fmt.Sprintf("Failed to load token: %v", err))
		}
		os.Exit(1)
	}

	if err := sender.Check(); err != nil {
		if *jsonOutput {
			c.outputJSON(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		} else {
			c.exitError(fmt.Sprintf("API check failed: %v", err))
		}
		os.Exit(1)
	}

	if *jsonOutput {
		c.outputJSON(map[string]interface{}{
			"success": true,
			"from":    config.FromAddress,
		})
	} else {
		fmt.Fprintln(c.stdout, "Gmail API connection OK!")
		fmt.Fprintf(c.stdout, "  From address: %s\n", config.FromAddress)
		fmt.Fprintf(c.stdout, "  Token path: %s\n", config.TokenPath)
	}
}

func (c *cliContext) gmailServer(args []string) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	port := fs.Int("port", 8087, "HTTP port")
	fs.Parse(args)

	config := gmail.DefaultConfig()

	server, err := gmail.NewServer(config, *port)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create server: %v", err))
	}

	if err := server.Start(); err != nil {
		c.exitError(fmt.Sprintf("Server error: %v", err))
	}
}

func (c *cliContext) printGmailUsage() {
	fmt.Fprintln(c.stdout, `Usage: google gmail <command> [arguments]

Commands:
  list [--max=10] [--query=TEXT]              List recent messages
  send --to=EMAIL --subject=SUBJ --body=TEXT   Send an email
  compose --to=EMAIL --subject=SUBJ --body=TEXT  Open Gmail compose
  check                                         Verify Gmail API connection
  server [--port=8087]                          Start webhook server

Send/Compose Options:
  --to          Recipient email address (required)
  --subject     Email subject (required)
  --body        Email body text (required unless --body-file)
  --body-file   Read body from file
  --mode        Send mode: api (default) or browser
  --signature   Override default signature

Options:
  --json    Output as JSON

Examples:
  google gmail send --to=user@example.com --subject="Test" --body="Hello"
  google gmail compose --to=user@example.com --subject="Review" --body="Please check"
  google gmail check
  google gmail server --port=8087`)
}
