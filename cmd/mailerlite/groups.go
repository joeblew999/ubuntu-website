package main

import "os"

// groups is mailerlite groups, the task mailerlite:groups: groups: list, create <name>, subscribers <id>, assign <id> <email>, unassign <id> <email>.
func groups() { os.Exit(run(os.Args[1:])) }
