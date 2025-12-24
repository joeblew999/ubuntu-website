// Command airspace manages FAA airspace data for the BVLOS demo.
//
// Usage:
//
//	airspace download             # Download all datasets
//	airspace download -dataset uas
//	airspace sync                 # Smart sync (only download if changed)
//	airspace status               # Show data file status
//
// Configuration:
//
//	The dataset definitions below MUST stay in sync with data/airspace/datasets.json
//	which is the single source of truth for all airspace constants.
//	See also: taskfiles/Taskfile.airspace.yml, layouts/fleet/airspace-demo.html
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// Dataset configuration
type Dataset struct {
	Name     string
	Filename string
	BaseURL  string
	// For FeatureServer APIs that require pagination
	IsPaginated bool
	PageSize    int
	// ETagURL is the URL to check for ETag/Last-Modified (for paginated APIs, this is the layer URL not the query URL)
	ETagURL string
}

var datasets = map[string]Dataset{
	"uas": {
		Name:        "UAS Facility Map",
		Filename:    "faa_uas_facility_map.geojson",
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    2000,
		// ETagURL points to the layer (not /query) which returns proper ETag/Last-Modified headers
		ETagURL: "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0",
	},
	"boundary": {
		Name:     "Airspace Boundary",
		Filename: "faa_airspace_boundary.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/67885972e4e940b2aa6d74024901c561/geojson?layers=0",
	},
	"sua": {
		Name:     "Special Use Airspace",
		Filename: "faa_special_use_airspace.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/dd0d1b726e504137ab3c41b21835d05b/geojson?layers=0",
	},
	"airports": {
		Name:     "Airports",
		Filename: "faa_airports.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/e747ab91a11045e8b3f8a3efd093d3b5/geojson?layers=0",
	},
	"navaids": {
		Name:     "Navigation Aids",
		Filename: "faa_navaids.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/990e238991b44dd08af27d7b43e70b92/geojson?layers=0",
	},
	"obstacles": {
		Name:     "Obstacles",
		Filename: "faa_obstacles.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/c6a62360338e408cb1512366ad61559e/geojson?layers=0",
	},
}

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
	case "status":
		runStatus()
	case "history":
		runHistory()
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
	fmt.Println("  download    Download FAA airspace data from ArcGIS APIs")
	fmt.Println("  sync        Smart sync (only download if source changed)")
	fmt.Println("  status      Show data file status and age")
	fmt.Println("  history     Show sync history and change patterns")
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
	fmt.Println("  airspace download                       # Download all datasets")
	fmt.Println("  airspace download -dataset uas          # Download only UAS Facility Map")
	fmt.Println("  airspace download -dataset airports     # Download only Airports")
	fmt.Println("  airspace sync                           # Sync only changed datasets")
	fmt.Println("  airspace sync -force                    # Force re-download all")
	fmt.Println("  airspace status                         # Show data status")
}

// ============================================================================
// Download Command
// ============================================================================

func runDownload() {
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	outputDir := fs.String("output", "static/airspace", "Output directory for GeoJSON files")
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
		outPath := filepath.Join(*outputDir, ds.Filename)

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
			fmt.Printf("  ✓ %s (%.1f MB)\n", ds.Filename, float64(info.Size())/(1024*1024))
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
	outputDir := fs.String("output", "static/airspace", "Data directory")
	fs.Parse(os.Args[1:])

	fmt.Println("Airspace Data Status")
	fmt.Println("====================")
	fmt.Println()

	datasetOrder := []string{"uas", "boundary", "sua", "airports", "navaids", "obstacles"}
	found := 0
	for _, key := range datasetOrder {
		ds := datasets[key]
		path := filepath.Join(*outputDir, ds.Filename)
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
	dataDir := fs.String("data-dir", "data/airspace", "Data directory")
	fs.Parse(os.Args[1:])

	historyFile := filepath.Join(*dataDir, "sync-history.json")
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
	outputDir := fs.String("output", "static/airspace", "Output directory for GeoJSON files")
	dataDir := fs.String("data-dir", "data/airspace", "Data directory for ETags and history")
	force := fs.Bool("force", false, "Force re-download even if unchanged")
	timeout := fs.Duration("timeout", 5*time.Minute, "HTTP request timeout")
	skipLarge := fs.Bool("skip-large", true, "Skip datasets >100MB (obstacles)")
	fs.Parse(os.Args[1:])

	etagFile := filepath.Join(*dataDir, "etags.json")
	historyFile := filepath.Join(*dataDir, "sync-history.json")
	resultFile := filepath.Join(*dataDir, "sync-result.json")

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
		outPath := filepath.Join(*outputDir, ds.Filename)
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
			fmt.Printf("  ✓ %s (%.1f MB)\n", ds.Filename, sizeMB)
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
