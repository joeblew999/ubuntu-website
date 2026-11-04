package main

import (
	"fmt"
	"path/filepath"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// history is airspace history, the task airspace:history.
func history() {
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
