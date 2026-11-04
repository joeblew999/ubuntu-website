package main

import "os"

// gmailSwitch is playwright gmail-switch, the task playwright:gmail-switch: Google's account chooser opened, to add or switch an account.
func gmailSwitch() { os.Exit(run(os.Args[1:])) }
