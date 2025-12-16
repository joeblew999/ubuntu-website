// Gmail settings automation - configure "Send mail as" via Playwright
package gmail

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
	"github.com/playwright-community/playwright-go"
)

// DefaultGmailProfileDir returns the default persistent profile directory for Gmail.
// Uses system Chrome profile since Playwright browsers are blocked by Google.
func DefaultGmailProfileDir() string {
	home, _ := os.UserHomeDir()
	return home + "/.playwright-profile-gmail-chrome"
}

// SMTPConfig holds SMTP relay configuration for "Send mail as"
type SMTPConfig struct {
	// Display name for the "From" field
	Name string
	// Email address to send from (your custom domain email)
	Email string
	// SMTP server hostname
	SMTPHost string
	// SMTP port (typically 587 for TLS or 465 for SSL)
	SMTPPort string
	// SMTP username (often your email or API key)
	SMTPUsername string
	// SMTP password
	SMTPPassword string
	// TreatAsAlias - if true, uses "Reply from the same address the message was sent to"
	TreatAsAlias bool
}

// SMTP2GOConfig returns pre-configured SMTP settings for SMTP2GO
func SMTP2GOConfig(name, email, username, password string) *SMTPConfig {
	return &SMTPConfig{
		Name:         name,
		Email:        email,
		SMTPHost:     "mail.smtp2go.com",
		SMTPPort:     "587",
		SMTPUsername: username,
		SMTPPassword: password,
		TreatAsAlias: true,
	}
}

// BrevoConfig returns pre-configured SMTP settings for Brevo
func BrevoConfig(name, email, username, password string) *SMTPConfig {
	return &SMTPConfig{
		Name:         name,
		Email:        email,
		SMTPHost:     "smtp-relay.brevo.com",
		SMTPPort:     "587",
		SMTPUsername: username,
		SMTPPassword: password,
		TreatAsAlias: true,
	}
}

// ResendConfig returns pre-configured SMTP settings for Resend
func ResendConfig(name, email, apiKey string) *SMTPConfig {
	return &SMTPConfig{
		Name:         name,
		Email:        email,
		SMTPHost:     "smtp.resend.com",
		SMTPPort:     "587",
		SMTPUsername: "resend",
		SMTPPassword: apiKey,
		TreatAsAlias: true,
	}
}

// ParseSMTPConfig parses CLI arguments into an SMTPConfig.
// Handles pre-configured providers (smtp2go, brevo, resend) or custom SMTP hosts.
func ParseSMTPConfig(name, email, provider, user, pass string) *SMTPConfig {
	switch strings.ToLower(provider) {
	case "smtp2go":
		return SMTP2GOConfig(name, email, user, pass)
	case "brevo":
		return BrevoConfig(name, email, user, pass)
	case "resend":
		return ResendConfig(name, email, pass) // Resend uses API key as password
	default:
		// Custom SMTP host
		return &SMTPConfig{
			Name:         name,
			Email:        email,
			SMTPHost:     provider,
			SMTPPort:     "587",
			SMTPUsername: user,
			SMTPPassword: pass,
			TreatAsAlias: true,
		}
	}
}

// SettingsAutomation handles Gmail settings configuration via Playwright
type SettingsAutomation struct {
	runner       *browser.PlaywrightRunner
	config       *SMTPConfig
	timeout      time.Duration
	profileDir   string // Persistent browser profile directory
	verbose      bool
	browserEngine browser.BrowserEngine // chromium, firefox, or webkit
}

// SettingsOption configures SettingsAutomation behavior
type SettingsOption func(*SettingsAutomation)

// WithProfile sets a custom profile directory for persistent sessions
func WithProfile(dir string) SettingsOption {
	return func(s *SettingsAutomation) {
		s.profileDir = dir
	}
}

// WithTimeout sets the operation timeout
func WithTimeout(d time.Duration) SettingsOption {
	return func(s *SettingsAutomation) {
		s.timeout = d
	}
}

