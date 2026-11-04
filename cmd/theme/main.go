// Command theme keeps the Hugo Plate theme, a Hugo Module vendored in _vendor/ (its exact commits in
// _vendor/modules.txt): it updates it from upstream, shows what's pinned, and makes assets/css/generated-theme.css from
// data/theme.json with upstream's own generator, themeGenerator.js, kept here unmodified. Each command is a file of
// its own and a task of .mise/tasks/theme.toml (theme:<command>).
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const (
	module     = "github.com/zeon-studio/hugoplate"
	modulesTxt = "_vendor/modules.txt"
	themeJSON  = "data/theme.json"
	themeCSS   = "assets/css/generated-theme.css"
	// Upstream's scripts/themeGenerator.js, as theme update syncs it from the pinned commit.
	generator = "cmd/theme/themeGenerator.js"
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: theme css | check | status | update [--ref <commit>]")
	}
	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "css":
		css()
	case "check":
		check()
	case "status":
		status()
	case "update":
		update(args)
	default:
		fail("no command " + strconv.Quote(cmd))
	}
}

// fail says what went wrong and exits 1.
func fail(msg string) {
	fmt.Println("✗ " + msg)
	os.Exit(1)
}

// run runs a tool with its output to ours.
func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}

// output is a tool's standard output.
func output(name string, args ...string) (string, error) {
	b, err := exec.Command(name, args...).Output()
	return string(b), err
}
