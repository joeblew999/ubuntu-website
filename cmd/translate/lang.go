package main

import "os"

// lang is translate lang <subcommand>, the task translate:lang: the configured languages (list, validate, add,
// remove, init).
func lang() { os.Exit(run(os.Args[1:])) }
