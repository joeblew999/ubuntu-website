// automation.go - Reusable browser automation framework with HIL support.
//
// This package provides composable patterns for browser automation:
// - HIL (Human-in-the-Loop): User handles auth/captcha, automation takes over
// - Actions: Reusable browser operations (click, fill, wait, etc.)
// - Sites: Site-specific automation modules (cloudflare, google, etc.)
package browser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

// =============================================================================
// Core Types
// =============================================================================

// Automation provides high-level browser automation with HIL support.
type Automation struct {
	runner  *PlaywrightRunner
	page    playwright.Page
	timeout time.Duration
	verbose bool
}

// AutomationConfig configures the automation runner.
type AutomationConfig struct {
	// Engine specifies which browser engine: chromium (default), firefox, webkit
	Engine BrowserEngine

	// Profile directory for persistent sessions (empty = ephemeral)
	Profile string

	// Headless mode (default: false for HIL)
	Headless bool

	// Timeout for operations
	Timeout time.Duration

	// Verbose logging
	Verbose bool

	// UseChrome uses system Chrome instead of bundled Chromium (Chromium engine only)
	UseChrome bool
}

// DefaultAutomationConfig returns sensible defaults for interactive automation.
// Uses Chromium by default.
func DefaultAutomationConfig() *AutomationConfig {
	return &AutomationConfig{
		Engine:   BrowserChromium,
		Headless: false,
		Timeout:  120 * time.Second,
		Verbose:  true,
	}
}

// WebKitAutomationConfig returns config for WebKit/Safari engine.
func WebKitAutomationConfig() *AutomationConfig {
	return &AutomationConfig{
		Engine:   BrowserWebKit,
		Headless: false,
		Timeout:  120 * time.Second,
		Verbose:  true,
	}
}

// FirefoxAutomationConfig returns config for Firefox engine.
func FirefoxAutomationConfig() *AutomationConfig {
	return &AutomationConfig{
		Engine:   BrowserFirefox,
		Headless: false,
		Timeout:  120 * time.Second,
		Verbose:  true,
	}
}

// NewAutomation creates a new automation instance.
func NewAutomation(config *AutomationConfig) *Automation {
	if config == nil {
		config = DefaultAutomationConfig()
	}
	return &Automation{
		timeout: config.Timeout,
		verbose: config.Verbose,
	}
}

// Start initializes the browser.
func (a *Automation) Start(config *AutomationConfig) error {
	if config == nil {
		config = DefaultAutomationConfig()
	}

	browserConfig := &PlaywrightConfig{
		Engine:      config.Engine,
		Headless:    config.Headless,
		Timeout:     float64(config.Timeout.Milliseconds()),
		UserDataDir: config.Profile,
		SlowMo:      50, // Small delay for stability
	}
	// Channel only applies to Chromium engine
	if config.UseChrome && (config.Engine == "" || config.Engine == BrowserChromium) {
		browserConfig.Channel = "chrome"
	}

	a.runner = NewRunner(browserConfig)
	if err := a.runner.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %w", err)
	}

	a.page = a.runner.Page()
	return nil
}

// Stop closes the browser.
func (a *Automation) Stop() {
	if a.runner != nil {
		a.runner.Stop()
	}
}

// Page returns the underlying Playwright page for direct access.
func (a *Automation) Page() playwright.Page {
	return a.page
}

// Runner returns the underlying PlaywrightRunner.
func (a *Automation) Runner() *PlaywrightRunner {
	return a.runner
}

