package airspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ETagStore holds ETags for each dataset to detect changes.
type ETagStore struct {
	ETags     map[string]string `json:"etags"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// SyncResult records the outcome of a sync run.
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

// DatasetSync records per-dataset sync metrics.
type DatasetSync struct {
	Status     string  `json:"status"` // "unchanged", "updated", "missing", "error"
	DurationMs int64   `json:"duration_ms,omitempty"`
	SizeBytes  int64   `json:"size_bytes,omitempty"`
	SizeMB     float64 `json:"size_mb,omitempty"`
	Features   int     `json:"features,omitempty"`
	ETag       string  `json:"etag,omitempty"`
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

// MaxHistoryRuns is the number of sync runs to keep in history.
const MaxHistoryRuns = 20

// SyncOptions configures sync behavior.
type SyncOptions struct {
	OutputDir string
	DataDir   string
	Force     bool
	SkipLarge bool
	Timeout   time.Duration
	Datasets  []string // Which datasets to sync (nil = default set)
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
	}
}

// Sync downloads FAA data with ETag-based change detection.
func Sync(opts SyncOptions) (*SyncResult, error) {
	syncStart := time.Now()

	etagFile := filepath.Join(opts.DataDir, FileSyncETags)
	historyFile := filepath.Join(opts.DataDir, FileSyncHistory)
	resultFile := filepath.Join(opts.DataDir, FileSyncResult)

	// Create directories
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output dir: %w", err)
	}
	if err := os.MkdirAll(opts.DataDir, 0755); err != nil {
		return nil, fmt.Errorf("creating data dir: %w", err)
	}

	// Load existing ETags and history
	store := LoadETags(etagFile)
	history := LoadSyncHistory(historyFile)
	client := &http.Client{Timeout: opts.Timeout}

	// Track results
	result := &SyncResult{
		Timestamp: time.Now().UTC(),
		Datasets:  make(map[string]DatasetSync),
	}
	var totalBytes int64

	for _, key := range opts.Datasets {
		dsStart := time.Now()
		ds := Datasets[key]
		outPath := filepath.Join(opts.OutputDir, ds.GeoJSON)
		dsResult := DatasetSync{}

		// Check if we need to download
		needsDownload := opts.Force
		var newETag string

		if !needsDownload {
			// Check if file exists locally
			if _, err := os.Stat(outPath); os.IsNotExist(err) {
				needsDownload = true
				dsResult.Status = "missing"
			}
		}

		if !needsDownload {
			// Check for changes using ETag or Last-Modified header
			checkURL := ds.BaseURL
			if ds.ETagURL != "" {
				checkURL = ds.ETagURL
			}
			newETag, needsDownload = CheckETag(client, checkURL, store.ETags[key])
			if needsDownload {
				dsResult.Status = "changed"
			} else {
				dsResult.Status = "unchanged"
				dsResult.ETag = store.ETags[key]
				dsResult.DurationMs = time.Since(dsStart).Milliseconds()
				result.Datasets[key] = dsResult
				result.Skipped++
				continue
			}
		}

		// Download the dataset
		var err error
		if ds.IsPaginated {
			err = DownloadPaginated(client, ds, outPath)
		} else {
			err = DownloadDirect(client, ds.BaseURL, outPath)
		}

		if err != nil {
			dsResult.Status = "error"
			dsResult.Error = err.Error()
			dsResult.DurationMs = time.Since(dsStart).Milliseconds()
			result.Datasets[key] = dsResult
			continue
		}

		// Update ETag for change tracking
		if newETag != "" {
			store.ETags[key] = newETag
			dsResult.ETag = newETag
		} else {
			// For paginated downloads, do HEAD to get current ETag
			checkURL := ds.BaseURL
			if ds.ETagURL != "" {
				checkURL = ds.ETagURL
			}
			if etag, _ := CheckETag(client, checkURL, ""); etag != "" {
				store.ETags[key] = etag
				dsResult.ETag = etag
			}
		}

		// Record metrics
		if info, err := os.Stat(outPath); err == nil {
			sizeBytes := info.Size()
			sizeMB := float64(sizeBytes) / (1024 * 1024)
			totalBytes += sizeBytes
			dsResult.SizeBytes = sizeBytes
			dsResult.SizeMB = sizeMB
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
	if err := SaveETags(etagFile, store); err != nil {
		return result, fmt.Errorf("saving etags: %w", err)
	}

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

// CheckETag does a HEAD request and compares ETag.
// Returns (newETag, needsDownload).
func CheckETag(client *http.Client, url, oldETag string) (string, bool) {
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

// LoadETags loads the ETag store from disk.
func LoadETags(path string) ETagStore {
	store := ETagStore{
		ETags: make(map[string]string),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return store
	}

	json.Unmarshal(data, &store)
	if store.ETags == nil {
		store.ETags = make(map[string]string)
	}
	return store
}

// SaveETags saves the ETag store to disk.
func SaveETags(path string, store ETagStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
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
