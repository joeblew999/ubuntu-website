package main

import "os"

// languages is autotranslate languages, the task autotranslate:languages: the languages the provider can translate to.
func languages() { os.Exit(run(os.Args[1:])) }
