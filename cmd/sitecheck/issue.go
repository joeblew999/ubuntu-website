package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/sitecheck"
)

// issue is sitecheck issue [--report <file>], the task sitecheck:issue, which the scheduled workflow runs: the HTTP
// check compared with the last run's (.sitecheck-state.json), as markdown for a GitHub issue, printed and, with
// --report, written to that file. It fails when it finds an issue, so the workflow opens one.
func issue(args []string) {
	var report string
	o := parse("issue", args, func(fs *flag.FlagSet) {
		fs.StringVar(&report, "report", "", "also write the markdown report to this file")
	})
	md, found, err := sitecheck.Report(o.url, o.nodes, o.wait, os.Stderr)
	if err != nil {
		fail(err.Error())
	}
	fmt.Println(md)
	if report != "" {
		if err := os.WriteFile(report, []byte(md+"\n"), 0o644); err != nil {
			fail(err.Error())
		}
	}
	if found {
		os.Exit(1)
	}
}
