package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// status is airspace status, the task airspace:status.
func status() {
	fmt.Println("Airspace Data Status")
	fmt.Println("====================")
	fmt.Println()

	found := 0
	for _, key := range airspace.AllDatasets {
		ds := airspace.Datasets[key]
		path := filepath.Join(airspace.DirGeoJSON, ds.GeoJSON)
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("  [%s] %s: NOT FOUND\n", key, ds.Name)
			continue
		}

		found++
		age := time.Since(info.ModTime())
		fmt.Printf("  [%s] %s: %.1f MB (%s old)\n",
			key, ds.Name, float64(info.Size())/(1024*1024), airspace.FormatAge(age))
	}

	fmt.Println()
	fmt.Printf("Found %d/%d datasets.\n", found, len(airspace.AllDatasets))
}
