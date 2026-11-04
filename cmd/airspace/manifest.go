package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// manifest is airspace manifest, the task airspace:manifest: the manifests again, for the last sync
// (sync_result.json) against what's published in R2 (internal/airspace doc.go, rule 5).
func manifest() {
	fmt.Println("Generating Airspace Manifests")
	fmt.Println("=============================")
	fmt.Println()

	published, err := airspace.FetchPublished(&http.Client{Timeout: time.Minute}, airspace.PublishedURL)
	if err != nil {
		fatal(err)
	}
	run := airspace.LoadSyncResult(filepath.Join(airspace.DirData, airspace.FileSyncResult))
	if err := airspace.GenerateManifests(airspace.DirData, airspace.DirPMTiles, airspace.DirGeoJSON, run, published); err != nil {
		fatal(err)
	}

	fmt.Printf("✓ %s\n", filepath.Join(airspace.DirData, airspace.FileManifest))
	fmt.Printf("✓ %s\n", filepath.Join(airspace.DirData, airspace.FileUSAManifest))
	fmt.Println("\nManifests updated.")
}
