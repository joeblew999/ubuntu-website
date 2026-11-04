package main

import (
	"fmt"
	"path/filepath"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// manifest is airspace manifest, the task airspace:manifest.
func manifest() {
	fmt.Println("Generating Airspace Manifests")
	fmt.Println("=============================")
	fmt.Println()

	if err := airspace.GenerateManifests(airspace.DirData, airspace.DirPMTiles, airspace.DirGeoJSON); err != nil {
		fatal(err)
	}

	fmt.Printf("✓ %s\n", filepath.Join(airspace.DirData, airspace.FileManifest))
	fmt.Printf("✓ %s\n", filepath.Join(airspace.DirData, airspace.FileUSAManifest))
	fmt.Println("\nManifests updated.")
}
