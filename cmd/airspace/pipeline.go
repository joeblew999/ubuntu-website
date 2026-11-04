package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
	"github.com/joeblew999/ubuntu-website/internal/airspace/tiler"
)

// pipeline is airspace pipeline, the task airspace:pipeline.
func pipeline() {
	fs := flag.NewFlagSet("pipeline", flag.ExitOnError)
	force := fs.Bool("force", false, "Force all steps even if no changes")
	tilerFlag := fs.String("tiler", "auto", "Tiler: auto, tippecanoe, gotiler")
	fs.Parse(os.Args[1:])

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   FAA Airspace Pipeline (Idempotent)   ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// Select tiler
	activeTiler, err := airspace.SelectTiler(*tilerFlag, tiler.New(), gotiler.New())
	if err != nil {
		fatal(err)
	}
	fmt.Printf("Tiler: %s\n\n", activeTiler.Name())

	// Run pipeline
	opts := airspace.PipelineOptions{
		Force:     *force,
		TilerName: *tilerFlag,
	}

	result, err := airspace.Pipeline(opts, activeTiler)
	if err != nil {
		fatal(err)
	}

	// Print summary
	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	if result.Skipped {
		fmt.Println("✓ No changes detected - pipeline complete (idempotent)")
		fmt.Println("  Skipped: tile generation, manifest update")
	} else {
		fmt.Printf("✓ Pipeline complete: %d synced, %d tiled\n",
			result.SyncResult.Updated, result.TileCount)
		fmt.Println("  Next: Run 'airspace upload' to push to R2")
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
