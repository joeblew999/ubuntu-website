package main

import "os"

// mute is youtube mute, the task youtube:mute: a video's sound removed.
func mute() { os.Exit(run(os.Args[1:])) }
