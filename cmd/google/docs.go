package main

import "os"

// docs is google docs, the task google:docs: Google Docs: get, create, append, insert, replace.
func docs() { os.Exit(run(os.Args[1:])) }
