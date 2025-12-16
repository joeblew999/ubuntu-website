package gsheets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// APIClient manages Sheets via Google Sheets API
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

// GetSpreadsheet gets spreadsheet metadata
func (c *APIClient) GetSpreadsheet(spreadsheetID string) (*Spreadsheet, error) {
	apiURL := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s?fields=spreadsheetId,properties(title,locale,timeZone),sheets(properties(title))",
		url.PathEscape(spreadsheetID))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp struct {
		SpreadsheetID string `json:"spreadsheetId"`
		Properties    struct {
			Title    string `json:"title"`
			Locale   string `json:"locale"`
			TimeZone string `json:"timeZone"`
		} `json:"properties"`
		Sheets []struct {
			Properties struct {
				Title string `json:"title"`
			} `json:"properties"`
		} `json:"sheets"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, err
	}

	sheetNames := make([]string, len(apiResp.Sheets))
	for i, s := range apiResp.Sheets {
		sheetNames[i] = s.Properties.Title
	}

	return &Spreadsheet{
		ID:         apiResp.SpreadsheetID,
		Title:      apiResp.Properties.Title,
		Locale:     apiResp.Properties.Locale,
		TimeZone:   apiResp.Properties.TimeZone,
		SheetNames: sheetNames,
	}, nil
}

// GetValues gets cell values from a range
func (c *APIClient) GetValues(spreadsheetID, rangeA1 string) (*GetResult, error) {
	apiURL := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s",
		url.PathEscape(spreadsheetID), url.PathEscape(rangeA1))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return &GetResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &GetResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &GetResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	var apiResp struct {
		Range  string          `json:"range"`
		Values [][]interface{} `json:"values"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return &GetResult{Success: false, Error: err.Error()}, err
	}

	return &GetResult{
		Success: true,
		Range:   apiResp.Range,
		Values:  apiResp.Values,
	}, nil
}

// UpdateValues updates cell values in a range
func (c *APIClient) UpdateValues(spreadsheetID, rangeA1 string, values [][]interface{}) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?valueInputOption=USER_ENTERED",
		url.PathEscape(spreadsheetID), url.PathEscape(rangeA1))

	body := map[string]interface{}{
		"range":  rangeA1,
		"values": values,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return &UpdateResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &UpdateResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &UpdateResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	var apiResp struct {
		UpdatedRange   string `json:"updatedRange"`
		UpdatedRows    int    `json:"updatedRows"`
		UpdatedColumns int    `json:"updatedColumns"`
		UpdatedCells   int    `json:"updatedCells"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return &UpdateResult{Success: false, Error: err.Error()}, err
	}

	return &UpdateResult{
		Success:        true,
		UpdatedRange:   apiResp.UpdatedRange,
		UpdatedRows:    apiResp.UpdatedRows,
		UpdatedColumns: apiResp.UpdatedColumns,
		UpdatedCells:   apiResp.UpdatedCells,
	}, nil
}

// AppendValues appends rows to a sheet
func (c *APIClient) AppendValues(spreadsheetID, rangeA1 string, values [][]interface{}) (*AppendResult, error) {
	apiURL := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s:append?valueInputOption=USER_ENTERED&insertDataOption=INSERT_ROWS",
		url.PathEscape(spreadsheetID), url.PathEscape(rangeA1))

	body := map[string]interface{}{
		"values": values,
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return &AppendResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &AppendResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &AppendResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	var apiResp struct {
		Updates struct {
			UpdatedRange string `json:"updatedRange"`
			UpdatedRows  int    `json:"updatedRows"`
		} `json:"updates"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return &AppendResult{Success: false, Error: err.Error()}, err
	}

	return &AppendResult{
		Success:      true,
		UpdatedRange: apiResp.Updates.UpdatedRange,
		UpdatedRows:  apiResp.Updates.UpdatedRows,
	}, nil
}

// Clear clears values from a range
func (c *APIClient) Clear(spreadsheetID, rangeA1 string) error {
	apiURL := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s:clear",
		url.PathEscape(spreadsheetID), url.PathEscape(rangeA1))

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
