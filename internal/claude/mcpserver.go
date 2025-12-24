// MCP Server registry and management.
//
// This provides a unified way to add/remove MCP servers to different Claude targets
// (VSCode, Desktop, CLI). Each MCP server (Google, GitHub, Filesystem, etc.) registers
// itself with its configuration and permissions.
//
// SAFETY: When modifying VSCode config (the environment we're running in), we use:
// - Backup before modifying (*.backup)
// - JSON validation before saving
// - Atomic write (temp file + rename)
//
// Usage:
//
//	// Register a server (typically done at init time by the server package)
//	claude.RegisterMCPServer(claude.MCPServerDef{
//	    Name:        "google",
//	    Command:     "google-mcp-server",
//	    Permissions: []string{"mcp__google__calendar_list", ...},
//	})
//
//	// Add server to a target
//	claude.AddMCPServer("google", claude.TargetVSCode, projectRoot)
//
//	// Remove server from a target
//	claude.RemoveMCPServer("google", claude.TargetVSCode, projectRoot)
package claude

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// backupFile creates a backup of a file before modifying it
// Returns the backup path, or empty string if file doesn't exist
func backupFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", nil // No file to backup
	}

	backupPath := path + ".backup." + time.Now().Format("20060102-150405")

	src, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file for backup: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(backupPath)
		return "", fmt.Errorf("failed to copy file for backup: %w", err)
	}

	return backupPath, nil
}

