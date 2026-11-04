package main

import "os"

// cloudflareEmailAdd is playwright cloudflare-email-add, the task playwright:cloudflare-email-add: an Email Routing destination address added.
func cloudflareEmailAdd() { os.Exit(run(os.Args[1:])) }
