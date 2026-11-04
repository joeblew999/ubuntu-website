package main

import "os"

// fade is youtube fade, the task youtube:fade: a fade in or out.
func fade() { os.Exit(run(os.Args[1:])) }
