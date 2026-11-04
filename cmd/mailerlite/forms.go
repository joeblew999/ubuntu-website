package main

import "os"

// forms is mailerlite forms, the task mailerlite:forms: forms: list, subscribers <id>.
func forms() { os.Exit(run(os.Args[1:])) }
