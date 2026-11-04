package main

import "os"

// add is vanityimport add, the task vanityimport:add: a package's vanity import page added: add <package-name>.
func add() { os.Exit(run(os.Args[1:])) }
