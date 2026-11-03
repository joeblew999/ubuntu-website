// Package browser provides shared browser automation utilities.
// Used by gmail, calendar, playwright CLI, and google-auth.
package browser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// OpenURL opens a URL in the default browser.
// Cross-platform: uses 'open' on macOS, 'xdg-open' on Linux, 'start' on Windows.
func OpenURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// OpenURLAndWait opens a URL and waits for the browser process to exit.
// Useful when you need to know when the user closes the browser.
func OpenURLAndWait(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", "-W", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "/wait", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// GetProjectRoot attempts to find the project root directory.
// Walks up from current directory looking for go.mod or .git.
func GetProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check for go.mod
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}
		// Check for .git (a directory, or a file in worktrees/submodules)
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding markers
			return "", fmt.Errorf("could not find project root (no go.mod or .git)")
		}
		dir = parent
	}
}

// fileExists checks if a file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
