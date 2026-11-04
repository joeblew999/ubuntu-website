// Command upgrade upgrades what the repo is built with, each step repeatable (with nothing new upstream it changes
// nothing): the tools mise.toml pins and Go (tools), the Go modules (modules), and Hugo (hugo). mise run upgrade:all runs
// every step but Hugo's, which goes through a Cloudflare preview build first. Each command is a file of its own and a
// task of .mise/tasks/upgrade.toml (upgrade:<command>); the theme (theme:update) and the workflows' actions
// (upgrade:actions, pinact) are tasks of their own.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: upgrade tools | modules [--force] [--ignore <regexp>] | hugo [--to <version>]")
	}
	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "tools":
		tools()
	case "modules":
		modules(args)
	case "hugo":
		hugo(args)
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
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

// output is a tool's standard output, trimmed.
func output(name string, args ...string) (string, error) {
	b, err := exec.Command(name, args...).Output()
	return strings.TrimSpace(string(b)), err
}