// validateJSON checks if data is valid JSON
func validateJSON(data []byte) error {
	var js json.RawMessage
	if err := json.Unmarshal(data, &js); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// safeWriteFile writes data to a file atomically (temp + rename)
// and validates JSON before writing
func safeWriteFile(path string, data []byte, perm os.FileMode) error {
	// Validate JSON first
	if err := validateJSON(data); err != nil {
		return err
	}

	// Write to temp file in same directory (for atomic rename)
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".mcp-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Clean up temp file on any error
	defer func() {
		if tmpPath != "" {
			os.Remove(tmpPath)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Set permissions before rename
	if err := os.Chmod(tmpPath, perm); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	tmpPath = "" // Prevent cleanup since rename succeeded
	return nil
}

// isVSCodeTarget returns true if modifying the current VSCode environment
func isVSCodeTarget(target Target) bool {
	return target == TargetVSCode || target == TargetClaude
}

// MCPServerDef defines an MCP server that can be added to Claude targets
type MCPServerDef struct {
	Name        string            // Server name (e.g., "google", "github", "filesystem")
	Command     string            // Command to run (e.g., "google-mcp-server")
	Args        []string          // Command arguments
	Env         map[string]string // Environment variables
	Permissions []string          // Required permissions for Claude Code (VSCode)
	AccountsDir string            // Directory where server stores tokens (e.g., ".google-mcp-accounts")
}

// mcpRegistry holds all registered MCP servers
var mcpRegistry = make(map[string]MCPServerDef)

// RegisterMCPServer adds an MCP server definition to the registry
func RegisterMCPServer(def MCPServerDef) {
	mcpRegistry[def.Name] = def
}

// GetMCPServer returns a registered MCP server definition
func GetMCPServer(name string) (MCPServerDef, bool) {
	def, ok := mcpRegistry[name]
	return def, ok
}

// ListMCPServers returns all registered MCP server names
func ListMCPServers() []string {
	names := make([]string, 0, len(mcpRegistry))
	for name := range mcpRegistry {
		names = append(names, name)
	}
	return names
}

// AddMCPServerResult contains the result of adding an MCP server
type AddMCPServerResult struct {
	ConfigPath     string // Path to the MCP config file that was modified
	SettingsPath   string // Path to the settings file that was modified (if any)
	ServerAdded    bool   // True if server was added (false if already existed)
	PermissionsSet bool   // True if permissions were added
	BackupPath     string // Path to backup file (if created for VSCode target)
}

// AddMCPServer adds an MCP server to a Claude target
// For VSCode targets, creates a backup before modifying to prevent breaking Claude Code
func AddMCPServer(serverName string, target Target, projectRoot string) (*AddMCPServerResult, error) {
	def, ok := GetMCPServer(serverName)
	if !ok {
		return nil, fmt.Errorf("unknown MCP server: %s", serverName)
	}

	result := &AddMCPServerResult{}

	// Get config path for target
	configPath, err := GetTargetConfigPath(target, projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}
	if configPath == "" {
		return nil, fmt.Errorf("target %s does not support MCP config", target)
	}

	result.ConfigPath = configPath

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// SAFETY: For VSCode targets, create backup before modifying
	// This protects the environment we're running in
	if isVSCodeTarget(target) {
		backupPath, err := backupFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to backup config (safety check): %w", err)
		}
		result.BackupPath = backupPath
	}

	// Load or create config
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Add server if not already present
	if !config.HasServer(def.Name) {
		config.AddServer(def.Name, Server{
			Command: def.Command,
			Args:    def.Args,
			Env:     def.Env,
		})

		// SAFETY: Use safe write for VSCode targets
		if isVSCodeTarget(target) {
			if err := SaveConfigSafe(configPath, config); err != nil {
				return nil, fmt.Errorf("failed to save config: %w", err)
			}
		} else {
			if err := SaveConfig(configPath, config); err != nil {
				return nil, fmt.Errorf("failed to save config: %w", err)
			}
		}
		result.ServerAdded = true
	}

	// Add permissions if target supports them
	targetInfo := GetTargetInfo(target)
	if targetInfo.HasPerms && len(def.Permissions) > 0 {
		if err := EnsureClaudeDir(projectRoot); err != nil {
			return nil, fmt.Errorf("failed to create .claude directory: %w", err)
		}

		settingsPath := GetSettingsPath(projectRoot)
		result.SettingsPath = settingsPath

		settings, err := LoadSettings(settingsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load settings: %w", err)
		}

		if settings.AddPermissions(def.Permissions) {
			if err := SaveSettings(settingsPath, settings); err != nil {
				return nil, fmt.Errorf("failed to save settings: %w", err)
			}
			result.PermissionsSet = true
		}
	}

	return result, nil
}

// RemoveMCPServerResult contains the result of removing an MCP server
type RemoveMCPServerResult struct {
	ConfigPath         string // Path to the MCP config file that was modified
	SettingsPath       string // Path to the settings file that was modified (if any)
	ServerRemoved      bool   // True if server was removed (false if didn't exist)
	PermissionsRemoved bool   // True if permissions were removed
	ConfigDeleted      bool   // True if config file was deleted (was last server)
	BackupPath         string // Path to backup file (if created for VSCode target)
}

// RemoveMCPServer removes an MCP server from a Claude target
// For VSCode targets, creates a backup before modifying to prevent breaking Claude Code
func RemoveMCPServer(serverName string, target Target, projectRoot string) (*RemoveMCPServerResult, error) {
	def, ok := GetMCPServer(serverName)
	if !ok {
		return nil, fmt.Errorf("unknown MCP server: %s", serverName)
	}

	result := &RemoveMCPServerResult{}

	// Get config path for target
	configPath, err := GetTargetConfigPath(target, projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}
	if configPath == "" {
		return nil, fmt.Errorf("target %s does not support MCP config", target)
	}

	result.ConfigPath = configPath

	// SAFETY: For VSCode targets, create backup before modifying
	// This protects the environment we're running in
	if isVSCodeTarget(target) {
		backupPath, err := backupFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to backup config (safety check): %w", err)
		}
		result.BackupPath = backupPath
	}

	// Load config
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Remove server if present
	if config.RemoveServer(def.Name) {
		result.ServerRemoved = true

		if config.IsEmpty() {
			// Delete the config file if no servers left
			os.Remove(configPath)
			result.ConfigDeleted = true
		} else {
			// SAFETY: Use safe write for VSCode targets
			if isVSCodeTarget(target) {
				if err := SaveConfigSafe(configPath, config); err != nil {
					return nil, fmt.Errorf("failed to save config: %w", err)
				}
			} else {
				if err := SaveConfig(configPath, config); err != nil {
					return nil, fmt.Errorf("failed to save config: %w", err)
				}
			}
		}
	}

	// Remove permissions if target supports them
	targetInfo := GetTargetInfo(target)
	if targetInfo.HasPerms && len(def.Permissions) > 0 {
		settingsPath := GetSettingsPath(projectRoot)
		result.SettingsPath = settingsPath

		settings, err := LoadSettings(settingsPath)
		if err == nil && settings != nil {
			if settings.RemovePermissions(def.Permissions) {
				SaveSettings(settingsPath, settings)
				result.PermissionsRemoved = true
			}
		}
	}

	return result, nil
}

// MCPServerStatus contains status information about an MCP server
type MCPServerStatus struct {
	Name            string
	Configured      bool   // Server is in MCP config
	ConfigPath      string // Where it's configured
	PermissionCount int    // Number of permissions set (for targets with perms)
	Installed       bool   // Binary is installed and in PATH
	AccountCount    int    // Number of authenticated accounts (if applicable)
}

// GetMCPServerStatus returns the status of an MCP server for a target
func GetMCPServerStatus(serverName string, target Target, projectRoot string) (*MCPServerStatus, error) {
	def, ok := GetMCPServer(serverName)
	if !ok {
		return nil, fmt.Errorf("unknown MCP server: %s", serverName)
	}

	status := &MCPServerStatus{Name: serverName}

	// Check config
	configPath, err := GetTargetConfigPath(target, projectRoot)
	if err == nil && configPath != "" {
		status.ConfigPath = configPath
		config, err := LoadConfig(configPath)
		if err == nil {
			status.Configured = config.HasServer(def.Name)
		}
	}

	// Check permissions
	targetInfo := GetTargetInfo(target)
	if targetInfo.HasPerms {
		settingsPath := GetSettingsPath(projectRoot)
		settings, err := LoadSettings(settingsPath)
		if err == nil && settings != nil {
			// Count matching permissions
			for _, perm := range def.Permissions {
				if settings.HasPermission(perm) {
					status.PermissionCount++
				}
			}
		}
	}

	// Check if binary is installed
	if def.Command != "" {
		if _, err := LookPath(def.Command); err == nil {
			status.Installed = true
		}
	}

	// Check accounts if applicable
	if def.AccountsDir != "" {
		home, _ := os.UserHomeDir()
		accounts, _ := filepath.Glob(filepath.Join(home, def.AccountsDir, "*.json"))
		status.AccountCount = len(accounts)
	}

	return status, nil
}

// LookPath is a wrapper for exec.LookPath to allow testing
var LookPath = func(file string) (string, error) {
	return lookPathImpl(file)
}

func lookPathImpl(file string) (string, error) {
	// Import exec inline to avoid circular deps in tests
	import_exec_LookPath := func(f string) (string, error) {
		// This will be replaced by the actual implementation
		return "", fmt.Errorf("not found")
	}
	_ = import_exec_LookPath

	// Use os/exec.LookPath
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		full := filepath.Join(dir, file)
		if info, err := os.Stat(full); err == nil && !info.IsDir() {
			return full, nil
		}
	}
	return "", fmt.Errorf("executable not found: %s", file)
}
