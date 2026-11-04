package main

import "os"

// gmailSettings is playwright gmail-settings, the task playwright:gmail-settings: Gmail's Settings, Accounts page opened, to set up by hand.
func gmailSettings() { os.Exit(run(os.Args[1:])) }
