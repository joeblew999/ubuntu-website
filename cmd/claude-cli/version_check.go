package main

import "os"

// versionCheck is claude-cli version-check, the task claude-cli:version-check: whether the installed version is at least the one given.
func versionCheck() { os.Exit(run(os.Args[1:])) }
