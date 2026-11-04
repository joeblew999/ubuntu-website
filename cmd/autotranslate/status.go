package main

import "os"

// status is autotranslate status, the task autotranslate:status: the provider's set-up and quota (DeepL's characters this month).
func status() { os.Exit(run(os.Args[1:])) }
