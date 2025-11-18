package env

import (
	"fmt"
	"os/exec"
	"sync"
)

// Global Hugo server management
var (
	hugoServerCmd  *exec.Cmd
	hugoServerMux  sync.Mutex
	hugoServerPort = 1313 // Default Hugo server port
)

// StartHugoServer starts a simple HTTP server for local preview of the built site
func StartHugoServer(mockMode bool) CommandOutput {
	hugoServerMux.Lock()
	defer hugoServerMux.Unlock()

	// Stop any existing server first
	if hugoServerCmd != nil {
		StopHugoServer()
	}

	if mockMode {
		localURL := fmt.Sprintf("http://localhost:%d", hugoServerPort)
		return CommandOutput{
			Output:   fmt.Sprintf("Starting preview server (mock mode)...\nServer running at %s", localURL),
			Error:    nil,
			LocalURL: localURL,
		}
	}

	// Use Hugo's built-in server to serve the built site
	// hugo server serves in production mode, no live reload
	localURL := fmt.Sprintf("http://localhost:%d", hugoServerPort)
	hugoServerCmd = exec.Command("hugo", "server", "-e", "production", "--disableLiveReload", "--port", fmt.Sprintf("%d", hugoServerPort))

	// Start the server in background
	if err := hugoServerCmd.Start(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to start preview server: %w", err),
		}
	}

	return CommandOutput{
		Output:   fmt.Sprintf("Preview server started at %s", localURL),
		Error:    nil,
		LocalURL: localURL,
	}
}

// StopHugoServer stops the running Hugo server
func StopHugoServer() CommandOutput {
	hugoServerMux.Lock()
	defer hugoServerMux.Unlock()

	if hugoServerCmd == nil {
		return CommandOutput{
			Output: "No Hugo server is running",
			Error:  nil,
		}
	}

	// Kill the server process
	if err := hugoServerCmd.Process.Kill(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to stop Hugo server: %w", err),
		}
	}

	hugoServerCmd = nil

	return CommandOutput{
		Output: "Hugo preview server stopped",
		Error:  nil,
	}
}

// BuildHugoSite runs `hugo --gc --minify` and returns streaming output
// Also starts a local preview server
func BuildHugoSite(mockMode bool) CommandOutput {
	if mockMode {
		localURL := fmt.Sprintf("http://localhost:%d", hugoServerPort)
		return CommandOutput{
			Output:   "Building Hugo site (mock mode)...\nBuild complete! (mock)\n\nStarting preview server...\nPreview server running at " + localURL,
			Error:    nil,
			LocalURL: localURL,
		}
	}

	// Build the site
	buildResult := runCommand("hugo", "--gc", "--minify")
	if buildResult.Error != nil {
		return buildResult
	}

	// Start preview server
	serverResult := StartHugoServer(mockMode)
	if serverResult.Error != nil {
		// Build succeeded but server failed - return build output with warning
		return CommandOutput{
			Output:   buildResult.Output + "\n\nWarning: Failed to start preview server: " + serverResult.Error.Error(),
			Error:    nil,
			LocalURL: "",
		}
	}

	// Combine build and server outputs
	return CommandOutput{
		Output:   buildResult.Output + "\n\n" + serverResult.Output,
		Error:    nil,
		LocalURL: serverResult.LocalURL,
	}
}
