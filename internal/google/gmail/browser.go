package gmail

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// BrowserSender sends emails via Playwright browser automation
type BrowserSender struct {
	config      *Config
	composeOnly bool // If true, opens compose but doesn't send
	headless    bool // If true, run browser headless (no UI) - for production use
}

// NewBrowserSender creates a new browser sender (visible mode for dev verification)
func NewBrowserSender(config *Config, composeOnly bool) *BrowserSender {
	return &BrowserSender{
		config:      config,
		composeOnly: composeOnly,
		headless:    false,
	}
}

// NewBrowserSenderWithOptions creates a browser sender with full control
func NewBrowserSenderWithOptions(config *Config, composeOnly, headless bool) *BrowserSender {
	return &BrowserSender{
		config:      config,
		composeOnly: composeOnly,
		headless:    headless,
	}
}

// SetHeadless sets headless mode (true = no browser UI, false = show browser)
func (s *BrowserSender) SetHeadless(headless bool) {
	s.headless = headless
}

// IsHeadless returns whether headless mode is enabled
func (s *BrowserSender) IsHeadless() bool {
	return s.headless
}

// Name returns the sender name
func (s *BrowserSender) Name() string {
	if s.composeOnly {
		return "compose"
	}
	if s.headless {
		return "browser-headless"
	}
	return "browser"
}

// Send opens Gmail in browser and sends/composes email
func (s *BrowserSender) Send(email *Email) (*SendResult, error) {
	// Set from address from config
	email.From = s.config.FromAddress

	// Validate
	if err := email.Validate(); err != nil {
		return &SendResult{
			Success: false,
			Error:   err.Error(),
			Mode:    s.Name(),
		}, err
	}

	// Build compose URL
	body := email.WithSignature(s.config.Signature)
	composeURL := buildGmailComposeURL(email.To, email.Subject, body)

	// Use playwright CLI if available, otherwise fall back to open command
	if err := s.openWithPlaywright(composeURL, email, body); err != nil {
		// Fallback to system open (compose only)
		if err := openInBrowser(composeURL); err != nil {
			return &SendResult{
				Success: false,
				Error:   err.Error(),
				Mode:    s.Name(),
			}, err
		}
		return &SendResult{
			Success: true,
			Mode:    "compose-fallback",
		}, nil
	}

	return &SendResult{
		Success: true,
		Mode:    s.Name(),
	}, nil
}

// openWithPlaywright uses the playwright CLI to automate Gmail
func (s *BrowserSender) openWithPlaywright(composeURL string, email *Email, body string) error {
	// Check if playwright binary exists
	playwrightBin := findPlaywrightBinary()
	if playwrightBin == "" {
		return fmt.Errorf("playwright binary not found")
	}

	// For compose-only mode, just open the URL
	if s.composeOnly {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, playwrightBin, "open", "--url", composeURL)
		return cmd.Run()
	}

	// For send mode, we need to:
	// 1. Open Gmail compose
	// 2. Wait for it to load
	// 3. Change "From" address to company email
	// 4. Click Send

	// This is complex with the CLI - for now, fall back to compose-only
	// Full automation would require the playwright Go library or MCP
	return fmt.Errorf("browser send mode requires MCP - use compose mode or API mode")
}

// buildGmailComposeURL builds the Gmail compose URL with pre-filled fields
func buildGmailComposeURL(to, subject, body string) string {
	params := url.Values{}
	params.Set("view", "cm")
	params.Set("fs", "1")
	params.Set("to", to)
	params.Set("su", subject)
	params.Set("body", body)

	return "https://mail.google.com/mail/u/0/?" + params.Encode()
}

// openInBrowser opens a URL in the default browser (uses shared browser package)
func openInBrowser(url string) error {
	return browser.OpenURL(url)
}

// findPlaywrightBinary looks for the playwright binary (uses shared browser package)
func findPlaywrightBinary() string {
	return browser.FindPlaywrightBinary()
}
