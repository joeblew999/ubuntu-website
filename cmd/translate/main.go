// translate manages the translations of the site's Hugo content: what changed in English since the last translation,
// what's missing, stale or orphaned in the other languages, the menus, and the languages themselves. Each command is
// a file of its own and a task of .mise/tasks/translate.toml; its subcommands and flags are internal/translate's, and
// mise passes them through (mise run translate:content status, mise run translate:content done).
//
//	translate content status | diff <file> | changed | next | done | missing | orphans | stale | clean
//	translate menu check | sync
//	translate lang list | validate | add <code> <name> <dirname> | remove <code> | init <code>
//	translate check      every check, as markdown for a GitHub issue (the weekly workflow)
//
// Flags, anywhere: --github-issue (markdown, exit 1 if action is needed), --force (no prompts), --version.
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/translate"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "content":
		content()
	case "menu":
		menu()
	case "lang":
		lang()
	case "check":
		check()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/translate's CLI on args.
func run(args []string) int {
	return translate.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr, os.Stdin)
}
