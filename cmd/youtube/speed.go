package main

import "os"

// speed is youtube speed, the task youtube:speed: a video sped up or slowed down.
func speed() { os.Exit(run(os.Args[1:])) }
