package main

import "os"

// update is vanityimport update, the task vanityimport:update: a package's page brought up to date (all of them, without a name).
func update() { os.Exit(run(os.Args[1:])) }
