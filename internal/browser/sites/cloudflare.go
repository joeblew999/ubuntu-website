// cloudflare.go - Cloudflare-specific browser automation.
//
// Provides automation for common Cloudflare tasks:
// - Email Routing configuration
// - API token creation
// - DNS management
// - Pages deployment management
package sites

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// CloudflareAutomation handles Cloudflare-specific automation.
type CloudflareAutomation struct {
	*browser.Automation
	domain string
}

// CloudflareConfig configures Cloudflare automation.
type CloudflareConfig struct {
	// Domain to manage (e.g., "ubuntusoftware.net")
	Domain string

	// Profile directory for persistent sessions
	Profile string

	// Timeout for operations
	Timeout time.Duration

	// Verbose logging
	Verbose bool
}

// DefaultCloudflareConfig returns default configuration.
func DefaultCloudflareConfig(domain string) *CloudflareConfig {
	return &CloudflareConfig{
		Domain:  domain,
		Profile: browser.DefaultProfileDir() + "-cloudflare",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}
}

// NewCloudflareAutomation creates a new Cloudflare automation instance.
func NewCloudflareAutomation(config *CloudflareConfig) *CloudflareAutomation {
	autoConfig := &browser.AutomationConfig{
		Profile: config.Profile,
		Timeout: config.Timeout,
		Verbose: config.Verbose,
	}

	return &CloudflareAutomation{
		Automation: browser.NewAutomation(autoConfig),
		domain:     config.Domain,
	}
}

// =============================================================================
// URLs
// =============================================================================

// DashboardURL returns the main dashboard URL.
func (c *CloudflareAutomation) DashboardURL() string {
	return "https://dash.cloudflare.com"
}

// EmailRoutingURL returns the Email Routing URL for the domain.
func (c *CloudflareAutomation) EmailRoutingURL() string {
	return fmt.Sprintf("https://dash.cloudflare.com/?to=/:account/%s/email/routing/routes", c.domain)
}

// DNSURL returns the DNS management URL for the domain.
func (c *CloudflareAutomation) DNSURL() string {
	return fmt.Sprintf("https://dash.cloudflare.com/?to=/:account/%s/dns", c.domain)
}

// PagesURL returns the Pages dashboard URL.
func (c *CloudflareAutomation) PagesURL() string {
	return "https://dash.cloudflare.com/?to=/:account/pages"
}

// APITokensURL returns the API tokens management URL.
func (c *CloudflareAutomation) APITokensURL() string {
	return "https://dash.cloudflare.com/profile/api-tokens"
}

// =============================================================================
// HIL Configurations
// =============================================================================

// LoginHIL returns HIL config for waiting for Cloudflare login.
func (c *CloudflareAutomation) LoginHIL() *browser.HILConfig {
	return &browser.HILConfig{
		ReadySelectors: []string{
			"[data-testid='zone-navigation']",
			"[data-testid='account-home']",
			"a[href*='/pages']",
			"text=Overview",
		},
		Timeout:     5 * time.Minute,
		WaitMessage: "Please log in to Cloudflare (complete any captcha)...",
	}
}

// EmailRoutingHIL returns HIL config for Email Routing page.
func (c *CloudflareAutomation) EmailRoutingHIL() *browser.HILConfig {
	return &browser.HILConfig{
		ReadySelectors: []string{
			"button:has-text('Create address')",
			"text=Create address",
			"[data-testid='create-address-button']",
		},
		Timeout:     5 * time.Minute,
		WaitMessage: "Please log in to Cloudflare (complete any captcha)...",
	}
}

// =============================================================================
// Email Routing
// =============================================================================

// EmailRoutingRule represents an email forwarding rule.
type EmailRoutingRule struct {
	FromAddress string // e.g., "contact" (without @domain)
	ToAddress   string // e.g., "user@gmail.com"
}

