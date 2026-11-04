package main

import "os"

// cloudflareToken is playwright cloudflare-token, the task playwright:cloudflare-token: Cloudflare's API token page opened, to make one.
func cloudflareToken() { os.Exit(run(os.Args[1:])) }
