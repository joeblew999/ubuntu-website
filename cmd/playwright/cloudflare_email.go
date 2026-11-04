package main

import "os"

// cloudflareEmail is playwright cloudflare-email, the task playwright:cloudflare-email: Cloudflare's Email Routing opened.
func cloudflareEmail() { os.Exit(run(os.Args[1:])) }
