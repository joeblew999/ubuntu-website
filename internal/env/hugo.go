package env

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Global Hugo server management
var (
	hugoServerCmd  *exec.Cmd
	hugoServerMux  sync.Mutex
	hugoServerPort = 1313 // Default Hugo server port
)

// GetLocalIP returns the non-loopback local IPv4 address for LAN access
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// Check the address type and if it is not a loopback
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			// Get IPv4 address
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// StartHugoServer starts a simple HTTPS server for local preview of the built site
// Uses Hugo's --tlsAuto flag for automatic certificate generation via mkcert
// Binds to 0.0.0.0 for LAN access (mobile testing)
func StartHugoServer(mockMode bool) CommandOutput {
	hugoServerMux.Lock()
	defer hugoServerMux.Unlock()

	// Stop any existing server first
	if hugoServerCmd != nil {
		stopHugoServerInternal()
	}

	// Detect LAN IP address for mobile testing
	lanIP := GetLocalIP()

	if mockMode {
		localURL := fmt.Sprintf("https://localhost:%d", hugoServerPort)
		lanURL := ""
		if lanIP != "" {
			lanURL = fmt.Sprintf("https://%s:%d", lanIP, hugoServerPort)
		}
		output := fmt.Sprintf("Starting preview server (mock mode)...\n  Local: %s\n  LAN:   %s", localURL, lanURL)
		return CommandOutput{
			Output:   output,
			Error:    nil,
			LocalURL: localURL,
			LANURL:   lanURL,
		}
	}

	// Generate HTTPS certificates with mkcert (includes LAN IP explicitly)
	// Store in temp directory and regenerate each time for simplicity
	tmpDir := os.TempDir()
	certFile := filepath.Join(tmpDir, "hugo-cert.pem")
	keyFile := filepath.Join(tmpDir, "hugo-key.pem")

	// Build mkcert arguments - explicitly include LAN IP for iOS Safari compatibility
	mkcertArgs := []string{
		"-cert-file", certFile,
		"-key-file", keyFile,
		"localhost", "127.0.0.1", "::1",
	}
	if lanIP != "" {
		mkcertArgs = append(mkcertArgs, lanIP) // Add LAN IP (e.g., 192.168.1.49)
	}

	// Run mkcert to generate certificates
	mkcertCmd := exec.Command("mkcert", mkcertArgs...)
	if err := mkcertCmd.Run(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to generate certificates with mkcert: %w (ensure mkcert is installed)", err),
		}
	}

	// Start Hugo server with generated certificates
	// Use development environment which enables relativeURLs for multi-hostname support
	hugoServerCmd = exec.Command("hugo", "server",
		"--environment", "development", // Loads config/development/config.toml with relativeURLs
		"--disableLiveReload",
		"--port", fmt.Sprintf("%d", hugoServerPort),
		"--tlsCertFile", certFile,  // Use mkcert-generated certificate
		"--tlsKeyFile", keyFile,     // Use mkcert-generated key
		"--bind", "0.0.0.0",         // Bind to all interfaces for LAN access
	)

	// Start the server in background
	if err := hugoServerCmd.Start(); err != nil {
		return CommandOutput{
			Output: "",
			Error:  fmt.Errorf("failed to start preview server: %w", err),
		}
	}

	// Build URLs for display
	localURL := fmt.Sprintf("https://localhost:%d", hugoServerPort)
	lanURL := ""
	if lanIP != "" {
		lanURL = fmt.Sprintf("https://%s:%d", lanIP, hugoServerPort)
	}

	output := fmt.Sprintf("Preview server started\n  Local: %s\n  LAN:   %s", localURL, lanURL)

	return CommandOutput{
		Output:   output,
		Error:    nil,
		LocalURL: localURL,
		LANURL:   lanURL,
	}
}

// stopHugoServerInternal stops the Hugo server without acquiring the mutex.
// This internal function should only be called when the caller already holds hugoServerMux.
func stopHugoServerInternal() CommandOutput {
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

// StopHugoServer stops the running Hugo server
func StopHugoServer() CommandOutput {
	hugoServerMux.Lock()
	defer hugoServerMux.Unlock()

	return stopHugoServerInternal()
}

// BuildHugoSite runs `hugo --gc --minify` and returns streaming output
// Also starts a local HTTPS preview server with LAN access
func BuildHugoSite(mockMode bool) CommandOutput {
	if mockMode {
		lanIP := GetLocalIP()
		localURL := fmt.Sprintf("https://localhost:%d", hugoServerPort)
		lanURL := ""
		if lanIP != "" {
			lanURL = fmt.Sprintf("https://%s:%d", lanIP, hugoServerPort)
		}
		output := fmt.Sprintf("Building Hugo site (mock mode)...\nBuild complete! (mock)\n\nStarting preview server...\nPreview server running\n  Local: %s\n  LAN:   %s", localURL, lanURL)
		return CommandOutput{
			Output:   output,
			Error:    nil,
			LocalURL: localURL,
			LANURL:   lanURL,
		}
	}

	// Build the site
	buildResult := runCommand("hugo", "--gc", "--minify")
	if buildResult.Error != nil {
		return buildResult
	}

	// Start preview server with HTTPS and LAN access
	serverResult := StartHugoServer(mockMode)
	if serverResult.Error != nil {
		// Build succeeded but server failed - return build output with warning
		return CommandOutput{
			Output:   buildResult.Output + "\n\nWarning: Failed to start preview server: " + serverResult.Error.Error(),
			Error:    nil,
			LocalURL: "",
			LANURL:   "",
		}
	}

	// Combine build and server outputs
	return CommandOutput{
		Output:   buildResult.Output + "\n\n" + serverResult.Output,
		Error:    nil,
		LocalURL: serverResult.LocalURL,
		LANURL:   serverResult.LANURL,
	}
}
