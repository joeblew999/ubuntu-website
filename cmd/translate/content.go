package main

import "os"

// content is translate content <subcommand>, the task translate:content: the English pages against their translations
// (status, diff, changed, next, done, missing, orphans, stale, clean).
func content() { os.Exit(run(os.Args[1:])) }
