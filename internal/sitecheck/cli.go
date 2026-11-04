// Package sitecheck checks site reachability from multiple global locations.
//
// Uses the check-host.net API to verify a URL is accessible from
// different geographic regions (US, EU, Asia, etc.).
//
// This file holds what cmd/sitecheck's commands call: Check, All and Report.
package sitecheck

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// noise: fewer failed nodes than this is check-host.net's own nodes failing now and then, not the site. A check passes
// below it, and the report raises no issue below it.
const noise = 3

// URL is what the checks check unless told otherwise: SITE_URL's robots.txt (the site's, without SITE_URL).
func URL() string {
	initConfig()
	return defaultURL
}

// Check runs one check (http, dns, tcp or redirect) of target from up to nodes global locations, waiting wait seconds
// for their results, and reports whether it passed.
func Check(checkType, target string, nodes, wait int, stdout, stderr io.Writer) bool {
	initConfig()
	return runCheck(checkType, target, nodes, wait, stdout, stderr)
}

// All runs every check, dns, tcp, redirect and http, one after the other, and reports whether all passed.
func All(target string, nodes, wait int, stdout, stderr io.Writer) bool {
	initConfig()
	allPassed := true
	for _, ct := range []string{"dns", "tcp", "redirect", "http"} {
		fmt.Fprintf(stdout, "=== %s Check ===\n", strings.ToUpper(ct))
		if !runCheck(ct, target, nodes, wait, stdout, stderr) {
			allPassed = false
		}
		fmt.Fprintln(stdout)
	}
	if allPassed {
		fmt.Fprintln(stdout, "Summary: All checks passed")
	}
	return allPassed
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

	return len(failures) < noise
}

// Report runs the HTTP check and compares it with the previous run's (.sitecheck-state.json), then saves this one's:
// a markdown report for a GitHub issue, and whether it found an issue (failures, more of them, or slower responses).
func Report(target string, nodes, wait int, stderr io.Writer) (report string, issues bool, err error) {
	initConfig()
	host := prepareHost("http", target)

	requestID, _, err := initiateCheck("http", host, nodes)
	if err != nil {
		return "", false, fmt.Errorf("failed to initiate check: %w", err)
	}

	time.Sleep(time.Duration(wait) * time.Second)

	results, err := getResults(requestID, "http")
	if err != nil {
		return "", false, fmt.Errorf("failed to get results: %w", err)
	}

	current := buildState("http", results)
	previous, _ := LoadState()
	report, issues = generateMarkdownReport(current, previous, results, requestID)

	if err := SaveState(current); err != nil {
		fmt.Fprintf(stderr, "Warning: failed to save state: %v\n", err)
	}
	return report, issues, nil
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

	// Issue 1: failures beyond the noise
	if current.FailedCount >= noise {
		issues = append(issues, fmt.Sprintf("**%d nodes failed** to reach the site", current.FailedCount))
		hasIssues = true
	}

	// Issue 2: New failures compared to previous
	if previous != nil {
		if previous.FailedCount < noise && current.FailedCount >= noise {
			issues = append(issues, fmt.Sprintf("**New failures detected** (was %d, within the noise)", previous.FailedCount))
			hasIssues = true
		} else if current.FailedCount >= previous.FailedCount+noise {
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
