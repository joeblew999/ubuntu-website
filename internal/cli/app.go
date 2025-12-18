// Package cli provides a shared CLI framework for Ubuntu Software tools.
//
// Usage:
//
//	app := cli.New("myapp", "v1.0.0")
//	app.Run(os.Args[1:], func(c *cli.Context) error {
//	    c.Println("Hello, world!")
//	    return nil
//	})
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

// App represents a CLI application with standard flags and output handling.
type App struct {
	Name    string
	Version string

	// Configuration
	config Config

	// Internal state
	flagSet *flag.FlagSet
	ctx     context.Context
	cancel  context.CancelFunc
}

// Config holds CLI configuration options.
type Config struct {
	// GitHubIssue enables markdown output for GitHub issues.
	GitHubIssue bool

	// Verbose enables verbose output.
	Verbose bool

	// Output writer (defaults to os.Stdout).
	Output io.Writer

	// ErrOutput writer (defaults to os.Stderr).
	ErrOutput io.Writer

	// Timeout for operations (defaults to 30s).
	Timeout time.Duration
}

// Context provides runtime context to command handlers.
type Context struct {
	*App
	Args []string
}

// New creates a new CLI application.
func New(name, version string) *App {
	return &App{
		Name:    name,
		Version: version,
		config: Config{
			Output:    os.Stdout,
			ErrOutput: os.Stderr,
			Timeout:   30 * time.Second,
		},
	}
}

// ParseFlags parses standard flags from args.
// Returns remaining args after flags.
func (a *App) ParseFlags(args []string) ([]string, error) {
	a.flagSet = flag.NewFlagSet(a.Name, flag.ContinueOnError)
	a.flagSet.SetOutput(io.Discard) // We handle errors ourselves

	var showVersion bool
	a.flagSet.BoolVar(&a.config.GitHubIssue, "github-issue", false, "Output markdown for GitHub issue")
	a.flagSet.BoolVar(&a.config.Verbose, "v", false, "Verbose output")
	a.flagSet.BoolVar(&showVersion, "version", false, "Show version")

	if err := a.flagSet.Parse(args); err != nil {
		return nil, err
	}

	if showVersion {
		fmt.Fprintf(a.config.Output, "%s %s\n", a.Name, a.Version)
		os.Exit(0)
	}

	return a.flagSet.Args(), nil
}

// Run executes the application with the given handler.
func (a *App) Run(args []string, handler func(*Context) error) error {
	remaining, err := a.ParseFlags(args)
	if err != nil {
		return err
	}

	a.ctx, a.cancel = context.WithTimeout(context.Background(), a.config.Timeout)
	defer a.cancel()

	return handler(&Context{
		App:  a,
		Args: remaining,
	})
}

// Context returns the application context.
func (a *App) Context() context.Context {
	if a.ctx == nil {
		a.ctx, a.cancel = context.WithTimeout(context.Background(), a.config.Timeout)
	}
	return a.ctx
}

// Cancel cancels the context.
func (a *App) Cancel() {
	if a.cancel != nil {
		a.cancel()
	}
}

// Config returns the current configuration.
func (a *App) Config() *Config {
	return &a.config
}

// GitHubIssue returns true if markdown output is enabled.
func (c *Context) GitHubIssue() bool {
	return c.config.GitHubIssue
}

// Verbose returns true if verbose output is enabled.
func (c *Context) Verbose() bool {
	return c.config.Verbose
}

// Output returns the output writer.
func (c *Context) Output() io.Writer {
	return c.config.Output
}

// SetOutput sets the output writer.
func (a *App) SetOutput(w io.Writer) {
	a.config.Output = w
}

// SetTimeout sets the context timeout.
func (a *App) SetTimeout(d time.Duration) {
	a.config.Timeout = d
}
