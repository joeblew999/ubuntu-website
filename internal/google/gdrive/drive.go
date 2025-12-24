// Package gdrive provides Google Drive file management via API.
//
// Supports:
//   - List files and folders
//   - Upload/download files
//   - Create folders
//   - File permissions and sharing
//
// Example:
//
//	config := gdrive.DefaultConfig()
//	client, err := gdrive.NewAPIClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	files, err := client.List("root", 10)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, f := range files.Files {
//	    fmt.Printf("%s (%s)\n", f.Name, f.ID)
//	}
package gdrive

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds Drive configuration
type Config struct {
	TokenPath string // Path to google-mcp-server token directory
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		TokenPath: filepath.Join(home, ".google-mcp-accounts"),
	}
}

// File represents a Google Drive file or folder
type File struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	MimeType     string   `json:"mimeType"`
	Size         string   `json:"size,omitempty"` // String because Google API returns it as string
	CreatedTime  string   `json:"createdTime,omitempty"`
	ModifiedTime string   `json:"modifiedTime,omitempty"`
	Parents      []string `json:"parents,omitempty"`
	WebViewLink  string   `json:"webViewLink,omitempty"`
}

// IsFolder returns true if the file is a folder
func (f *File) IsFolder() bool {
	return f.MimeType == MimeTypeFolder
}

// ListResult contains the result of a list operation
type ListResult struct {
	Success       bool    `json:"success"`
	Files         []*File `json:"files,omitempty"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
	Error         string  `json:"error,omitempty"`
}

// UploadResult contains the result of an upload operation
type UploadResult struct {
	Success bool   `json:"success"`
	File    *File  `json:"file,omitempty"`
	Error   string `json:"error,omitempty"`
}

// DownloadResult contains the result of a download operation
type DownloadResult struct {
	Success  bool   `json:"success"`
	Content  []byte `json:"content,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Common MIME types
const (
	MimeTypeFolder       = "application/vnd.google-apps.folder"
	MimeTypeDocument     = "application/vnd.google-apps.document"
	MimeTypeSpreadsheet  = "application/vnd.google-apps.spreadsheet"
	MimeTypePresentation = "application/vnd.google-apps.presentation"
	MimeTypePDF          = "application/pdf"
	MimeTypeText         = "text/plain"
	MimeTypeJSON         = "application/json"
)

// Validate validates the config
func (c *Config) Validate() error {
	if c.TokenPath == "" {
		return fmt.Errorf("token path is required")
	}
	return nil
}
