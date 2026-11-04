package main

import "os"

// cleanup is claude-cli cleanup, the task claude-cli:cleanup: the old npm/bun installations removed.
func cleanup() { os.Exit(run(os.Args[1:])) }
