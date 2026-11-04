// Command sitecheck checks that the site is reachable from around the world, from check-host.net's nodes (no account,
// no secret): over HTTP, DNS and TCP, the apex domain's redirect to www, all of them at once (mise run sitecheck:all),
// and the report the scheduled workflow (.github/workflows/monitor-sitecheck.yml, mise run sitecheck:issue) opens a
// GitHub issue with. SITE_URL checks another site. And one with no network, links: every link of the built site
// (public/) lands, in root:ci. Each command is a file of its own and a task of
// .mise/tasks/sitecheck.toml (sitecheck:<command>); the checks are internal/sitecheck.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joeblew999/ubuntu-website/internal/sitecheck"
)

func main() {
	if len(os.Args) < 2 {
		fail("usage: sitecheck http | dns | tcp | redirect | all | issue [--url <url>] [--nodes <n>] [--wait <seconds>] | links [--dir <dir>]")
	}
	switch cmd, args := os.Args[1], os.Args[2:]; cmd {
	case "http":
		http(args)
	case "dns":
		dns(args)
	case "tcp":
		tcp(args)
	case "redirect":
		redirect(args)
	case "all":
		all(args)
	case "issue":
		issue(args)
	case "links":
		links(args)
	default:
		fail("no command " + strconv.Quote(cmd))
	}
}

// fail says what went wrong and exits 1.
func fail(msg string) {
	fmt.Fprintln(os.Stderr, "✗ "+msg)
	os.Exit(1)
}

// options are the flags every command takes: what to check, from how many nodes, and how long to wait for them.
type options struct {
	url         string
	nodes, wait int
}

// parse reads a command's flags: the ones every command takes, and its own, which own adds to fs.
func parse(name string, args []string, own func(fs *flag.FlagSet)) options {
	var o options
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	fs.StringVar(&o.url, "url", sitecheck.URL(), "URL to check (for HTTP) or domain (for DNS/TCP)")
	fs.IntVar(&o.nodes, "nodes", sitecheck.DefaultNodes, "most global nodes to check from")
	fs.IntVar(&o.wait, "wait", sitecheck.DefaultWait, "seconds to wait for their results")
	if own != nil {
		own(fs)
	}
	fs.Parse(args)
	return o
}

// check runs one check with a command's flags, and exits 1 if it fails.
func check(checkType string, args []string) {
	o := parse(checkType, args, nil)
	if !sitecheck.Check(checkType, o.url, o.nodes, o.wait, os.Stdout, os.Stderr) {
		os.Exit(1)
	}
}
