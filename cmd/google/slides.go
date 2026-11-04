package main

import "os"

// slides is google slides, the task google:slides: Google Slides: get, create, add, delete.
func slides() { os.Exit(run(os.Args[1:])) }
