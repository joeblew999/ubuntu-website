package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// sync is airspace sync, the task airspace:sync.
func sync() {
	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	force := fs.Bool("force", false, "Force re-download even if unchanged")
	skipLarge := fs.Bool("skip-large", true, "Skip datasets >100MB (obstacles)")
	timeout := fs.Duration("timeout", 5*time.Minute, "HTTP request timeout")
	fs.Parse(os.Args[1:])

	fmt.Println("Syncing FAA Airspace Data")
	fmt.Println("=========================")
	if *force {
		fmt.Println("Mode: FORCE (re-downloading all)")
	} else {
		fmt.Println("Mode: Smart (ETag-based diff)")
	}
	fmt.Println()

	opts := airspace.DefaultSyncOptions()
	opts.Force = *force
	opts.SkipLarge = *skipLarge
	opts.Timeout = *timeout

	if !*skipLarge {
		opts.Datasets = airspace.AllDatasets
	}

	result, err := airspace.Sync(opts)
	if err != nil {
		fatal(err)
	}

	// Print dataset results
	for key, ds := range result.Datasets {
		switch ds.Status {
		case "updated":
			fmt.Printf("[%s] ✓ updated (%.1f MB)\n", key, ds.SizeMB)
		case "unchanged":
			fmt.Printf("[%s] unchanged\n", key)
		case "missing":
			fmt.Printf("[%s] downloaded (was missing)\n", key)
		case "error":
			fmt.Printf("[%s] ERROR: %s\n", key, ds.Error)
		}
	}

	fmt.Println()
	fmt.Printf("Done: %d updated, %d unchanged in %s\n", result.Updated, result.Skipped, result.Duration)
	if result.HasChanges {
		fmt.Printf("Downloaded: %.1f MB total\n", result.TotalSizeMB)
		fmt.Println("\nRun 'airspace pipeline' to generate tiles.")
	}
}
