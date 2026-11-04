package main

import "os"

// detect is claude-cli detect, the task claude-cli:detect: every installation found (--json).
func detect() { os.Exit(run(os.Args[1:])) }
