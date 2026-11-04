package main

import "os"

// releases is mailerlite releases, the task mailerlite:releases: the GitHub releases emails link to: latest, list, urls.
func releases() { os.Exit(run(os.Args[1:])) }