// AddEmailRoutingRule adds an email forwarding rule using HIL.
// The user must log in and pass captcha, then automation adds the rule.
func (c *CloudflareAutomation) AddEmailRoutingRule(rule *EmailRoutingRule) error {
	// Start browser
	if err := c.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-cloudflare",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	defer c.Stop()

	// Extract custom address part (before @) if full email provided
	customAddress := rule.FromAddress
	if strings.Contains(customAddress, "@") {
		customAddress = strings.Split(customAddress, "@")[0]
	}

	c.Log("╔══════════════════════════════════════════════════════════════╗")
	c.Log("║    Cloudflare Email Routing - Add Forwarding Rule (HIL)      ║")
	c.Log("╚══════════════════════════════════════════════════════════════╝")
	c.Log("")
	c.Log("  From: %s@%s", customAddress, c.domain)
	c.Log("  To:   %s", rule.ToAddress)
	c.Log("")

	// Run HIL flow
	return c.RunHILFlow(c.EmailRoutingURL(), c.EmailRoutingHIL(), func() error {
		// Click "Create address" button
		c.Log("Clicking 'Create address'...")
		if err := c.ClickFirst(
			"button:has-text('Create address')",
			"text=Create address",
		); err != nil {
			return fmt.Errorf("could not click Create address: %w", err)
		}

		c.Sleep(1 * time.Second)

		// Fill custom address
		c.Log("Filling custom address: %s", customAddress)
		if err := c.Fill("input[placeholder*='custom'], input[name*='address'], input[type='text']", customAddress); err != nil {
			return fmt.Errorf("could not fill custom address: %w", err)
		}

		c.Sleep(500 * time.Millisecond)

		// Select destination
		c.Log("Setting destination: %s", rule.ToAddress)
		if c.Exists("select, [role='combobox']") {
			if err := c.SelectOption("select, [role='combobox']", rule.ToAddress); err != nil {
				// Try clicking dropdown and selecting
				c.Click("[role='combobox']")
				c.Sleep(300 * time.Millisecond)
				c.ClickFirst(fmt.Sprintf("text=%s", rule.ToAddress))
			}
		} else if c.Exists("input[placeholder*='destination']") {
			c.Fill("input[placeholder*='destination']", rule.ToAddress)
		}

		c.Sleep(500 * time.Millisecond)

		// Click Save
		c.Log("Clicking Save...")
		if err := c.ClickFirst(
			"button:has-text('Save')",
			"button:has-text('Create')",
			"button[type='submit']",
		); err != nil {
			return fmt.Errorf("could not click Save: %w", err)
		}

		c.Sleep(2 * time.Second)

		c.Log("")
		c.Log("✅ Email forwarding rule created!")
		c.Log("   %s@%s → %s", customAddress, c.domain, rule.ToAddress)

		return nil
	})
}

// OpenEmailRouting opens the Email Routing page for manual configuration.
func (c *CloudflareAutomation) OpenEmailRouting() error {
	if err := c.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-cloudflare",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't defer Stop() - we want browser to stay open

	c.Log("Opening Cloudflare Email Routing...")
	c.Log("  Domain: %s", c.domain)
	c.Log("")
	c.Log("Instructions:")
	c.Log("  1. Log in if needed")
	c.Log("  2. Click 'Create address' to add forwarding rule")
	c.Log("  3. Enter: Custom address → Destination email")
	c.Log("")

	return c.Navigate(c.EmailRoutingURL())
}

// =============================================================================
// API Tokens
// =============================================================================

// APITokenPermission represents a permission for an API token.
type APITokenPermission struct {
	Category   string // e.g., "Zone", "Account"
	Permission string // e.g., "Email Routing Rules"
	Access     string // e.g., "Edit", "Read"
}

// CreateAPITokenGuidance opens the API tokens page with guidance.
func (c *CloudflareAutomation) CreateAPITokenGuidance(permissions []APITokenPermission) error {
	if err := c.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-cloudflare",
		Timeout: 10 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't defer Stop() - we want browser to stay open for token creation

	c.Log("╔══════════════════════════════════════════════════════════════╗")
	c.Log("║         Cloudflare API Token Creation Guide                  ║")
	c.Log("╚══════════════════════════════════════════════════════════════╝")
	c.Log("")
	c.Log("Steps:")
	c.Log("  1. Log in to Cloudflare (browser will open)")
	c.Log("  2. Click 'Create Token'")
	c.Log("  3. Use 'Create Custom Token' at the bottom")
	c.Log("  4. Add these permissions:")
	for _, p := range permissions {
		c.Log("     - %s > %s > %s", p.Category, p.Permission, p.Access)
	}
	c.Log("  5. Set Zone Resources: Include > Specific zone > %s", c.domain)
	c.Log("  6. Click 'Continue to summary' > 'Create Token'")
	c.Log("  7. COPY THE TOKEN (shown only once!)")
	c.Log("")

	return c.Navigate(c.APITokensURL())
}

// EmailRoutingTokenPermissions returns the permissions needed for Email Routing.
func EmailRoutingTokenPermissions() []APITokenPermission {
	return []APITokenPermission{
		{Category: "Zone", Permission: "Email Routing Rules", Access: "Edit"},
		{Category: "Zone", Permission: "Zone", Access: "Read"},
		{Category: "Account", Permission: "Email Routing Addresses", Access: "Edit"},
	}
}

