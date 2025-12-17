// Package cfanalytics provides Cloudflare Web Analytics fetching and reporting.
//
// This file contains the CLI entry point. The main.go in cmd/cfanalytics
// just imports and calls Run().
package cfanalytics

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"
)

// CLIOptions holds global CLI flags
type CLIOptions struct {
	WebhookURL  string
	Days        int
	Verbose     bool
	GithubIssue bool
	Version     bool
}

// Run is the main entry point for the analytics CLI.
// Returns exit code (0 = success, 1 = error or changes detected in github-issue mode).
func Run(args []string, version string, stdout, stderr io.Writer) int {
	// Parse flags
	fs := flag.NewFlagSet("analytics", flag.ContinueOnError)
	fs.SetOutput(stderr)

	opts := &CLIOptions{}
	fs.StringVar(&opts.WebhookURL, "webhook", "", "Webhook URL to post changes (Slack/Discord)")
	fs.IntVar(&opts.Days, "days", 7, "Number of days to analyze")
	fs.BoolVar(&opts.Verbose, "v", false, "Verbose output")
	fs.BoolVar(&opts.GithubIssue, "github-issue", false, "Output markdown for GitHub Issue (exits 1 if changes detected)")
	fs.BoolVar(&opts.Version, "version", false, "Print version and exit")

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if opts.Version {
		fmt.Fprintf(stdout, "analytics %s\n", version)
		return 0
	}

	token := os.Getenv("CLOUDFLARE_API_TOKEN")
	if token == "" {
		fmt.Fprintln(stderr, "Error: CLOUDFLARE_API_TOKEN environment variable not set")
		fmt.Fprintln(stderr, "")
		fmt.Fprintln(stderr, "Create a token at: https://dash.cloudflare.com/profile/api-tokens")
		fmt.Fprintln(stderr, "Required permissions: Account Analytics:Read")
		return 1
	}

	// Calculate date range
	until := time.Now().UTC().Truncate(24 * time.Hour)
	since := until.AddDate(0, 0, -opts.Days)

	if opts.Verbose {
		fmt.Fprintf(stdout, "Fetching analytics for %s to %s...\n", since.Format("2006-01-02"), until.Format("2006-01-02"))
	}

	// Fetch current analytics
	current, err := FetchAnalytics(token, since, until)
	if err != nil {
		fmt.Fprintf(stderr, "Error fetching analytics: %v\n", err)
		return 1
	}
	current.Period = fmt.Sprintf("%s to %s", since.Format("Jan 2"), until.Format("Jan 2"))

	// Load previous state
	previous, err := LoadState()
	if err != nil && opts.Verbose {
		fmt.Fprintln(stdout, "No previous state found (first run)")
	}

	// Generate report
	report := GenerateReport(current, previous)

	// GitHub Issue mode: output markdown and exit with code based on changes
	if opts.GithubIssue {
		fmt.Fprintln(stdout, GenerateMarkdownReport(current, previous, report))
		// Save state for next run
		if err := SaveState(current); err != nil {
			fmt.Fprintf(stderr, "Warning: failed to save state: %v\n", err)
		}
		if report.HasChanges {
			return 1 // Signal to workflow that issue should be created
		}
		return 0
	}

	// Print report
	fmt.Fprintln(stdout, report.Summary)
	if len(report.Changes) > 0 {
		fmt.Fprintln(stdout, "\nSignificant Changes:")
		for _, change := range report.Changes {
			fmt.Fprintf(stdout, "  %s\n", change)
		}
	}

	// Save current state for next comparison
	if err := SaveState(current); err != nil {
		fmt.Fprintf(stderr, "Warning: failed to save state: %v\n", err)
	}

	// Post to webhook if configured and there are changes
	if opts.WebhookURL != "" && len(report.Changes) > 0 {
		if err := PostToWebhook(opts.WebhookURL, report); err != nil {
			fmt.Fprintf(stderr, "Warning: failed to post to webhook: %v\n", err)
		} else if opts.Verbose {
			fmt.Fprintln(stdout, "Posted to webhook")
		}
	}

	return 0
}

