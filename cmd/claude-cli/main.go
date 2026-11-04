// Command claude-cli manages the Claude CLI on this machine, Anthropic's native binary (a machine-wide tool, not a
// project pin): which installations there are, old npm/bun ones cleaned up, the native one installed or upgraded to at
// least a version (mise run claude-cli:ensure), and the MCP servers Claude Code runs. It's claude-cli, not claude, so it
// never shadows Anthropic's claude on a task's PATH (bin/ comes first). Each command is a file of its own and a task of
// .mise/tasks/claude-cli.toml; its arguments are internal/claude's, passed through.
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/claude"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "status":
		status()
	case "check":
		check()
	case "detect":
		detect()
	case "cleanup":
		cleanup()
	case "version-check":
		versionCheck()
	case "ensure":
		ensure()
	case "mcp":
		mcp()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/claude's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return claude.NewCLI(version, os.Stdout, os.Stderr).Run(args)
}
