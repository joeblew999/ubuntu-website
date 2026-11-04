package main

import "os"

// oauth is playwright oauth, the task playwright:oauth: an OAuth flow in the browser, the code (and token, query) from its callback printed as JSON.
func oauth() { os.Exit(run(os.Args[1:])) }
