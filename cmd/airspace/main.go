// Command airspace manages the FAA airspace data behind the BVLOS demo (layouts/fleet/airspace-demo.html): FAA
// ArcGIS → GeoJSON → PMTiles → R2, only the datasets the FAA changed since what R2 has (the sync rules: the
// internal/airspace package doc). Each command is a file of its own and a task of .mise/tasks/airspace.toml:
//
//	mise run airspace:pipeline    sync → tile → manifest, each only if something changed (local dev)
//	mise run airspace:upload      the tiles and manifests to R2 (production)
//	mise run airspace:sync        download only the datasets that changed
//	mise run airspace:tile        GeoJSON to PMTiles
//	mise run airspace:manifest    the manifests
//	mise run airspace:status      the local data files
package main

import (
	"fmt"
	"os"
)

// fatal prints an error and exits.
func fatal(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...) // Remove subcommand from args

	switch cmd {
	case "pipeline":
		pipeline()
	case "sync":
		sync()
	case "tile":
		tile()
	case "manifest":
		manifest()
	case "upload":
		upload()
	case "status":
		status()
	case "history":
		history()
	case "summary":
		summary()
	case "check":
		check()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: airspace <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  pipeline    Full pipeline: sync → tile (if changed) → manifest")
	fmt.Println("  sync        Smart sync FAA data (only download if source changed)")
	fmt.Println("  tile        Convert GeoJSON to PMTiles")
	fmt.Println("  manifest    Generate manifest files with file sizes and metadata")
	fmt.Println("  upload      Upload PMTiles to Cloudflare R2 (manifest-driven)")
	fmt.Println("  status      Show data file status and age")
	fmt.Println("  history     Show sync history and change patterns")
	fmt.Println("  check       Output sync result for GitHub Actions")
	fmt.Println("  summary     Generate GitHub Actions step summary")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -force              Force all steps even if no changes")
	fmt.Println("  -tiler <name>       Tiler: auto, tippecanoe, gotiler (default: auto)")
	fmt.Println("  -dataset <name>     Process single dataset (uas, boundary, sua, airports, navaids)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  airspace pipeline                  # Full idempotent pipeline")
	fmt.Println("  airspace pipeline -tiler gotiler   # Pipeline with pure Go tiler")
	fmt.Println("  airspace sync                      # Sync only changed datasets")
	fmt.Println("  airspace sync -force               # Force re-download all")
	fmt.Println("  airspace tile -dataset uas         # Convert single dataset")
	fmt.Println("  airspace upload                    # Upload to R2")
	fmt.Println("  airspace upload -test              # Test R2 endpoints")
}

// ============================================================================
// Pipeline Command
// ============================================================================

// ============================================================================
// Sync Command
// ============================================================================

// ============================================================================
// Tile Command
// ============================================================================

// ============================================================================
// Manifest Command
// ============================================================================

// ============================================================================
// Status Command
// ============================================================================

// ============================================================================
// History Command
// ============================================================================

// ============================================================================
// Upload Command
// ============================================================================

// ============================================================================
// CI Helpers
// ============================================================================
