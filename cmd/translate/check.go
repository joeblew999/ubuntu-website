package main

import (
	"fmt"
	"os"
)

// check is translate check, the task translate:check, which the weekly workflow (.github/workflows/translate.yml)
// runs: every check (missing, stale and orphaned pages, the menus) as one markdown report for a GitHub issue. It
// fails if any finds something to do.
func check() {
	fmt.Print("# Translation Status Report\n\n")
	found := false
	for _, c := range [][]string{{"content", "missing"}, {"content", "stale"}, {"content", "orphans"}, {"menu", "check"}} {
		if run(append([]string{"--github-issue"}, c...)) != 0 {
			found = true
		}
	}
	if found {
		os.Exit(1)
	}
}
