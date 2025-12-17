// checkhost.go - check-host.net API client.
package sitecheck

import (
	"encoding/json"
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
	// Fallback values if SITE_URL env var not set
	fallbackSiteURL = "https://www.ubuntusoftware.net"

	defaultNodes = 56                       // Use all available nodes
	defaultWait  = 8                        // Seconds to wait for global responses
	apiBase      = "https://check-host.net"
	StateFile    = ".sitecheck-state.json"
)

// Package-level config derived from SITE_URL environment variable
var (
	defaultURL  string // Full URL to check (e.g., https://www.example.com/robots.txt)
	defaultHost string // Host portion (e.g., www.example.com)
	apexURL     string // Apex domain URL for redirect check (e.g., http://example.com)
)

// initConfig initializes config from environment variables
func initConfig() {
	siteURL := os.Getenv("SITE_URL")
	if siteURL == "" {
		siteURL = fallbackSiteURL
	}
	// Ensure no trailing slash
	siteURL = strings.TrimSuffix(siteURL, "/")

	// Derive check URL (robots.txt for HTTP check)
	defaultURL = siteURL + "/robots.txt"

	// Extract host from URL
	if parsed, err := url.Parse(siteURL); err == nil {
		defaultHost = parsed.Host
		// Derive apex URL (remove www. prefix, use http for redirect check)
		apexHost := strings.TrimPrefix(defaultHost, "www.")
		apexURL = "http://" + apexHost
	} else {
		defaultHost = "www.ubuntusoftware.net"
		apexURL = "http://ubuntusoftware.net"
	}
}

// Check type to API endpoint mapping
var checkEndpoints = map[string]string{
	"http":     "check-http",
	"dns":      "check-dns",
	"tcp":      "check-tcp",
	"redirect": "check-http", // Uses HTTP check but expects 301/302
}

// CheckResponse is the initial response from check-host.net
type CheckResponse struct {
	OK            int              `json:"ok"`
	RequestID     string           `json:"request_id"`
	Nodes         map[string][]any `json:"nodes"`
	PermanentLink string           `json:"permanent_link"`
}

// Result represents a single check result from a node
type Result struct {
	Node    string
	Success bool
	Time    float64 // seconds
	Status  string  // HTTP status, DNS records, or error
	IP      string
	Pending bool
}

// State represents stored check state for comparison
type State struct {
	Timestamp     time.Time `json:"timestamp"`
	CheckType     string    `json:"check_type"`
	TotalNodes    int       `json:"total_nodes"`
	OKCount       int       `json:"ok_count"`
	FailedCount   int       `json:"failed_count"`
	FailedNodes   []string  `json:"failed_nodes"`
	AvgResponseMS float64   `json:"avg_response_ms"`
	MaxResponseMS float64   `json:"max_response_ms"`
}

// prepareHost converts the URL to the appropriate format for each check type
func prepareHost(checkType, targetURL string) string {
	switch checkType {
	case "dns":
		// DNS check needs just the domain
		return extractDomain(targetURL)
	case "tcp":
		// TCP check needs domain:port
		return extractDomain(targetURL) + ":443"
	case "redirect":
		// Redirect check always uses apex domain
		return apexURL
	default:
		// HTTP check needs full URL
		return targetURL
	}
}

// extractDomain pulls the domain from a URL
func extractDomain(targetURL string) string {
	// If it's already just a domain, return it
	if !strings.Contains(targetURL, "://") {
		return strings.Split(targetURL, ":")[0] // Remove port if present
	}
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return defaultHost
	}
	return parsed.Host
}

func initiateCheck(checkType, host string, maxNodes int) (string, map[string][]any, error) {
	endpoint := checkEndpoints[checkType]
	apiURL := fmt.Sprintf("%s/%s?host=%s&max_nodes=%d",
		apiBase, endpoint, url.QueryEscape(host), maxNodes)

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

func getResults(requestID, checkType string) ([]Result, error) {
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

		// Parse based on check type
		switch checkType {
		case "dns":
			parseDNSResult(&r, raw)
		case "tcp":
			parseTCPResult(&r, raw)
		case "redirect":
			parseRedirectResult(&r, raw)
		default:
			parseHTTPResult(&r, raw)
		}

		results = append(results, r)
	}

	return results, nil
}

// sortResults sorts results by node name for consistent output
func sortResults(results []Result) {
	sort.Slice(results, func(i, j int) bool {
		return results[i].Node < results[j].Node
	})
}

