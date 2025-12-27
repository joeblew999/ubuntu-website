// Command airspace manages FAA airspace data for the BVLOS demo.
//
// Usage:
//
//	airspace download             # Download all datasets
//	airspace download -dataset uas
//	airspace sync                 # Smart sync (only download if changed)
//	airspace status               # Show data file status
//	airspace pipeline             # Full pipeline: sync → tile → manifest
//	airspace manifest             # Generate manifest files
//	airspace tile                 # Convert GeoJSON to PMTiles (requires tippecanoe)
//
// Configuration:
//
//	All file paths are defined as constants below - no external dependencies.
//	See also: layouts/fleet/airspace-demo.html
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
)

// =============================================================================
// File Path Constants - Single source of truth for all file locations
// =============================================================================

const (
	// Directory paths
	DirGeoJSON   = "static/airspace"          // GeoJSON output directory
	DirPMTiles   = "static/airspace/tiles"    // PMTiles output directory
	DirData      = "data/airspace"            // Data/metadata directory

	// Data files (in DirData) - use underscores for Hugo data access compatibility
	FileSyncETags   = "sync_etags.json"       // ETag cache for change detection
	FileSyncResult  = "sync_result.json"      // Last sync result (for pipeline idempotency)
	FileSyncHistory = "sync_history.json"     // Rolling sync history

	// Manifest files (in DirData, copied to static) - prefix with "manifest_" for alphabetical grouping
	FileManifest    = "manifest.json"         // Global manifest
	FileUSAManifest = "manifest_usa.json"     // USA regional manifest

	// GeoJSON files (in DirGeoJSON)
	GeoJSONBoundary = "faa_airspace_boundary.geojson"
	GeoJSONSUA      = "faa_special_use_airspace.geojson"
	GeoJSONUAS      = "faa_uas_facility_map.geojson"
	GeoJSONAirports = "faa_airports.geojson"
	GeoJSONNavaids  = "faa_navaids.geojson"
	GeoJSONObstacles = "faa_obstacles.geojson"

	// PMTiles files (in DirPMTiles)
	PMTilesBoundary = "faa_airspace_boundary.pmtiles"
	PMTilesSUA      = "faa_special_use_airspace.pmtiles"
	PMTilesUAS      = "faa_uas_facility_map.pmtiles"
	PMTilesAirports = "faa_airports.pmtiles"
	PMTilesNavaids  = "faa_navaids.pmtiles"
	PMTilesObstacles = "faa_obstacles.pmtiles"
	PMTilesCombined = "faa_airspace_combined.pmtiles"

	// PMTiles layer names (used in tippecanoe and map rendering)
	LayerBoundary = "boundary"
	LayerSUA      = "sua"
	LayerUAS      = "uas"
	LayerAirports = "airports"
	LayerNavaids  = "navaids"
	LayerObstacles = "obstacles"
)

// Dataset configuration
type Dataset struct {
	Name        string
	Key         string // Dataset key (uas, boundary, etc.)
	GeoJSON     string // GeoJSON filename (uses constants above)
	PMTiles     string // PMTiles filename (uses constants above)
	Layer       string // PMTiles layer name
	BaseURL     string
	IsPaginated bool   // For FeatureServer APIs that require pagination
	PageSize    int
	ETagURL     string // URL to check for ETag/Last-Modified (for paginated APIs)
}

// TileConfig holds tippecanoe settings per dataset
type TileConfig struct {
	MinZoom          int
	MaxZoom          int
	DropDensest      bool   // --drop-densest-as-needed
	NoFeatureLimit   bool   // --no-feature-limit
	NoTileSizeLimit  bool   // --no-tile-size-limit
	ReduceRate       int    // -r rate (1 = no reduction)
}

var datasets = map[string]Dataset{
	"uas": {
		Name:        "UAS Facility Map",
		Key:         "uas",
		GeoJSON:     GeoJSONUAS,
		PMTiles:     PMTilesUAS,
		Layer:       LayerUAS,
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    2000,
		ETagURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0",
	},
	"boundary": {
		Name:    "Airspace Boundary",
		Key:     "boundary",
		GeoJSON: GeoJSONBoundary,
		PMTiles: PMTilesBoundary,
		Layer:   LayerBoundary,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/67885972e4e940b2aa6d74024901c561/geojson?layers=0",
	},
	"sua": {
		Name:    "Special Use Airspace",
		Key:     "sua",
		GeoJSON: GeoJSONSUA,
		PMTiles: PMTilesSUA,
		Layer:   LayerSUA,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/dd0d1b726e504137ab3c41b21835d05b/geojson?layers=0",
	},
	"airports": {
		Name:    "Airports",
		Key:     "airports",
		GeoJSON: GeoJSONAirports,
		PMTiles: PMTilesAirports,
		Layer:   LayerAirports,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/e747ab91a11045e8b3f8a3efd093d3b5/geojson?layers=0",
	},
	"navaids": {
		Name:    "Navigation Aids",
		Key:     "navaids",
		GeoJSON: GeoJSONNavaids,
		PMTiles: PMTilesNavaids,
		Layer:   LayerNavaids,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/990e238991b44dd08af27d7b43e70b92/geojson?layers=0",
	},
	"obstacles": {
		Name:    "Obstacles",
		Key:     "obstacles",
		GeoJSON: GeoJSONObstacles,
		PMTiles: PMTilesObstacles,
		Layer:   LayerObstacles,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/c6a62360338e408cb1512366ad61559e/geojson?layers=0",
	},
}

