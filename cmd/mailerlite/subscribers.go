package main

import "os"

// subscribers is mailerlite subscribers, the task mailerlite:subscribers: subscribers: list, count, get <email>, add <email> [name], delete <email>.
func subscribers() { os.Exit(run(os.Args[1:])) }
