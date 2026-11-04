package main

import (
	"os"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// check is airspace check, the task airspace:check.
func check() {
	airspace.PrintCheckOutput(os.Stdout)
}
