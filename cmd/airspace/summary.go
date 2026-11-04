package main

import (
	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// summary is airspace summary, the task airspace:summary.
func summary() {
	w, cleanup := airspace.GetSummaryWriter()
	defer cleanup()
	airspace.PrintSummary(w)
}
