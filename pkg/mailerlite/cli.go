package mailerlite

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// CLI provides the command-line interface for MailerLite operations.
type CLI struct {
	client      *Client
	ctx         context.Context
	out         io.Writer
	githubIssue bool
	verbose     bool
}

// CLIConfig holds configuration for the CLI.
type CLIConfig struct {
	// APIKey is the MailerLite API key. If empty, uses MAILERLITE_API_KEY env var.
	APIKey string

	// GitHubIssue enables markdown output for GitHub issues.
	GitHubIssue bool

	// Verbose enables verbose output.
	Verbose bool

	// Output is the writer for output. Defaults to os.Stdout.
	Output io.Writer

	// Timeout is the context timeout. Defaults to 30 seconds.
	Timeout time.Duration
}

// NewCLI creates a new CLI instance.
func NewCLI(cfg *CLIConfig) (*CLI, context.CancelFunc, error) {
	if cfg == nil {
		cfg = &CLIConfig{}
	}

	// Get API key
	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("MAILERLITE_API_KEY")
	}

	// Create client (may be nil for commands that don't need it)
	var client *Client
	if apiKey != "" {
		client = NewClient(apiKey)
	}

	// Set defaults
	out := cfg.Output
	if out == nil {
		out = os.Stdout
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	cli := &CLI{
		client:      client,
		ctx:         ctx,
		out:         out,
		githubIssue: cfg.GitHubIssue,
		verbose:     cfg.Verbose,
	}

	return cli, cancel, nil
}

// Run executes a CLI command with the given arguments.
func (c *CLI) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	cmd := args[0]
	subArgs := args[1:]

	// Commands that don't require API key
	noAPIKeyCommands := map[string]bool{
		"releases": true,
		"open":     true,
		"docs":     true,
	}

	if c.client == nil && !noAPIKeyCommands[cmd] {
		return fmt.Errorf("MAILERLITE_API_KEY required. Get your key from: %s", APIKeyURL)
	}

	switch cmd {
	case "subscribers":
		return c.handleSubscribers(subArgs)
	case "groups":
		return c.handleGroups(subArgs)
	case "forms":
		return c.handleForms(subArgs)
	case "webhooks":
		return c.handleWebhooks(subArgs)
	case "automations":
		return c.handleAutomations(subArgs)
	case "stats":
		return c.handleStats()
	case "open":
		return c.handleOpen(subArgs)
	case "releases":
		return c.handleReleases(subArgs)
	case "server":
		return c.handleServer(subArgs)
	case "docs":
		return c.handleDocs()
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// handleDocs outputs package documentation for registry publishing.
func (c *CLI) handleDocs() error {
	PackageDoc().Write(c.out)
	return nil
}

// printf writes formatted output.
func (c *CLI) printf(format string, a ...interface{}) {
	fmt.Fprintf(c.out, format, a...)
}

// println writes a line.
func (c *CLI) println(a ...interface{}) {
	fmt.Fprintln(c.out, a...)
}
