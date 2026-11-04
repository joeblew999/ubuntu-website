// Command vanityimport manages the site's Go vanity import pages (/pkg/<name>: www.ubuntusoftware.net/pkg/<name> as a
// Go module path). Each command is a file of its own and a task of .mise/tasks/vanityimport.toml; its arguments and
// flags are internal/vanityimport's, passed through (mise run vanityimport:add <package-name>).
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/vanityimport"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		os.Exit(run(nil))
	}
	switch os.Args[1] {
	case "add":
		add()
	case "list":
		list()
	case "update":
		update()
	case "info":
		info()
	case "docs":
		docs()
	default:
		os.Exit(run(os.Args[1:]))
	}
}

// run is internal/vanityimport's CLI on args (its flags anywhere, internal/cli.Interspersed).
func run(args []string) int {
	return vanityimport.Run(append([]string{os.Args[0]}, args...), version, os.Stdout, os.Stderr)
}
