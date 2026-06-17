// Package cfemail sends transactional email via the Cloudflare Email Service
// REST API (https://api.cloudflare.com/client/v4/accounts/{id}/email/sending/send).
//
// It is a thin, dependency-free wrapper around net/http intended to replace
// SMTP2GO/Gmail as the transactional sending relay. It does NOT manage
// subscriber lists, groups, or campaigns — that remains MailerLite's job.
package cfemail

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	// apiBase is the Cloudflare API v4 base URL.
	apiBase = "https://api.cloudflare.com/client/v4"

	// DefaultFrom is the canonical sending address for the company.
	// Hardcoded as the default to prevent accidental wrong-sender mistakes,
	// mirroring the Gmail tool convention. Override with --from if needed.
	DefaultFrom = "gerard.webb@ubuntusoftware.net"

	// APITokenURL is where to create a token with Email Sending: Edit.
	APITokenURL = "https://dash.cloudflare.com/profile/api-tokens"
)

// Client talks to the Cloudflare Email Service REST API.
type Client struct {
	token      string
	accountID  string
	httpClient *http.Client
}

// NewClient creates a client with an explicit token and account ID.
func NewClient(token, accountID string) *Client {
	return &Client{
		token:      token,
		accountID:  accountID,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// NewClientFromEnv builds a client from CLOUDFLARE_API_TOKEN and CF_ACCOUNT_ID,
// matching the convention used by internal/cfanalytics.
func NewClientFromEnv() (*Client, error) {
	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		return nil, errors.New("CLOUDFLARE_API_TOKEN not set (create one with 'Email Sending: Edit' at " + APITokenURL + ")")
	}
	accountID := os.Getenv("CF_ACCOUNT_ID")
	if accountID == "" {
		return nil, errors.New("CF_ACCOUNT_ID not set")
	}
	return NewClient(token, accountID), nil
}

// Address is a recipient or sender, optionally with a display name.
type Address struct {
	Email string
	Name  string
}

// marshalable renders an Address as a JSON string or {address,name} object,
// matching the Cloudflare API which accepts either form.
func (a Address) marshalable() any {
	if a.Name == "" {
		return a.Email
	}
	return map[string]string{"address": a.Email, "name": a.Name}
}

// Message is an outbound email. At least one of HTML or Text must be set.
type Message struct {
	From    Address
	To      []Address
	CC      []Address
	BCC     []Address
	ReplyTo string
	Subject string
	HTML    string
	Text    string
	Headers map[string]string
}

// SendResult is the per-recipient delivery status returned by the API.
type SendResult struct {
	Delivered        []string `json:"delivered"`
	Queued           []string `json:"queued"`
	PermanentBounces []string `json:"permanent_bounces"`
}

// apiResponse is the standard Cloudflare API envelope.
type apiResponse struct {
	Success  bool       `json:"success"`
	Errors   []apiError `json:"errors"`
	Messages []apiError `json:"messages"`
	Result   SendResult `json:"result"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Send delivers a single message and returns the per-recipient status.
func (c *Client) Send(ctx context.Context, m Message) (*SendResult, error) {
	if m.From.Email == "" {
		m.From.Email = DefaultFrom
	}
	if len(m.To) == 0 {
		return nil, errors.New("at least one recipient (to) is required")
	}
	if m.HTML == "" && m.Text == "" {
		return nil, errors.New("message must have html or text body")
	}
	total := len(m.To) + len(m.CC) + len(m.BCC)
	if total > 50 {
		return nil, fmt.Errorf("too many recipients: %d (to+cc+bcc must be <= 50)", total)
	}

	body := map[string]any{
		"from":    m.From.marshalable(),
		"to":      addrList(m.To),
		"subject": m.Subject,
	}
	if len(m.CC) > 0 {
		body["cc"] = addrList(m.CC)
	}
	if len(m.BCC) > 0 {
		body["bcc"] = addrList(m.BCC)
	}
	if m.ReplyTo != "" {
		body["reply_to"] = m.ReplyTo
	}
	if m.HTML != "" {
		body["html"] = m.HTML
	}
	if m.Text != "" {
		body["text"] = m.Text
	}
	if len(m.Headers) > 0 {
		body["headers"] = m.Headers
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/%s/email/sending/send", apiBase, c.accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed apiResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(raw))
	}

	if !parsed.Success || len(parsed.Errors) > 0 {
		var msgs []string
		for _, e := range parsed.Errors {
			msgs = append(msgs, fmt.Sprintf("[%d] %s", e.Code, e.Message))
		}
		if len(msgs) == 0 {
			msgs = append(msgs, fmt.Sprintf("HTTP %d", resp.StatusCode))
		}
		return nil, fmt.Errorf("cloudflare email send failed: %s", strings.Join(msgs, "; "))
	}

	return &parsed.Result, nil
}

// addrList converts a slice of Address into the API's JSON representation.
func addrList(addrs []Address) []any {
	out := make([]any, len(addrs))
	for i, a := range addrs {
		out[i] = a.marshalable()
	}
	return out
}

// ParseAddress parses "Name <email@host>" or "email@host" into an Address.
func ParseAddress(s string) Address {
	s = strings.TrimSpace(s)
	if i := strings.LastIndex(s, "<"); i >= 0 && strings.HasSuffix(s, ">") {
		name := strings.TrimSpace(s[:i])
		email := strings.TrimSpace(s[i+1 : len(s)-1])
		return Address{Email: email, Name: name}
	}
	return Address{Email: s}
}

// ParseAddressList splits a comma-separated recipient string into Addresses.
func ParseAddressList(s string) []Address {
	var out []Address
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, ParseAddress(part))
		}
	}
	return out
}
