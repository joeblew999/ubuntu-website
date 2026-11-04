package main

import "os"

// stats is mailerlite stats, the task mailerlite:stats: the account's statistics.
func stats() { os.Exit(run(os.Args[1:])) }
