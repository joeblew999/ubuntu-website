package main

import "os"

// gmailSendas is playwright gmail-sendas, the task playwright:gmail-sendas: Gmail's Send mail as set up with an SMTP relay: gmail-sendas <name> <email> <provider> <user> <pass>.
func gmailSendas() { os.Exit(run(os.Args[1:])) }
