// Package vanity provides functionality for managing Go vanity import packages
// in the Hugo content directory.
package vanityimport

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Package represents a Go package with vanity import metadata.
type Package struct {
	Title            string         `yaml:"title"`
	ImportPath       string         `yaml:"import_path"`
	RepoURL          string         `yaml:"repo_url"`
	Description      string         `yaml:"description"`
	Version          string         `yaml:"version"`
	DocumentationURL string         `yaml:"documentation_url,omitempty"`
	License          string         `yaml:"license"`
	Author           string         `yaml:"author"`
	CreatedAt        time.Time      `yaml:"created_at"`
	UpdatedAt        time.Time      `yaml:"updated_at"`
	HasBinary        bool           `yaml:"has_binary,omitempty"`
	BinaryName       string         `yaml:"binary_name,omitempty"`   // CLI binary name (e.g., "mailerlite")
	TaskfilePath     string         `yaml:"taskfile_path,omitempty"` // Path in repo (e.g., "taskfiles/Taskfile.mailerlite.yml")
	Process          *ProcessConfig `yaml:"process,omitempty"`       // Process configuration for servers
}

// ProcessConfig defines how a package runs as a long-running process.
// This enables automatic process-compose.yaml generation from package metadata.
type ProcessConfig struct {
	// Command is the command to run (e.g., "task mailerlite:server")
	Command string `yaml:"command,omitempty"`

	// Port is the HTTP port the service listens on (e.g., 8086)
	Port int `yaml:"port,omitempty"`

	// HealthPath is the health check endpoint (e.g., "/health")
	HealthPath string `yaml:"health_path,omitempty"`

	// Disabled means the process is defined but not started by default
	Disabled bool `yaml:"disabled,omitempty"`

	// DependsOn lists processes that must start before this one
	DependsOn []string `yaml:"depends_on,omitempty"`

	// Namespace groups related processes (e.g., "servers", "workers")
	Namespace string `yaml:"namespace,omitempty"`
}

// DefaultContentDir is the default location for package content files.
const DefaultContentDir = "content/english/pkg"

// ReadPackage reads a package from a Hugo markdown file.
func ReadPackage(filePath string) (*Package, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Extract frontmatter
	frontmatter, err := extractFrontmatter(data)
	if err != nil {
		return nil, err
	}

	var pkg Package
	if err := yaml.Unmarshal(frontmatter, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	return &pkg, nil
}

// WritePackage writes a package to a Hugo markdown file.
func WritePackage(filePath string, pkg *Package) error {
	// Marshal package to YAML
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(pkg); err != nil {
		return fmt.Errorf("failed to encode package: %w", err)
	}

	// Create markdown content with frontmatter
	content := fmt.Sprintf("---\n%s---\n", buf.String())

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ListPackages returns all package files in the content directory.
func ListPackages(contentDir string) ([]string, error) {
	var packages []string

	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read content directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".md") && name != "_index.md" {
			packages = append(packages, filepath.Join(contentDir, name))
		}
	}

	return packages, nil
}

// GetPackage reads a package by name from the content directory.
func GetPackage(contentDir, name string) (*Package, error) {
	filePath := filepath.Join(contentDir, name+".md")
	return ReadPackage(filePath)
}

// extractFrontmatter extracts YAML frontmatter from markdown content.
func extractFrontmatter(data []byte) ([]byte, error) {
	content := string(data)

	// Check for frontmatter delimiters
	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("no frontmatter found")
	}

	// Find the closing delimiter
	endIndex := strings.Index(content[4:], "\n---")
	if endIndex == -1 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}

	// Extract frontmatter content (without delimiters)
	frontmatter := content[4 : endIndex+4]
	return []byte(frontmatter), nil
}
