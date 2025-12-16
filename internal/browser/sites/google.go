// google.go - Google services browser automation.
//
// Provides automation for Google services:
// - Gmail settings (Send mail as, filters, etc.)
// - Google Calendar
// - Google Cloud Console
// - OAuth flows
package sites

import (
	"fmt"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// GoogleAutomation handles Google-specific automation.
type GoogleAutomation struct {
	*browser.Automation
	account string // Optional: specific Google account to use
}

// GoogleConfig configures Google automation.
type GoogleConfig struct {
	// Account email (optional, for account selection)
	Account string

	// Profile directory for persistent sessions
	Profile string

	// Timeout for operations
	Timeout time.Duration

	// Verbose logging
	Verbose bool
}

// DefaultGoogleConfig returns default configuration.
func DefaultGoogleConfig() *GoogleConfig {
	return &GoogleConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}
}

// NewGoogleAutomation creates a new Google automation instance.
func NewGoogleAutomation(config *GoogleConfig) *GoogleAutomation {
	if config == nil {
		config = DefaultGoogleConfig()
	}

	autoConfig := &browser.AutomationConfig{
		Profile: config.Profile,
		Timeout: config.Timeout,
		Verbose: config.Verbose,
	}

	return &GoogleAutomation{
		Automation: browser.NewAutomation(autoConfig),
		account:    config.Account,
	}
}

// =============================================================================
// URLs
// =============================================================================

// GmailURL returns the Gmail inbox URL.
func (g *GoogleAutomation) GmailURL() string {
	return "https://mail.google.com"
}

// GmailSettingsURL returns the Gmail settings URL (Accounts & Import tab).
func (g *GoogleAutomation) GmailSettingsURL() string {
	return "https://mail.google.com/mail/u/0/#settings/accounts"
}

// GmailComposeURL returns a Gmail compose URL with pre-filled fields.
func (g *GoogleAutomation) GmailComposeURL(to, subject, body string) string {
	return fmt.Sprintf("https://mail.google.com/mail/u/0/?view=cm&fs=1&to=%s&su=%s&body=%s",
		to, subject, body)
}

// CalendarURL returns the Google Calendar URL.
func (g *GoogleAutomation) CalendarURL() string {
	return "https://calendar.google.com"
}

// DriveURL returns the Google Drive URL.
func (g *GoogleAutomation) DriveURL() string {
	return "https://drive.google.com"
}

// CloudConsoleURL returns the Google Cloud Console URL.
func (g *GoogleAutomation) CloudConsoleURL() string {
	return "https://console.cloud.google.com"
}

// OAuthConsentURL returns the OAuth consent screen URL.
func (g *GoogleAutomation) OAuthConsentURL() string {
	return "https://console.cloud.google.com/apis/credentials/consent"
}

// CredentialsURL returns the API credentials URL.
func (g *GoogleAutomation) CredentialsURL() string {
	return "https://console.cloud.google.com/apis/credentials"
}

// =============================================================================
// HIL Configurations
// =============================================================================

// LoginHIL returns HIL config for waiting for Google login.
func (g *GoogleAutomation) LoginHIL() *browser.HILConfig {
	return &browser.HILConfig{
		ReadySelectors: []string{
			"[data-ogsr-up]",                    // Signed in indicator
			"a[href*='SignOutOptions']",         // Sign out link (means signed in)
			"[aria-label='Google Account']",     // Account avatar
			"img[data-profileimagelarge='true']", // Profile image
		},
		Timeout:     5 * time.Minute,
		WaitMessage: "Please sign in to your Google account...",
	}
}

// GmailSettingsHIL returns HIL config for Gmail settings page.
func (g *GoogleAutomation) GmailSettingsHIL() *browser.HILConfig {
	return &browser.HILConfig{
		ReadySelectors: []string{
			"text=Send mail as",
			"text=Add another email address",
			"[data-tooltip*='Send mail as']",
			"div[data-section-id='accounts']",
		},
		Timeout:     5 * time.Minute,
		WaitMessage: "Please sign in to Gmail...",
	}
}

// CalendarHIL returns HIL config for Google Calendar.
func (g *GoogleAutomation) CalendarHIL() *browser.HILConfig {
	return &browser.HILConfig{
		ReadySelectors: []string{
			"[data-view='day']",
			"[data-view='week']",
			"[data-view='month']",
			"[aria-label='Create']",
		},
		Timeout:     5 * time.Minute,
		WaitMessage: "Please sign in to Google Calendar...",
	}
}

// =============================================================================
// Gmail: Send Mail As
// =============================================================================

// SMTPConfig holds SMTP relay configuration for Gmail "Send mail as".
type SMTPConfig struct {
	Name         string // Display name (e.g., "Gerard Webb")
	Email        string // Email address to send from
	SMTPHost     string // SMTP server hostname
	SMTPPort     string // SMTP port (usually "587")
	SMTPUsername string // SMTP username
	SMTPPassword string // SMTP password
	TreatAsAlias bool   // Whether to treat as alias
}

