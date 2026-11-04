package main

import "os"

// check is claude-cli check, the task claude-cli:check: whether the native binary is installed.
func check() { os.Exit(run(os.Args[1:])) }
