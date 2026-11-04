package main

import "os"

// docs is mailerlite docs, the task mailerlite:docs: the CLI's docs.
func docs() { os.Exit(run(os.Args[1:])) }
