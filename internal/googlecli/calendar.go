package googlecli

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joeblew999/ubuntu-website/internal/google/gcal"
)

func (c *cliContext) handleCalendar(args []string) {
	if len(args) < 1 {
		c.printCalendarUsage()
		return
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "list":
		c.calendarList(cmdArgs)
	case "today":
		c.calendarToday(cmdArgs)
	case "create":
		c.calendarCreate(cmdArgs)
	case "compose":
		c.calendarCompose(cmdArgs)
	case "check":
		c.calendarCheck(cmdArgs)
	case "open":
		c.calendarOpen(cmdArgs)
	case "server":
		c.calendarServer(cmdArgs)
	default:
		fmt.Fprintf(c.stderr, "Unknown calendar command: %s\n", cmd)
		c.printCalendarUsage()
		os.Exit(1)
	}
}

func (c *cliContext) calendarList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	startTime := fs.String("start", "", "Start time (default: now)")
	endTime := fs.String("end", "", "End time (default: end of day)")
	maxResults := fs.Int("max", 10, "Maximum events")
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	now := time.Now()
	start := now
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	if *startTime != "" {
		var err error
		start, err = parseTime(*startTime)
		if err != nil {
			c.exitError(fmt.Sprintf("Invalid start time: %v", err))
		}
	}
	if *endTime != "" {
		var err error
		end, err = parseTime(*endTime)
		if err != nil {
			c.exitError(fmt.Sprintf("Invalid end time: %v", err))
		}
	}

	config := gcal.DefaultConfig()
	client, err := gcal.NewAPIClient(config)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create API client: %v", err))
	}

	result, err := client.List(start, end, *maxResults)
	if err != nil {
		c.exitError(fmt.Sprintf("List failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
	} else {
		if len(result.Events) == 0 {
			fmt.Fprintln(c.stdout, "No events found.")
			return
		}
		fmt.Fprintf(c.stdout, "Events (%d):\n", len(result.Events))
		for _, event := range result.Events {
			fmt.Fprintf(c.stdout, "\n  %s\n", event.Title)
			fmt.Fprintf(c.stdout, "    %s - %s\n", gcal.FormatEventTime(event.Start), gcal.FormatEventTime(event.End))
			if event.Location != "" {
				fmt.Fprintf(c.stdout, "    Location: %s\n", event.Location)
			}
		}
	}
}

func (c *cliContext) calendarToday(args []string) {
	fs := flag.NewFlagSet("today", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	config := gcal.DefaultConfig()
	client, err := gcal.NewAPIClient(config)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create API client: %v", err))
	}

	result, err := client.List(start, end, 20)
	if err != nil {
		c.exitError(fmt.Sprintf("List failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
	} else {
		fmt.Fprintf(c.stdout, "Today's Events - %s\n", gcal.FormatEventDate(now))
		if len(result.Events) == 0 {
			fmt.Fprintln(c.stdout, "\nNo events today.")
			return
		}
		for _, event := range result.Events {
			fmt.Fprintf(c.stdout, "\n  %s\n", event.Title)
			fmt.Fprintf(c.stdout, "    %s - %s\n", event.Start.Format("3:04 PM"), event.End.Format("3:04 PM"))
			if event.Location != "" {
				fmt.Fprintf(c.stdout, "    Location: %s\n", event.Location)
			}
		}
	}
}

func (c *cliContext) calendarCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	title := fs.String("title", "", "Event title")
	startTime := fs.String("start", "", "Start time")
	endTime := fs.String("end", "", "End time")
	description := fs.String("description", "", "Event description")
	location := fs.String("location", "", "Event location")
	attendees := fs.String("attendees", "", "Comma-separated attendee emails")
	mode := fs.String("mode", "api", "Create mode: api or browser")
	headless := fs.Bool("headless", false, "Run browser headless")
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	if *title == "" {
		c.exitError("--title is required")
	}
	if *startTime == "" {
		c.exitError("--start is required")
	}
	if *endTime == "" {
		c.exitError("--end is required")
	}

	start, err := parseTime(*startTime)
	if err != nil {
		c.exitError(fmt.Sprintf("Invalid start time: %v", err))
	}
	end, err := parseTime(*endTime)
	if err != nil {
		c.exitError(fmt.Sprintf("Invalid end time: %v", err))
	}

	var attendeeList []string
	if *attendees != "" {
		attendeeList = strings.Split(*attendees, ",")
		for i := range attendeeList {
			attendeeList[i] = strings.TrimSpace(attendeeList[i])
		}
	}

	config := gcal.DefaultConfig()
	event := &gcal.Event{
		Title:       *title,
		Description: *description,
		Location:    *location,
		Start:       start,
		End:         end,
		Attendees:   attendeeList,
	}

	var creator gcal.Creator
	switch strings.ToLower(*mode) {
	case "api":
		client, err := gcal.NewAPIClient(config)
		if err != nil {
			c.exitError(fmt.Sprintf("Failed to create API client: %v", err))
		}
		creator = client
	case "browser":
		creator = gcal.NewBrowserClientWithOptions(config, false, *headless)
	default:
		c.exitError(fmt.Sprintf("Invalid mode: %s (use 'api' or 'browser')", *mode))
	}

	result, err := creator.Create(event)
	if err != nil {
		c.exitError(fmt.Sprintf("Create failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
	} else {
		fmt.Fprintln(c.stdout, "Event created successfully!")
		fmt.Fprintf(c.stdout, "  Title: %s\n", *title)
		fmt.Fprintf(c.stdout, "  Start: %s\n", gcal.FormatEventTime(start))
		fmt.Fprintf(c.stdout, "  End: %s\n", gcal.FormatEventTime(end))
		fmt.Fprintf(c.stdout, "  Mode: %s\n", result.Mode)
		if result.EventID != "" {
			fmt.Fprintf(c.stdout, "  Event ID: %s\n", result.EventID)
		}
		if result.Link != "" {
			fmt.Fprintf(c.stdout, "  Link: %s\n", result.Link)
		}
	}
}

func (c *cliContext) calendarCompose(args []string) {
	fs := flag.NewFlagSet("compose", flag.ExitOnError)
	title := fs.String("title", "", "Event title")
	startTime := fs.String("start", "", "Start time")
	endTime := fs.String("end", "", "End time")
	description := fs.String("description", "", "Event description")
	location := fs.String("location", "", "Event location")
	attendees := fs.String("attendees", "", "Comma-separated attendee emails")
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	if *title == "" {
		c.exitError("--title is required")
	}
	if *startTime == "" {
		c.exitError("--start is required")
	}
	if *endTime == "" {
		c.exitError("--end is required")
	}

	start, err := parseTime(*startTime)
	if err != nil {
		c.exitError(fmt.Sprintf("Invalid start time: %v", err))
	}
	end, err := parseTime(*endTime)
	if err != nil {
		c.exitError(fmt.Sprintf("Invalid end time: %v", err))
	}

	var attendeeList []string
	if *attendees != "" {
		attendeeList = strings.Split(*attendees, ",")
		for i := range attendeeList {
			attendeeList[i] = strings.TrimSpace(attendeeList[i])
		}
	}

	config := gcal.DefaultConfig()
	event := &gcal.Event{
		Title:       *title,
		Description: *description,
		Location:    *location,
		Start:       start,
		End:         end,
		Attendees:   attendeeList,
	}

	creator := gcal.NewBrowserClient(config, true)
	result, err := creator.Create(event)
	if err != nil {
		c.exitError(fmt.Sprintf("Compose failed: %v", err))
	}

	if *jsonOutput {
		c.outputJSON(result)
	} else {
		fmt.Fprintln(c.stdout, "Calendar compose opened!")
		fmt.Fprintf(c.stdout, "  Title: %s\n", *title)
		fmt.Fprintf(c.stdout, "  Start: %s\n", gcal.FormatEventTime(start))
		fmt.Fprintf(c.stdout, "  End: %s\n", gcal.FormatEventTime(end))
		fmt.Fprintln(c.stdout, "\nReview and click Save in the browser.")
	}
}

func (c *cliContext) calendarCheck(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	jsonOutput := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	config := gcal.DefaultConfig()

	client, err := gcal.NewAPIClient(config)
	if err != nil {
		if *jsonOutput {
			c.outputJSON(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		} else {
			c.exitError(fmt.Sprintf("Failed to load token: %v", err))
		}
		os.Exit(1)
	}

	if err := client.Check(); err != nil {
		if *jsonOutput {
			c.outputJSON(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		} else {
			c.exitError(fmt.Sprintf("API check failed: %v", err))
		}
		os.Exit(1)
	}

	if *jsonOutput {
		c.outputJSON(map[string]interface{}{
			"success":  true,
			"calendar": config.CalendarID,
		})
	} else {
		fmt.Fprintln(c.stdout, "Calendar API connection OK!")
		fmt.Fprintf(c.stdout, "  Calendar: %s\n", config.CalendarID)
		fmt.Fprintf(c.stdout, "  Token path: %s\n", config.TokenPath)
	}
}

func (c *cliContext) calendarOpen(args []string) {
	view := ""
	if len(args) > 0 {
		view = args[0]
	}

	if err := gcal.OpenCalendar(view); err != nil {
		c.exitError(fmt.Sprintf("Failed to open calendar: %v", err))
	}
	fmt.Fprintln(c.stdout, "Opening Google Calendar...")
}

func (c *cliContext) calendarServer(args []string) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)
	port := fs.Int("port", 8088, "HTTP port")
	fs.Parse(args)

	config := gcal.DefaultConfig()

	server, err := gcal.NewServer(config, *port)
	if err != nil {
		c.exitError(fmt.Sprintf("Failed to create server: %v", err))
	}

	if err := server.Start(); err != nil {
		c.exitError(fmt.Sprintf("Server error: %v", err))
	}
}

func (c *cliContext) printCalendarUsage() {
	fmt.Fprintln(c.stdout, `Usage: google calendar <command> [arguments]

Commands:
  list [--start=TIME] [--end=TIME] [--max=N]  List events
  today                                        List today's events
  create --title=T --start=T --end=T           Create calendar event
  compose --title=T --start=T --end=T          Open calendar to create
  check                                        Verify API connection
  open [day|week|month|agenda]                 Open calendar in browser
  server [--port=8088]                         Start webhook server

Create Options:
  --title        Event title (required)
  --start        Start time (required)
  --end          End time (required)
  --description  Event description
  --location     Event location
  --attendees    Comma-separated emails
  --mode         Create mode: api (default) or browser
  --headless     Run browser headless

Options:
  --json    Output as JSON

Time Formats:
  RFC3339: 2024-12-13T14:00:00+07:00
  Relative: "today 2pm", "tomorrow 10am", "+1h"

Examples:
  google calendar list --max=5
  google calendar today
  google calendar create --title="Meeting" --start="tomorrow 2pm" --end="tomorrow 3pm"
  google calendar open week`)
}

// parseTime parses various time formats
func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	formats := []string{
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}
	for _, format := range formats {
		if t, err := time.ParseInLocation(format, s, time.Local); err == nil {
			return t, nil
		}
	}

	now := time.Now()
	s = strings.ToLower(strings.TrimSpace(s))

	if strings.HasPrefix(s, "today ") {
		timeStr := strings.TrimPrefix(s, "today ")
		return parseTimeOfDay(now, timeStr)
	}
	if strings.HasPrefix(s, "tomorrow ") {
		timeStr := strings.TrimPrefix(s, "tomorrow ")
		return parseTimeOfDay(now.AddDate(0, 0, 1), timeStr)
	}

	if strings.HasPrefix(s, "+") {
		d, err := time.ParseDuration(s[1:])
		if err == nil {
			return now.Add(d), nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized time format: %s", s)
}

func parseTimeOfDay(date time.Time, timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(strings.ToLower(timeStr))

	var hour, minute int
	var isPM bool

	if strings.HasSuffix(timeStr, "pm") {
		isPM = true
		timeStr = strings.TrimSuffix(timeStr, "pm")
	} else if strings.HasSuffix(timeStr, "am") {
		timeStr = strings.TrimSuffix(timeStr, "am")
	}

	if strings.Contains(timeStr, ":") {
		_, err := fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		_, err := fmt.Sscanf(timeStr, "%d", &hour)
		if err != nil {
			return time.Time{}, err
		}
	}

	if isPM && hour < 12 {
		hour += 12
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location()), nil
}
