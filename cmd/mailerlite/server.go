package main

import "os"

// server is mailerlite server, the task mailerlite:server: the webhook server for Web3Forms' submissions (PORT=8086, GROUP_ID= to assign a group).
func server() { os.Exit(run(os.Args[1:])) }