// Common SMTP provider configurations
var (
	// SMTP2GOHost is the SMTP2GO server hostname.
	SMTP2GOHost = "mail.smtp2go.com"
	// BrevoHost is the Brevo (Sendinblue) server hostname.
	BrevoHost = "smtp-relay.brevo.com"
	// ResendHost is the Resend server hostname.
	ResendHost = "smtp.resend.com"
)

// ConfigureSendMailAs configures Gmail "Send mail as" using HIL.
// This automates the process of adding a custom email address to Gmail.
func (g *GoogleAutomation) ConfigureSendMailAs(smtp *SMTPConfig) error {
	if err := g.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	defer g.Stop()

	g.Log("🔧 Gmail 'Send mail as' Configuration")
	g.Log("=====================================")
	g.Log("  Name:     %s", smtp.Name)
	g.Log("  Email:    %s", smtp.Email)
	g.Log("  SMTP:     %s:%s", smtp.SMTPHost, smtp.SMTPPort)
	g.Log("  Username: %s", smtp.SMTPUsername)
	g.Log("")

	// Navigate to Gmail settings
	return g.RunHILFlow(g.GmailSettingsURL(), g.GmailSettingsHIL(), func() error {
		// Click "Add another email address"
		g.Log("Clicking 'Add another email address'...")
		if err := g.ClickFirst(
			"text=Add another email address",
			"span:has-text('Add another email')",
		); err != nil {
			return fmt.Errorf("could not find 'Add another email address' button: %w", err)
		}

		g.Sleep(2 * time.Second)

		// A popup window should appear - need to handle it
		// For now, provide guidance
		g.Log("")
		g.Log("📝 A popup window should appear. Please fill in:")
		g.Log("   Name: %s", smtp.Name)
		g.Log("   Email: %s", smtp.Email)
		g.Log("")
		g.Log("   Click 'Next Step'")
		g.Log("")
		g.Log("   SMTP Server: %s", smtp.SMTPHost)
		g.Log("   Port: %s", smtp.SMTPPort)
		g.Log("   Username: %s", smtp.SMTPUsername)
		g.Log("   Password: (enter your SMTP password)")
		g.Log("")
		g.Log("   Click 'Add Account'")
		g.Log("")
		g.Log("⏳ Waiting for you to complete the popup...")

		// Wait for the modal to close or success indication
		g.Sleep(30 * time.Second)

		return nil
	})
}

// OpenGmailSettings opens Gmail settings page for manual configuration.
func (g *GoogleAutomation) OpenGmailSettings() error {
	if err := g.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't Stop() - keep browser open

	g.Log("Opening Gmail Settings > Accounts...")
	g.Log("")
	g.Log("Instructions for 'Send mail as':")
	g.Log("  1. Sign in to Gmail if needed")
	g.Log("  2. Find 'Send mail as' section")
	g.Log("  3. Click 'Add another email address'")
	g.Log("  4. Enter your custom email details")
	g.Log("")

	return g.Navigate(g.GmailSettingsURL())
}

// =============================================================================
// Gmail: Compose
// =============================================================================

// ComposeEmail opens Gmail compose with pre-filled fields.
func (g *GoogleAutomation) ComposeEmail(to, subject, body string) error {
	if err := g.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't Stop() - keep browser open for user to send

	g.Log("Opening Gmail compose...")
	return g.Navigate(g.GmailComposeURL(to, subject, body))
}

// =============================================================================
// Google Calendar
// =============================================================================

// CalendarEvent represents a calendar event.
type CalendarEvent struct {
	Title       string
	Description string
	Location    string
	Start       time.Time
	End         time.Time
	Attendees   []string
}

// OpenCalendar opens Google Calendar.
func (g *GoogleAutomation) OpenCalendar() error {
	if err := g.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't Stop() - keep browser open

	g.Log("Opening Google Calendar...")
	return g.Navigate(g.CalendarURL())
}

// =============================================================================
// Google Cloud Console
// =============================================================================

// OpenCloudConsole opens Google Cloud Console.
func (g *GoogleAutomation) OpenCloudConsole() error {
	if err := g.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't Stop() - keep browser open

	g.Log("Opening Google Cloud Console...")
	return g.Navigate(g.CloudConsoleURL())
}

// CreateOAuthCredentialsGuidance opens the credentials page with guidance.
func (g *GoogleAutomation) CreateOAuthCredentialsGuidance() error {
	if err := g.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-google",
		Timeout: 10 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't Stop() - keep browser open

	g.Log("╔══════════════════════════════════════════════════════════════╗")
	g.Log("║         Google OAuth Credentials Guide                       ║")
	g.Log("╚══════════════════════════════════════════════════════════════╝")
	g.Log("")
	g.Log("Steps:")
	g.Log("  1. Sign in to Google Cloud Console")
	g.Log("  2. Select or create a project")
	g.Log("  3. Click 'Create Credentials' > 'OAuth client ID'")
	g.Log("  4. Application type: 'Desktop app'")
	g.Log("  5. Name: 'Ubuntu Software CLI' (or any name)")
	g.Log("  6. Click 'Create'")
	g.Log("  7. COPY the Client ID and Client Secret")
	g.Log("")

	return g.Navigate(g.CredentialsURL())
}
