package main

import "os"

// audio is youtube audio, the task youtube:audio: a video's sound extracted.
func audio() { os.Exit(run(os.Args[1:])) }
