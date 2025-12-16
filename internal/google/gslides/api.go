package gslides

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// APIClient manages Slides via Google Slides API
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

// Get gets a presentation by ID
func (c *APIClient) Get(presentationID string) (*Presentation, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s",
		url.PathEscape(presentationID))

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

	var pres Presentation
	if err := json.Unmarshal(respBody, &pres); err != nil {
		return nil, err
	}

	return &pres, nil
}

// Create creates a new presentation
func (c *APIClient) Create(title string) (*CreateResult, error) {
	apiURL := "https://slides.googleapis.com/v1/presentations"

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

	var pres Presentation
	if err := json.Unmarshal(respBody, &pres); err != nil {
		return &CreateResult{Success: false, Error: err.Error()}, err
	}

	return &CreateResult{
		Success:      true,
		Presentation: &pres,
	}, nil
}

// AddSlide adds a new slide to the presentation
func (c *APIClient) AddSlide(presentationID string, insertionIndex int) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s:batchUpdate",
		url.PathEscape(presentationID))

	body := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"createSlide": map[string]interface{}{
					"insertionIndex": insertionIndex,
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
		}, fmt.Errorf("add slide failed")
	}

	return &UpdateResult{Success: true}, nil
}

// DeleteSlide removes a slide from the presentation
func (c *APIClient) DeleteSlide(presentationID, slideID string) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s:batchUpdate",
		url.PathEscape(presentationID))

	body := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"deleteObject": map[string]interface{}{
					"objectId": slideID,
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
		}, fmt.Errorf("delete slide failed")
	}

	return &UpdateResult{Success: true}, nil
}

// InsertText inserts text into a shape on a slide
func (c *APIClient) InsertText(presentationID, objectID string, text string) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s:batchUpdate",
		url.PathEscape(presentationID))

	body := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"insertText": map[string]interface{}{
					"objectId":       objectID,
					"insertionIndex": 0,
					"text":           text,
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
		}, fmt.Errorf("insert text failed")
	}

	return &UpdateResult{Success: true}, nil
}

// CreateTextBox creates a text box on a slide
func (c *APIClient) CreateTextBox(presentationID, slideID string, x, y, width, height float64) (*UpdateResult, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s:batchUpdate",
		url.PathEscape(presentationID))

	// Dimensions in EMU (English Metric Units) - 914400 EMU = 1 inch
	emuPerPoint := 12700.0

	body := map[string]interface{}{
		"requests": []map[string]interface{}{
			{
				"createShape": map[string]interface{}{
					"shapeType": "TEXT_BOX",
					"elementProperties": map[string]interface{}{
						"pageObjectId": slideID,
						"size": map[string]interface{}{
							"width":  map[string]interface{}{"magnitude": width * emuPerPoint, "unit": "EMU"},
							"height": map[string]interface{}{"magnitude": height * emuPerPoint, "unit": "EMU"},
						},
						"transform": map[string]interface{}{
							"scaleX":     1,
							"scaleY":     1,
							"translateX": x * emuPerPoint,
							"translateY": y * emuPerPoint,
							"unit":       "EMU",
						},
					},
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
		}, fmt.Errorf("create text box failed")
	}

	return &UpdateResult{Success: true}, nil
}
