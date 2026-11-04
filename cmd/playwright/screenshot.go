package main

import "os"

// screenshot is playwright screenshot, the task playwright:screenshot: a page's screenshot.
func screenshot() { os.Exit(run(os.Args[1:])) }
