// Package sitecheck checks site reachability from multiple global locations.
//
// Uses the check-host.net API to verify a URL is accessible from
// different geographic regions (US, EU, Asia, etc.).
//
// This file contains the CLI entry point. The main.go in cmd/sitecheck
// just imports and calls Run().
package sitecheck

import (
	"flag"
	"fmt"
	"io"
	"strings"
	"time"
)

// Run is the main entry point for the sitecheck CLI.
// Returns exit code (0 = success, 1 = error or issues detected).
func Run(args []string, version string, stdout, stderr io.Writer) int {
	// Initialize config from environment
	initConfig()

	fs := flag.NewFlagSet("sitecheck", flag.ContinueOnError)
	fs.SetOutput(stderr)

	urlFlag := fs.String("url", defaultURL, "URL to check (for HTTP) or domain (for DNS/TCP)")
	typeFlag := fs.String("type", "http", "Check type: http, dns, tcp, redirect, or all")
	nodesFlag := fs.Int("nodes", defaultNodes, "Maximum number of global nodes to check from")
	waitFlag := fs.Int("wait", defaultWait, "Seconds to wait for results")
	githubIssue := fs.Bool("github-issue", false, "Output markdown for GitHub Issue (exits 1 if issues detected)")
	ver := fs.Bool("version", false, "Print version and exit")

	if err := fs.Parse(args[1:]); err != nil {
		return 1
	}

	if *ver {
		fmt.Fprintf(stdout, "sitecheck %s\n", version)
		return 0
	}

	checkType := strings.ToLower(*typeFlag)

	// GitHub Issue mode: run HTTP check and output markdown
	if *githubIssue {
		return runGitHubIssueMode(*urlFlag, *nodesFlag, *waitFlag, stdout, stderr)
	}

	if checkType == "all" {
		// Run all checks sequentially
		allPassed := true
		for _, ct := range []string{"dns", "tcp", "redirect", "http"} {
			fmt.Fprintf(stdout, "=== %s Check ===\n", strings.ToUpper(ct))
			passed := runCheck(ct, *urlFlag, *nodesFlag, *waitFlag, stdout, stderr)
			if !passed {
				allPassed = false
			}
			fmt.Fprintln(stdout)
		}
		if !allPassed {
			return 1
		}
		fmt.Fprintln(stdout, "Summary: All checks passed")
		return 0
	}

	if _, ok := checkEndpoints[checkType]; !ok {
		fmt.Fprintf(stderr, "Unknown check type: %s (use http, dns, tcp, redirect, or all)\n", checkType)
		return 1
	}

	if !runCheck(checkType, *urlFlag, *nodesFlag, *waitFlag, stdout, stderr) {
		return 1
	}
	return 0
}

// runCheck executes a single check type and returns true if passed
func runCheck(checkType, targetURL string, maxNodes, waitSecs int, stdout, stderr io.Writer) bool {
	// Prepare the host parameter based on check type
	host := prepareHost(checkType, targetURL)

	fmt.Fprintf(stdout, "Checking %s from %d global locations...\n\n", host, maxNodes)

	// Initiate the check
	requestID, nodes, err := initiateCheck(checkType, host, maxNodes)
	if err != nil {
		fmt.Fprintf(stderr, "Failed to initiate check: %v\n", err)
		return false
	}

	fmt.Fprintf(stdout, "Waiting %d seconds for %d nodes...\n", waitSecs, len(nodes))

	time.Sleep(time.Duration(waitSecs) * time.Second)

	// Get results
	results, err := getResults(requestID, checkType)
	if err != nil {
		fmt.Fprintf(stderr, "Failed to get results: %v\n", err)
		return false
	}

	// Sort and process results
	sortResults(results)

	// Count and collect failures, track response times
	var failures []Result
	var times []float64
	pending := 0
	for _, r := range results {
		if r.Pending {
			pending++
		} else if !r.Success {
			failures = append(failures, r)
		} else if r.Time > 0 {
			times = append(times, r.Time)
		}
	}

	ok := len(results) - len(failures) - pending

	// Only print failures (if any)
	if len(failures) > 0 {
		fmt.Fprintln(stdout, "Failures:")
		for _, r := range failures {
			fmt.Fprintf(stdout, "  ✗ %s: %s\n", r.Node, r.Status)
		}
		fmt.Fprintln(stdout)
	}

	// Summary line
	fmt.Fprintf(stdout, "✓ %d/%d nodes OK", ok, len(results))
	if len(failures) > 0 {
		fmt.Fprintf(stdout, ", %d failed", len(failures))
	}
	if pending > 0 {
		fmt.Fprintf(stdout, ", %d pending", pending)
	}
	// Show response times if available
	if len(times) > 0 {
		var sum, max float64
		for _, t := range times {
			sum += t
			if t > max {
				max = t
			}
		}
		avg := sum / float64(len(times))
		fmt.Fprintf(stdout, " (avg %.0fms, max %.0fms)", avg*1000, max*1000)
	}
	fmt.Fprintf(stdout, " - %s/check-report/%s\n", apiBase, requestID)

	// Exit 1 only if 3+ failures (1-2 is noise)
	return len(failures) < 3
}

