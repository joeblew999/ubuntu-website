// sitecheck checks site reachability from multiple global locations.
//
// Uses the check-host.net API to verify a URL is accessible from
// different geographic regions (US, EU, Asia, etc.).
//
// Usage:
//
//	go run cmd/sitecheck/main.go                           # Check production URL
//	go run cmd/sitecheck/main.go -url https://example.com  # Check custom URL
//	task site:check                                        # Via Taskfile
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	defaultURL   = "https://www.ubuntusoftware.net"
	defaultNodes = 56 // Use all available nodes
	defaultWait  = 8  // Seconds to wait for global responses
	apiBase      = "https://check-host.net"
)

// CheckResponse is the initial response from check-host.net
type CheckResponse struct {
	OK        int               `json:"ok"`
	RequestID string            `json:"request_id"`
	Nodes     map[string][]any  `json:"nodes"`
	PermanentLink string        `json:"permanent_link"`
}

// Result represents a single check result from a node
type Result struct {
	Node     string
	Success  bool
	Time     float64 // seconds
	Status   string  // HTTP status or error
	IP       string
	Pending  bool
}

func main() {
	urlFlag := flag.String("url", defaultURL, "URL to check")
	nodesFlag := flag.Int("nodes", defaultNodes, "Maximum number of global nodes to check from")
	waitFlag := flag.Int("wait", defaultWait, "Seconds to wait for results")
	flag.Parse()

	fmt.Printf("Checking %s from %d global locations...\n\n", *urlFlag, *nodesFlag)

	// Initiate the check
	requestID, nodes, err := initiateCheck(*urlFlag, *nodesFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initiate check: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Request ID: %s\n", requestID)
	fmt.Printf("Checking from: %s\n", strings.Join(nodeNames(nodes), ", "))
	fmt.Printf("Waiting %d seconds for results...\n\n", *waitFlag)

	time.Sleep(time.Duration(*waitFlag) * time.Second)

	// Get results
	results, err := getResults(requestID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get results: %v\n", err)
		os.Exit(1)
	}

	// Sort by node name for consistent output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Node < results[j].Node
	})

	// Count and collect failures
	var failures []Result
	pending := 0
	for _, r := range results {
		if r.Pending {
			pending++
		} else if !r.Success {
			failures = append(failures, r)
		}
	}

	ok := len(results) - len(failures) - pending

	// Only print failures (if any)
	if len(failures) > 0 {
		fmt.Println("Failures:")
		for _, r := range failures {
			fmt.Printf("  ✗ %s: %s\n", r.Node, r.Status)
		}
		fmt.Println()
	}

	// Summary line
	fmt.Printf("✓ %d/%d nodes OK", ok, len(results))
	if len(failures) > 0 {
		fmt.Printf(", %d failed", len(failures))
	}
	if pending > 0 {
		fmt.Printf(", %d pending", pending)
	}
	fmt.Println()
	fmt.Printf("Full report: %s/check-report/%s\n", apiBase, requestID)

	// Exit 1 only if 3+ failures (1-2 is noise)
	if len(failures) >= 3 {
		os.Exit(1)
	}
}

func initiateCheck(targetURL string, maxNodes int) (string, map[string][]any, error) {
	apiURL := fmt.Sprintf("%s/check-http?host=%s&max_nodes=%d",
		apiBase, url.QueryEscape(targetURL), maxNodes)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	var checkResp CheckResponse
	if err := json.Unmarshal(body, &checkResp); err != nil {
		return "", nil, fmt.Errorf("failed to parse response: %w\n%s", err, string(body))
	}

	if checkResp.RequestID == "" {
		return "", nil, fmt.Errorf("no request ID in response: %s", string(body))
	}

	return checkResp.RequestID, checkResp.Nodes, nil
}

func getResults(requestID string) ([]Result, error) {
	apiURL := fmt.Sprintf("%s/check-result/%s", apiBase, requestID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the dynamic JSON structure
	var rawResults map[string]json.RawMessage
	if err := json.Unmarshal(body, &rawResults); err != nil {
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	var results []Result
	for node, raw := range rawResults {
		r := Result{Node: node}

		// Check if null (pending)
		if string(raw) == "null" {
			r.Pending = true
			results = append(results, r)
			continue
		}

		// Parse the array structure: [[status, time_or_error, status_text, http_code, ip]]
		var nodeResult [][]any
		if err := json.Unmarshal(raw, &nodeResult); err != nil {
			r.Status = "parse error"
			results = append(results, r)
			continue
		}

		if len(nodeResult) == 0 || len(nodeResult[0]) < 3 {
			r.Status = "incomplete data"
			results = append(results, r)
			continue
		}

		data := nodeResult[0]

		// First element is success indicator (1 = success)
		if status, ok := data[0].(float64); ok && status == 1 {
			r.Success = true
			// Second element is response time in seconds
			if t, ok := data[1].(float64); ok {
				r.Time = t
			}
			// Third element is status text
			if s, ok := data[2].(string); ok {
				r.Status = s
			}
			// Fourth element is HTTP code (can be string or number)
			if len(data) > 3 {
				switch v := data[3].(type) {
				case string:
					r.Status = v
				case float64:
					r.Status = fmt.Sprintf("%d", int(v))
				}
			}
			// Fifth element is IP
			if len(data) > 4 {
				if ip, ok := data[4].(string); ok {
					r.IP = ip
				}
			}
		} else {
			// Failure - second element contains error message
			r.Success = false
			if len(data) > 2 {
				if errMsg, ok := data[2].(string); ok {
					r.Status = errMsg
				}
			}
			if r.Status == "" {
				r.Status = "unknown error"
			}
		}

		results = append(results, r)
	}

	return results, nil
}

func nodeNames(nodes map[string][]any) []string {
	names := make([]string, 0, len(nodes))
	for name := range nodes {
		// Extract country code from node name (e.g., "us1.node.check-host.net" -> "us1")
		short := strings.Split(name, ".")[0]
		names = append(names, short)
	}
	sort.Strings(names)
	return names
}
