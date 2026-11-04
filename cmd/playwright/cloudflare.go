package main

import "os"

// cloudflare is playwright cloudflare, the task playwright:cloudflare: the Cloudflare dashboard opened.
func cloudflare() { os.Exit(run(os.Args[1:])) }
