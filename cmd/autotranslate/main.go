// Command autotranslate provides automatic translation of Hugo markdown content
// using external translation APIs (DeepL or Claude).
//
// Usage:
//
//	autotranslate [flags] <command> [args]
//
// Commands:
//
//	file     Translate a single file
//	missing  Translate all missing files for a language
//	batch    Translate multiple files
//	status   Show translation quota/usage
//
// Examples:
//
//	# Translate a single file to Vietnamese using DeepL
//	autotranslate file content/english/blog/post.md vi
//
//	# Translate using Claude (no per-character cost if you have subscription)
//	autotranslate --provider=claude missing vi
//
//	# Dry-run to see what would be translated
//	autotranslate missing vi --dry-run
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/autotranslate"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	exitCode := autotranslate.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