// Log prints a message if verbose mode is enabled.
func (a *Automation) Log(format string, args ...interface{}) {
	if a.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// =============================================================================
// Navigation
// =============================================================================

// Navigate goes to a URL.
func (a *Automation) Navigate(url string) error {
	a.Log("Navigating to: %s", url)
	_, err := a.page.Goto(url)
	return err
}

// WaitForURL waits until the URL matches a pattern.
func (a *Automation) WaitForURL(pattern string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for URL pattern: %s", pattern)
		default:
			currentURL := a.page.URL()
			if strings.Contains(currentURL, pattern) {
				return nil
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// =============================================================================
// Element Interaction
// =============================================================================

// WaitForSelector waits for an element matching the selector.
func (a *Automation) WaitForSelector(selector string, timeout time.Duration) (playwright.Locator, error) {
	a.Log("Waiting for: %s", selector)
	locator := a.page.Locator(selector)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for selector: %s", selector)
		default:
			count, _ := locator.Count()
			if count > 0 {
				return locator, nil
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// WaitForAnySelector waits for any of the selectors to appear.
// Returns the first matching selector and its locator.
func (a *Automation) WaitForAnySelector(selectors []string, timeout time.Duration) (string, playwright.Locator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return "", nil, fmt.Errorf("timeout waiting for any of %d selectors", len(selectors))
		default:
			for _, sel := range selectors {
				locator := a.page.Locator(sel)
				count, _ := locator.Count()
				if count > 0 {
					return sel, locator, nil
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// Click clicks on an element.
func (a *Automation) Click(selector string) error {
	a.Log("Clicking: %s", selector)
	locator := a.page.Locator(selector).First()
	return locator.Click()
}

// ClickFirst clicks the first element matching any of the selectors.
func (a *Automation) ClickFirst(selectors ...string) error {
	for _, sel := range selectors {
		locator := a.page.Locator(sel).First()
		count, _ := locator.Count()
		if count > 0 {
			a.Log("Clicking: %s", sel)
			return locator.Click()
		}
	}
	return fmt.Errorf("no matching element found for any selector")
}

// Fill types text into an input.
func (a *Automation) Fill(selector, text string) error {
	a.Log("Filling: %s", selector)
	locator := a.page.Locator(selector).First()
	return locator.Fill(text)
}

// Type types text character by character (useful for triggering events).
func (a *Automation) Type(selector, text string) error {
	a.Log("Typing into: %s", selector)
	locator := a.page.Locator(selector).First()
	return locator.Type(text)
}

// SelectOption selects an option from a dropdown.
func (a *Automation) SelectOption(selector string, value string) error {
	a.Log("Selecting: %s in %s", value, selector)
	locator := a.page.Locator(selector).First()
	_, err := locator.SelectOption(playwright.SelectOptionValues{Values: &[]string{value}})
	return err
}

// GetText gets the text content of an element.
func (a *Automation) GetText(selector string) (string, error) {
	locator := a.page.Locator(selector).First()
	return locator.TextContent()
}

// IsVisible checks if an element is visible.
func (a *Automation) IsVisible(selector string) bool {
	locator := a.page.Locator(selector).First()
	visible, _ := locator.IsVisible()
	return visible
}

// Exists checks if an element exists in the DOM.
func (a *Automation) Exists(selector string) bool {
	locator := a.page.Locator(selector)
	count, _ := locator.Count()
	return count > 0
}

// =============================================================================
// HIL (Human-in-the-Loop) Support
// =============================================================================

// HILConfig configures a Human-in-the-Loop flow.
type HILConfig struct {
	// ReadySelectors - selectors that indicate the user has completed manual steps
	// (e.g., login complete, captcha passed)
	ReadySelectors []string

	// Timeout for waiting for user to complete manual steps
	Timeout time.Duration

	// Message to show while waiting
	WaitMessage string
}

// DefaultHILConfig returns default HIL configuration.
func DefaultHILConfig() *HILConfig {
	return &HILConfig{
		Timeout:     5 * time.Minute,
		WaitMessage: "Complete manual steps in the browser...",
	}
}

// WaitForHIL waits for the user to complete manual steps (login, captcha, etc).
// Returns when any of the ReadySelectors are found.
func (a *Automation) WaitForHIL(config *HILConfig) (string, error) {
	if config == nil {
		config = DefaultHILConfig()
	}

	if len(config.ReadySelectors) == 0 {
		return "", fmt.Errorf("no ReadySelectors specified")
	}

	a.Log("🔐 %s", config.WaitMessage)
	a.Log("   Waiting for: %v", config.ReadySelectors)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout waiting for HIL completion")
		default:
			for _, sel := range config.ReadySelectors {
				locator := a.page.Locator(sel)
				count, _ := locator.Count()
				if count > 0 {
					a.Log("✓ Manual steps complete (found: %s)", sel)
					return sel, nil
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// RunHILFlow runs a complete HIL automation flow:
// 1. Navigate to URL
// 2. Wait for user to complete manual steps
// 3. Execute automation actions
func (a *Automation) RunHILFlow(url string, hilConfig *HILConfig, actions func() error) error {
	// Navigate
	if err := a.Navigate(url); err != nil {
		return fmt.Errorf("navigation failed: %w", err)
	}

	// Wait for HIL completion
	if _, err := a.WaitForHIL(hilConfig); err != nil {
		return fmt.Errorf("HIL wait failed: %w", err)
	}

	// Execute automation
	a.Log("🤖 Automation taking over...")
	return actions()
}

// =============================================================================
// Utility Functions
// =============================================================================

// Screenshot takes a screenshot.
func (a *Automation) Screenshot(path string) error {
	a.Log("Taking screenshot: %s", path)
	_, err := a.page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
	return err
}

// Sleep pauses for a duration.
func (a *Automation) Sleep(d time.Duration) {
	time.Sleep(d)
}

// Reload reloads the current page.
func (a *Automation) Reload() error {
	_, err := a.page.Reload()
	return err
}

// CurrentURL returns the current page URL.
func (a *Automation) CurrentURL() string {
	return a.page.URL()
}

// Title returns the page title.
func (a *Automation) Title() (string, error) {
	return a.page.Title()
}
