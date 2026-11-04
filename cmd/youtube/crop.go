package main

import "os"

// crop is youtube crop, the task youtube:crop: a video cropped.
func crop() { os.Exit(run(os.Args[1:])) }
