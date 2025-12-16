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

// FindPlaywrightBinary locates the Playwright CLI binary.
// Searches in order:
// 1. PLAYWRIGHT_BIN environment variable
// 2. .playwright-mcp/node_modules/.bin/playwright (project-local)
// 3. node_modules/.bin/playwright (npm-local)
// 4. PATH lookup
func FindPlaywrightBinary() string {
	// 1. Environment variable override
	if bin := os.Getenv("PLAYWRIGHT_BIN"); bin != "" {
		if _, err := os.Stat(bin); err == nil {
			return bin
		}
	}

	// 2. Project-local .playwright-mcp directory
	if bin := ".playwright-mcp/node_modules/.bin/playwright"; fileExists(bin) {
		return bin
	}

	// 3. npm-local node_modules
	if bin := "node_modules/.bin/playwright"; fileExists(bin) {
		return bin
	}

	// 4. PATH lookup
	if path, err := exec.LookPath("playwright"); err == nil {
		return path
	}

	// Default fallback (will fail if not found, but gives clear error)
	return "playwright"
}

// FindPlaywrightBinaryInDir locates Playwright relative to a specific directory.
// Useful when running from different working directories.
func FindPlaywrightBinaryInDir(dir string) string {
	// 1. Environment variable override
	if bin := os.Getenv("PLAYWRIGHT_BIN"); bin != "" {
		if _, err := os.Stat(bin); err == nil {
			return bin
		}
	}

	// 2. Project-local .playwright-mcp directory
	bin := filepath.Join(dir, ".playwright-mcp/node_modules/.bin/playwright")
	if fileExists(bin) {
		return bin
	}

	// 3. npm-local node_modules
	bin = filepath.Join(dir, "node_modules/.bin/playwright")
	if fileExists(bin) {
		return bin
	}

	// 4. PATH lookup
	if path, err := exec.LookPath("playwright"); err == nil {
		return path
	}

	return "playwright"
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
		// Check for .git
		if fileExists(filepath.Join(dir, ".git")) {
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
