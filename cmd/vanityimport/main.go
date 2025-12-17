// vanityimport manages Go vanity import packages for www.ubuntusoftware.net/pkg/.
//
// Usage:
//
//	vanityimport add <package-name>   Add a new package
//	vanityimport list                 List all packages
//	vanityimport update [name]        Update package metadata
//	vanityimport info <name>          Show package info
//	vanityimport docs <name>          Generate documentation
package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/vanity"
)

// version is set via ldflags at build time
var version = "dev"

func main() {
	exitCode := vanity.Run(os.Args, version, os.Stdout, os.Stderr)
	os.Exit(exitCode)
}
