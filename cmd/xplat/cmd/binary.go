package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

// BinaryCmd is the parent command for binary operations
var BinaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Binary management commands",
	Long: `Commands for managing binary tool installation.

Provides a centralized way to install binaries that:
- First checks if the binary already exists (PATH or install dir)
- Builds from local source if Go is available
- Downloads from GitHub releases as fallback`,
}

// BinaryInstallCmd installs a binary (build or download)
var BinaryInstallCmd = &cobra.Command{
	Use:   "install <name> <version> <repo>",
	Short: "Install a binary (build from source or download)",
	Long: `Install a binary tool, using the best available strategy:

1. Check if binary exists in PATH - skip if found
2. Check if binary exists in install dir - skip if found
3. Build from source if Go is available AND --source is provided
4. Download from GitHub release as fallback

Arguments:
  name      Binary name (e.g., "analytics")
  version   Version tag (e.g., "v0.1.0")
  repo      GitHub repo (e.g., "joeblew999/ubuntu-website")

Examples:
  xplat binary install analytics v0.1.0 joeblew999/ubuntu-website --source ./cmd/analytics
  xplat binary install sitecheck v0.1.0 joeblew999/ubuntu-website
  xplat binary install analytics v0.1.0 joeblew999/ubuntu-website --force`,
	Args: cobra.ExactArgs(3),
	RunE: runBinaryInstall,
}

var (
	binarySource string
	binaryDir    string
	binaryForce  bool
)

func init() {
	BinaryInstallCmd.Flags().StringVar(&binarySource, "source", "", "Local source path for building (e.g., ./cmd/analytics)")
	BinaryInstallCmd.Flags().StringVar(&binaryDir, "dir", "", "Install directory (default: ~/.local/bin or ~/bin on Windows)")
	BinaryInstallCmd.Flags().BoolVar(&binaryForce, "force", false, "Force reinstall even if binary exists")

	BinaryCmd.AddCommand(BinaryInstallCmd)
}

func runBinaryInstall(cmd *cobra.Command, args []string) error {
	name := args[0]
	version := args[1]
	repo := args[2]

	// Default install directory
	installDir := binaryDir
	if installDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		if runtime.GOOS == "windows" {
			installDir = filepath.Join(home, "bin")
		} else {
			installDir = filepath.Join(home, ".local", "bin")
		}
	}

	// Binary extension for Windows
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	binPath := filepath.Join(installDir, name+ext)

	// Check if binary exists (unless --force)
	if !binaryForce {
		// Check PATH
		if path, err := exec.LookPath(name + ext); err == nil {
			fmt.Printf("OK: %s found at %s\n", name, path)
			return nil
		}
		// Check install directory
		if _, err := os.Stat(binPath); err == nil {
			fmt.Printf("OK: %s found at %s\n", name, binPath)
			return nil
		}
	}

	// Create install directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Strategy 1: Build from source if Go available AND source path provided
	if binarySource != "" {
		if _, err := exec.LookPath("go"); err == nil {
			fmt.Printf("DEBUG: binarySource (raw): %q\n", binarySource)

			// Clean the path to handle mixed separators (e.g., D:\foo/bar → D:\foo\bar)
			// MUST clean BEFORE checking IsAbs, because mixed separators confuse IsAbs on Windows
			sourcePath := filepath.Clean(binarySource)
			fmt.Printf("DEBUG: after Clean: %q\n", sourcePath)
			fmt.Printf("DEBUG: IsAbs: %v\n", filepath.IsAbs(sourcePath))

			// Only join with cwd if path is relative
			// After Clean, IsAbs works correctly on Windows
			if !filepath.IsAbs(sourcePath) {
				if cwd, err := os.Getwd(); err == nil {
					fmt.Printf("DEBUG: cwd: %q\n", cwd)
					sourcePath = filepath.Join(cwd, sourcePath)
					fmt.Printf("DEBUG: after Join: %q\n", sourcePath)
				}
			}

			info, statErr := os.Stat(sourcePath)
			if statErr != nil {
				fmt.Printf("DEBUG: source path stat failed: %s (error: %v)\n", sourcePath, statErr)
			} else if !info.IsDir() {
				fmt.Printf("DEBUG: source path is not a directory: %s\n", sourcePath)
			}

			if statErr == nil && info.IsDir() {
				fmt.Printf("Building %s from source...\n", name)
				fmt.Printf("    Source: %s\n", sourcePath)
				buildCmd := exec.Command("go", "build", "-o", binPath, sourcePath)
				buildCmd.Stdout = os.Stdout
				buildCmd.Stderr = os.Stderr
				if err := buildCmd.Run(); err != nil {
					return fmt.Errorf("build failed: %w", err)
				}
				fmt.Printf("OK: %s built from source\n", name)
				fmt.Printf("    Installed to: %s\n", binPath)
				return nil
			}
		}
	}

	// Strategy 2: Download from GitHub release
	fmt.Printf("Downloading %s %s from GitHub...\n", name, version)

	// Build download URL
	// Format: https://github.com/REPO/releases/download/NAME-VERSION/NAME-OS-ARCH[.exe]
	binaryName := fmt.Sprintf("%s-%s-%s%s", name, runtime.GOOS, runtime.GOARCH, ext)
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s-%s/%s",
		repo, name, version, binaryName)

	fmt.Printf("URL: %s\n", url)

	// Download binary
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d\nRelease %s-%s may not exist yet.\nInstall Go and use --source to build from source.", resp.StatusCode, name, version)
	}

	// Create output file
	out, err := os.Create(binPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy content
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		os.Remove(binPath) // Clean up partial download
		return fmt.Errorf("download incomplete: %w", err)
	}

	// Make executable (no-op on Windows)
	if err := os.Chmod(binPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	fmt.Printf("OK: %s %s installed (%d bytes)\n", name, version, written)
	fmt.Printf("    Installed to: %s\n", binPath)

	return nil
}
