package main

import "os"

// gmail is google gmail, the task google:gmail: Gmail: list, send, compose, server, check (From is always gerard.webb@ubuntusoftware.net).
func gmail() { os.Exit(run(os.Args[1:])) }
