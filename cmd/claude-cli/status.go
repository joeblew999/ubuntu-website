package main

import "os"

// status is claude-cli status, the task claude-cli:status: the Claude CLI installations found, and what needs doing (--dev, --ci).
func status() { os.Exit(run(os.Args[1:])) }
