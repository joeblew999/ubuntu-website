package main

import "os"

// open is playwright open, the task playwright:open: a URL opened in a Playwright browser.
func open() { os.Exit(run(os.Args[1:])) }
