package main

import "os"

// sheets is google sheets, the task google:sheets: Google Sheets: info, get, set, append, clear.
func sheets() { os.Exit(run(os.Args[1:])) }