// tileConfigs holds tippecanoe configuration per dataset
var tileConfigs = map[string]TileConfig{
	"boundary": {MinZoom: -1, MaxZoom: -1, DropDensest: false},  // -zg (auto zoom)
	"sua":      {MinZoom: -1, MaxZoom: -1, DropDensest: false},  // -zg
	"uas":      {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"airports": {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"navaids":  {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"obstacles": {MinZoom: -1, MaxZoom: -1, DropDensest: true},  // -zg --drop-densest-as-needed
}

// datasetOrder defines the processing order (consistent ordering)
var datasetOrder = []string{"uas", "boundary", "sua", "airports", "navaids"}

// datasetOrderWithObstacles includes obstacles (large file, skipped by default)
var datasetOrderWithObstacles = []string{"uas", "boundary", "sua", "airports", "navaids", "obstacles"}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...) // Remove subcommand from args

	switch cmd {
	case "download":
		runDownload()
	case "sync":
		runSync()
	case "tile":
		runTile()
	case "manifest":
		runManifest()
	case "pipeline":
		runPipeline()
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
	fmt.Println("  sync        Smart sync FAA data (only download if source changed)")
	fmt.Println("  tile        Convert GeoJSON to PMTiles")
	fmt.Println("  manifest    Generate manifest files with file sizes and metadata")
	fmt.Println("  pipeline    Full pipeline: sync → tile (if changed) → manifest")
	fmt.Println("  download    Download FAA airspace data (use sync instead)")
	fmt.Println("  status      Show data file status and age")
	fmt.Println("  history     Show sync history and change patterns")
	fmt.Println("  check       Output sync result for GitHub Actions")
	fmt.Println("  summary     Generate GitHub Actions step summary")
	fmt.Println()
	fmt.Println("Tiler Options (for tile command):")
	fmt.Println("  -tiler auto       Auto-detect (tippecanoe if available, else gotiler)")
	fmt.Println("  -tiler tippecanoe Use tippecanoe (external binary)")
	fmt.Println("  -tiler gotiler    Use pure Go tiler (no dependencies)")
	fmt.Println()
	fmt.Println("Datasets:")
	fmt.Println("  uas         UAS Facility Map (LAANC ceiling altitudes)")
	fmt.Println("  boundary    Airspace Boundary (Class B/C/D/E)")
	fmt.Println("  sua         Special Use Airspace (MOAs, Restricted, Prohibited)")
	fmt.Println("  airports    Airports (US, PR, VI)")
	fmt.Println("  navaids     Navigation Aids (VOR, NDB, etc.)")
	fmt.Println("  obstacles   Obstacles (towers, buildings, etc.)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  airspace pipeline                       # Full idempotent pipeline")
	fmt.Println("  airspace pipeline -tiler gotiler        # Pipeline with pure Go tiler")
	fmt.Println("  airspace sync                           # Sync only changed datasets")
	fmt.Println("  airspace sync -force                    # Force re-download all")
	fmt.Println("  airspace tile                           # Convert using auto-detected tiler")
	fmt.Println("  airspace tile -tiler gotiler            # Convert using pure Go tiler")
	fmt.Println("  airspace tile -dataset uas              # Convert single dataset")
	fmt.Println("  airspace manifest                       # Generate manifest files")
	fmt.Println("  airspace status                         # Show data status")
}

// ============================================================================
// Download Command
// ============================================================================

func runDownload() {
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	outputDir := fs.String("output", DirGeoJSON, "Output directory for GeoJSON files")
	datasetFlag := fs.String("dataset", "", "Specific dataset (uas, boundary, sua). Empty = all")
	timeout := fs.Duration("timeout", 5*time.Minute, "HTTP request timeout")
	fs.Parse(os.Args[1:])

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Dataset order for consistent processing
	datasetOrder := []string{"uas", "boundary", "sua", "airports", "navaids", "obstacles"}

	// Determine which datasets to download
	toDownload := make([]Dataset, 0)
	if *datasetFlag != "" {
		ds, ok := datasets[*datasetFlag]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown dataset: %s (valid: %v)\n", *datasetFlag, datasetOrder)
			os.Exit(1)
		}
		toDownload = append(toDownload, ds)
	} else {
		// Download all in consistent order
		for _, key := range datasetOrder {
			toDownload = append(toDownload, datasets[key])
		}
	}

	client := &http.Client{Timeout: *timeout}

	for _, ds := range toDownload {
		fmt.Printf("Downloading %s...\n", ds.Name)
		outPath := filepath.Join(*outputDir, ds.GeoJSON)

		var err error
		if ds.IsPaginated {
			err = downloadPaginated(client, ds, outPath)
		} else {
			err = downloadDirect(client, ds.BaseURL, outPath)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", ds.Name, err)
			os.Exit(1)
		}

		// Report file size
		if info, err := os.Stat(outPath); err == nil {
			fmt.Printf("  ✓ %s (%.1f MB)\n", ds.GeoJSON, float64(info.Size())/(1024*1024))
		}
	}

	fmt.Println()
	fmt.Println("Done. Run 'task r2:airspace:upload' to sync to R2.")
}

// downloadDirect downloads a file directly (no pagination)
func downloadDirect(client *http.Client, url, outPath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// downloadPaginated handles ArcGIS FeatureServer pagination
func downloadPaginated(client *http.Client, ds Dataset, outPath string) error {
	type FeatureCollection struct {
		Type     string        `json:"type"`
		Features []interface{} `json:"features"`
	}

	collection := FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]interface{}, 0),
	}

	offset := 0
	for {
		params := url.Values{}
		params.Set("where", "1=1")
		params.Set("outFields", "*")
		params.Set("f", "geojson")
		params.Set("resultRecordCount", fmt.Sprintf("%d", ds.PageSize))
		params.Set("resultOffset", fmt.Sprintf("%d", offset))

		queryURL := ds.BaseURL + "?" + params.Encode()

		resp, err := client.Get(queryURL)
		if err != nil {
			return fmt.Errorf("fetch page at offset %d: %w", offset, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("HTTP %d at offset %d", resp.StatusCode, offset)
		}

		var page FeatureCollection
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			resp.Body.Close()
			return fmt.Errorf("decode page at offset %d: %w", offset, err)
		}
		resp.Body.Close()

		collection.Features = append(collection.Features, page.Features...)
		fmt.Printf("  ... fetched %d features (total: %d)\n", len(page.Features), len(collection.Features))

		if len(page.Features) < ds.PageSize {
			break
		}
		offset += ds.PageSize
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(collection)
}

