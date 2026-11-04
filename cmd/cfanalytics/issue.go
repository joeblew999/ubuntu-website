package main

import "os"

// issue is cfanalytics issue, the task cfanalytics:issue: the weekly report as markdown for a GitHub issue; fails when traffic changed significantly.
func issue() { os.Exit(run(append([]string{"-github-issue"}, os.Args[2:]...))) }
