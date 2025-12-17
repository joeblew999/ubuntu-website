// translate provides translation workflow management for Hugo multilingual content.
//
// Usage:
//
//	translate content status          Show what English files changed since last translation
//	translate content diff <file>     Show git diff for specific file since checkpoint
//	translate content changed         Show detailed changes for all files
//	translate content next            Show next file to translate with progress
//	translate content done            Mark translations complete (update checkpoint)
//	translate content missing         Show files missing in target languages
//	translate content orphans         Show target files with no English source
//	translate content stale           Show potentially outdated translations
//	translate content clean           Delete orphaned files (prompts unless -force)
//
//	translate menu check              Validate menu files for broken links and sync issues
//	translate menu sync               Generate translated menu files from English
//
//	translate lang list               Show configured languages and detect stray directories
//	translate lang validate           Check translator config matches Hugo config
//	translate lang add <code> <name> <dirname>   Add a new target language
//	translate lang remove <code>      Remove a language (prompts unless -force)
//	translate lang init <code>        Initialize content directory for configured language
//
// Flags:
//
//	-github-issue    Output markdown for GitHub Issue (exit 1 if action needed)
//	-force           Skip confirmation prompts (for CI)
//	-version         Print version and exit
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/translator"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	exitCode := translator.Run(os.Args, version, os.Stdout, os.Stderr, os.Stdin)
	os.Exit(exitCode)
}
