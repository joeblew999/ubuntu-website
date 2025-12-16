package gdrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/joeblew999/ubuntu-website/internal/googleauth"
)

// APIClient manages Drive via Google Drive API
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

// List lists files in a folder
func (c *APIClient) List(folderID string, maxResults int) (*ListResult, error) {
	if folderID == "" {
		folderID = "root"
	}
	if maxResults <= 0 {
		maxResults = 10
	}

	params := url.Values{}
	params.Set("q", fmt.Sprintf("'%s' in parents and trashed = false", folderID))
	params.Set("pageSize", fmt.Sprintf("%d", maxResults))
	params.Set("fields", "nextPageToken, files(id, name, mimeType, size, createdTime, modifiedTime, parents, webViewLink)")
	params.Set("orderBy", "folder, name")

	apiURL := "https://www.googleapis.com/drive/v3/files?" + params.Encode()

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return &ListResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &ListResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &ListResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("API error: %s", string(respBody))
	}

	var apiResp struct {
		Files         []*File `json:"files"`
		NextPageToken string  `json:"nextPageToken"`
	}

	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return &ListResult{Success: false, Error: err.Error()}, err
	}

	return &ListResult{
		Success:       true,
		Files:         apiResp.Files,
		NextPageToken: apiResp.NextPageToken,
	}, nil
}

// Get gets file metadata by ID
func (c *APIClient) Get(fileID string) (*File, error) {
	params := url.Values{}
	params.Set("fields", "id, name, mimeType, size, createdTime, modifiedTime, parents, webViewLink")

	apiURL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s?%s",
		url.PathEscape(fileID), params.Encode())

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

	var file File
	if err := json.Unmarshal(respBody, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

// Download downloads file content
func (c *APIClient) Download(fileID string) (*DownloadResult, error) {
	apiURL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s?alt=media",
		url.PathEscape(fileID))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return &DownloadResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &DownloadResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return &DownloadResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(body)),
		}, fmt.Errorf("download failed")
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return &DownloadResult{Success: false, Error: err.Error()}, err
	}

	return &DownloadResult{
		Success:  true,
		Content:  content,
		MimeType: resp.Header.Get("Content-Type"),
	}, nil
}

// Upload uploads a file to Drive
func (c *APIClient) Upload(name string, content []byte, mimeType, parentID string) (*UploadResult, error) {
	// Create metadata
	metadata := map[string]interface{}{
		"name": name,
	}
	if parentID != "" {
		metadata["parents"] = []string{parentID}
	}

	metadataJSON, _ := json.Marshal(metadata)

	// Create multipart request
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add metadata part
	metadataPart, _ := writer.CreateFormField("metadata")
	metadataPart.Write(metadataJSON)

	// Add file part
	filePart, _ := writer.CreateFormFile("file", name)
	filePart.Write(content)

	writer.Close()

	apiURL := "https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart&fields=id,name,mimeType,webViewLink"

	req, err := http.NewRequest("POST", apiURL, &buf)
	if err != nil {
		return &UploadResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &UploadResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("upload failed")
	}

	var file File
	if err := json.Unmarshal(respBody, &file); err != nil {
		return &UploadResult{Success: false, Error: err.Error()}, err
	}

	return &UploadResult{
		Success: true,
		File:    &file,
	}, nil
}

// CreateFolder creates a new folder
func (c *APIClient) CreateFolder(name, parentID string) (*UploadResult, error) {
	metadata := map[string]interface{}{
		"name":     name,
		"mimeType": MimeTypeFolder,
	}
	if parentID != "" {
		metadata["parents"] = []string{parentID}
	}

	metadataJSON, _ := json.Marshal(metadata)

	apiURL := "https://www.googleapis.com/drive/v3/files?fields=id,name,mimeType,webViewLink"

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(metadataJSON))
	if err != nil {
		return &UploadResult{Success: false, Error: err.Error()}, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &UploadResult{Success: false, Error: err.Error()}, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return &UploadResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, fmt.Errorf("create folder failed")
	}

	var file File
	if err := json.Unmarshal(respBody, &file); err != nil {
		return &UploadResult{Success: false, Error: err.Error()}, err
	}

	return &UploadResult{
		Success: true,
		File:    &file,
	}, nil
}

// Delete deletes a file (moves to trash)
func (c *APIClient) Delete(fileID string) error {
	apiURL := fmt.Sprintf("https://www.googleapis.com/drive/v3/files/%s",
		url.PathEscape(fileID))

	req, err := http.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