// parseHTTPResult parses HTTP check results
func parseHTTPResult(r *Result, raw json.RawMessage) {
	// Parse the array structure: [[status, time_or_error, status_text, http_code, ip]]
	var nodeResult [][]any
	if err := json.Unmarshal(raw, &nodeResult); err != nil {
		r.Status = "parse error"
		return
	}

	if len(nodeResult) == 0 || len(nodeResult[0]) < 3 {
		r.Status = "incomplete data"
		return
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
}

// parseRedirectResult parses HTTP check results expecting a 301/302 redirect
func parseRedirectResult(r *Result, raw json.RawMessage) {
	var nodeResult [][]any
	if err := json.Unmarshal(raw, &nodeResult); err != nil {
		r.Status = "parse error"
		return
	}

	if len(nodeResult) == 0 || len(nodeResult[0]) < 3 {
		r.Status = "incomplete data"
		return
	}

	data := nodeResult[0]

	// First element is success indicator (1 = HTTP request succeeded)
	if status, ok := data[0].(float64); ok && status == 1 {
		// Get the HTTP status code from element 3
		var httpCode int
		if len(data) > 3 {
			switch v := data[3].(type) {
			case string:
				// Try to parse string as int
				fmt.Sscanf(v, "%d", &httpCode)
			case float64:
				httpCode = int(v)
			}
		}

		// For redirect check: 301/302 = success, anything else = failure
		if httpCode == 301 || httpCode == 302 {
			r.Success = true
			r.Status = fmt.Sprintf("%d redirect", httpCode)
		} else if httpCode == 200 {
			r.Success = false
			r.Status = "200 (no redirect)"
		} else {
			r.Success = false
			r.Status = fmt.Sprintf("%d (expected 301/302)", httpCode)
		}

		// Get response time
		if t, ok := data[1].(float64); ok {
			r.Time = t
		}
		// Get IP
		if len(data) > 4 {
			if ip, ok := data[4].(string); ok {
				r.IP = ip
			}
		}
	} else {
		// HTTP request failed
		r.Success = false
		if len(data) > 2 {
			if errMsg, ok := data[2].(string); ok {
				r.Status = errMsg
			}
		}
		if r.Status == "" {
			r.Status = "request failed"
		}
	}
}

// parseDNSResult parses DNS check results
// Format: {"A":["104.21.x.x","172.67.x.x"],"AAAA":["2606:4700:..."],"TTL":300}
func parseDNSResult(r *Result, raw json.RawMessage) {
	// DNS results are wrapped in an array
	var wrapper []json.RawMessage
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		r.Status = "parse error"
		return
	}

	if len(wrapper) == 0 {
		r.Status = "no data"
		return
	}

	// Parse the actual DNS result object
	var dnsResult map[string]any
	if err := json.Unmarshal(wrapper[0], &dnsResult); err != nil {
		r.Status = "parse error"
		return
	}

	// Check for A records
	var ips []string
	if aRecords, ok := dnsResult["A"].([]any); ok {
		for _, ip := range aRecords {
			if ipStr, ok := ip.(string); ok {
				ips = append(ips, ipStr)
			}
		}
	}

	if len(ips) > 0 {
		r.Success = true
		r.IP = strings.Join(ips, ", ")
		r.Status = fmt.Sprintf("A: %s", r.IP)

		// Add AAAA if present
		if aaaaRecords, ok := dnsResult["AAAA"].([]any); ok && len(aaaaRecords) > 0 {
			r.Status += " (+AAAA)"
		}
	} else {
		r.Success = false
		r.Status = "no A records"
	}
}

// parseTCPResult parses TCP check results
// Format: [{"address":"ip","time":0.123}] or [{"error":"message"}] on failure
func parseTCPResult(r *Result, raw json.RawMessage) {
	var nodeResult []map[string]any
	if err := json.Unmarshal(raw, &nodeResult); err != nil {
		r.Status = "parse error"
		return
	}

	if len(nodeResult) == 0 {
		r.Status = "incomplete data"
		return
	}

	data := nodeResult[0]

	// Check for error field (failure case)
	if errMsg, ok := data["error"].(string); ok {
		r.Success = false
		r.Status = errMsg
		return
	}

	// Success case: object with address and time
	r.Success = true
	if addr, ok := data["address"].(string); ok {
		r.IP = addr
	}
	if t, ok := data["time"].(float64); ok {
		r.Time = t
		r.Status = fmt.Sprintf("%.0fms", t*1000)
	} else {
		r.Status = "connected"
	}
}

// buildState creates a State from check results
func buildState(checkType string, results []Result) *State {
	state := &State{
		Timestamp:  time.Now().UTC(),
		CheckType:  checkType,
		TotalNodes: len(results),
	}

	var times []float64
	for _, r := range results {
		if r.Pending {
			continue
		}
		if r.Success {
			state.OKCount++
			if r.Time > 0 {
				times = append(times, r.Time*1000) // Convert to ms
			}
		} else {
			state.FailedCount++
			state.FailedNodes = append(state.FailedNodes, r.Node)
		}
	}

	// Calculate response time stats
	if len(times) > 0 {
		var sum float64
		for _, t := range times {
			sum += t
			if t > state.MaxResponseMS {
				state.MaxResponseMS = t
			}
		}
		state.AvgResponseMS = sum / float64(len(times))
	}

	return state
}

// LoadState loads the previous state from disk.
func LoadState() (*State, error) {
	data, err := os.ReadFile(StateFile)
	if err != nil {
		return nil, err
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

// SaveState saves the current state to disk.
func SaveState(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StateFile, data, 0644)
}
