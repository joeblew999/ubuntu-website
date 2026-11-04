package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/sitecheck"
)

// links is sitecheck links [--dir <dir>], the task sitecheck:links, in root:ci after the production build: every link
// of the built site (public/) lands on a page or file the build made. No network: it checks the build, not the live
// site. It fails when one doesn't, naming it and the pages that link to it.
func links(args []string) {
	fs := flag.NewFlagSet("links", flag.ExitOnError)
	dir := fs.String("dir", "public", "the built site")
	fs.Parse(args)
	missing, err := sitecheck.Links(*dir, os.Stdout)
	if err != nil {
		fail(err.Error())
	}
	if missing > 0 {
		fail(fmt.Sprintf("%d links land nowhere", missing))
	}
}
