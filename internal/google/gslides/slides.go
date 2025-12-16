// Package gslides provides Google Slides presentation management via API.
//
// Supports:
//   - Get presentation metadata
//   - Create new presentations
//   - Add/remove slides
//   - Insert text and shapes
//
// Example:
//
//	config := gslides.DefaultConfig()
//	client, err := gslides.NewAPIClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	pres, err := client.Get(presentationID)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Title: %s (%d slides)\n", pres.Title, len(pres.Slides))
package gslides

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds Slides configuration
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

// Presentation represents a Google Slides presentation
type Presentation struct {
	ID     string   `json:"presentationId"`
	Title  string   `json:"title"`
	Locale string   `json:"locale,omitempty"`
	Slides []*Slide `json:"slides,omitempty"`
}

// Slide represents a single slide
type Slide struct {
	ObjectID      string         `json:"objectId"`
	PageElements  []*PageElement `json:"pageElements,omitempty"`
	SlideProperties *SlideProperties `json:"slideProperties,omitempty"`
}

// SlideProperties contains slide metadata
type SlideProperties struct {
	LayoutObjectID string `json:"layoutObjectId,omitempty"`
}

// PageElement represents an element on a slide
type PageElement struct {
	ObjectID string `json:"objectId"`
	Size     *Size  `json:"size,omitempty"`
	Transform *Transform `json:"transform,omitempty"`
	Shape    *Shape `json:"shape,omitempty"`
}

// Size represents dimensions
type Size struct {
	Width  *Dimension `json:"width,omitempty"`
	Height *Dimension `json:"height,omitempty"`
}

// Dimension represents a measurement
type Dimension struct {
	Magnitude float64 `json:"magnitude"`
	Unit      string  `json:"unit"` // EMU, PT, etc.
}

// Transform represents position/rotation
type Transform struct {
	ScaleX     float64 `json:"scaleX,omitempty"`
	ScaleY     float64 `json:"scaleY,omitempty"`
	TranslateX float64 `json:"translateX,omitempty"`
	TranslateY float64 `json:"translateY,omitempty"`
	Unit       string  `json:"unit,omitempty"`
}

// Shape represents a shape element
type Shape struct {
	ShapeType string `json:"shapeType,omitempty"`
	Text      *TextContent `json:"text,omitempty"`
}

// TextContent represents text within a shape
type TextContent struct {
	TextElements []*TextElement `json:"textElements,omitempty"`
}

// TextElement represents a text run
type TextElement struct {
	TextRun *TextRun `json:"textRun,omitempty"`
}

// TextRun represents styled text
type TextRun struct {
	Content string `json:"content"`
}

// CreateResult contains the result of a create operation
type CreateResult struct {
	Success      bool          `json:"success"`
	Presentation *Presentation `json:"presentation,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// UpdateResult contains the result of an update operation
type UpdateResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// Validate validates the config
func (c *Config) Validate() error {
	if c.TokenPath == "" {
		return fmt.Errorf("token path is required")
	}
	return nil
}
