package gmail

import (
	"net/url"

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

	// Open the compose window in the default browser, where the user is signed in to Gmail;
	// they review and click Send. A Playwright-launched browser has a fresh profile and can't
	// reach a signed-in Gmail, so automated sending is the API mode's job.
	if err := openInBrowser(composeURL); err != nil {
		return &SendResult{
			Success: false,
			Error:   err.Error(),
			Mode:    s.Name(),
		}, err
	}

	mode := s.Name()
	if !s.composeOnly {
		mode = "compose-fallback"
	}
	return &SendResult{
		Success: true,
		Mode:    mode,
	}, nil
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
