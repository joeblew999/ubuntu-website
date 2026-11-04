package main

import "os"

// gif is youtube gif, the task youtube:gif: a video made into a GIF.
func gif() { os.Exit(run(os.Args[1:])) }
