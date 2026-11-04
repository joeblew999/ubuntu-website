package main

import "os"

// menu is translate menu <subcommand>, the task translate:menu: the menus checked for broken links (check), or
// generated from English (sync).
func menu() { os.Exit(run(os.Args[1:])) }
