package mailerlite

import (
	"flag"
	"fmt"
	"io"

	"github.com/joeblew999/ubuntu-website/internal/cli"
)

// Main is the mailerlite command line on args (args[0], the program): its flags (-github-issue, -v, -version), which
// may come anywhere (cli.Interspersed), then its command. It returns the exit code.
func Main(args []string, version string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mailerlite", flag.ContinueOnError)
	fs.SetOutput(stderr)
	githubIssue := fs.Bool("github-issue", false, "Output markdown for GitHub issue")
	verbose := fs.Bool("v", false, "Verbose output")
	showVersion := fs.Bool("version", false, "Show version")
	if err := fs.Parse(cli.Ordered(fs, args[1:])); err != nil {
		return 1
	}
	if *showVersion {
		fmt.Fprintf(stdout, "mailerlite %s\n", version)
		return 0
	}
	if fs.NArg() == 0 {
		PrintUsage()
		return 1
	}
	c, cancel, err := NewCLI(&CLIConfig{GitHubIssue: *githubIssue, Verbose: *verbose})
	if err != nil {
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return 1
	}
	defer cancel()
	if err := c.Run(fs.Args()); err != nil {
		fmt.Fprintf(stderr, "ERROR: %v\n", err)
		return 1
	}
	return 0
}
