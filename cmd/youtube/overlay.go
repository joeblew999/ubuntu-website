package main

import "os"

// overlay is youtube overlay, the task youtube:overlay: an image laid over a video (-image, -pos).
func overlay() { os.Exit(run(os.Args[1:])) }
