package main

import "os"

// compress is youtube compress, the task youtube:compress: a video compressed.
func compress() { os.Exit(run(os.Args[1:])) }
