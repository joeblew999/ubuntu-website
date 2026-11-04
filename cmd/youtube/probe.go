package main

import "os"

// probe is youtube probe, the task youtube:probe: a video's streams, size and length.
func probe() { os.Exit(run(os.Args[1:])) }