// runGitHubIssueMode runs HTTP check in GitHub Issue mode
// Outputs markdown report and returns exit code (1 if issues detected)
func runGitHubIssueMode(targetURL string, maxNodes, waitSecs int, stdout, stderr io.Writer) int {
	host := prepareHost("http", targetURL)

	// Initiate HTTP check
	requestID, _, err := initiateCheck("http", host, maxNodes)
	if err != nil {
		fmt.Fprintf(stderr, "Failed to initiate check: %v\n", err)
		return 1
	}

	time.Sleep(time.Duration(waitSecs) * time.Second)

	// Get results
	results, err := getResults(requestID, "http")
	if err != nil {
		fmt.Fprintf(stderr, "Failed to get results: %v\n", err)
		return 1
	}

	// Build current state
	current := buildState("http", results)

	// Load previous state
	previous, _ := LoadState()

	// Generate markdown report
	report, hasIssues := generateMarkdownReport(current, previous, results, requestID)
	fmt.Fprintln(stdout, report)

	// Save current state
	if err := SaveState(current); err != nil {
		fmt.Fprintf(stderr, "Warning: failed to save state: %v\n", err)
	}

	// Exit 1 if issues detected (triggers GitHub Issue creation)
	if hasIssues {
		return 1
	}
	return 0
}

// generateMarkdownReport creates markdown output and returns whether issues were detected
func generateMarkdownReport(current, previous *State, results []Result, requestID string) (string, bool) {
	var sb strings.Builder
	hasIssues := false

	sb.WriteString("## Site Reachability Check\n\n")
	sb.WriteString(fmt.Sprintf("**Target:** %s\n", defaultURL))
	sb.WriteString(fmt.Sprintf("**Time:** %s\n\n", current.Timestamp.Format("2006-01-02 15:04 UTC")))

	// Check for issues
	var issues []string

	// Issue 1: Any failures
	if current.FailedCount > 0 {
		issues = append(issues, fmt.Sprintf("**%d nodes failed** to reach the site", current.FailedCount))
		hasIssues = true
	}

	// Issue 2: New failures compared to previous
	if previous != nil {
		if previous.FailedCount == 0 && current.FailedCount > 0 {
			issues = append(issues, "**New failures detected** (was 100% reachable)")
			hasIssues = true
		} else if current.FailedCount >= previous.FailedCount+3 {
			issues = append(issues, fmt.Sprintf("**Failure count increased** from %d to %d", previous.FailedCount, current.FailedCount))
			hasIssues = true
		}

		// Issue 3: Response time degradation (>50% increase)
		if previous.AvgResponseMS > 0 && current.AvgResponseMS > previous.AvgResponseMS*1.5 {
			issues = append(issues, fmt.Sprintf("**Response time degraded** from %.0fms to %.0fms avg (+%.0f%%)",
				previous.AvgResponseMS, current.AvgResponseMS,
				(current.AvgResponseMS-previous.AvgResponseMS)/previous.AvgResponseMS*100))
			hasIssues = true
		}
	}

	// Issues section
	if len(issues) > 0 {
		sb.WriteString("### Issues Detected\n")
		for _, issue := range issues {
			sb.WriteString(fmt.Sprintf("- %s\n", issue))
		}
		sb.WriteString("\n")
	}

	// Results summary
	sb.WriteString("### Results\n\n")
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Nodes Checked | %d |\n", current.TotalNodes))
	sb.WriteString(fmt.Sprintf("| Successful | %d |\n", current.OKCount))
	sb.WriteString(fmt.Sprintf("| Failed | %d |\n", current.FailedCount))
	if current.AvgResponseMS > 0 {
		sb.WriteString(fmt.Sprintf("| Avg Response | %.0fms |\n", current.AvgResponseMS))
		sb.WriteString(fmt.Sprintf("| Max Response | %.0fms |\n", current.MaxResponseMS))
	}

	// Comparison with previous
	if previous != nil {
		sb.WriteString("\n### Comparison with Previous Check\n\n")
		sb.WriteString("| Metric | Previous | Current | Change |\n")
		sb.WriteString("|--------|----------|---------|--------|\n")
		sb.WriteString(fmt.Sprintf("| Failed Nodes | %d | %d | %+d |\n",
			previous.FailedCount, current.FailedCount, current.FailedCount-previous.FailedCount))
		if previous.AvgResponseMS > 0 && current.AvgResponseMS > 0 {
			change := (current.AvgResponseMS - previous.AvgResponseMS) / previous.AvgResponseMS * 100
			sb.WriteString(fmt.Sprintf("| Avg Response | %.0fms | %.0fms | %+.0f%% |\n",
				previous.AvgResponseMS, current.AvgResponseMS, change))
		}
	}

	// Failed nodes detail
	if current.FailedCount > 0 {
		sb.WriteString("\n### Failed Nodes\n")
		for _, r := range results {
			if !r.Success && !r.Pending {
				sb.WriteString(fmt.Sprintf("- **%s**: %s\n", r.Node, r.Status))
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\n---\n[Full Report](%s/check-report/%s) | *Generated by site monitor workflow*\n", apiBase, requestID))

	return sb.String(), hasIssues
}
