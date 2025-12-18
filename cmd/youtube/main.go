package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/youtube"
)

// version is set via ldflags at build time.
var version = "dev"

func main() {
	exitCode := youtube.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
