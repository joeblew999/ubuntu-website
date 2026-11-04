package main

import "os"

// automations is mailerlite automations, the task mailerlite:automations: automations: list, get <id>.
func automations() { os.Exit(run(os.Args[1:])) }
