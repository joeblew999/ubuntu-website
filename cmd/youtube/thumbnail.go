package main

import "os"

// thumbnail is youtube thumbnail, the task youtube:thumbnail: a poster frame as JPG (-time).
func thumbnail() { os.Exit(run(os.Args[1:])) }
