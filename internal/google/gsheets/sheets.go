// Package gsheets provides Google Sheets spreadsheet management via API.
//
// Supports:
//   - Read cell values and ranges
//   - Write/update cell values
//   - Append rows
//   - Get spreadsheet metadata
//
// Example:
//
//	config := gsheets.DefaultConfig()
//	client, err := gsheets.NewAPIClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	values, err := client.GetValues(spreadsheetID, "Sheet1!A1:D10")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, row := range values.Values {
//	    fmt.Println(row)
//	}
package gsheets

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds Sheets configuration
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

// Spreadsheet represents spreadsheet metadata
type Spreadsheet struct {
	ID         string   `json:"spreadsheetId"`
	Title      string   `json:"title"`
	Locale     string   `json:"locale,omitempty"`
	TimeZone   string   `json:"timeZone,omitempty"`
	SheetNames []string `json:"sheetNames,omitempty"`
}

// ValueRange represents a range of cell values
type ValueRange struct {
	Range  string          `json:"range"`
	Values [][]interface{} `json:"values"`
}

// GetResult contains the result of a get values operation
type GetResult struct {
	Success bool        `json:"success"`
	Range   string      `json:"range,omitempty"`
	Values  [][]interface{} `json:"values,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// UpdateResult contains the result of an update operation
type UpdateResult struct {
	Success         bool   `json:"success"`
	UpdatedRange    string `json:"updatedRange,omitempty"`
	UpdatedRows     int    `json:"updatedRows,omitempty"`
	UpdatedColumns  int    `json:"updatedColumns,omitempty"`
	UpdatedCells    int    `json:"updatedCells,omitempty"`
	Error           string `json:"error,omitempty"`
}

// AppendResult contains the result of an append operation
type AppendResult struct {
	Success      bool   `json:"success"`
	UpdatedRange string `json:"updatedRange,omitempty"`
	UpdatedRows  int    `json:"updatedRows,omitempty"`
	Error        string `json:"error,omitempty"`
}

// Validate validates the config
func (c *Config) Validate() error {
	if c.TokenPath == "" {
		return fmt.Errorf("token path is required")
	}
	return nil
}
