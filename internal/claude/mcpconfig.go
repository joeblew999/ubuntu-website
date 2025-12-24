// MCP (Model Context Protocol) configuration management.
// Used for managing .vscode/mcp.json and .claude/settings.json files.
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SchemaURL is the JSON schema for MCP config files
const SchemaURL = "https://modelcontextprotocol.io/schema/config.json"

// Config represents the mcp.json file structure
type Config struct {
	Schema     string            `json:"$schema,omitempty"`
	MCPServers map[string]Server `json:"mcpServers"`
}

// Server represents an MCP server configuration
type Server struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

// Location represents a config file location
type Location struct {
	Dir  string // Subdirectory (empty = project root)
	File string // Filename
}

// Standard config locations
var Locations = map[string]Location{
	"vscode":  {Dir: ".vscode", File: "mcp.json"},
	"project": {Dir: "", File: ".mcp.json"},
	"claude":  {Dir: ".claude", File: "mcp.json"},
}

// LoadConfig loads an mcp.json file, returning empty config if not exists
func LoadConfig(mcpFile string) (*Config, error) {
	config := &Config{
		Schema:     SchemaURL,
		MCPServers: make(map[string]Server),
	}

	data, err := os.ReadFile(mcpFile)
	if os.IsNotExist(err) {
		return config, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", mcpFile, err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", mcpFile, err)
	}

	if config.MCPServers == nil {
		config.MCPServers = make(map[string]Server)
	}

	// Ensure schema is set
	if config.Schema == "" {
		config.Schema = SchemaURL
	}

	return config, nil
}

// SaveConfig saves the config to mcp.json with pretty formatting
func SaveConfig(mcpFile string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	if err := os.WriteFile(mcpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", mcpFile, err)
	}

	return nil
}

// SaveConfigSafe saves config using atomic write (temp file + rename)
// with JSON validation. Use this for VSCode configs to prevent corruption.
func SaveConfigSafe(mcpFile string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	// Use safe write (validates JSON, atomic rename)
	return safeWriteFile(mcpFile, data, 0644)
}

// AddServer adds a server to the config
func (c *Config) AddServer(name string, server Server) {
	c.MCPServers[name] = server
}

// RemoveServer removes a server from the config
func (c *Config) RemoveServer(name string) bool {
	if _, exists := c.MCPServers[name]; exists {
		delete(c.MCPServers, name)
		return true
	}
	return false
}

// HasServer checks if a server exists
func (c *Config) HasServer(name string) bool {
	_, exists := c.MCPServers[name]
	return exists
}

// IsEmpty returns true if no servers are configured
func (c *Config) IsEmpty() bool {
	return len(c.MCPServers) == 0
}

// GetConfigPath returns the full path to an MCP config file
func GetConfigPath(projectRoot, locationName string) (string, error) {
	loc, ok := Locations[locationName]
	if !ok {
		return "", fmt.Errorf("unknown location: %s (valid: vscode, project, claude)", locationName)
	}

	if loc.Dir != "" {
		return filepath.Join(projectRoot, loc.Dir, loc.File), nil
	}
	return filepath.Join(projectRoot, loc.File), nil
}

// EnsureDir creates the directory for the config file if needed
func EnsureDir(mcpFile string) error {
	dir := filepath.Dir(mcpFile)
	return os.MkdirAll(dir, 0755)
}
