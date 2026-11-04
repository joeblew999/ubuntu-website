package main

import "os"

// drive is google drive, the task google:drive: Google Drive: list, get, download, upload, mkdir, rm.
func drive() { os.Exit(run(os.Args[1:])) }
