package main

import "os"

// list is youtube list, the task youtube:list: the videos downloaded.
func list() { os.Exit(run(os.Args[1:])) }
