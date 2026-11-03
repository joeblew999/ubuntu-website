package airspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// PrintCheckOutput prints sync result in key=value format for GitHub Actions.
func PrintCheckOutput(w io.Writer) {
	resultPath := filepath.Join(DirData, FileSyncResult)
	result := LoadSyncResult(resultPath)

	fmt.Fprintf(w, "has_changes=%t\n", result.HasChanges)
	fmt.Fprintf(w, "updated_count=%d\n", result.Updated)
}

// PrintSummary prints a markdown summary for GitHub Actions step summary.
func PrintSummary(w io.Writer) {
	fmt.Fprintln(w, "## Airspace Sync Complete")
	fmt.Fprintln(w)

	// Sync result
	resultPath := filepath.Join(DirData, FileSyncResult)
	if data, err := os.ReadFile(resultPath); err == nil {
		fmt.Fprintln(w, "### Sync Result")
		fmt.Fprintln(w, "```json")
		fmt.Fprintln(w, string(data))
		fmt.Fprintln(w, "```")
		fmt.Fprintln(w)
	}

	// PMTiles files
	fmt.Fprintln(w, "### PMTiles Files")
	fmt.Fprintln(w, "| File | Size |")
	fmt.Fprintln(w, "|------|------|")

	files, _ := filepath.Glob(filepath.Join(DirPMTiles, "*.pmtiles"))
	if len(files) == 0 {
		fmt.Fprintln(w, "| _No PMTiles found_ | - |")
	} else {
		for _, f := range files {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}
			sizeMB := float64(info.Size()) / 1024 / 1024
			fmt.Fprintf(w, "| %s | %.1f MB |\n", filepath.Base(f), sizeMB)
		}
	}
	fmt.Fprintln(w)
}

// GetSummaryWriter returns the appropriate writer for GitHub Actions summary.
// If GITHUB_STEP_SUMMARY is set, returns a file writer; otherwise returns stdout.
func GetSummaryWriter() (io.Writer, func()) {
	if summaryPath := os.Getenv("GITHUB_STEP_SUMMARY"); summaryPath != "" {
		f, err := os.OpenFile(summaryPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			return f, func() { f.Close() }
		}
	}
	return os.Stdout, func() {}
}
