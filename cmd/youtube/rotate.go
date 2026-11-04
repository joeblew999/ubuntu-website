package main

import "os"

// rotate is youtube rotate, the task youtube:rotate: a video rotated.
func rotate() { os.Exit(run(os.Args[1:])) }
