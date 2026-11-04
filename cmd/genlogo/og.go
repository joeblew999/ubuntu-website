package main

import "os"

// og is genlogo og, the task genlogo:og: the Open Graph image, static/images/og-image.png (social cards).
func og() { os.Exit(run(append([]string{"-asset", "og"}, os.Args[2:]...))) }
