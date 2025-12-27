// Command airspace manages FAA airspace data for the BVLOS demo.
//
// Usage:
//
//	airspace pipeline             # Full pipeline: sync → tile → manifest
//	airspace sync                 # Smart sync (only download if changed)
//	airspace tile                 # Convert GeoJSON to PMTiles
//	airspace manifest             # Generate manifest files
//	airspace status               # Show data file status
//	airspace download             # Download FAA data (use sync instead)
//
// See also: layouts/fleet/airspace-demo.html
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
	"github.com/joeblew999/ubuntu-website/internal/airspace/tiler"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...) // Remove subcommand from args

	switch cmd {
	case "pipeline":
		runPipeline()
	case "sync":
		runSync()
	case "tile":
		runTile()
	case "manifest":
		runManifest()
	case "download":
		runDownload()
	case "status":
		runStatus()
	case "history":
		runHistory()
	case "summary":
		runSummary()
	case "check":
		runCheck()
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
	fmt.Println("  download    Download FAA airspace data (use sync instead)")
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
}

// ============================================================================
// Pipeline Command
// ============================================================================

func runPipeline() {
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Tiler: %s\n\n", activeTiler.Name())

	// Run pipeline
	opts := airspace.PipelineOptions{
		Force:     *force,
		TilerName: *tilerFlag,
	}

	result, err := airspace.Pipeline(opts, activeTiler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pipeline error: %v\n", err)
		os.Exit(1)
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
		fmt.Println("  Next: Run 'task r2:airspace:upload' to push to R2")
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// ============================================================================
// Sync Command
// ============================================================================

func runSync() {
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
		fmt.Fprintf(os.Stderr, "Sync error: %v\n", err)
		os.Exit(1)
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

// ============================================================================
// Tile Command
// ============================================================================

func runTile() {
	fs := flag.NewFlagSet("tile", flag.ExitOnError)
	datasetFlag := fs.String("dataset", "", "Specific dataset to tile (empty = all)")
	force := fs.Bool("force", false, "Force regenerate even if up-to-date")
	tilerFlag := fs.String("tiler", "auto", "Tiler: auto, tippecanoe, gotiler")
	fs.Parse(os.Args[1:])

	// Select tiler
	activeTiler, err := airspace.SelectTiler(*tilerFlag, tiler.New(), gotiler.New())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converting GeoJSON to PMTiles (using %s)\n", activeTiler.Name())
	fmt.Println("=========================================")
	fmt.Println()

	if *datasetFlag != "" {
		// Single dataset
		if err := airspace.TileOne(activeTiler, *datasetFlag, *force); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("[%s] ✓ done\n", *datasetFlag)
	} else {
		// All datasets
		count, err := airspace.TileAll(activeTiler, *force)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("\nDone: %d datasets tiled\n", count)
	}
}

// ============================================================================
// Manifest Command
// ============================================================================

func runManifest() {
	fmt.Println("Generating Airspace Manifests")
	fmt.Println("=============================")
	fmt.Println()

	if err := airspace.GenerateManifests(airspace.DirData, airspace.DirPMTiles, airspace.DirGeoJSON); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ %s\n", filepath.Join(airspace.DirData, airspace.FileManifest))
	fmt.Printf("✓ %s\n", filepath.Join(airspace.DirData, airspace.FileUSAManifest))
	fmt.Println("\nManifests updated.")
}

// ============================================================================
// Download Command (legacy - use sync instead)
// ============================================================================

func runDownload() {
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	outputDir := fs.String("output", airspace.DirGeoJSON, "Output directory")
	datasetFlag := fs.String("dataset", "", "Specific dataset (empty = all)")
	timeout := fs.Duration("timeout", 5*time.Minute, "HTTP timeout")
	fs.Parse(os.Args[1:])

	fmt.Println("Downloading FAA Data")
	fmt.Println("====================")
	fmt.Println("Note: Use 'airspace sync' for incremental updates.")
	fmt.Println()

	client := &http.Client{Timeout: *timeout}

	var datasets []string
	if *datasetFlag != "" {
		datasets = []string{*datasetFlag}
	} else {
		datasets = airspace.AllDatasets
	}

	if err := airspace.Download(client, *outputDir, datasets); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDone.")
}

// ============================================================================
// Status Command
// ============================================================================

func runStatus() {
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

// ============================================================================
// History Command
// ============================================================================

func runHistory() {
	historyFile := filepath.Join(airspace.DirData, airspace.FileSyncHistory)
	history := airspace.LoadSyncHistory(historyFile)

	fmt.Println("FAA Airspace Sync History")
	fmt.Println("=========================")
	fmt.Println()

	if history.TotalRuns == 0 {
		fmt.Println("No sync history yet. Run 'airspace sync' to start tracking.")
		return
	}

	fmt.Println("Summary:")
	fmt.Printf("  Total sync runs: %d\n", history.TotalRuns)
	fmt.Printf("  Changes detected: %d times (%.1f%%)\n", history.ChangeCount,
		float64(history.ChangeCount)/float64(history.TotalRuns)*100)
	fmt.Printf("  Last change: %s\n", airspace.FormatTimeSince(history.LastChange))
	fmt.Printf("  Average duration: %s\n", history.AvgDuration)
	fmt.Println()

	fmt.Println("Recent Runs (newest first):")
	fmt.Println("----------------------------")
	for i, run := range history.Runs {
		if i >= 10 {
			fmt.Printf("  ... and %d more\n", len(history.Runs)-10)
			break
		}
		status := "no changes"
		if run.HasChanges {
			status = fmt.Sprintf("%d updated, %.1f MB", run.Updated, run.TotalSizeMB)
		}
		fmt.Printf("  %s  %s  (%s)\n",
			run.Timestamp.Format("2006-01-02 15:04"),
			status,
			run.Duration)
	}
}

// ============================================================================
// CI Helpers
// ============================================================================

func runCheck() {
	resultPath := filepath.Join(airspace.DirData, airspace.FileSyncResult)
	result := airspace.LoadSyncResult(resultPath)

	fmt.Printf("has_changes=%t\n", result.HasChanges)
	fmt.Printf("updated_count=%d\n", result.Updated)
}

func runSummary() {
	var out *os.File
	if summaryPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryPath != "" {
		f, err := os.OpenFile(summaryPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			out = os.Stdout
		} else {
			defer f.Close()
			out = f
		}
	} else {
		out = os.Stdout
	}

	fmt.Fprintln(out, "## Airspace Sync Complete")
	fmt.Fprintln(out)

	// Sync result
	resultPath := filepath.Join(airspace.DirData, airspace.FileSyncResult)
	if data, err := os.ReadFile(resultPath); err == nil {
		fmt.Fprintln(out, "### Sync Result")
		fmt.Fprintln(out, "```json")
		fmt.Fprintln(out, string(data))
		fmt.Fprintln(out, "```")
		fmt.Fprintln(out)
	}

	// PMTiles files
	fmt.Fprintln(out, "### PMTiles Files")
	fmt.Fprintln(out, "| File | Size |")
	fmt.Fprintln(out, "|------|------|")

	files, _ := filepath.Glob(filepath.Join(airspace.DirPMTiles, "*.pmtiles"))
	if len(files) == 0 {
		fmt.Fprintln(out, "| _No PMTiles found_ | - |")
	} else {
		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}
			sizeMB := float64(info.Size()) / 1024 / 1024
			fmt.Fprintf(out, "| %s | %.1f MB |\n", filepath.Base(f), sizeMB)
		}
	}
	fmt.Fprintln(out)
}
