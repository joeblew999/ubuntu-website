package gdocs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// APIClient manages Docs via Google Docs API
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

// Get gets a document by ID
func (c *APIClient) Get(documentID string) (*Document, error) {
	apiURL := fmt.Sprintf("https://docs.googleapis.com/v1/documents/%s",
		url.PathEscape(documentID))

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

	var doc Document
	if err := json.Unmarshal(respBody, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// Create creates a new document
func (c *APIClient) Create(title string) (*CreateResult, error) {
	apiURL := "https://docs.googleapis.com/v1/documents"

	body := map[string]string{"title": title}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return &CreateResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &CreateResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &CreateResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("create failed")
	}

	var doc Document
	if err := json.Unmarshal(respBody, &doc); err != nil {
		return &CreateResult{Success: false, Error: err.Error()}, err
	}

	return &CreateResult{
		Success:  true,
		Document: &doc,
	}, nil
}

// InsertText inserts text at a specific index
func (c *APIClient) InsertText(documentID string, index int, text string) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://docs.googleapis.com/v1/documents/%s:batchUpdate",
		url.PathEscape(documentID))

	body := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"insertText": map[string]interface{}{
					"location": map[string]int{
						"index": index,
					},
					"text": text,
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
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
		}, fmt.Errorf("insert failed")
	}

	var apiResp struct {
		DocumentID string `json:"documentId"`
		Replies    []struct {
			InsertText struct{} `json:"insertText"`
		} `json:"replies"`
		WriteControl struct {
			RequiredRevisionID string `json:"requiredRevisionId"`
		} `json:"writeControl"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return &UpdateResult{Success: false, Error: err.Error()}, err
	}

	return &UpdateResult{
		Success:  true,
		Revision: apiResp.WriteControl.RequiredRevisionID,
	}, nil
}

// AppendText appends text to the end of the document
func (c *APIClient) AppendText(documentID, text string) (*UpdateResult, error) {
	// First get the document to find the end index
	doc, err := c.Get(documentID)
	if err != nil {
		return &UpdateResult{Success: false, Error: err.Error()}, err
	}

	// Find the end index (last element's end index - 1 to insert before final newline)
	endIndex := 1
	if doc.Body != nil && len(doc.Body.Content) > 0 {
		lastElem := doc.Body.Content[len(doc.Body.Content)-1]
		endIndex = lastElem.EndIndex - 1
		if endIndex < 1 {
			endIndex = 1
		}
	}

	return c.InsertText(documentID, endIndex, text)
}

// ReplaceText replaces all occurrences of a string
func (c *APIClient) ReplaceText(documentID, find, replace string) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://docs.googleapis.com/v1/documents/%s:batchUpdate",
		url.PathEscape(documentID))

	body := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"replaceAllText": map[string]interface{}{
					"containsText": map[string]interface{}{
						"text":      find,
						"matchCase": true,
					},
					"replaceText": replace,
				},
			},
		},
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
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
		}, fmt.Errorf("replace failed")
	}

	return &UpdateResult{Success: true}, nil
}
