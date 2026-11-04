package main

import "os"

// ensure is claude-cli ensure, the task claude-cli:ensure: the native CLI installed with Anthropic's installer for this OS, or upgraded, to at least a version (default 2.0.0).
func ensure() { os.Exit(run(os.Args[1:])) }
