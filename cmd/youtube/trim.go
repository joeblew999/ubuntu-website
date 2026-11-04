package main

import "os"

// trim is youtube trim, the task youtube:trim: a video trimmed (-start, -end).
func trim() { os.Exit(run(os.Args[1:])) }
