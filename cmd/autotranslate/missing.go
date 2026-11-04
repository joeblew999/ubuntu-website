package main

import "os"

// missing is autotranslate missing, the task autotranslate:missing: every page missing in a language translated: missing <lang>.
func missing() { os.Exit(run(os.Args[1:])) }
