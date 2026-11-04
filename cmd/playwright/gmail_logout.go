package main

import "os"

// gmailLogout is playwright gmail-logout, the task playwright:gmail-logout: signed out of Gmail in the Playwright profile.
func gmailLogout() { os.Exit(run(os.Args[1:])) }
