// translate provides translation workflow management for Hugo multilingual content.
//
// Usage:
//
//	translate status              Show what English files changed since last translation
//	translate diff <file>         Show git diff for specific file since checkpoint
//	translate missing             Show files missing in target languages
//	translate stale               Show potentially outdated translations
//	translate orphans             Show target files with no English source
//	translate clean               Delete orphaned translation files
//	translate done                Update checkpoint tag to current commit
//	translate next                Show next file to translate with progress
//	translate changed             Show detailed changes for all files
//	translate validate            Check translator config matches Hugo config
//
// Flags:
//
//	-github-issue    Output markdown for GitHub Issue (exit 1 if action needed)
//	-version         Print version and exit
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/translator"
)

const version = "0.1.0"

func main() {
	// Global flags
	githubIssue := flag.Bool("github-issue", false, "Output markdown for GitHub Issue")
	ver := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *ver {
		fmt.Printf("translate v%s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		usage()
		os.Exit(1)
	}

	// Create translator instance (no API key needed for status commands)
	t, err := translator.NewChecker()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	cmd := flag.Arg(0)
	var exitCode int

	switch cmd {
	case "status":
		exitCode = t.Status(*githubIssue)
	case "diff":
		file := flag.Arg(1)
		if file == "" {
			fmt.Fprintln(os.Stderr, "Error: diff requires a file argument")
			fmt.Fprintln(os.Stderr, "Usage: translate diff <file>")
			os.Exit(1)
		}
		exitCode = t.Diff(file)
	case "missing":
		exitCode = t.Missing(*githubIssue)
	case "stale":
		exitCode = t.Stale(*githubIssue)
	case "orphans":
		exitCode = t.Orphans(*githubIssue)
	case "clean":
		exitCode = t.Clean()
	case "done":
		exitCode = t.Done()
	case "next":
		exitCode = t.Next()
	case "changed":
		exitCode = t.Changed()
	case "validate":
		exitCode = t.Validate()
	case "langs":
		exitCode = t.Langs()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func usage() {
	fmt.Fprintf(os.Stderr, `translate - Translation workflow for Hugo multilingual content

Usage:
  translate <command> [flags]

Commands:
  status      Show what English files changed since last translation
  diff <file> Show git diff for specific file since checkpoint
  missing     Show files missing in target languages
  stale       Show potentially outdated translations (target < 50%% of source)
  orphans     Show target files with no English source (should be deleted)
  clean       Delete orphaned translation files
  done        Update checkpoint tag to current commit
  next        Show next file to translate with progress
  changed     Show detailed changes for all files
  validate    Check translator config matches Hugo config
  langs       Show configured languages and detect stray directories

Flags:
  -github-issue  Output markdown for GitHub Issue (exit 1 if action needed)
  -version       Print version and exit

Examples:
  translate status                    # See what changed
  translate diff blog/my-post.md      # See specific file changes
  translate missing -github-issue     # CI mode: exit 1 if missing files
  translate done                      # Mark translations complete

`)
}
