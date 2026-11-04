package main

import "os"

// concat is youtube concat, the task youtube:concat: videos joined.
func concat() { os.Exit(run(os.Args[1:])) }
