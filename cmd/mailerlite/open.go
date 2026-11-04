package main

import "os"

// open is mailerlite open, the task mailerlite:open: the MailerLite dashboard opened (or a page of it).
func open() { os.Exit(run(os.Args[1:])) }
