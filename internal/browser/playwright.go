// Playwright integration for browser automation.
// Provides both CLI-based and Go library-based Playwright usage.
package browser

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/playwright-community/playwright-go"
)

// BrowserEngine represents the browser engine to use.
type BrowserEngine string

const (
	// BrowserChromium uses Chromium (default, most compatible)
	BrowserChromium BrowserEngine = "chromium"
	// BrowserFirefox uses Firefox
	BrowserFirefox BrowserEngine = "firefox"
	// BrowserWebKit uses WebKit (Safari engine)
	BrowserWebKit BrowserEngine = "webkit"
)

// PlaywrightConfig holds configuration for Playwright automation.
type PlaywrightConfig struct {
	// Engine specifies which browser engine: chromium (default), firefox, webkit
	Engine BrowserEngine

	// Headless runs browser without visible window (default: true in production)
	Headless bool

	// SlowMo adds delay between actions in ms (useful for debugging)
	SlowMo float64

	// Timeout for operations in ms (default: 30000)
	Timeout float64

	// UserDataDir for persistent browser profile (empty = temp profile)
	UserDataDir string

	// Channel specifies browser channel: "chrome", "msedge", "chromium" (default)
	// Only applies to Chromium engine.
	Channel string
}

// DefaultConfig returns sensible defaults for Playwright.
func DefaultConfig() *PlaywrightConfig {
	return &PlaywrightConfig{
		Engine:   BrowserChromium,
		Headless: os.Getenv("BROWSER_HEADLESS") != "false",
		SlowMo:   0,
		Timeout:  30000,
		Channel:  "", // Use default Chromium
	}
}

// WebKitConfig returns config for WebKit/Safari engine.
func WebKitConfig() *PlaywrightConfig {
	return &PlaywrightConfig{
		Engine:   BrowserWebKit,
		Headless: os.Getenv("BROWSER_HEADLESS") != "false",
		SlowMo:   0,
		Timeout:  30000,
	}
}

// FirefoxConfig returns config for Firefox engine.
func FirefoxConfig() *PlaywrightConfig {
	return &PlaywrightConfig{
		Engine:   BrowserFirefox,
		Headless: os.Getenv("BROWSER_HEADLESS") != "false",
		SlowMo:   0,
		Timeout:  30000,
	}
}

// InteractiveConfig returns config for visible browser (debugging/demos).
func InteractiveConfig() *PlaywrightConfig {
	return &PlaywrightConfig{
		Headless: false,
		SlowMo:   100,
		Timeout:  60000,
		Channel:  "",
	}
}

// PersistentConfig returns config that reuses browser profile (keeps logins).
// The profileDir should be a directory path where browser data is stored.
// Common locations:
//   - ~/.playwright-profile (custom, recommended)
//   - ~/Library/Application Support/Google/Chrome (system Chrome - be careful!)
func PersistentConfig(profileDir string) *PlaywrightConfig {
	return &PlaywrightConfig{
		Headless:    false,
		SlowMo:      100,
		Timeout:     60000,
		UserDataDir: profileDir,
		Channel:     "chrome", // Use system Chrome for best compatibility
	}
}

// DefaultProfileDir returns the default persistent profile directory.
func DefaultProfileDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.playwright-profile"
}

// PlaywrightRunner manages Playwright browser automation.
type PlaywrightRunner struct {
	config     *PlaywrightConfig
	pw         *playwright.Playwright
	browser    playwright.Browser
	context    playwright.BrowserContext
	page       playwright.Page
}

// getBrowserType returns the appropriate browser type based on config.
func (r *PlaywrightRunner) getBrowserType() playwright.BrowserType {
	switch r.config.Engine {
	case BrowserFirefox:
		return r.pw.Firefox
	case BrowserWebKit:
		return r.pw.WebKit
	default:
		return r.pw.Chromium
	}
}

// NewRunner creates a new Playwright runner with the given config.
func NewRunner(config *PlaywrightConfig) *PlaywrightRunner {
	if config == nil {
		config = DefaultConfig()
	}
	return &PlaywrightRunner{config: config}
}

// Start initializes Playwright and launches a browser.
// If UserDataDir is set, uses persistent context (reuses cookies/logins).
func (r *PlaywrightRunner) Start() error {
	var err error

	// Install playwright if needed
	if err = playwright.Install(); err != nil {
		return fmt.Errorf("failed to install playwright: %w", err)
	}

	// Start playwright
	r.pw, err = playwright.Run()
	if err != nil {
		return fmt.Errorf("failed to start playwright: %w", err)
	}

	// If UserDataDir is set, use persistent context (keeps cookies/logins)
	if r.config.UserDataDir != "" {
		return r.startPersistent()
	}

	// Otherwise use regular ephemeral browser
	return r.startEphemeral()
}

