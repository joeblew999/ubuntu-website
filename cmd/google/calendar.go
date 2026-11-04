package main

import "os"

// calendar is google calendar, the task google:calendar: Google Calendar: today, list, create, compose, open, server, check.
func calendar() { os.Exit(run(os.Args[1:])) }
