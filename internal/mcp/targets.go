// Claude target definitions - where MCP configs live for different Claude products
//
// This file defines configuration locations for:
// - Claude Code (VSCode extension) - developers
// - Claude Desktop (standalone app) - end users
// - Claude CLI (command line) - developers/power users
// - Claude Cloud (API-only) - no local config
//
// Each target has different:
// - Config file locations
// - Permission models (Claude-specific permissions vs API-based)
// - Installation methods
package mcp

import (
	"os"
	"path/filepath"
	"runtime"
)

// Target represents a Claude product/installation target
type Target string

const (
	// TargetVSCode is Claude Code (VSCode extension) - project-level config
	// Config: .vscode/mcp.json (project) or ~/.claude/settings.json (user)
	// Permissions: settings.json with allow/deny lists
	// Used by: Developers using VSCode
	TargetVSCode Target = "vscode"

	// TargetProject is project-level MCP config (CLI)
	// Config: .mcp.json in project root
	// Used by: Claude CLI, shared across team
	TargetProject Target = "project"

	// TargetClaude is Claude folder config
	// Config: .claude/mcp.json in project root
	// Used by: Claude Code project-specific settings
	TargetClaude Target = "claude"

	// TargetDesktop is Claude Desktop (standalone app)
	// Config: ~/Library/Application Support/Claude/claude_desktop_config.json (macOS)
	//         %APPDATA%\Claude\claude_desktop_config.json (Windows)
	//         ~/.config/Claude/claude_desktop_config.json (Linux)
	// Permissions: None (app handles permissions via OAuth)
	// Used by: End users
	TargetDesktop Target = "desktop"

	// TargetUserGlobal is user-level global Claude settings
	// Config: ~/.claude/settings.json or ~/.claude.json
	// Used by: All Claude Code projects for this user
	TargetUserGlobal Target = "user"

	// TargetCloud represents Claude Cloud (Anthropic API)
	// Config: None local - API-based configuration
	// Permissions: API keys and organization policies
	// Used by: Server deployments, CI/CD, automation
	TargetCloud Target = "cloud"
)

// TargetInfo contains information about a Claude target
type TargetInfo struct {
	Name        string // Human-readable name
	Description string // What this target is for
	ConfigDir   string // Subdirectory relative to project/home (empty = root)
	ConfigFile  string // Config filename
	HasPerms    bool   // Whether this target uses permission files
	PermsFile   string // Permission file (if HasPerms)
	ForDevs     bool   // Whether this is primarily for developers
	ForUsers    bool   // Whether this is for end users
}

// GetTargetInfo returns information about a Claude target
func GetTargetInfo(t Target) TargetInfo {
	switch t {
	case TargetVSCode:
		return TargetInfo{
			Name:        "Claude Code (VSCode)",
			Description: "VSCode extension - project-level MCP config",
			ConfigDir:   ".vscode",
			ConfigFile:  "mcp.json",
			HasPerms:    true,
			PermsFile:   ".claude/settings.json",
			ForDevs:     true,
			ForUsers:    false,
		}
	case TargetProject:
		return TargetInfo{
			Name:        "Project MCP",
			Description: "Project-level MCP config for CLI",
			ConfigDir:   "",
			ConfigFile:  ".mcp.json",
			HasPerms:    false, // CLI manages permissions differently
			ForDevs:     true,
			ForUsers:    false,
		}
	case TargetClaude:
		return TargetInfo{
			Name:        "Claude Folder",
			Description: "Claude folder project config",
			ConfigDir:   ".claude",
			ConfigFile:  "mcp.json",
			HasPerms:    true,
			PermsFile:   ".claude/settings.json",
			ForDevs:     true,
			ForUsers:    false,
		}
	case TargetDesktop:
		return TargetInfo{
			Name:        "Claude Desktop",
			Description: "Claude Desktop standalone application",
			ConfigDir:   getDesktopConfigDir(),
			ConfigFile:  "claude_desktop_config.json",
			HasPerms:    false, // Desktop uses OAuth, not file permissions
			ForDevs:     false,
			ForUsers:    true,
		}
	case TargetUserGlobal:
		return TargetInfo{
			Name:        "User Global",
			Description: "User-level global Claude settings",
			ConfigDir:   ".claude",
			ConfigFile:  "settings.json",
			HasPerms:    true,
			PermsFile:   ".claude/settings.json",
			ForDevs:     true,
			ForUsers:    true,
		}
	case TargetCloud:
		return TargetInfo{
			Name:        "Claude Cloud",
			Description: "Anthropic API - no local config",
			ConfigDir:   "",
			ConfigFile:  "",
			HasPerms:    false, // API-based permissions
			ForDevs:     true,
			ForUsers:    true,
		}
	default:
		return TargetInfo{
			Name:        "Unknown",
			Description: "Unknown target",
		}
	}
}

// getDesktopConfigDir returns the Claude Desktop config directory for the current OS
func getDesktopConfigDir() string {
	switch runtime.GOOS {
	case "darwin":
		return "Library/Application Support/Claude"
	case "windows":
		// Note: This will be joined with APPDATA, not home
		return "Claude"
	case "linux":
		return ".config/Claude"
	default:
		return ".config/Claude"
	}
}

// GetDesktopConfigPath returns the full path to Claude Desktop config
func GetDesktopConfigPath() (string, error) {
	var basePath string

	switch runtime.GOOS {
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		basePath = filepath.Join(home, "Library", "Application Support", "Claude")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		basePath = filepath.Join(appData, "Claude")
	case "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		basePath = filepath.Join(home, ".config", "Claude")
	default:
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		basePath = filepath.Join(home, ".config", "Claude")
	}

	return filepath.Join(basePath, "claude_desktop_config.json"), nil
}

// GetUserGlobalConfigPath returns the path to user-level global Claude config
func GetUserGlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

// GetTargetConfigPath returns the config file path for a target
// For project-based targets, projectRoot is required
// For user-based targets (desktop, user), projectRoot is ignored
func GetTargetConfigPath(t Target, projectRoot string) (string, error) {
	info := GetTargetInfo(t)

	switch t {
	case TargetDesktop:
		return GetDesktopConfigPath()
	case TargetUserGlobal:
		return GetUserGlobalConfigPath()
	case TargetCloud:
		return "", nil // No local config
	default:
		// Project-based targets
		if info.ConfigDir != "" {
			return filepath.Join(projectRoot, info.ConfigDir, info.ConfigFile), nil
		}
		return filepath.Join(projectRoot, info.ConfigFile), nil
	}
}

// AllTargets returns all defined targets
func AllTargets() []Target {
	return []Target{
		TargetVSCode,
		TargetProject,
		TargetClaude,
		TargetDesktop,
		TargetUserGlobal,
		TargetCloud,
	}
}

// DevTargets returns targets primarily for developers
func DevTargets() []Target {
	return []Target{
		TargetVSCode,
		TargetProject,
		TargetClaude,
		TargetUserGlobal,
	}
}

// UserTargets returns targets primarily for end users
func UserTargets() []Target {
	return []Target{
		TargetDesktop,
		TargetUserGlobal,
		TargetCloud,
	}
}

// IsDesktopInstalled checks if Claude Desktop is installed
func IsDesktopInstalled() bool {
	configPath, err := GetDesktopConfigPath()
	if err != nil {
		return false
	}

	// Check if config file or parent directory exists
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); err == nil {
		return true
	}

	return false
}
