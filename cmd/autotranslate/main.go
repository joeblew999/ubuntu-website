// Command autotranslate machine-translates the site's Hugo markdown (DeepL, Claude's API, or the Claude CLI), keeping
// front matter, shortcodes, code and links as they are. Each command is a file of its own and a task of
// .mise/tasks/autotranslate.toml; its arguments and flags are internal/autotranslate's, passed through:
//
//	mise run autotranslate:missing vi                        every page missing in Vietnamese (DeepL)
//	mise run autotranslate:missing vi --provider claude-cli  the same, with the Claude CLI
//	mise run autotranslate:file content/english/blog/post.md de --dry-run
//
// Flags, anywhere: --provider (deepl, claude, claude-cli), --dry-run, --verbose, --api-key, --bundle.
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/autotranslate"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "file":
		file()
	case "missing":
		missing()
	case "arb":
		arb()
	case "arb-status":
		arbStatus()
	case "languages":
		languages()
	case "status":
		status()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/autotranslate's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return autotranslate.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
