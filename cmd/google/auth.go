package main

import "os"

// auth is google auth, the task google:auth: signing in to Google: login, open, add, remove, status, check, guide.
func auth() { os.Exit(run(os.Args[1:])) }
