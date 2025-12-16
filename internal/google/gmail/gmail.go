// Package gmail provides email sending via Gmail API and browser automation.
package gmail

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Config holds gmail sender configuration
type Config struct {
	// FromAddress is the sender email (always company email)
	FromAddress string
	// Signature is appended to emails
	Signature string
	// TokenPath is path to Google OAuth tokens
	TokenPath string
}

// DefaultConfig returns the standard configuration
func DefaultConfig() *Config {
	return &Config{
		FromAddress: "gerard.webb@ubuntusoftware.net",
		Signature:   "Ubuntu Software Local AI",
		TokenPath:   "~/.google-mcp-accounts",
	}
}

// Email represents an email to send
type Email struct {
	To      string
	Subject string
	Body    string
	From    string // Set automatically from config
}

// SendResult contains the result of sending an email
type SendResult struct {
	Success   bool   `json:"success"`
	MessageID string `json:"message_id,omitempty"`
	Error     string `json:"error,omitempty"`
	Mode      string `json:"mode"` // "api" or "browser"
}

// MessageSummary is a lightweight Gmail message for listing
type MessageSummary struct {
	ID       string    `json:"id"`
	ThreadID string    `json:"thread_id,omitempty"`
	Subject  string    `json:"subject,omitempty"`
	From     string    `json:"from,omitempty"`
	Snippet  string    `json:"snippet,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Unread   bool      `json:"unread,omitempty"`
	Link     string    `json:"link,omitempty"`
}

// ListResult contains the result of a list operation
type ListResult struct {
	Success       bool              `json:"success"`
	Messages      []*MessageSummary `json:"messages,omitempty"`
	NextPageToken string            `json:"nextPageToken,omitempty"`
	Error         string            `json:"error,omitempty"`
}

// Validate checks if the email is valid
func (e *Email) Validate() error {
	if e.To == "" {
		return fmt.Errorf("to address is required")
	}
	if !isValidEmail(e.To) {
		return fmt.Errorf("invalid to address: %s", e.To)
	}
	if e.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if e.Body == "" {
		return fmt.Errorf("body is required")
	}
	return nil
}

// WithSignature returns the body with signature appended
func (e *Email) WithSignature(signature string) string {
	if signature == "" {
		return e.Body
	}
	return fmt.Sprintf("%s\n\nBest,\n%s", strings.TrimSpace(e.Body), signature)
}

// isValidEmail checks if an email address is valid
func isValidEmail(email string) bool {
	// Simple regex for email validation
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// Sender interface for different send modes
type Sender interface {
	Send(email *Email) (*SendResult, error)
	Name() string
}
