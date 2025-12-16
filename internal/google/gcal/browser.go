package gcal

import (
	"fmt"
	"net/url"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/browser"
)

// BrowserClient manages calendar via browser automation
type BrowserClient struct {
	config      *Config
	composeOnly bool // If true, just open calendar, don't create
	headless    bool // If true, run browser headless (no UI) - for production use
}

// NewBrowserClient creates a new browser client (visible mode for dev verification)
func NewBrowserClient(config *Config, composeOnly bool) *BrowserClient {
	return &BrowserClient{
		config:      config,
		composeOnly: composeOnly,
		headless:    false,
	}
}

// NewBrowserClientWithOptions creates a browser client with full control
func NewBrowserClientWithOptions(config *Config, composeOnly, headless bool) *BrowserClient {
	return &BrowserClient{
		config:      config,
		composeOnly: composeOnly,
		headless:    headless,
	}
}

// SetHeadless sets headless mode (true = no browser UI, false = show browser)
func (c *BrowserClient) SetHeadless(headless bool) {
	c.headless = headless
}

// IsHeadless returns whether headless mode is enabled
func (c *BrowserClient) IsHeadless() bool {
	return c.headless
}

// Name returns the client name
func (c *BrowserClient) Name() string {
	if c.composeOnly {
		return "compose"
	}
	if c.headless {
		return "browser-headless"
	}
	return "browser"
}

// Create opens Google Calendar to create an event
func (c *BrowserClient) Create(event *Event) (*CreateResult, error) {
	if err := event.Validate(); err != nil {
		return &CreateResult{
			Success: false,
			Error:   err.Error(),
			Mode:    c.Name(),
		}, err
	}

	// Build Google Calendar URL
	calURL := buildCalendarURL(event)

	// Headless mode: use Playwright for automation without UI
	// Visible mode: open in default browser for user verification
	if c.headless && !c.composeOnly {
		// For headless browser automation, recommend using API mode instead
		// Browser automation in headless mode is complex due to Google auth
		// The API mode is more reliable for programmatic calendar operations
		return &CreateResult{
			Success: false,
			Error:   "headless browser mode: use --mode=api for programmatic calendar access",
			Mode:    c.Name(),
		}, fmt.Errorf("headless browser mode not supported - use API mode instead")
	}

	// Open in browser (visible mode for dev verification or compose mode)
	if err := browser.OpenURL(calURL); err != nil {
		return &CreateResult{
			Success: false,
			Error:   fmt.Sprintf("failed to open browser: %v", err),
			Mode:    c.Name(),
		}, err
	}

	return &CreateResult{
		Success: true,
		Link:    calURL,
		Mode:    c.Name(),
	}, nil
}

// buildCalendarURL builds a Google Calendar event creation URL
func buildCalendarURL(event *Event) string {
	params := url.Values{}
	params.Set("action", "TEMPLATE")
	params.Set("text", event.Title)

	// Format dates for Google Calendar URL
	// Format: YYYYMMDDTHHmmss/YYYYMMDDTHHmmss
	dateFormat := "20060102T150405"
	dates := event.Start.Format(dateFormat) + "/" + event.End.Format(dateFormat)
	params.Set("dates", dates)

	if event.Description != "" {
		params.Set("details", event.Description)
	}
	if event.Location != "" {
		params.Set("location", event.Location)
	}
	if len(event.Attendees) > 0 {
		// Join attendees with comma
		attendees := ""
		for i, a := range event.Attendees {
			if i > 0 {
				attendees += ","
			}
			attendees += a
		}
		params.Set("add", attendees)
	}

	return "https://calendar.google.com/calendar/render?" + params.Encode()
}

// OpenCalendar opens Google Calendar in the browser
func OpenCalendar(view string) error {
	var calURL string
	switch view {
	case "day":
		calURL = "https://calendar.google.com/calendar/r/day"
	case "week":
		calURL = "https://calendar.google.com/calendar/r/week"
	case "month":
		calURL = "https://calendar.google.com/calendar/r/month"
	case "agenda":
		calURL = "https://calendar.google.com/calendar/r/agenda"
	default:
		calURL = "https://calendar.google.com/calendar/r"
	}

	return browser.OpenURL(calURL)
}

// OpenEventCreate opens the event creation page
func OpenEventCreate() error {
	return browser.OpenURL("https://calendar.google.com/calendar/r/eventedit")
}

// FormatEventTime formats a time for display
func FormatEventTime(t time.Time) string {
	return t.Format("Mon Jan 2, 3:04 PM")
}

// FormatEventDate formats a date for display
func FormatEventDate(t time.Time) string {
	return t.Format("Mon Jan 2, 2006")
}