// WithVerbose enables verbose logging
func WithVerbose(v bool) SettingsOption {
	return func(s *SettingsAutomation) {
		s.verbose = v
	}
}

// WithBrowserEngine sets the browser engine (chromium, firefox, webkit)
// Firefox is recommended for Google login as Chromium is often blocked.
func WithBrowserEngine(engine browser.BrowserEngine) SettingsOption {
	return func(s *SettingsAutomation) {
		s.browserEngine = engine
	}
}

// NewSettingsAutomation creates a new Gmail settings automation
func NewSettingsAutomation(config *SMTPConfig, opts ...SettingsOption) *SettingsAutomation {
	s := &SettingsAutomation{
		config:        config,
		timeout:       120 * time.Second,
		profileDir:   DefaultGmailProfileDir(), // Use persistent profile by default
		verbose:       true,
		browserEngine: browser.BrowserChromium, // Use Chromium with "chrome" channel for system Chrome
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// log prints a message if verbose mode is enabled
func (s *SettingsAutomation) log(format string, args ...interface{}) {
	if s.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// ConfigureSendAs automates the Gmail "Send mail as" setup using HIL pattern.
// The browser opens, user logs in if needed, then automation takes over.
func (s *SettingsAutomation) ConfigureSendAs() error {
	// Use persistent profile to preserve login sessions
	// Use "chrome" channel to launch real system Chrome (less likely to be blocked by Google)
	browserConfig := &browser.PlaywrightConfig{
		Engine:      s.browserEngine,
		Channel:     "chrome", // Use system Chrome installation
		Headless:    false,
		SlowMo:      100,
		Timeout:     float64(s.timeout.Milliseconds()),
		UserDataDir: s.profileDir,
	}

	s.log("🌐 Starting browser with persistent profile: %s", s.profileDir)
	s.runner = browser.NewRunner(browserConfig)
	if err := s.runner.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %w", err)
	}
	defer s.runner.Stop()

	page := s.runner.Page()

	// Step 1: Navigate to Gmail Settings > Accounts
	s.log("📧 Opening Gmail Settings > Accounts...")
	if _, err := page.Goto("https://mail.google.com/mail/u/0/#settings/accounts"); err != nil {
		return fmt.Errorf("failed to open Gmail settings: %w", err)
	}

	// Step 2: HIL - Wait for user to log in if needed
	s.log("🔐 Waiting for Gmail settings page to load...")
	s.log("   If prompted, please log in to your Google account")

	// Selectors that indicate we're on the settings page (logged in)
	settingsReadySelectors := []string{
		"text=Add another email address",
		"text=Send mail as",
		"span:has-text('Add another email address')",
		"div[aria-label='Settings']",
	}

	// Wait for settings page (HIL pattern - user handles login)
	if err := s.waitForAnySelector(page, settingsReadySelectors, s.timeout); err != nil {
		page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("/tmp/gmail-hil-timeout.png"),
		})
		return fmt.Errorf("timeout waiting for Gmail settings (login may be needed) - screenshot saved to /tmp/gmail-hil-timeout.png: %w", err)
	}
	s.log("✅ Gmail settings page loaded")

	// Small delay to ensure page is fully interactive
	time.Sleep(1 * time.Second)

	// Step 3: Click "Add another email address"
	s.log("🔍 Looking for 'Add another email address' link...")

	addEmailSelectors := []string{
		"text=Add another email address",
		"span:has-text('Add another email address')",
		"a:has-text('Add another')",
	}

	var clicked bool
	for _, selector := range addEmailSelectors {
		if err := page.Click(selector, playwright.PageClickOptions{
			Timeout: playwright.Float(5000),
		}); err == nil {
			clicked = true
			s.log("✅ Clicked 'Add another email address'")
			break
		}
	}

	if !clicked {
		page.Screenshot(playwright.PageScreenshotOptions{
			Path: playwright.String("/tmp/gmail-settings-debug.png"),
		})
		return fmt.Errorf("could not find 'Add another email address' link - screenshot saved to /tmp/gmail-settings-debug.png")
	}

	// Step 4: Wait for and find popup window
	s.log("⏳ Waiting for popup window...")
	time.Sleep(2 * time.Second)

	popup, err := s.findPopup(page)
	if err != nil {
		return err
	}

	// Step 5: Fill in the form - Page 1 (Name and Email)
	if err := s.fillNameAndEmail(popup); err != nil {
		return err
	}

	// Click Next
	s.log("➡️ Clicking Next...")
	if err := s.clickNext(popup); err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	// Step 6: Fill in SMTP settings - Page 2
	if err := s.fillSMTPSettings(popup); err != nil {
		return err
	}

	// Step 7: Click Add Account / Send Verification
	s.log("📤 Submitting SMTP configuration...")
	if err := s.clickAddAccount(popup); err != nil {
		return err
	}

	// Wait for result
	time.Sleep(3 * time.Second)

	// Check for errors
	errorText, _ := popup.TextContent(".errormsg")
	if errorText != "" {
		return fmt.Errorf("Gmail returned error: %s", errorText)
	}

	s.log("")
	s.log("✅ SMTP configuration submitted!")
	s.log("📬 Gmail will send a verification email to: %s", s.config.Email)
	s.log("👉 Click the link in that email to complete setup")
	s.log("")
	s.log("Keeping browser open for 30 seconds so you can verify...")

	// Keep browser open for user to see result
	time.Sleep(30 * time.Second)

	return nil
}

