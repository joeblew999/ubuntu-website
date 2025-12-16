// commands.go - High-level browser commands for CLI use
//
// These functions provide simple wrappers around browser automation
// that can be called from CLI tools or other packages.
package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CLIOAuthConfig configures an OAuth flow for CLI use.
// Unlike OAuthConfig which is for full token exchange,
// this just opens a URL and waits for a callback code.
type CLIOAuthConfig struct {
	URL     string
	Port    int
	Timeout time.Duration
	Browser *PlaywrightConfig
}

// CLIOAuthResult contains the result of a CLI OAuth flow
type CLIOAuthResult struct {
	Code  string            `json:"code,omitempty"`
	Token string            `json:"token,omitempty"`
	Error string            `json:"error,omitempty"`
	Query map[string]string `json:"query,omitempty"`
}

// RunCLIOAuthFlow runs an OAuth flow with callback server for CLI use.
// This opens a browser to the given URL and waits for a callback with code.
func RunCLIOAuthFlow(config *CLIOAuthConfig) (*CLIOAuthResult, error) {
	fmt.Println("Starting OAuth flow...")
	fmt.Printf("  URL: %s\n", config.URL)
	fmt.Printf("  Port: %d\n", config.Port)
	fmt.Printf("  Timeout: %s\n", config.Timeout)
	fmt.Println()

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Start callback server
	server := NewCallbackServer(config.Port)
	server.SuccessHTML = `<!DOCTYPE html>
<html>
<head><title>Success</title></head>
<body style="font-family: -apple-system, sans-serif; padding: 40px; text-align: center;">
<h1 style="color: #22c55e;">Authentication Complete</h1>
<p>You can close this window and return to the terminal.</p>
</body>
</html>`

	if err := server.Start(); err != nil {
		return nil, fmt.Errorf("could not start callback server: %w", err)
	}
	defer server.Stop(ctx)

	// Start browser
	runner := NewRunner(config.Browser)
	if err := runner.Start(); err != nil {
		return nil, fmt.Errorf("could not start browser: %w", err)
	}
	defer runner.Stop()

	fmt.Println("Opening browser...")
	if err := runner.Navigate(config.URL); err != nil {
		return nil, fmt.Errorf("could not navigate: %w", err)
	}

	fmt.Println("Waiting for callback...")

	code, err := server.WaitForCode(ctx)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout waiting for callback")
		}
		return nil, err
	}

	result := &CLIOAuthResult{
		Code:  code,
		Query: map[string]string{"code": code},
	}

	// Output result as JSON
	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))

	return result, nil
}

// OpenURLInPlaywright opens a URL in a Playwright browser and waits.
// Use this when you need the Playwright browser (for automation, debugging).
// For opening in the default system browser, use OpenURL instead.
func OpenURLInPlaywright(url string, timeout time.Duration, config *PlaywrightConfig) error {
	fmt.Printf("Opening %s...\n", url)

	runner := NewRunner(config)
	if err := runner.Start(); err != nil {
		return fmt.Errorf("could not start browser: %w", err)
	}
	defer runner.Stop()

	if err := runner.Navigate(url); err != nil {
		return fmt.Errorf("could not navigate: %w", err)
	}

	fmt.Println("Browser open. Press Ctrl+C to close.")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	<-ctx.Done()

	return nil
}

// TakeScreenshot takes a screenshot of a URL
func TakeScreenshot(url, filename string, config *PlaywrightConfig) error {
	fmt.Printf("Taking screenshot of %s...\n", url)

	// Force headless for screenshots
	config.Headless = true

	runner := NewRunner(config)
	if err := runner.Start(); err != nil {
		return fmt.Errorf("could not start browser: %w", err)
	}
	defer runner.Stop()

	if err := runner.Navigate(url); err != nil {
		return fmt.Errorf("could not navigate: %w", err)
	}

	// Wait for page to load
	time.Sleep(2 * time.Second)

	// Ensure filename has extension
	if !strings.HasSuffix(filename, ".png") && !strings.HasSuffix(filename, ".jpg") {
		filename = filename + ".png"
	}

	if err := runner.Screenshot(filename); err != nil {
		return fmt.Errorf("could not take screenshot: %w", err)
	}

	fmt.Printf("Screenshot saved to %s\n", filename)
	return nil
}

// InstallBrowsers installs the specified browser engine(s)
func InstallBrowsers(engine BrowserEngine) error {
	switch engine {
	case BrowserFirefox:
		fmt.Println("Installing Firefox browser...")
		return InstallFirefox()
	case BrowserWebKit:
		fmt.Println("Installing WebKit (Safari engine) browser...")
		return InstallWebKit()
	case BrowserChromium:
		fmt.Println("Installing Chromium browser...")
		return InstallChromium()
	default:
		fmt.Println("Installing all Playwright browsers...")
		return InstallAllBrowsers()
	}
}
