package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
	"github.com/joeblew999/ubuntu-website/internal/airspace/tiler"
)

// tile is airspace tile, the task airspace:tile.
func tile() {
	fs := flag.NewFlagSet("tile", flag.ExitOnError)
	datasetFlag := fs.String("dataset", "", "Specific dataset to tile (empty = all)")
	force := fs.Bool("force", false, "Force regenerate even if up-to-date")
	tilerFlag := fs.String("tiler", "auto", "Tiler: auto, tippecanoe, gotiler")
	fs.Parse(os.Args[1:])

	// Select tiler
	activeTiler, err := airspace.SelectTiler(*tilerFlag, tiler.New(), gotiler.New())
	if err != nil {
		fatal(err)
	}

	fmt.Printf("Converting GeoJSON to PMTiles (using %s)\n", activeTiler.Name())
	fmt.Println("=========================================")
	fmt.Println()

	if *datasetFlag != "" {
		// Single dataset
		if err := airspace.TileOne(activeTiler, *datasetFlag, *force); err != nil {
			fatal(err)
		}
		fmt.Printf("[%s] ✓ done\n", *datasetFlag)
	} else {
		// All datasets
		count, err := airspace.TileAll(activeTiler, *force)
		if err != nil {
			fatal(err)
		}
		fmt.Printf("\nDone: %d datasets tiled\n", count)
	}
}
