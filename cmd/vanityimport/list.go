package main

import "os"

// list is vanityimport list, the task vanityimport:list: the packages with a vanity import page.
func list() { os.Exit(run(os.Args[1:])) }
