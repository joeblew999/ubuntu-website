package main

import (
	"flag"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// upload is airspace upload, the task airspace:upload.
func upload() {
	fs := flag.NewFlagSet("upload", flag.ExitOnError)
	testOnly := fs.Bool("test", false, "Test endpoints only (don't upload)")
	fs.Parse(os.Args[1:])

	if *testOnly {
		// Test mode - check if endpoints are accessible
		if err := airspace.TestR2Endpoints(); err != nil {
			fatal(err)
		}
		return
	}

	// Upload mode
	if err := airspace.UploadToR2(); err != nil {
		fatal(err)
	}
}
