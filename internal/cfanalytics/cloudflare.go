// cloudflare.go - Cloudflare Web Analytics API client.
package cfanalytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	cfGraphQLEndpoint = "https://api.cloudflare.com/client/v4/graphql"
	StateFile         = ".analytics-state.json"
	ChangeThreshold   = 0.20 // 20% change triggers alert

	// Default values (fallbacks if env vars not set)
	defaultAccountTag = "7384af54e33b8a54ff240371ea368440"
	defaultSiteTag    = "4c28a6bfb5514996914a603c999d5c79"
)

// State represents the stored analytics state from previous run
type State struct {
	Timestamp time.Time        `json:"timestamp"`
	Period    string           `json:"period"`
	Visits    int64            `json:"visits"`
	PageViews int64            `json:"pageviews"`
	TopPages  map[string]int64 `json:"top_pages"`
	Countries map[string]int64 `json:"countries"`
}

// Report contains the generated analytics report
type Report struct {
	Summary    string
	Changes    []string
	HasChanges bool
}

// GraphQL query for Cloudflare Web Analytics
const analyticsQuery = `
query WebAnalytics($accountTag: string!, $filter: AccountRumPageloadEventsAdaptiveGroupsFilter_InputObject!) {
  viewer {
    accounts(filter: {accountTag: $accountTag}) {
      rumPageloadEventsAdaptiveGroups(
        filter: $filter
        limit: 5000
      ) {
        sum {
          visits
        }
        count
        dimensions {
          requestPath
          countryName
        }
      }
    }
  }
}
`

// GraphQL response structures
type graphQLResponse struct {
	Data   responseData   `json:"data"`
	Errors []graphQLError `json:"errors,omitempty"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type responseData struct {
	Viewer viewer `json:"viewer"`
}

type viewer struct {
	Accounts []account `json:"accounts"`
}

type account struct {
	RumGroups []rumGroup `json:"rumPageloadEventsAdaptiveGroups"`
}

type rumGroup struct {
	Sum        sumData    `json:"sum"`
	Dimensions dimensions `json:"dimensions"`
	Count      int64      `json:"count"`
}

type sumData struct {
	Visits int64 `json:"visits"`
}

type dimensions struct {
	RequestPath string `json:"requestPath"`
	CountryName string `json:"countryName"`
}

// GetConfig returns Cloudflare account and site tags from environment variables,
// falling back to defaults for backward compatibility.
func GetConfig() (accountTag, siteTag string) {
	accountTag = os.Getenv("CF_ACCOUNT_ID")
	if accountTag == "" {
		accountTag = defaultAccountTag
	}
	siteTag = os.Getenv("CF_WEB_ANALYTICS_SITE_TAG")
	if siteTag == "" {
		siteTag = defaultSiteTag
	}
	return
}

// FetchAnalytics retrieves analytics data from Cloudflare for the given date range.
func FetchAnalytics(token string, since, until time.Time) (*State, error) {
	// Get config from environment (with fallbacks)
	accountTag, siteTag := GetConfig()

	// Build GraphQL request with proper filter structure
	filter := map[string]any{
		"AND": []map[string]any{
			{
				"datetime_geq": since.Format(time.RFC3339),
				"datetime_leq": until.Format(time.RFC3339),
			},
			{"bot": 0}, // Exclude bots
			{
				"OR": []map[string]any{
					{"siteTag": siteTag},
				},
			},
		},
	}

	reqBody := map[string]any{
		"query": analyticsQuery,
		"variables": map[string]any{
			"accountTag": accountTag,
			"filter":     filter,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", cfGraphQLEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var gqlResp graphQLResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		errMsg := gqlResp.Errors[0].Message
		if strings.Contains(errMsg, "not authorized") {
			return nil, fmt.Errorf("not authorized - ensure your API token has 'Account Analytics:Read' permission\nCreate/edit token at: https://dash.cloudflare.com/profile/api-tokens")
		}
		return nil, fmt.Errorf("GraphQL error: %s", errMsg)
	}

	// Aggregate results
	state := &State{
		Timestamp: time.Now().UTC(),
		TopPages:  make(map[string]int64),
		Countries: make(map[string]int64),
	}

	if len(gqlResp.Data.Viewer.Accounts) == 0 {
		return state, nil // No data
	}

	for _, group := range gqlResp.Data.Viewer.Accounts[0].RumGroups {
		state.Visits += group.Sum.Visits
		state.PageViews += group.Count // Count is pageviews in this API

		if group.Dimensions.RequestPath != "" {
			state.TopPages[group.Dimensions.RequestPath] += group.Count
		}
		if group.Dimensions.CountryName != "" {
			state.Countries[group.Dimensions.CountryName] += group.Count
		}
	}

	return state, nil
}

// LoadState loads the previous analytics state from disk.
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

// SaveState saves the current analytics state to disk.
func SaveState(state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(StateFile, data, 0644)
}

// PostToWebhook posts the report to a Slack/Discord webhook.
func PostToWebhook(url string, report Report) error {
	// Build webhook payload (works for Slack/Discord)
	payload := map[string]any{
		"text": fmt.Sprintf("*Analytics Alert*\n%s\n\n*Changes:*\n%s",
			report.Summary,
			strings.Join(report.Changes, "\n")),
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// kv is a key-value pair for sorting maps.
type kv struct {
	Key   string
	Value int64
}

// sortMapByValue returns the top N entries from a map, sorted by value descending.
func sortMapByValue(m map[string]int64, limit int) []kv {
	var sorted []kv
	for k, v := range m {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}
