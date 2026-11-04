package main

import "os"

// file is autotranslate file, the task autotranslate:file: one page translated: file <path> <lang>.
func file() { os.Exit(run(os.Args[1:])) }
