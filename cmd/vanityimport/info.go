package main

import "os"

// info is vanityimport info, the task vanityimport:info: a package's vanity import details: info <package-name>.
func info() { os.Exit(run(os.Args[1:])) }
