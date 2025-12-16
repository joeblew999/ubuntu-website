// Package gcal provides Google Calendar management via API and browser automation.
package gcal

import (
	"fmt"
	"time"
)

// Config holds calendar configuration
type Config struct {
	// CalendarID is the calendar to use (default: primary)
	CalendarID string
	// TokenPath is path to Google OAuth tokens
	TokenPath string
	// DefaultDuration for events without end time
	DefaultDuration time.Duration
}

// DefaultConfig returns the standard configuration
func DefaultConfig() *Config {
	return &Config{
		CalendarID:      "primary",
		TokenPath:       "~/.google-mcp-accounts",
		DefaultDuration: 1 * time.Hour,
	}
}

// Event represents a calendar event
type Event struct {
	Title       string
	Description string
	Location    string
	Start       time.Time
	End         time.Time
	Attendees   []string
}

// Validate checks if the event is valid
func (e *Event) Validate() error {
	if e.Title == "" {
		return fmt.Errorf("title is required")
	}
	if e.Start.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if e.End.IsZero() {
		return fmt.Errorf("end time is required")
	}
	if e.End.Before(e.Start) {
		return fmt.Errorf("end time must be after start time")
	}
	return nil
}

// CreateResult contains the result of creating an event
type CreateResult struct {
	Success bool   `json:"success"`
	EventID string `json:"event_id,omitempty"`
	Link    string `json:"link,omitempty"`
	Error   string `json:"error,omitempty"`
	Mode    string `json:"mode"` // "api" or "browser"
}

// ListResult contains the result of listing events
type ListResult struct {
	Success bool           `json:"success"`
	Events  []*EventSummary `json:"events,omitempty"`
	Error   string         `json:"error,omitempty"`
}

// EventSummary is a simplified event for listing
type EventSummary struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Location    string    `json:"location,omitempty"`
	Description string    `json:"description,omitempty"`
	Link        string    `json:"link,omitempty"`
}

// Creator interface for different create modes
type Creator interface {
	Create(event *Event) (*CreateResult, error)
	Name() string
}

// Lister interface for listing events
type Lister interface {
	List(start, end time.Time, maxResults int) (*ListResult, error)
}
