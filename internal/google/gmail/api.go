package gmail

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// APISender sends emails via Gmail API
type APISender struct {
	config *Config
	token  string
}

// NewAPISender creates a new API sender
func NewAPISender(config *Config) (*APISender, error) {
	token, err := googleauth.LoadAccessToken(config.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load access token: %w", err)
	}
	return &APISender{
		config: config,
		token:  token,
	}, nil
}

// Name returns the sender name
func (s *APISender) Name() string {
	return "api"
}

// Send sends an email via Gmail API
func (s *APISender) Send(email *Email) (*SendResult, error) {
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

	// Build RFC 2822 message
	body := email.WithSignature(s.config.Signature)
	msg := buildRFC2822Message(email.From, email.To, email.Subject, body)

	// Base64 URL encode
	encoded := base64.URLEncoding.EncodeToString([]byte(msg))

	// Build request body
	reqBody := map[string]string{
		"raw": encoded,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return &SendResult{
			Success: false,
			Error:   err.Error(),
			Mode:    s.Name(),
		}, err
	}

	// Make API request
	req, err := http.NewRequest("POST", "https://gmail.googleapis.com/gmail/v1/users/me/messages/send", bytes.NewBuffer(jsonBody))
	if err != nil {
		return &SendResult{
			Success: false,
			Error:   err.Error(),
			Mode:    s.Name(),
		}, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &SendResult{
			Success: false,
			Error:   err.Error(),
			Mode:    s.Name(),
		}, err
	}
	defer resp.Body.Close()

	// Read response
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &SendResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
			Mode:    s.Name(),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	// Parse response for message ID
	var result struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return &SendResult{
			Success: true,
			Mode:    s.Name(),
		}, nil
	}

	return &SendResult{
		Success:   true,
		MessageID: result.ID,
		Mode:      s.Name(),
	}, nil
}

// Check verifies the API token is valid
func (s *APISender) Check() error {
	req, err := http.NewRequest("GET", "https://gmail.googleapis.com/gmail/v1/users/me/profile", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API check failed: %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// List returns recent Gmail messages using the API
func (s *APISender) List(maxResults int, query string) (*ListResult, error) {
	if maxResults <= 0 {
		maxResults = 10
	}

	params := url.Values{}
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	if strings.TrimSpace(query) != "" {
		params.Set("q", query)
	}

	listURL := "https://gmail.googleapis.com/gmail/v1/users/me/messages?" + params.Encode()

	req, err := http.NewRequest("GET", listURL, nil)
	if err != nil {
		return &ListResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &ListResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return &ListResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(body)),
		}, fmt.Errorf("API error: %s", string(body))
	}

	var listResp struct {
		Messages []struct {
			ID       string `json:"id"`
			ThreadID string `json:"threadId"`
		} `json:"messages"`
		NextPageToken string `json:"nextPageToken"`
	}

	if err := json.Unmarshal(body, &listResp); err != nil {
		return &ListResult{Success: false, Error: err.Error()}, err
	}

	messages := make([]*MessageSummary, 0, len(listResp.Messages))
	for _, msg := range listResp.Messages {
		summary, err := s.fetchMessageMetadata(msg.ID)
		if err != nil {
			return &ListResult{Success: false, Error: err.Error()}, err
		}
		summary.ThreadID = msg.ThreadID
		summary.Link = messageLink(msg.ID)
		messages = append(messages, summary)
	}

	return &ListResult{
		Success:       true,
		Messages:      messages,
		NextPageToken: listResp.NextPageToken,
	}, nil
}

// fetchMessageMetadata retrieves metadata for a single message
func (s *APISender) fetchMessageMetadata(messageID string) (*MessageSummary, error) {
	params := url.Values{}
	params.Set("format", "metadata")
	params.Add("metadataHeaders", "Subject")
	params.Add("metadataHeaders", "From")
	params.Add("metadataHeaders", "Date")

	messageURL := fmt.Sprintf("https://gmail.googleapis.com/gmail/v1/users/me/messages/%s?%s",
		url.PathEscape(messageID), params.Encode())

	req, err := http.NewRequest("GET", messageURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var msgResp struct {
		ID       string   `json:"id"`
		ThreadID string   `json:"threadId"`
		Snippet  string   `json:"snippet"`
		LabelIDs []string `json:"labelIds"`
		Payload  struct {
			Headers []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"headers"`
		} `json:"payload"`
	}

	if err := json.Unmarshal(body, &msgResp); err != nil {
		return nil, err
	}

	headers := normalizeHeaders(msgResp.Payload.Headers)
	summary := &MessageSummary{
		ID:      msgResp.ID,
		From:    headers["from"],
		Subject: headers["subject"],
		Snippet: strings.TrimSpace(msgResp.Snippet),
		Unread:  containsLabel(msgResp.LabelIDs, "UNREAD"),
	}

	if dateStr, ok := headers["date"]; ok {
		if parsed, err := mail.ParseDate(dateStr); err == nil {
			summary.Date = parsed
		}
	}

	return summary, nil
}

// buildRFC2822Message builds a properly formatted email message
func buildRFC2822Message(from, to, subject, body string) string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	return msg.String()
}

// normalizeHeaders converts Gmail header list to a case-insensitive map
func normalizeHeaders(headers []struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		key := strings.ToLower(h.Name)
		result[key] = h.Value
	}
	return result
}

// containsLabel checks if the label list contains the target label
func containsLabel(labels []string, target string) bool {
	for _, label := range labels {
		if strings.EqualFold(label, target) {
			return true
		}
	}
	return false
}

// messageLink builds a Gmail web URL for a message ID
func messageLink(messageID string) string {
	return fmt.Sprintf("https://mail.google.com/mail/u/0/#inbox/%s", messageID)
}