// startPersistent launches browser with persistent profile (reuses logins).
func (r *PlaywrightRunner) startPersistent() error {
	persistentOpts := playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(r.config.Headless),
	}
	if r.config.SlowMo > 0 {
		persistentOpts.SlowMo = playwright.Float(r.config.SlowMo)
	}
	// Channel only applies to Chromium
	if r.config.Channel != "" && r.config.Engine == BrowserChromium {
		persistentOpts.Channel = playwright.String(r.config.Channel)
	}

	browserType := r.getBrowserType()
	var err error
	r.context, err = browserType.LaunchPersistentContext(r.config.UserDataDir, persistentOpts)
	if err != nil {
		r.pw.Stop()
		return fmt.Errorf("failed to launch persistent %s browser: %w", r.config.Engine, err)
	}

	// Get first page or create one
	pages := r.context.Pages()
	if len(pages) > 0 {
		r.page = pages[0]
	} else {
		r.page, err = r.context.NewPage()
		if err != nil {
			r.context.Close()
			r.pw.Stop()
			return fmt.Errorf("failed to create page: %w", err)
		}
	}

	r.page.SetDefaultTimeout(r.config.Timeout)
	return nil
}

// startEphemeral launches browser with fresh profile (no saved state).
func (r *PlaywrightRunner) startEphemeral() error {
	launchOpts := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(r.config.Headless),
	}
	if r.config.SlowMo > 0 {
		launchOpts.SlowMo = playwright.Float(r.config.SlowMo)
	}
	// Channel only applies to Chromium
	if r.config.Channel != "" && r.config.Engine == BrowserChromium {
		launchOpts.Channel = playwright.String(r.config.Channel)
	}

	browserType := r.getBrowserType()
	var err error
	r.browser, err = browserType.Launch(launchOpts)
	if err != nil {
		r.pw.Stop()
		return fmt.Errorf("failed to launch %s browser: %w", r.config.Engine, err)
	}

	r.context, err = r.browser.NewContext()
	if err != nil {
		r.browser.Close()
		r.pw.Stop()
		return fmt.Errorf("failed to create context: %w", err)
	}

	r.page, err = r.context.NewPage()
	if err != nil {
		r.context.Close()
		r.browser.Close()
		r.pw.Stop()
		return fmt.Errorf("failed to create page: %w", err)
	}

	r.page.SetDefaultTimeout(r.config.Timeout)
	return nil
}

// Stop closes the browser and cleans up Playwright resources.
func (r *PlaywrightRunner) Stop() {
	if r.page != nil {
		r.page.Close()
	}
	if r.context != nil {
		r.context.Close()
	}
	if r.browser != nil {
		r.browser.Close()
	}
	if r.pw != nil {
		r.pw.Stop()
	}
}

// Page returns the current page for direct manipulation.
func (r *PlaywrightRunner) Page() playwright.Page {
	return r.page
}

// Context returns the browser context.
func (r *PlaywrightRunner) Context() playwright.BrowserContext {
	return r.context
}

// Browser returns the browser instance.
func (r *PlaywrightRunner) Browser() playwright.Browser {
	return r.browser
}

// Navigate goes to a URL and waits for load.
func (r *PlaywrightRunner) Navigate(url string) error {
	_, err := r.page.Goto(url)
	return err
}

// Screenshot takes a screenshot and saves it to the given path.
func (r *PlaywrightRunner) Screenshot(path string) error {
	_, err := r.page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
	return err
}

// WaitForSelector waits for an element to appear.
func (r *PlaywrightRunner) WaitForSelector(selector string) error {
	_, err := r.page.WaitForSelector(selector)
	return err
}

// Click clicks on an element.
func (r *PlaywrightRunner) Click(selector string) error {
	return r.page.Click(selector)
}

// Fill types text into an input field.
func (r *PlaywrightRunner) Fill(selector, text string) error {
	return r.page.Fill(selector, text)
}

// GetText gets the text content of an element.
func (r *PlaywrightRunner) GetText(selector string) (string, error) {
	return r.page.TextContent(selector)
}

// =============================================================================
// CLI-based Playwright operations (for when Go library isn't suitable)
// =============================================================================

// RunPlaywrightCLI executes a playwright CLI command.
func RunPlaywrightCLI(args ...string) error {
	bin := FindPlaywrightBinary()
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunPlaywrightCLIWithOutput executes a playwright CLI command and returns output.
func RunPlaywrightCLIWithOutput(args ...string) (string, error) {
	bin := FindPlaywrightBinary()
	cmd := exec.Command(bin, args...)
	output, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)), err
}

// InstallPlaywrightBrowsers installs Playwright browsers via CLI.
// Installs the specified browser engine, or all supported browsers if empty.
func InstallPlaywrightBrowsers(engines ...BrowserEngine) error {
	if len(engines) == 0 {
		// Install all browsers
		return RunPlaywrightCLI("install")
	}

	for _, engine := range engines {
		if err := RunPlaywrightCLI("install", string(engine)); err != nil {
			return fmt.Errorf("failed to install %s: %w", engine, err)
		}
	}
	return nil
}

// InstallChromium installs only Chromium browser.
func InstallChromium() error {
	return InstallPlaywrightBrowsers(BrowserChromium)
}

// InstallFirefox installs only Firefox browser.
func InstallFirefox() error {
	return InstallPlaywrightBrowsers(BrowserFirefox)
}

// InstallWebKit installs only WebKit (Safari engine) browser.
func InstallWebKit() error {
	return InstallPlaywrightBrowsers(BrowserWebKit)
}

// InstallAllBrowsers installs all supported browsers.
func InstallAllBrowsers() error {
	return InstallPlaywrightBrowsers()
}

// PlaywrightVersion returns the installed Playwright version.
func PlaywrightVersion() (string, error) {
	return RunPlaywrightCLIWithOutput("--version")
}
