package airspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SyncResult records the outcome of a sync run (sync_result.json): which datasets this run downloaded, which tile,
// manifest and upload then act on (doc.go, rules 5 and 6).
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

	// Published is the manifest in R2 the run compared with (doc.go, rule 2), nil if R2 had none.
	Published *ManifestUSA `json:"-"`
}

// DatasetSync records per-dataset sync metrics.
type DatasetSync struct {
	Status     string  `json:"status"` // "unchanged" (at the published version), "updated" (downloaded), "error"
	DurationMs int64   `json:"duration_ms,omitempty"`
	SizeBytes  int64   `json:"size_bytes,omitempty"`
	SizeMB     float64 `json:"size_mb,omitempty"`
	Features   int     `json:"features,omitempty"`
	Version    string  `json:"version,omitempty"` // The FAA's version of it (doc.go, rule 1)
	Error      string  `json:"error,omitempty"`
}

// SyncHistory maintains a rolling log of sync runs.
type SyncHistory struct {
	LastChange    time.Time    `json:"last_change"`
	Runs          []SyncResult `json:"runs"`
	ChangeCount   int          `json:"change_count"`
	TotalRuns     int          `json:"total_runs"`
	AvgDuration   string       `json:"avg_duration"`
	AvgDurationMs int64        `json:"avg_duration_ms"`
}

// SyncOptions configures sync behavior.
type SyncOptions struct {
	OutputDir string
	DataDir   string
	Force     bool
	SkipLarge bool
	Timeout   time.Duration
	Datasets  []string // Which datasets to sync (nil = default set)

	PublishedURL string // The manifest in R2 that records what's published (doc.go, rule 2)
}

// DefaultSyncOptions returns sensible defaults.
func DefaultSyncOptions() SyncOptions {
	return SyncOptions{
		OutputDir: DirGeoJSON,
		DataDir:   DirData,
		Force:     false,
		SkipLarge: true,
		Timeout:   5 * time.Minute,
		Datasets:  DatasetOrder,

		PublishedURL: PublishedURL,
	}
}

// Sync downloads the datasets whose FAA version isn't the one published in R2, or every one with Force (doc.go, rules
// 1 to 4), and records the run: sync_result.json, which tile, manifest and upload read, and sync_history.json. A
// dataset it couldn't version or download is in the result as "error", and Pipeline fails on it.
func Sync(opts SyncOptions) (*SyncResult, error) {
	syncStart := time.Now()

	historyFile := filepath.Join(opts.DataDir, FileSyncHistory)
	resultFile := filepath.Join(opts.DataDir, FileSyncResult)

	// Create directories
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output dir: %w", err)
	}
	if err := os.MkdirAll(opts.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("creating data dir: %w", err)
	}

	history := LoadSyncHistory(historyFile)
	client := &http.Client{Timeout: opts.Timeout}

	// Rule 2: what R2 serves, and the versions it was built from.
	published, err := FetchPublished(client, opts.PublishedURL)
	if err != nil {
		return nil, err
	}

	// Track results
	result := &SyncResult{
		Timestamp: time.Now().UTC(),
		Datasets:  make(map[string]DatasetSync),
		Published: published,
	}
	var totalBytes int64

	for _, key := range opts.Datasets {
		dsStart := time.Now()
		ds := Datasets[key]
		outPath := filepath.Join(opts.OutputDir, ds.GeoJSON)
		failed := func(err error) {
			result.Datasets[key] = DatasetSync{Status: "error", Error: err.Error(), DurationMs: time.Since(dsStart).Milliseconds()}
		}

		// Rule 1: the FAA's version of it.
		version, err := SourceVersion(client, ds)
		if err != nil {
			failed(err) // rule 4
			continue
		}

		// Rule 3: at the published version, R2 already serves it, whether or not this machine has a file of it.
		if !opts.Force && version == published.publishedVersion(key) {
			result.Datasets[key] = DatasetSync{Status: "unchanged", Version: version, DurationMs: time.Since(dsStart).Milliseconds()}
			result.Skipped++
			continue
		}

		// Changed, never published, or forced: downloaded whole (rule 1: the FAA has no delta).
		if ds.IsPaginated {
			err = DownloadPaginated(client, ds, outPath)
		} else {
			err = DownloadDirect(client, ds.BaseURL, outPath)
		}
		if err != nil {
			failed(err) // rule 4
			continue
		}

		dsResult := DatasetSync{Status: "updated", Version: version}
		if info, err := os.Stat(outPath); err == nil {
			sizeBytes := info.Size()
			totalBytes += sizeBytes
			dsResult.SizeBytes = sizeBytes
			dsResult.SizeMB = float64(sizeBytes) / (1024 * 1024)
		}
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

	// Update and save history
	if result.HasChanges {
		history.LastChange = result.Timestamp
		history.ChangeCount++
	}
	history.TotalRuns++
	history.Runs = append([]SyncResult{*result}, history.Runs...)
	if len(history.Runs) > MaxHistoryRuns {
		history.Runs = history.Runs[:MaxHistoryRuns]
	}
	// Calculate average duration
	var totalMs int64
	for _, r := range history.Runs {
		totalMs += r.DurationMs
	}
	history.AvgDurationMs = totalMs / int64(len(history.Runs))
	history.AvgDuration = time.Duration(history.AvgDurationMs * int64(time.Millisecond)).String()

	if err := SaveSyncHistory(historyFile, history); err != nil {
		return result, fmt.Errorf("saving history: %w", err)
	}

	// Save current result (for pipeline idempotency)
	if err := SaveSyncResult(resultFile, *result); err != nil {
		return result, fmt.Errorf("saving result: %w", err)
	}

	return result, nil
}

// LoadSyncHistory loads sync history from disk.
func LoadSyncHistory(path string) SyncHistory {
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

// SaveSyncHistory saves sync history to disk.
func SaveSyncHistory(path string, history SyncHistory) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// SaveSyncResult saves the current sync result.
func SaveSyncResult(path string, result SyncResult) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// LoadSyncResult loads the last sync result.
func LoadSyncResult(path string) SyncResult {
	var result SyncResult
	data, err := os.ReadFile(path)
	if err != nil {
		return result
	}
	json.Unmarshal(data, &result)
	return result
}

// FormatTimeSince formats a duration since a time.
func FormatTimeSince(t time.Time) string {
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

// FormatAge formats a duration as a human-readable age.
func FormatAge(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%d min", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}
