package main

import "os"

// webhooks is mailerlite webhooks, the task mailerlite:webhooks: webhooks: list, create <url> [events], delete <id>.
func webhooks() { os.Exit(run(os.Args[1:])) }