// =============================================================================
// Dashboard Navigation
// =============================================================================

// OpenDashboard opens the Cloudflare dashboard to a specific page.
// page can be: "", "email", "pages", "dns", or a full URL.
func (c *CloudflareAutomation) OpenDashboard(page string) error {
	if err := c.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-cloudflare",
		Timeout: 5 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't defer Stop() - we want browser to stay open

	var url string
	switch page {
	case "email":
		url = c.EmailRoutingURL()
	case "pages":
		url = c.PagesURL()
	case "dns":
		url = c.DNSURL()
	case "":
		url = c.DashboardURL()
	default:
		// Treat as URL if starts with http, otherwise as unknown page
		if strings.HasPrefix(page, "http") {
			url = page
		} else {
			url = c.DashboardURL()
		}
	}

	c.Log("Opening Cloudflare Dashboard...")
	c.Log("  URL: %s", url)
	c.Log("")

	if err := c.Navigate(url); err != nil {
		return err
	}

	c.Log("Browser open. Log in if needed.")
	c.Log("Press Ctrl+C to close when done.")

	// Block forever - user closes with Ctrl+C
	select {}
}

// =============================================================================
// API Token Setup
// =============================================================================

// SetupAPIToken guides user through creating an API token and saves it.
func (c *CloudflareAutomation) SetupAPIToken(envFile string) error {
	if err := c.Start(&browser.AutomationConfig{
		Profile: browser.DefaultProfileDir() + "-cloudflare",
		Timeout: 10 * time.Minute,
		Verbose: true,
	}); err != nil {
		return err
	}
	// Note: Don't defer Stop() - we want browser to stay open

	c.Log("╔══════════════════════════════════════════════════════════════╗")
	c.Log("║         Cloudflare API Token Setup                           ║")
	c.Log("╚══════════════════════════════════════════════════════════════╝")
	c.Log("")
	c.Log("Steps:")
	c.Log("  1. Log in to Cloudflare (browser will open)")
	c.Log("  2. Click 'Create Token'")
	c.Log("  3. Use 'Create Custom Token' at the bottom")
	c.Log("  4. Add permissions:")
	c.Log("     - Zone > Email Routing Rules > Edit")
	c.Log("     - Zone > Zone > Read")
	c.Log("     - Account > Email Routing Addresses > Edit")
	c.Log("  5. Set Zone Resources: Include > Specific zone > %s", c.domain)
	c.Log("  6. Click 'Continue to summary' > 'Create Token'")
	c.Log("  7. COPY THE TOKEN (shown only once!)")
	c.Log("")

	if err := c.Navigate(c.APITokensURL()); err != nil {
		return err
	}

	c.Log("Browser opened to API Tokens page.")
	c.Log("")
	c.Log("After creating your token, paste it below.")
	c.Log("(Press Ctrl+C to cancel)")
	c.Log("")

	// Read the token from user
	fmt.Print("Paste your API token: ")
	reader := bufio.NewReader(os.Stdin)
	token, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("could not read token: %w", err)
	}
	token = strings.TrimSpace(token)

	if token == "" {
		return fmt.Errorf("no token provided")
	}

	// Verify the token works
	c.Log("")
	c.Log("Verifying token...")

	verified, err := verifyCloudflareToken(token)
	if err != nil {
		c.Log("Warning: Could not verify token: %v", err)
	} else if verified {
		c.Log("✓ Token verified successfully!")
	}

	// Save to .env file
	c.Log("")
	c.Log("Saving token to %s...", envFile)

	if err := appendToEnvFile(envFile, "CLOUDFLARE_API_TOKEN", token); err != nil {
		return fmt.Errorf("could not save token: %w", err)
	}

	c.Log("✓ Token saved!")
	c.Log("")
	c.Log("You can now use Cloudflare API commands:")
	c.Log("  task cf:email:status")
	c.Log("  task cf:email:add FROM=x TO=y")

	return nil
}

// verifyCloudflareToken tests if a token is valid
func verifyCloudflareToken(token string) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200, nil
}

// appendToEnvFile appends or updates a key=value in .env file
func appendToEnvFile(filename, key, value string) error {
	content, err := os.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false
	newLines := make([]string, 0, len(lines)+1)

	for _, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			newLines = append(newLines, fmt.Sprintf("%s=%s", key, value))
			found = true
		} else {
			newLines = append(newLines, line)
		}
	}

	if !found {
		newLines = append(newLines, fmt.Sprintf("%s=%s", key, value))
	}

	return os.WriteFile(filename, []byte(strings.Join(newLines, "\n")), 0644)
}
