package main

import "os"

// arbStatus is autotranslate arb-status, the task autotranslate:arb-status: the ARB strings' state per language.
func arbStatus() { os.Exit(run(os.Args[1:])) }
