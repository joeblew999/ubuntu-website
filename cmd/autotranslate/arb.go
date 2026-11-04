package main

import "os"

// arb is autotranslate arb, the task autotranslate:arb: a tokibundle's ARB strings translated: arb <lang>.
func arb() { os.Exit(run(os.Args[1:])) }
