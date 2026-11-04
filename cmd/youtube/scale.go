package main

import "os"

// scale is youtube scale, the task youtube:scale: a video scaled.
func scale() { os.Exit(run(os.Args[1:])) }
