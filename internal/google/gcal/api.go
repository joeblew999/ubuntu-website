package gcal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// APIClient manages calendar via Google Calendar API
type APIClient struct {
	config *Config
	token  string
}

// NewAPIClient creates a new API client
func NewAPIClient(config *Config) (*APIClient, error) {
	token, err := googleauth.LoadAccessToken(config.TokenPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load access token: %w", err)
	}
	return &APIClient{
		config: config,
		token:  token,
	}, nil
}

// Name returns the client name
func (c *APIClient) Name() string {
	return "api"
}

// Create creates a new calendar event via API
func (c *APIClient) Create(event *Event) (*CreateResult, error) {
	if err := event.Validate(); err != nil {
		return &CreateResult{
			Success: false,
			Error:   err.Error(),
			Mode:    c.Name(),
		}, err
	}

	// Build API request body
	reqBody := map[string]interface{}{
		"summary": event.Title,
		"start": map[string]string{
			"dateTime": event.Start.Format(time.RFC3339),
			"timeZone": event.Start.Location().String(),
		},
		"end": map[string]string{
			"dateTime": event.End.Format(time.RFC3339),
			"timeZone": event.End.Location().String(),
		},
	}

	if event.Description != "" {
		reqBody["description"] = event.Description
	}
	if event.Location != "" {
		reqBody["location"] = event.Location
	}
	if len(event.Attendees) > 0 {
		attendees := make([]map[string]string, len(event.Attendees))
		for i, email := range event.Attendees {
			attendees[i] = map[string]string{"email": email}
		}
		reqBody["attendees"] = attendees
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return &CreateResult{
			Success: false,
			Error:   err.Error(),
			Mode:    c.Name(),
		}, err
	}

	// Make API request
	apiURL := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events",
		url.PathEscape(c.config.CalendarID))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return &CreateResult{
			Success: false,
			Error:   err.Error(),
			Mode:    c.Name(),
		}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &CreateResult{
			Success: false,
			Error:   err.Error(),
			Mode:    c.Name(),
		}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &CreateResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
			Mode:    c.Name(),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	// Parse response
	var result struct {
		ID       string `json:"id"`
		HTMLLink string `json:"htmlLink"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return &CreateResult{
			Success: true,
			Mode:    c.Name(),
		}, nil
	}

	return &CreateResult{
		Success: true,
		EventID: result.ID,
		Link:    result.HTMLLink,
		Mode:    c.Name(),
	}, nil
}

// List lists calendar events within a time range
func (c *APIClient) List(start, end time.Time, maxResults int) (*ListResult, error) {
	if maxResults <= 0 {
		maxResults = 10
	}

	params := url.Values{}
	params.Set("timeMin", start.Format(time.RFC3339))
	params.Set("timeMax", end.Format(time.RFC3339))
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	params.Set("singleEvents", "true")
	params.Set("orderBy", "startTime")

	apiURL := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events?%s",
		url.PathEscape(c.config.CalendarID), params.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return &ListResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &ListResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &ListResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	// Parse response
	var apiResp struct {
		Items []struct {
			ID          string `json:"id"`
			Summary     string `json:"summary"`
			Description string `json:"description"`
			Location    string `json:"location"`
			HTMLLink    string `json:"htmlLink"`
			Start       struct {
				DateTime string `json:"dateTime"`
				Date     string `json:"date"`
			} `json:"start"`
			End struct {
				DateTime string `json:"dateTime"`
				Date     string `json:"date"`
			} `json:"end"`
		} `json:"items"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return &ListResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	events := make([]*EventSummary, 0, len(apiResp.Items))
	for _, item := range apiResp.Items {
		event := &EventSummary{
			ID:          item.ID,
			Title:       item.Summary,
			Description: item.Description,
			Location:    item.Location,
			Link:        item.HTMLLink,
		}

		// Parse start time
		if item.Start.DateTime != "" {
			event.Start, _ = time.Parse(time.RFC3339, item.Start.DateTime)
		} else if item.Start.Date != "" {
			event.Start, _ = time.Parse("2006-01-02", item.Start.Date)
		}

		// Parse end time
		if item.End.DateTime != "" {
			event.End, _ = time.Parse(time.RFC3339, item.End.DateTime)
		} else if item.End.Date != "" {
			event.End, _ = time.Parse("2006-01-02", item.End.Date)
		}

		events = append(events, event)
	}

	return &ListResult{
		Success: true,
		Events:  events,
	}, nil
}

// Check verifies the API token is valid
func (c *APIClient) Check() error {
	apiURL := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s",
		url.PathEscape(c.config.CalendarID))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

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
