package main

import "os"

// logoSvg is genlogo logo-svg, the task genlogo:logo-svg: the logo SVGs, assets/images/logo.svg and logo-darkmode.svg.
func logoSvg() { os.Exit(run(append([]string{"-asset", "logo-svg"}, os.Args[2:]...))) }
