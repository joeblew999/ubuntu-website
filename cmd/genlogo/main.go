// Command genlogo generates the site's branding from Go, the source of truth: the logo SVGs (light and dark), the
// favicon, the Open Graph image, the email logo, and the Bluesky avatar and banner. Never edit those files by hand.
// Each asset is a command of its own (a file each) and a task of .mise/tasks/genlogo.toml; mise run genlogo:all makes
// every one. Flags, anywhere: --dir (the project root, default .).
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/genlogo"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "all":
		all()
	case "logo-svg":
		logoSvg()
	case "favicon":
		favicon()
	case "og":
		og()
	case "email":
		email()
	case "avatar":
		avatar()
	case "banner":
		banner()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/genlogo's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return genlogo.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
