// Google MCP configuration tool for Claude Code
//
// Manages MCP server configuration in various locations.
//
// Usage:
//
//	google-auth add [-location=vscode|project|claude]
//	google-auth remove [-location=vscode|project|claude]
//	google-auth status [-location=vscode|project|claude]
//
// Locations:
//
//	vscode  - .vscode/mcp.json (default, VSCode extension)
//	project - .mcp.json (CLI project-level)
//	claude  - .claude/mcp.json (Claude folder)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// MCPConfig represents the mcp.json file structure
type MCPConfig struct {
	Schema     string               `json:"$schema,omitempty"`
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

// MCPServer represents an MCP server configuration
type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
}

const (
	schemaURL    = "https://modelcontextprotocol.io/schema/config.json"
	serverName   = "google"
	serverCmd    = "google-mcp-server"
	envClientID  = "${GOOGLE_CLIENT_ID}"
	envSecretID  = "${GOOGLE_CLIENT_SECRET}"
)

// Location configs
var locations = map[string]struct {
	dir  string // subdirectory (empty = project root)
	file string // filename
}{
	"vscode":  {dir: ".vscode", file: "mcp.json"},
	"project": {dir: "", file: ".mcp.json"},
	"claude":  {dir: ".claude", file: "mcp.json"},
}

func main() {
	// Define flags
	location := flag.String("location", "vscode", "Config location: vscode, project, or claude")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	cmd := args[0]

	// Validate location
	loc, ok := locations[*location]
	if !ok {
		fmt.Printf("Unknown location: %s\n", *location)
		fmt.Println("Valid locations: vscode, project, claude")
		os.Exit(1)
	}

	// Find project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Printf("Error finding project root: %v\n", err)
		os.Exit(1)
	}

	// Build paths
	var targetDir, mcpFile string
	if loc.dir != "" {
		targetDir = filepath.Join(projectRoot, loc.dir)
		mcpFile = filepath.Join(targetDir, loc.file)
	} else {
		targetDir = projectRoot
		mcpFile = filepath.Join(projectRoot, loc.file)
	}

	switch cmd {
	case "add":
		if err := addServer(targetDir, mcpFile, loc.dir != ""); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "remove":
		if err := removeServer(mcpFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "status":
		if err := showStatus(mcpFile, *location); err != nil {
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
	fmt.Println("Manages MCP server configuration for Claude Code.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  google-auth [-location=LOC] add      Add google server")
	fmt.Println("  google-auth [-location=LOC] remove   Remove google server")
	fmt.Println("  google-auth [-location=LOC] status   Show config status")
	fmt.Println("")
	fmt.Println("Locations (-location flag):")
	fmt.Println("  vscode  - .vscode/mcp.json (default, VSCode extension)")
	fmt.Println("  project - .mcp.json (CLI project-level)")
	fmt.Println("  claude  - .claude/mcp.json (Claude folder)")
	fmt.Println("")
	fmt.Println("Environment variables are referenced as ${VAR} - no secrets stored.")
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

// loadConfig loads the mcp.json file, returning empty config if not exists
func loadConfig(mcpFile string) (*MCPConfig, error) {
	config := &MCPConfig{
		Schema:     schemaURL,
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

	// Ensure schema is set
	if config.Schema == "" {
		config.Schema = schemaURL
	}

	return config, nil
}

// saveConfig saves the config to mcp.json with pretty formatting
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

// addServer adds the google MCP server to mcp.json
func addServer(targetDir, mcpFile string, createDir bool) error {
	// Ensure target directory exists if needed
	if createDir {
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

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

	fmt.Println("✅ Google MCP server added")
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

// removeServer removes the google MCP server from mcp.json
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

	fmt.Println("✅ Google MCP server removed")
	fmt.Printf("File: %s\n", mcpFile)
	fmt.Println("")
	fmt.Println("Restart Claude Code to apply changes.")

	return nil
}

// showStatus shows the current mcp.json configuration
func showStatus(mcpFile, location string) error {
	fmt.Println("=== Google MCP Configuration Status ===")
	fmt.Printf("Location: %s\n", location)
	fmt.Println("")

	// Check if file exists
	if _, err := os.Stat(mcpFile); os.IsNotExist(err) {
		fmt.Printf("❌ No config found at %s\n", mcpFile)
		fmt.Println("")
		fmt.Printf("Run: google-auth -location=%s add\n", location)
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
		fmt.Printf("Run: google-auth -location=%s add\n", location)
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
