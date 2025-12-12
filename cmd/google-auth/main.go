// Google MCP configuration tool for Claude Code
//
// Manages the .mcp.json file in the project root to add/remove
// the Google MCP server configuration.
//
// Usage:
//
//	google-auth add      # Add google server to .mcp.json
//	google-auth remove   # Remove google server from .mcp.json
//	google-auth status   # Show current .mcp.json config
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// MCPConfig represents the .mcp.json file structure
type MCPConfig struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

const (
	mcpFileName  = ".mcp.json"
	serverName   = "google"
	serverCmd    = "google-mcp-server"
	envClientID  = "${GOOGLE_CLIENT_ID}"
	envSecretID  = "${GOOGLE_CLIENT_SECRET}"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	// Find project root (where .mcp.json should live)
	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Printf("Error finding project root: %v\n", err)
		os.Exit(1)
	}

	mcpFile := filepath.Join(projectRoot, mcpFileName)

	switch cmd {
	case "add":
		if err := addServer(mcpFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "remove":
		if err := removeServer(mcpFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "status":
		if err := showStatus(mcpFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Google MCP Configuration Tool")
	fmt.Println("")
	fmt.Println("Manages the .mcp.json file for Claude Code project-level MCP servers.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  google-auth add      Add google server to .mcp.json")
	fmt.Println("  google-auth remove   Remove google server from .mcp.json")
	fmt.Println("  google-auth status   Show current .mcp.json config")
	fmt.Println("")
	fmt.Println("The tool writes to .mcp.json in the project root (safe, version-controlled).")
	fmt.Println("Environment variables are referenced as ${VAR} - no secrets stored in file.")
}

// findProjectRoot looks for the project root by finding go.mod or .git
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check for go.mod or .git
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root, use current working directory
			return os.Getwd()
		}
		dir = parent
	}
}

// loadConfig loads the .mcp.json file, returning empty config if not exists
func loadConfig(mcpFile string) (*MCPConfig, error) {
	config := &MCPConfig{
		MCPServers: make(map[string]MCPServer),
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
		config.MCPServers = make(map[string]MCPServer)
	}

	return config, nil
}

// saveConfig saves the config to .mcp.json with pretty formatting
func saveConfig(mcpFile string, config *MCPConfig) error {
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

// addServer adds the google MCP server to .mcp.json
func addServer(mcpFile string) error {
	config, err := loadConfig(mcpFile)
	if err != nil {
		return err
	}

	// Check if already exists
	if _, exists := config.MCPServers[serverName]; exists {
		fmt.Printf("✅ Google MCP server already configured in %s\n", mcpFile)
		return nil
	}

	// Add the server
	config.MCPServers[serverName] = MCPServer{
		Command: serverCmd,
		Args:    []string{},
		Env: map[string]string{
			"GOOGLE_CLIENT_ID":     envClientID,
			"GOOGLE_CLIENT_SECRET": envSecretID,
		},
	}

	if err := saveConfig(mcpFile, config); err != nil {
		return err
	}

	fmt.Println("✅ Google MCP server added to .mcp.json")
	fmt.Println("")
	fmt.Printf("File: %s\n", mcpFile)
	fmt.Println("")
	fmt.Println("The config uses ${GOOGLE_CLIENT_ID} and ${GOOGLE_CLIENT_SECRET}")
	fmt.Println("environment variables - no secrets stored in the file.")
	fmt.Println("")
	fmt.Println("Make sure your .env file has these variables set.")
	fmt.Println("Restart Claude Code to pick up the new MCP server.")

	return nil
}

// removeServer removes the google MCP server from .mcp.json
func removeServer(mcpFile string) error {
	config, err := loadConfig(mcpFile)
	if err != nil {
		return err
	}

	// Check if exists
	if _, exists := config.MCPServers[serverName]; !exists {
		fmt.Printf("ℹ️  Google MCP server not found in %s\n", mcpFile)
		return nil
	}

	// Remove the server
	delete(config.MCPServers, serverName)

	// If no servers left, remove the file
	if len(config.MCPServers) == 0 {
		if err := os.Remove(mcpFile); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", mcpFile, err)
		}
		fmt.Println("✅ Google MCP server removed")
		fmt.Printf("   Deleted %s (no servers remaining)\n", mcpFile)
		return nil
	}

	if err := saveConfig(mcpFile, config); err != nil {
		return err
	}

	fmt.Println("✅ Google MCP server removed from .mcp.json")
	fmt.Println("")
	fmt.Println("Restart Claude Code to apply changes.")

	return nil
}

// showStatus shows the current .mcp.json configuration
func showStatus(mcpFile string) error {
	fmt.Println("=== Google MCP Configuration Status ===")
	fmt.Println("")

	// Check if file exists
	if _, err := os.Stat(mcpFile); os.IsNotExist(err) {
		fmt.Printf("❌ No .mcp.json found at %s\n", mcpFile)
		fmt.Println("")
		fmt.Println("Run: google-auth add")
		return nil
	}

	config, err := loadConfig(mcpFile)
	if err != nil {
		return err
	}

	fmt.Printf("File: %s\n", mcpFile)
	fmt.Println("")

	// Check for google server
	if server, exists := config.MCPServers[serverName]; exists {
		fmt.Println("✅ Google MCP server configured:")
		fmt.Printf("   Command: %s\n", server.Command)
		if len(server.Env) > 0 {
			fmt.Println("   Environment:")
			for k, v := range server.Env {
				fmt.Printf("     %s: %s\n", k, v)
			}
		}
	} else {
		fmt.Println("❌ Google MCP server not configured")
		fmt.Println("")
		fmt.Println("Run: google-auth add")
	}

	// Show other servers
	otherCount := 0
	for name := range config.MCPServers {
		if name != serverName {
			otherCount++
		}
	}
	if otherCount > 0 {
		fmt.Println("")
		fmt.Printf("Other MCP servers in config: %d\n", otherCount)
		for name := range config.MCPServers {
			if name != serverName {
				fmt.Printf("   - %s\n", name)
			}
		}
	}

	return nil
}
