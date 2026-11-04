package main

import "os"

// info is youtube info, the task youtube:info: a video's details, before downloading it.
func info() { os.Exit(run(os.Args[1:])) }
