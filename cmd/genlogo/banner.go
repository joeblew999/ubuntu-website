package main

import "os"

// banner is genlogo banner, the task genlogo:banner: the Bluesky banner, static/images/bluesky-banner.png.
func banner() { os.Exit(run(append([]string{"-asset", "banner"}, os.Args[2:]...))) }
