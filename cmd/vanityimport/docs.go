package main

import "os"

// docs is vanityimport docs, the task vanityimport:docs: the packages' docs pages.
func docs() { os.Exit(run(os.Args[1:])) }
