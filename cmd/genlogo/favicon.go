package main

import "os"

// favicon is genlogo favicon, the task genlogo:favicon: the favicon, assets/images/favicon.png.
func favicon() { os.Exit(run(append([]string{"-asset", "favicon"}, os.Args[2:]...))) }
