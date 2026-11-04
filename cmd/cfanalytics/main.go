// Command cfanalytics reports the site's traffic from Cloudflare Web Analytics, against the last report
// (.analytics-state.json): mise run cfanalytics:report by hand, and cfanalytics:issue, the weekly workflow's, as markdown
// for a GitHub issue when traffic changed a lot. Each command is a file of its own and a task of
// .mise/tasks/cfanalytics.toml. Flags, anywhere: --days, -v, --webhook (Slack/Discord). UBUNTU_WEBSITE_CF_TOKEN from .env.
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/cfanalytics"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "report":
		report()
	case "issue":
		issue()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/cfanalytics's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return cfanalytics.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