// GenerateReport creates an analytics report comparing current and previous state.
func GenerateReport(current, previous *State) Report {
	report := Report{}

	// Build summary
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Analytics Report (%s)\n", current.Period))
	sb.WriteString(strings.Repeat("=", 40) + "\n")
	sb.WriteString(fmt.Sprintf("Visits:     %d\n", current.Visits))
	sb.WriteString(fmt.Sprintf("Page Views: %d\n", current.PageViews))

	// Top pages
	if len(current.TopPages) > 0 {
		sb.WriteString("\nTop Pages:\n")
		pages := sortMapByValue(current.TopPages, 5)
		for _, p := range pages {
			sb.WriteString(fmt.Sprintf("  %s: %d\n", p.Key, p.Value))
		}
	}

	// Top countries
	if len(current.Countries) > 0 {
		sb.WriteString("\nTop Countries:\n")
		countries := sortMapByValue(current.Countries, 5)
		for _, c := range countries {
			sb.WriteString(fmt.Sprintf("  %s: %d\n", c.Key, c.Value))
		}
	}

	report.Summary = sb.String()

	// Compare with previous if available
	if previous != nil {
		// Check visits change
		if change := percentChange(previous.Visits, current.Visits); math.Abs(change) >= ChangeThreshold*100 {
			direction := "increased"
			if change < 0 {
				direction = "decreased"
			}
			report.Changes = append(report.Changes,
				fmt.Sprintf("Visits %s %.0f%% (%d -> %d)", direction, math.Abs(change), previous.Visits, current.Visits))
			report.HasChanges = true
		}

		// Check pageviews change
		if change := percentChange(previous.PageViews, current.PageViews); math.Abs(change) >= ChangeThreshold*100 {
			direction := "increased"
			if change < 0 {
				direction = "decreased"
			}
			report.Changes = append(report.Changes,
				fmt.Sprintf("Page views %s %.0f%% (%d -> %d)", direction, math.Abs(change), previous.PageViews, current.PageViews))
			report.HasChanges = true
		}
	}

	return report
}

// GenerateMarkdownReport creates a markdown-formatted report for GitHub Issues.
func GenerateMarkdownReport(current, previous *State, report Report) string {
	var sb strings.Builder

	sb.WriteString("## Analytics Change Detected\n\n")
	sb.WriteString(fmt.Sprintf("**Period:** %s\n\n", current.Period))

	// Changes section
	if len(report.Changes) > 0 {
		sb.WriteString("### Changes\n")
		for _, change := range report.Changes {
			sb.WriteString(fmt.Sprintf("- **%s**\n", change))
		}
		sb.WriteString("\n")
	}

	// Comparison table
	sb.WriteString("### Current Stats\n\n")
	sb.WriteString("| Metric | Previous | Current | Change |\n")
	sb.WriteString("|--------|----------|---------|--------|\n")

	if previous != nil {
		visitChange := percentChange(previous.Visits, current.Visits)
		pvChange := percentChange(previous.PageViews, current.PageViews)
		sb.WriteString(fmt.Sprintf("| Visits | %d | %d | %+.0f%% |\n", previous.Visits, current.Visits, visitChange))
		sb.WriteString(fmt.Sprintf("| Page Views | %d | %d | %+.0f%% |\n", previous.PageViews, current.PageViews, pvChange))
	} else {
		sb.WriteString(fmt.Sprintf("| Visits | - | %d | (first run) |\n", current.Visits))
		sb.WriteString(fmt.Sprintf("| Page Views | - | %d | (first run) |\n", current.PageViews))
	}

	// Top pages
	if len(current.TopPages) > 0 {
		sb.WriteString("\n### Top Pages\n")
		pages := sortMapByValue(current.TopPages, 5)
		for i, p := range pages {
			sb.WriteString(fmt.Sprintf("%d. `%s` - %d views\n", i+1, p.Key, p.Value))
		}
	}

	// Top countries
	if len(current.Countries) > 0 {
		sb.WriteString("\n### Top Countries\n")
		countries := sortMapByValue(current.Countries, 5)
		for i, c := range countries {
			sb.WriteString(fmt.Sprintf("%d. %s - %d\n", i+1, c.Key, c.Value))
		}
	}

	sb.WriteString("\n---\n*Generated by analytics change detection workflow*\n")

	return sb.String()
}

func percentChange(old, new int64) float64 {
	if old == 0 {
		if new > 0 {
			return 100
		}
		return 0
	}
	return float64(new-old) / float64(old) * 100
}
