// genlogo generates Ubuntu Software logo assets using Go graphics.
//
// SINGLE SOURCE OF TRUTH for all logo/branding assets.
//
// Usage:
//
//	go run cmd/genlogo/main.go -asset all      # Generate everything
//	go run cmd/genlogo/main.go -asset favicon  # Generate specific asset
//	task genlogo:all                           # Via Taskfile
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/genlogo"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	exitCode := genlogo.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
