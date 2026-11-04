package main

import "os"

// install is playwright install, the task playwright:install: the browsers Playwright drives, installed.
func install() { os.Exit(run(os.Args[1:])) }
