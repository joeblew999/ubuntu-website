// Package gdocs provides Google Docs document management via API.
//
// Supports:
//   - Get document content
//   - Create new documents
//   - Insert/update text
//   - Basic formatting
//
// Example:
//
//	config := gdocs.DefaultConfig()
//	client, err := gdocs.NewAPIClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	doc, err := client.Get(documentID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Title: %s\n", doc.Title)
package gdocs

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds Docs configuration
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

// Document represents a Google Doc
type Document struct {
	ID       string `json:"documentId"`
	Title    string `json:"title"`
	Body     *Body  `json:"body,omitempty"`
	Revision string `json:"revisionId,omitempty"`
}

// Body represents document body content
type Body struct {
	Content []*StructuralElement `json:"content,omitempty"`
}

// StructuralElement represents a document element
type StructuralElement struct {
	StartIndex int        `json:"startIndex"`
	EndIndex   int        `json:"endIndex"`
	Paragraph  *Paragraph `json:"paragraph,omitempty"`
}

// Paragraph represents a paragraph
type Paragraph struct {
	Elements []*ParagraphElement `json:"elements,omitempty"`
}

// ParagraphElement represents an element within a paragraph
type ParagraphElement struct {
	StartIndex int      `json:"startIndex"`
	EndIndex   int      `json:"endIndex"`
	TextRun    *TextRun `json:"textRun,omitempty"`
}

// TextRun represents a run of text
type TextRun struct {
	Content string `json:"content"`
}

// CreateResult contains the result of a create operation
type CreateResult struct {
	Success  bool      `json:"success"`
	Document *Document `json:"document,omitempty"`
	Error    string    `json:"error,omitempty"`
}

// UpdateResult contains the result of an update operation
type UpdateResult struct {
	Success  bool   `json:"success"`
	Revision string `json:"revision,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Validate validates the config
func (c *Config) Validate() error {
	if c.TokenPath == "" {
		return fmt.Errorf("token path is required")
	}
	return nil
}

// GetText extracts plain text from a document
func (d *Document) GetText() string {
	if d.Body == nil {
		return ""
	}

	var text string
	for _, elem := range d.Body.Content {
		if elem.Paragraph != nil {
			for _, pe := range elem.Paragraph.Elements {
				if pe.TextRun != nil {
					text += pe.TextRun.Content
				}
			}
		}
	}
	return text
}