// ============================================================================
// Status Command
// ============================================================================

func runStatus() {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	outputDir := fs.String("output", DirGeoJSON, "Data directory")
	fs.Parse(os.Args[1:])

	fmt.Println("Airspace Data Status")
	fmt.Println("====================")
	fmt.Println()

	datasetOrder := []string{"uas", "boundary", "sua", "airports", "navaids", "obstacles"}
	found := 0
	for _, key := range datasetOrder {
		ds := datasets[key]
		path := filepath.Join(*outputDir, ds.GeoJSON)
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("  [%s] %s: NOT FOUND\n", key, ds.Name)
			continue
		}

		found++
		age := time.Since(info.ModTime())
		fmt.Printf("  [%s] %s: %.1f MB (%s old)\n",
			key, ds.Name, float64(info.Size())/(1024*1024), formatAge(age))
	}

	fmt.Println()
	fmt.Printf("Found %d/%d datasets.\n", found, len(datasetOrder))
	fmt.Println("Run 'airspace download' to refresh data.")
}

func formatAge(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%d min", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}

// ============================================================================
// History Command - Show sync history for cron tuning
// ============================================================================

func runHistory() {
	fs := flag.NewFlagSet("history", flag.ExitOnError)
	dataDir := fs.String("data-dir", DirData, "Data directory")
	fs.Parse(os.Args[1:])

	historyFile := filepath.Join(*dataDir, FileSyncHistory)
	history := loadSyncHistory(historyFile)

	fmt.Println("FAA Airspace Sync History")
	fmt.Println("=========================")
	fmt.Println()

	if history.TotalRuns == 0 {
		fmt.Println("No sync history yet. Run 'airspace sync' to start tracking.")
		return
	}

	// Summary stats
	fmt.Println("Summary:")
	fmt.Printf("  Total sync runs: %d\n", history.TotalRuns)
	fmt.Printf("  Changes detected: %d times (%.1f%%)\n", history.ChangeCount,
		float64(history.ChangeCount)/float64(history.TotalRuns)*100)
	fmt.Printf("  Last change: %s\n", formatTimeSince(history.LastChange))
	fmt.Printf("  Average duration: %s\n", history.AvgDuration)
	fmt.Println()

	// Recent runs
	fmt.Println("Recent Runs (newest first):")
	fmt.Println("----------------------------")
	for i, run := range history.Runs {
		if i >= 10 {
			fmt.Printf("  ... and %d more (see sync-history.json)\n", len(history.Runs)-10)
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
	fmt.Println()

	// Cron recommendation
	fmt.Println("Cron Recommendation:")
	if history.ChangeCount == 0 {
		fmt.Println("  No changes detected yet - FAA data appears stable.")
		fmt.Println("  Weekly sync (current) is appropriate.")
	} else {
		daysSinceChange := time.Since(history.LastChange).Hours() / 24
		if daysSinceChange < 7 {
			fmt.Println("  Recent change detected - consider daily sync temporarily.")
		} else if daysSinceChange < 28 {
			fmt.Println("  Change within AIRAC cycle - weekly sync is appropriate.")
		} else {
			fmt.Println("  No recent changes - weekly sync is sufficient.")
		}
	}
	fmt.Printf("  Expected sync time: %s (no changes) to ~5-10 min (full download)\n", history.AvgDuration)
}

// ============================================================================
// Sync Command - Smart download with ETag tracking
// ============================================================================

// ETagStore holds ETags for each dataset to detect changes
type ETagStore struct {
	ETags     map[string]string `json:"etags"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// SyncResult records the outcome of a sync run for history/metrics
type SyncResult struct {
	Timestamp   time.Time              `json:"timestamp"`
	Duration    string                 `json:"duration"`
	DurationMs  int64                  `json:"duration_ms"`
	Updated     int                    `json:"updated"`
	Skipped     int                    `json:"skipped"`
	HasChanges  bool                   `json:"has_changes"`
	Datasets    map[string]DatasetSync `json:"datasets"`
	TotalBytes  int64                  `json:"total_bytes"`
	TotalSizeMB float64                `json:"total_size_mb"`
}

// DatasetSync records per-dataset sync metrics
type DatasetSync struct {
	Status     string  `json:"status"` // "unchanged", "updated", "missing", "error"
	DurationMs int64   `json:"duration_ms,omitempty"`
	SizeBytes  int64   `json:"size_bytes,omitempty"`
	SizeMB     float64 `json:"size_mb,omitempty"`
	Features   int     `json:"features,omitempty"`
	ETag       string  `json:"etag,omitempty"`
	Error      string  `json:"error,omitempty"`
}

// SyncHistory maintains a rolling log of sync runs for cron tuning
type SyncHistory struct {
	LastChange  time.Time    `json:"last_change"`           // When FAA data last changed
	Runs        []SyncResult `json:"runs"`                  // Recent sync runs (keep last 20)
	ChangeCount int          `json:"change_count"`          // Total changes detected
	TotalRuns   int          `json:"total_runs"`            // Total sync runs
	AvgDuration string       `json:"avg_duration"`          // Average run duration
	AvgDurationMs int64      `json:"avg_duration_ms"`
}

const maxHistoryRuns = 20

func runSync() {
	syncStart := time.Now()

	fs := flag.NewFlagSet("sync", flag.ExitOnError)
	outputDir := fs.String("output", DirGeoJSON, "Output directory for GeoJSON files")
	dataDir := fs.String("data-dir", DirData, "Data directory for ETags and history")
	force := fs.Bool("force", false, "Force re-download even if unchanged")
	timeout := fs.Duration("timeout", 5*time.Minute, "HTTP request timeout")
	skipLarge := fs.Bool("skip-large", true, "Skip datasets >100MB (obstacles)")
	fs.Parse(os.Args[1:])

	etagFile := filepath.Join(*dataDir, FileSyncETags)
	historyFile := filepath.Join(*dataDir, FileSyncHistory)
	resultFile := filepath.Join(*dataDir, FileSyncResult)

	// Create directories
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data dir: %v\n", err)
		os.Exit(1)
	}

	// Load existing ETags and history
	store := loadETags(etagFile)
	history := loadSyncHistory(historyFile)
	client := &http.Client{Timeout: *timeout}

	// Datasets to sync (exclude obstacles by default - too large for GeoJSON)
	datasetOrder := []string{"uas", "boundary", "sua", "airports", "navaids"}
	if !*skipLarge {
		datasetOrder = append(datasetOrder, "obstacles")
	}

	fmt.Println("Syncing FAA Airspace Data")
	fmt.Println("=========================")
	if *force {
		fmt.Println("Mode: FORCE (re-downloading all)")
	} else {
		fmt.Println("Mode: Smart (ETag-based diff)")
	}
	fmt.Println()

	// Track results
	result := SyncResult{
		Timestamp: time.Now().UTC(),
		Datasets:  make(map[string]DatasetSync),
	}
	var totalBytes int64

	for _, key := range datasetOrder {
		dsStart := time.Now()
		ds := datasets[key]
		outPath := filepath.Join(*outputDir, ds.GeoJSON)
		dsResult := DatasetSync{}

		// Check if we need to download
		needsDownload := *force
		var newETag string

		if !needsDownload {
			// Check if file exists locally
			if _, err := os.Stat(outPath); os.IsNotExist(err) {
				needsDownload = true
				dsResult.Status = "missing"
				fmt.Printf("[%s] %s: MISSING\n", key, ds.Name)
			}
		}

		if !needsDownload {
			// Check for changes using ETag or Last-Modified header
			checkURL := ds.BaseURL
			if ds.ETagURL != "" {
				checkURL = ds.ETagURL
			}
			newETag, needsDownload = checkETag(client, checkURL, store.ETags[key])
			if needsDownload {
				dsResult.Status = "changed"
				fmt.Printf("[%s] %s: CHANGED (source updated)\n", key, ds.Name)
			} else {
				dsResult.Status = "unchanged"
				dsResult.ETag = store.ETags[key]
				dsResult.DurationMs = time.Since(dsStart).Milliseconds()
				result.Datasets[key] = dsResult
				fmt.Printf("[%s] %s: unchanged\n", key, ds.Name)
				result.Skipped++
				continue
			}
		}

		// Download the dataset
		fmt.Printf("  Downloading %s...\n", ds.Name)
		var err error
		if ds.IsPaginated {
			err = downloadPaginated(client, ds, outPath)
		} else {
			err = downloadDirect(client, ds.BaseURL, outPath)
		}

		if err != nil {
			dsResult.Status = "error"
			dsResult.Error = err.Error()
			dsResult.DurationMs = time.Since(dsStart).Milliseconds()
			result.Datasets[key] = dsResult
			fmt.Fprintf(os.Stderr, "  ERROR: %v\n", err)
			continue
		}

		// Update ETag/Last-Modified for change tracking
		if newETag != "" {
			store.ETags[key] = newETag
			dsResult.ETag = newETag
		} else {
			// For paginated downloads, do HEAD to get current ETag
			checkURL := ds.BaseURL
			if ds.ETagURL != "" {
				checkURL = ds.ETagURL
			}
			if etag, _ := checkETag(client, checkURL, ""); etag != "" {
				store.ETags[key] = etag
				dsResult.ETag = etag
			}
		}

		// Report file size and record metrics
		if info, err := os.Stat(outPath); err == nil {
			sizeBytes := info.Size()
			sizeMB := float64(sizeBytes) / (1024 * 1024)
			totalBytes += sizeBytes
			dsResult.SizeBytes = sizeBytes
			dsResult.SizeMB = sizeMB
			fmt.Printf("  ✓ %s (%.1f MB)\n", ds.GeoJSON, sizeMB)
		}

		dsResult.Status = "updated"
		dsResult.DurationMs = time.Since(dsStart).Milliseconds()
		result.Datasets[key] = dsResult
		result.Updated++
	}

	// Finalize result
	syncDuration := time.Since(syncStart)
	result.Duration = syncDuration.Round(time.Millisecond).String()
	result.DurationMs = syncDuration.Milliseconds()
	result.HasChanges = result.Updated > 0
	result.TotalBytes = totalBytes
	result.TotalSizeMB = float64(totalBytes) / (1024 * 1024)

	// Save ETags
	store.UpdatedAt = time.Now()
	saveETags(etagFile, store)

	// Update and save history
	if result.HasChanges {
		history.LastChange = result.Timestamp
		history.ChangeCount++
	}
	history.TotalRuns++
	history.Runs = append([]SyncResult{result}, history.Runs...)
	if len(history.Runs) > maxHistoryRuns {
		history.Runs = history.Runs[:maxHistoryRuns]
	}
	// Calculate average duration
	var totalMs int64
	for _, r := range history.Runs {
		totalMs += r.DurationMs
	}
	history.AvgDurationMs = totalMs / int64(len(history.Runs))
	history.AvgDuration = time.Duration(history.AvgDurationMs * int64(time.Millisecond)).String()
	saveSyncHistory(historyFile, history)

	// Save current result (for pipeline idempotency)
	saveSyncResult(resultFile, result)

	// Print summary
	fmt.Println()
	fmt.Printf("Done: %d updated, %d unchanged in %s\n", result.Updated, result.Skipped, result.Duration)
	if result.HasChanges {
		fmt.Printf("Downloaded: %.1f MB total\n", result.TotalSizeMB)
	}
	fmt.Println()

	// Print history summary
	fmt.Println("Sync History:")
	fmt.Printf("  Total runs: %d\n", history.TotalRuns)
	fmt.Printf("  Changes detected: %d times\n", history.ChangeCount)
	fmt.Printf("  Last change: %s\n", formatTimeSince(history.LastChange))
	fmt.Printf("  Average duration: %s\n", history.AvgDuration)
	fmt.Println()

	if result.HasChanges {
		fmt.Println("Changes detected. Run 'airspace pipeline' to generate tiles and upload.")
	} else {
		fmt.Println("No changes. Pipeline can skip tile/manifest/upload steps.")
	}
}

func formatTimeSince(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	d := time.Since(t)
	if d < time.Hour {
		return fmt.Sprintf("%d min ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1f hours ago", d.Hours())
	}
	return fmt.Sprintf("%.1f days ago (%s)", d.Hours()/24, t.Format("2006-01-02"))
}

func loadSyncHistory(path string) SyncHistory {
	history := SyncHistory{
		Runs: make([]SyncResult, 0),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return history
	}
	json.Unmarshal(data, &history)
	if history.Runs == nil {
		history.Runs = make([]SyncResult, 0)
	}
	return history
}

func saveSyncHistory(path string, history SyncHistory) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func saveSyncResult(path string, result SyncResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// checkETag does a HEAD request and compares ETag
// Returns (newETag, needsDownload)
func checkETag(client *http.Client, url, oldETag string) (string, bool) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", true // Download on error
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", true // Download on error
	}
	defer resp.Body.Close()

	newETag := resp.Header.Get("ETag")
	if newETag == "" {
		// No ETag support, check Last-Modified instead
		newETag = resp.Header.Get("Last-Modified")
	}

	if newETag == "" {
		return "", true // No way to check, always download
	}

	return newETag, newETag != oldETag
}

func loadETags(path string) ETagStore {
	store := ETagStore{
		ETags: make(map[string]string),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return store // Return empty store if file doesn't exist
	}

	json.Unmarshal(data, &store)
	if store.ETags == nil {
		store.ETags = make(map[string]string)
	}
	return store
}

func saveETags(path string, store ETagStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ============================================================================
// Tile Command - Convert GeoJSON to PMTiles
// ============================================================================

func runTile() {
	fs := flag.NewFlagSet("tile", flag.ExitOnError)
	datasetFlag := fs.String("dataset", "", "Specific dataset to tile (empty = all)")
	force := fs.Bool("force", false, "Force regenerate even if PMTiles is newer than GeoJSON")
	tilerFlag := fs.String("tiler", "auto", "Tiler to use: auto, tippecanoe, gotiler")
	fs.Parse(os.Args[1:])

	// Determine which tiler to use
	useTippecanoe := false
	useGoTiler := false

	switch *tilerFlag {
	case "tippecanoe":
		if _, err := exec.LookPath("tippecanoe"); err != nil {
			fmt.Fprintln(os.Stderr, "Error: tippecanoe not found in PATH")
			fmt.Fprintln(os.Stderr, "Install with: brew install tippecanoe (macOS) or apt install tippecanoe (Ubuntu)")
			os.Exit(1)
		}
		useTippecanoe = true
	case "gotiler", "go":
		useGoTiler = true
	case "auto", "":
		// Auto-detect: prefer tippecanoe if available, fall back to gotiler
		if _, err := exec.LookPath("tippecanoe"); err == nil {
			useTippecanoe = true
		} else {
			useGoTiler = true
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown tiler: %s (valid: auto, tippecanoe, gotiler)\n", *tilerFlag)
		os.Exit(1)
	}

	if useTippecanoe {
		fmt.Println("Using: tippecanoe (external)")
	} else if useGoTiler {
		fmt.Println("Using: gotiler (pure Go)")
	}

	// Create tiles directory
	if err := os.MkdirAll(DirPMTiles, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tiles dir: %v\n", err)
		os.Exit(1)
	}

	// Determine which datasets to tile
	toTile := datasetOrder
	if *datasetFlag != "" {
		if _, ok := datasets[*datasetFlag]; !ok {
			fmt.Fprintf(os.Stderr, "Unknown dataset: %s\n", *datasetFlag)
			os.Exit(1)
		}
		toTile = []string{*datasetFlag}
	}

	fmt.Println("Converting GeoJSON to PMTiles")
	fmt.Println("=============================")
	fmt.Println()

	// Create gotiler instance if needed
	var goTiler *gotiler.GoTiler
	if useGoTiler {
		goTiler = gotiler.New()
	}

	tiled := 0
	skipped := 0
	for _, key := range toTile {
		ds := datasets[key]
		geoJSONPath := filepath.Join(DirGeoJSON, ds.GeoJSON)
		pmTilesPath := filepath.Join(DirPMTiles, ds.PMTiles)

		// Check if GeoJSON exists
		geoJSONInfo, err := os.Stat(geoJSONPath)
		if err != nil {
			fmt.Printf("[%s] SKIP: GeoJSON not found (%s)\n", key, ds.GeoJSON)
			skipped++
			continue
		}

		// Check if PMTiles is newer than GeoJSON (skip if up-to-date)
		if !*force {
			if pmTilesInfo, err := os.Stat(pmTilesPath); err == nil {
				if pmTilesInfo.ModTime().After(geoJSONInfo.ModTime()) {
					fmt.Printf("[%s] unchanged (PMTiles newer than GeoJSON)\n", key)
					skipped++
					continue
				}
			}
		}

		// Convert using selected tiler
		fmt.Printf("[%s] Converting %s → %s\n", key, ds.GeoJSON, ds.PMTiles)

		var tileErr error
		if useTippecanoe {
			tileErr = runTippecanoe(key, geoJSONPath, pmTilesPath, ds.Layer)
		} else if useGoTiler {
			cfg := tileConfigs[key]
			tileConfig := airspace.TileConfig{
				MinZoom: cfg.MinZoom,
				MaxZoom: cfg.MaxZoom,
				Layer:   ds.Layer,
			}
			// Use sensible defaults for gotiler
			if tileConfig.MinZoom < 0 {
				tileConfig.MinZoom = 0
			}
			if tileConfig.MaxZoom < 0 {
				tileConfig.MaxZoom = 10
			}
			tileErr = goTiler.Tile(geoJSONPath, pmTilesPath, tileConfig)
		}

		if tileErr != nil {
			fmt.Fprintf(os.Stderr, "[%s] ERROR: %v\n", key, tileErr)
			continue
		}

		// Report size
		if info, err := os.Stat(pmTilesPath); err == nil {
			fmt.Printf("[%s] ✓ %.1f MB\n", key, float64(info.Size())/(1024*1024))
		}
		tiled++
	}

	fmt.Println()
	fmt.Printf("Done: %d tiled, %d skipped\n", tiled, skipped)
}

func runTippecanoe(key, inputPath, outputPath, layer string) error {
	config := tileConfigs[key]

	args := []string{
		"-o", outputPath,
		"--layer=" + layer,
		"--force",
	}

	// Zoom settings
	if config.MinZoom >= 0 && config.MaxZoom >= 0 {
		args = append(args, fmt.Sprintf("-Z%d", config.MinZoom))
		args = append(args, fmt.Sprintf("-z%d", config.MaxZoom))
	} else {
		args = append(args, "-zg") // Auto-detect zoom
	}

	// Feature reduction
	if config.ReduceRate > 0 {
		args = append(args, fmt.Sprintf("-r%d", config.ReduceRate))
	}
	if config.DropDensest {
		args = append(args, "--drop-densest-as-needed")
	}
	if config.NoFeatureLimit {
		args = append(args, "--no-feature-limit")
	}
	if config.NoTileSizeLimit {
		args = append(args, "--no-tile-size-limit")
	}

	args = append(args, inputPath)

	cmd := exec.Command("tippecanoe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ============================================================================
// Manifest Command - Generate manifest files
// ============================================================================

// ManifestGlobal is the top-level manifest structure
type ManifestGlobal struct {
	Version int                       `json:"version"`
	Updated string                    `json:"updated"`
	Regions map[string]ManifestRegion `json:"regions"`
	Notes   map[string]string         `json:"notes,omitempty"`
}

type ManifestRegion struct {
	Name          string   `json:"name"`
	BBox          []float64 `json:"bbox"`
	TilesPath     string   `json:"tiles_path"`
	ManifestFile  string   `json:"manifest_file"`
	DefaultLayers []string `json:"default_layers"`
}

// ManifestUSA is the USA regional manifest structure
type ManifestUSA struct {
	Region  string                    `json:"region"`
	Name    string                    `json:"name"`
	Version int                       `json:"version"`
	Updated string                    `json:"updated"`
	BBox    []float64                 `json:"bbox"`
	Layers  map[string]ManifestLayer  `json:"layers"`
	Source  ManifestSource            `json:"source"`
}

type ManifestLayer struct {
	Name           string          `json:"name"`
	File           string          `json:"file"`
	PMTilesLayer   string          `json:"pmtiles_layer"`
	GeomType       string          `json:"geom_type"`       // polygon, point, line
	SizeMB         float64         `json:"size_mb"`
	Features       int             `json:"features"`
	ZoomRange      []int           `json:"zoom_range"`
	DefaultVisible bool            `json:"default_visible"`
	RenderRules    []RenderRule    `json:"render_rules"`    // Ordered paint rules
	Legend         []LegendEntry   `json:"legend,omitempty"` // For UI display
}

// RenderRule defines how to style features (evaluated in order, first match wins)
type RenderRule struct {
	FilterProp  string  `json:"filter_prop,omitempty"`  // Property to filter on (e.g., "CLASS", "TYPE_CODE")
	FilterValue string  `json:"filter_value,omitempty"` // Value to match (empty = default/fallback)
	Fill        string  `json:"fill"`
	Stroke      string  `json:"stroke,omitempty"`
	Opacity     float64 `json:"opacity,omitempty"`
	Width       float64 `json:"width,omitempty"`
	Radius      float64 `json:"radius,omitempty"` // For points
}

// LegendEntry for UI layer toggles
type LegendEntry struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

type ManifestSource struct {
	Authority   string            `json:"authority"`
	URLs        map[string]string `json:"urls"`
	UpdateCycle string            `json:"update_cycle"`
}

func runManifest() {
	fs := flag.NewFlagSet("manifest", flag.ExitOnError)
	fs.Parse(os.Args[1:])

	timestamp := time.Now().UTC().Format(time.RFC3339)

	fmt.Println("Generating Airspace Manifests")
	fmt.Println("=============================")
	fmt.Printf("Timestamp: %s\n\n", timestamp)

	// Collect metrics for each layer
	layerMetrics := make(map[string]struct {
		SizeMB   float64
		Features int
	})

	for _, key := range datasetOrder {
		ds := datasets[key]
		pmTilesPath := filepath.Join(DirPMTiles, ds.PMTiles)
		geoJSONPath := filepath.Join(DirGeoJSON, ds.GeoJSON)

		var sizeMB float64
		var features int

		// Get PMTiles size
		if info, err := os.Stat(pmTilesPath); err == nil {
			sizeMB = float64(info.Size()) / (1024 * 1024)
		}

		// Count features in GeoJSON (approximate)
		features = countGeoJSONFeatures(geoJSONPath)

		layerMetrics[key] = struct {
			SizeMB   float64
			Features int
		}{sizeMB, features}

		fmt.Printf("  [%s] %.1f MB, %d features\n", key, sizeMB, features)
	}

	// Create global manifest
	globalManifest := ManifestGlobal{
		Version: 1,
		Updated: timestamp,
		Regions: map[string]ManifestRegion{
			"usa": {
				Name:          "United States",
				BBox:          []float64{-125, 24, -66, 50},
				TilesPath:     "tiles",
				ManifestFile:  FileUSAManifest,
				DefaultLayers: []string{"boundary", "sua"},
			},
		},
		Notes: map[string]string{
			"bbox_format":    "[west, south, east, north]",
			"tiles_path":     "Relative to /airspace/ in R2",
			"future_regions": "europe, canada, australia, japan",
		},
	}

	// Create USA manifest with all layers
	usaManifest := ManifestUSA{
		Region:  "usa",
		Name:    "United States",
		Version: 1,
		Updated: timestamp,
		BBox:    []float64{-125, 24, -66, 50},
		Layers:  make(map[string]ManifestLayer),
		Source: ManifestSource{
			Authority:   "FAA",
			URLs:        map[string]string{"adds": "https://adds-faa.opendata.arcgis.com", "udds": "https://udds-faa.opendata.arcgis.com"},
			UpdateCycle: "28-day AIRAC",
		},
	}

	// Add layer definitions with full render rules
	usaManifest.Layers["boundary"] = ManifestLayer{
		Name:           "Airspace Boundary",
		File:           PMTilesBoundary,
		PMTilesLayer:   LayerBoundary,
		GeomType:       "polygon",
		SizeMB:         layerMetrics["boundary"].SizeMB,
		Features:       layerMetrics["boundary"].Features,
		ZoomRange:      []int{4, 14},
		DefaultVisible: true,
		RenderRules: []RenderRule{
			{FilterProp: "CLASS", FilterValue: "A", Fill: "#0066cc", Stroke: "#0066cc", Opacity: 0.15, Width: 1},
			{FilterProp: "CLASS", FilterValue: "C", Fill: "#cc00cc", Stroke: "#cc00cc", Opacity: 0.2, Width: 2},
			{FilterProp: "CLASS", FilterValue: "D", Fill: "#0099cc", Stroke: "#0099cc", Opacity: 0.15, Width: 2},
			{FilterProp: "CLASS", FilterValue: "E", Fill: "#00cc99", Stroke: "#00cc99", Opacity: 0.1, Width: 1},
			{FilterProp: "CLASS", FilterValue: "G", Fill: "#999999", Stroke: "#999999", Opacity: 0.05, Width: 1},
			{Fill: "#666666", Stroke: "#666666", Opacity: 0.1, Width: 1}, // fallback
		},
		Legend: []LegendEntry{
			{Label: "Class A", Color: "#0066cc"},
			{Label: "Class C", Color: "#cc00cc"},
			{Label: "Class D", Color: "#0099cc"},
			{Label: "Class E", Color: "#00cc99"},
		},
	}

	usaManifest.Layers["sua"] = ManifestLayer{
		Name:           "Special Use Airspace",
		File:           PMTilesSUA,
		PMTilesLayer:   LayerSUA,
		GeomType:       "polygon",
		SizeMB:         layerMetrics["sua"].SizeMB,
		Features:       layerMetrics["sua"].Features,
		ZoomRange:      []int{4, 14},
		DefaultVisible: true,
		RenderRules: []RenderRule{
			{FilterProp: "TYPE_CODE", FilterValue: "R", Fill: "#cc0000", Stroke: "#cc0000", Opacity: 0.3, Width: 2},
			{FilterProp: "TYPE_CODE", FilterValue: "P", Fill: "#ff0000", Stroke: "#ff0000", Opacity: 0.4, Width: 2},
			{FilterProp: "TYPE_CODE", FilterValue: "MOA", Fill: "#ff9900", Stroke: "#ff9900", Opacity: 0.2, Width: 1},
			{FilterProp: "TYPE_CODE", FilterValue: "A", Fill: "#ffcc00", Stroke: "#ffcc00", Opacity: 0.2, Width: 1},
			{FilterProp: "TYPE_CODE", FilterValue: "W", Fill: "#996600", Stroke: "#996600", Opacity: 0.15, Width: 1},
			{Fill: "#666666", Stroke: "#666666", Opacity: 0.1, Width: 1}, // fallback
		},
		Legend: []LegendEntry{
			{Label: "Restricted", Color: "#cc0000"},
			{Label: "Prohibited", Color: "#ff0000"},
			{Label: "MOA", Color: "#ff9900"},
			{Label: "Alert", Color: "#ffcc00"},
			{Label: "Warning", Color: "#996600"},
		},
	}

	usaManifest.Layers["laanc"] = ManifestLayer{
		Name:           "LAANC/UAS Facility Map",
		File:           PMTilesUAS,
		PMTilesLayer:   LayerUAS,
		GeomType:       "polygon",
		SizeMB:         layerMetrics["uas"].SizeMB,
		Features:       layerMetrics["uas"].Features,
		ZoomRange:      []int{6, 14},
		DefaultVisible: false,
		RenderRules: []RenderRule{
			{Fill: "#ffff00", Stroke: "#cc9900", Opacity: 0.5, Width: 1},
		},
		Legend: []LegendEntry{
			{Label: "LAANC Grid", Color: "#ffff00"},
		},
	}

	usaManifest.Layers["airports"] = ManifestLayer{
		Name:           "Airports",
		File:           PMTilesAirports,
		PMTilesLayer:   LayerAirports,
		GeomType:       "point",
		SizeMB:         layerMetrics["airports"].SizeMB,
		Features:       layerMetrics["airports"].Features,
		ZoomRange:      []int{0, 10},
		DefaultVisible: false,
		RenderRules: []RenderRule{
			{Fill: "#00ff00", Stroke: "#006600", Width: 1, Radius: 5},
		},
		Legend: []LegendEntry{
			{Label: "Airport", Color: "#00ff00"},
		},
	}

	usaManifest.Layers["navaids"] = ManifestLayer{
		Name:           "Navigation Aids",
		File:           PMTilesNavaids,
		PMTilesLayer:   LayerNavaids,
		GeomType:       "point",
		SizeMB:         layerMetrics["navaids"].SizeMB,
		Features:       layerMetrics["navaids"].Features,
		ZoomRange:      []int{0, 10},
		DefaultVisible: false,
		RenderRules: []RenderRule{
			{Fill: "#ff00ff", Stroke: "#660066", Width: 1, Radius: 4},
		},
		Legend: []LegendEntry{
			{Label: "VOR/NDB", Color: "#ff00ff"},
		},
	}

	// Ensure data directory exists
	if err := os.MkdirAll(DirData, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data dir: %v\n", err)
		os.Exit(1)
	}

	// Write manifests
	globalPath := filepath.Join(DirData, FileManifest)
	usaPath := filepath.Join(DirData, FileUSAManifest)

	if err := writeJSON(globalPath, globalManifest); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing global manifest: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n✓ %s\n", globalPath)

	if err := writeJSON(usaPath, usaManifest); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing USA manifest: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ %s\n", usaPath)

	// Copy to static directory for local dev
	staticGlobal := filepath.Join(DirGeoJSON, FileManifest)
	staticUSA := filepath.Join(DirGeoJSON, FileUSAManifest)

	if err := copyFile(globalPath, staticGlobal); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not copy to static: %v\n", err)
	} else {
		fmt.Printf("✓ %s (copy)\n", staticGlobal)
	}

	if err := copyFile(usaPath, staticUSA); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not copy to static: %v\n", err)
	} else {
		fmt.Printf("✓ %s (copy)\n", staticUSA)
	}

	fmt.Println("\nManifests updated.")
}

func countGeoJSONFeatures(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	// Count '"type":"Feature"' occurrences (same as shell grep)
	return strings.Count(string(data), `"type":"Feature"`)
}

func writeJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// ============================================================================
// Pipeline Command - Full idempotent pipeline: sync → tile → manifest
// ============================================================================

func runPipeline() {
	fs := flag.NewFlagSet("pipeline", flag.ExitOnError)
	force := fs.Bool("force", false, "Force all steps even if no changes")
	tilerFlag := fs.String("tiler", "auto", "Tiler to use: auto, tippecanoe, gotiler")
	fs.Parse(os.Args[1:])

	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║   FAA Airspace Pipeline (Idempotent)   ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// Step 1: Sync
	fmt.Println("▶ Step 1: Sync FAA Data")
	fmt.Println("------------------------")
	runSyncInternal(*force)

	// Check if changes were detected
	resultPath := filepath.Join(DirData, FileSyncResult)
	result := loadSyncResult(resultPath)

	if !result.HasChanges && !*force {
		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("✓ No changes detected - pipeline complete (idempotent)")
		fmt.Println("  Skipped: tile generation, manifest update")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		return
	}

	fmt.Println()
	fmt.Printf("▶ Step 2: Generate PMTiles (%d datasets changed)\n", result.Updated)
	fmt.Println("------------------------------------------------")
	runTileInternal(false, *tilerFlag) // Don't force - rely on file timestamps

	fmt.Println()
	fmt.Println("▶ Step 3: Update Manifests")
	fmt.Println("---------------------------")
	runManifestInternal()

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✓ Pipeline complete")
	fmt.Println("  Next: Run 'task r2:airspace:upload' to push to R2")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func loadSyncResult(path string) SyncResult {
	var result SyncResult
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}
	json.Unmarshal(data, &result)
	return result
}

// Internal versions that don't parse args (for pipeline use)
func runSyncInternal(force bool) {
	// Save current args and restore after
	oldArgs := os.Args
	if force {
		os.Args = []string{"airspace", "-force"}
	} else {
		os.Args = []string{"airspace"}
	}
	runSync()
	os.Args = oldArgs
}

func runTileInternal(force bool, tiler string) {
	oldArgs := os.Args
	args := []string{"airspace"}
	if force {
		args = append(args, "-force")
	}
	if tiler != "" && tiler != "auto" {
		args = append(args, "-tiler", tiler)
	}
	os.Args = args
	runTile()
	os.Args = oldArgs
}

func runManifestInternal() {
	oldArgs := os.Args
	os.Args = []string{"airspace"}
	runManifest()
	os.Args = oldArgs
}

// runCheck outputs sync result as GitHub Actions outputs.
// Usage in workflow: eval $(go run ./cmd/airspace check)
func runCheck() {
	resultPath := filepath.Join(DirData, FileSyncResult)
	data, err := os.ReadFile(resultPath)
	if err != nil {
		fmt.Println("has_changes=false")
		fmt.Println("updated_count=0")
		return
	}

	var result struct {
		HasChanges bool `json:"has_changes"`
		Updated    int  `json:"updated"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("has_changes=false")
		fmt.Println("updated_count=0")
		return
	}

	fmt.Printf("has_changes=%t\n", result.HasChanges)
	fmt.Printf("updated_count=%d\n", result.Updated)
}

// runSummary outputs GitHub Actions step summary markdown.
// Writes to GITHUB_STEP_SUMMARY if set, otherwise stdout.
func runSummary() {
	var out *os.File
	if summaryPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryPath != "" {
		f, err := os.OpenFile(summaryPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: cannot write to GITHUB_STEP_SUMMARY: %v\n", err)
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
	resultPath := filepath.Join(DirData, FileSyncResult)
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

	files, _ := filepath.Glob(filepath.Join(DirPMTiles, "*.pmtiles"))
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
