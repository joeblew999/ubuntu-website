package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/sitecheck"
)

// all is sitecheck all, the task sitecheck:all: every check, dns, tcp, redirect and http, one after the other. It fails
// if any does.
func all(args []string) {
	o := parse("all", args, nil)
	if !sitecheck.All(o.url, o.nodes, o.wait, os.Stdout, os.Stderr) {
		os.Exit(1)
	}
}