// waitForAnySelector waits for any of the given selectors to appear
func (s *SettingsAutomation) waitForAnySelector(page playwright.Page, selectors []string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for selectors")
		default:
			for _, sel := range selectors {
				locator := page.Locator(sel)
				count, _ := locator.Count()
				if count > 0 {
					return nil
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// findPopup finds the Gmail "Add email" popup window
func (s *SettingsAutomation) findPopup(mainPage playwright.Page) (playwright.Page, error) {
	pages := s.runner.Context().Pages()
	var popup playwright.Page

	for _, p := range pages {
		url := p.URL()
		// Gmail popup URLs contain "view=cf" or similar
		if strings.Contains(url, "view=cf") || strings.Contains(url, "view=cm") {
			popup = p
			s.log("📝 Found popup window: %s", url)
			break
		}
	}

	if popup == nil {
		// Check for popup by excluding main page
		mainURL := mainPage.URL()
		for _, p := range pages {
			url := p.URL()
			if url != mainURL && url != "about:blank" && !strings.Contains(url, "#settings") {
				popup = p
				s.log("📝 Found popup window (by exclusion): %s", url)
				break
			}
		}
	}

	if popup == nil {
		// The form might be inline in Gmail's newer UI
		s.log("📝 No separate popup found, using main page (inline form)")
		popup = mainPage
	}

	return popup, nil
}

// fillNameAndEmail fills the first page of the form
func (s *SettingsAutomation) fillNameAndEmail(popup playwright.Page) error {
	s.log("📝 Filling name: %s", s.config.Name)

	// Try multiple selectors for name field
	nameSelectors := []string{"input[name='cfn']", "#cfn", "input[aria-label*='Name']"}
	if err := s.fillField(popup, nameSelectors, s.config.Name); err != nil {
		return fmt.Errorf("failed to fill name field: %w", err)
	}

	s.log("📝 Filling email: %s", s.config.Email)
	emailSelectors := []string{"input[name='cfa']", "#cfa", "input[aria-label*='Email']"}
	if err := s.fillField(popup, emailSelectors, s.config.Email); err != nil {
		return fmt.Errorf("failed to fill email field: %w", err)
	}

	// Handle "Treat as alias" checkbox if needed
	if s.config.TreatAsAlias {
		s.log("☑️ 'Treat as alias' should be checked by default")
	}

	return nil
}

// fillField tries multiple selectors to fill a field
func (s *SettingsAutomation) fillField(page playwright.Page, selectors []string, value string) error {
	for _, sel := range selectors {
		if err := page.Fill(sel, value, playwright.PageFillOptions{
			Timeout: playwright.Float(5000),
		}); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no matching field found")
}

// clearAndFillField clears a field first (triple-click + delete) then fills
func (s *SettingsAutomation) clearAndFillField(page playwright.Page, selector, value string) error {
	locator := page.Locator(selector).First()

	// Triple-click to select all
	if err := locator.Click(playwright.LocatorClickOptions{
		ClickCount: playwright.Int(3),
	}); err != nil {
		return err
	}

	// Type the new value (replaces selected text)
	return locator.Fill(value)
}

// fillSMTPSettings fills the SMTP configuration page
func (s *SettingsAutomation) fillSMTPSettings(popup playwright.Page) error {
	// SMTP Host - MUST clear first (Gmail pre-fills "gmail.com")
	s.log("📝 Filling SMTP Host: %s", s.config.SMTPHost)
	hostSelector := "input[name='smtpServerHostName']"

	// Clear and fill the host field
	if err := s.clearAndFillField(popup, hostSelector, s.config.SMTPHost); err != nil {
		// Try alternate selector
		if err := s.clearAndFillField(popup, "#smtpServerHostName", s.config.SMTPHost); err != nil {
			return fmt.Errorf("failed to fill SMTP host: %w", err)
		}
	}

	// SMTP Port - clear and fill
	s.log("📝 Filling SMTP Port: %s", s.config.SMTPPort)
	portSelector := "input[name='smtpServerPort']"
	if err := s.clearAndFillField(popup, portSelector, s.config.SMTPPort); err != nil {
		s.log("⚠️ Could not set port, using default")
	}

	// SMTP Username
	s.log("📝 Filling SMTP Username: %s", s.config.SMTPUsername)
	if err := popup.Fill("input[name='smtpServerUsername']", s.config.SMTPUsername); err != nil {
		return fmt.Errorf("failed to fill SMTP username: %w", err)
	}

	// SMTP Password
	s.log("📝 Filling SMTP Password: ****")
	if err := popup.Fill("input[name='smtpServerPassword']", s.config.SMTPPassword); err != nil {
		return fmt.Errorf("failed to fill SMTP password: %w", err)
	}

	// Select TLS (port 587) - usually a radio button
	s.log("🔒 Selecting TLS security...")
	tlsSelectors := []string{
		"input[value='tls']",
		"label:has-text('TLS')",
		"text=Secured connection using TLS",
	}
	for _, sel := range tlsSelectors {
		if err := popup.Click(sel, playwright.PageClickOptions{
			Timeout: playwright.Float(2000),
		}); err == nil {
			break
		}
	}

	return nil
}

// clickNext clicks the Next button on the form
func (s *SettingsAutomation) clickNext(popup playwright.Page) error {
	nextSelectors := []string{
		"text=Next Step",
		"input[type='submit']",
		"button:has-text('Next')",
	}

	for _, sel := range nextSelectors {
		if err := popup.Click(sel, playwright.PageClickOptions{
			Timeout: playwright.Float(3000),
		}); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to click Next button")
}

// clickAddAccount clicks the Add Account button
func (s *SettingsAutomation) clickAddAccount(popup playwright.Page) error {
	addSelectors := []string{
		"text=Add Account",
		"input[type='submit']",
		"button:has-text('Add')",
	}

	for _, sel := range addSelectors {
		if err := popup.Click(sel, playwright.PageClickOptions{
			Timeout: playwright.Float(3000),
		}); err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to click Add Account button")
}

// Logout opens the Google account logout page to switch accounts
func (s *SettingsAutomation) Logout() error {
	browserConfig := &browser.PlaywrightConfig{
		Engine:      s.browserEngine,
		Channel:     "chrome", // Use system Chrome
		Headless:    false,
		SlowMo:      100,
		Timeout:     float64(s.timeout.Milliseconds()),
		UserDataDir: s.profileDir,
	}

	s.log("🌐 Starting browser with persistent profile: %s", s.profileDir)
	s.runner = browser.NewRunner(browserConfig)
	if err := s.runner.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %w", err)
	}
	// Don't defer Stop() - keep browser open for user

	s.log("🚪 Opening Google Account logout page...")
	if _, err := s.runner.Page().Goto("https://accounts.google.com/Logout"); err != nil {
		return fmt.Errorf("failed to open logout page: %w", err)
	}

	s.log("✅ Logged out of Google accounts")
	s.log("👉 You can now log in with a different account")
	s.log("")
	s.log("Browser will stay open. Press Ctrl+C to close.")
	select {}
}

// SwitchAccount opens the Google account chooser to add/switch accounts
func (s *SettingsAutomation) SwitchAccount() error {
	browserConfig := &browser.PlaywrightConfig{
		Engine:      s.browserEngine,
		Channel:     "chrome", // Use system Chrome
		Headless:    false,
		SlowMo:      100,
		Timeout:     float64(s.timeout.Milliseconds()),
		UserDataDir: s.profileDir,
	}

	s.log("🌐 Starting browser with persistent profile: %s", s.profileDir)
	s.runner = browser.NewRunner(browserConfig)
	if err := s.runner.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %w", err)
	}
	// Don't defer Stop() - keep browser open for user

	s.log("👥 Opening Google Account chooser...")
	// This URL shows all logged-in accounts and option to add another
	if _, err := s.runner.Page().Goto("https://accounts.google.com/AccountChooser"); err != nil {
		return fmt.Errorf("failed to open account chooser: %w", err)
	}

	s.log("✅ Account chooser opened")
	s.log("👉 Select an existing account or click 'Use another account' to add one")
	s.log("")
	s.log("Browser will stay open. Press Ctrl+C to close.")
	select {}
}

// OpenSettingsPage just opens the Gmail settings page for manual configuration
func (s *SettingsAutomation) OpenSettingsPage() error {
	// Use persistent profile to preserve login sessions
	// Use system Chrome since Playwright browsers are blocked by Google
	browserConfig := &browser.PlaywrightConfig{
		Engine:      s.browserEngine,
		Channel:     "chrome", // Use system Chrome
		Headless:    false,
		SlowMo:      100,
		Timeout:     float64(s.timeout.Milliseconds()),
		UserDataDir: s.profileDir,
	}

	s.log("🌐 Starting browser with persistent profile: %s", s.profileDir)
	s.runner = browser.NewRunner(browserConfig)
	if err := s.runner.Start(); err != nil {
		return fmt.Errorf("failed to start browser: %w", err)
	}
	// Don't defer Stop() - keep browser open

	s.log("📧 Opening Gmail Settings > Accounts...")
	if _, err := s.runner.Page().Goto("https://mail.google.com/mail/u/0/#settings/accounts"); err != nil {
		return fmt.Errorf("failed to open Gmail settings: %w", err)
	}

	s.log("✅ Gmail settings page opened")
	s.log("👉 Click 'Add another email address' to add your custom domain")
	s.log("")
	s.log("SMTP Settings for SMTP2GO:")
	s.log("  Host: mail.smtp2go.com")
	s.log("  Port: 587")
	s.log("  Username: (from SMTP2GO dashboard)")
	s.log("  Password: (from SMTP2GO dashboard)")
	s.log("  Security: TLS")

	// Keep browser open indefinitely
	s.log("")
	s.log("Browser will stay open. Press Ctrl+C to close.")
	select {}
}
