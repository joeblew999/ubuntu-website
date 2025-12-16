// Claude settings management for .claude/settings.json
package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SettingsFile is the standard Claude settings filename
const SettingsFile = "settings.json"

// ClaudeSettings represents the .claude/settings.json file structure
type ClaudeSettings struct {
	Permissions                Permissions       `json:"permissions"`
	EnableAllProjectMcpServers bool              `json:"enableAllProjectMcpServers"`
	Env                        map[string]string `json:"env,omitempty"`
}

// Permissions represents the permissions block in settings.json
type Permissions struct {
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

// LoadSettings loads the settings.json file, returning empty settings if not exists
func LoadSettings(settingsPath string) (*ClaudeSettings, error) {
	settings := &ClaudeSettings{
		Permissions: Permissions{
			Allow: []string{},
			Deny:  []string{},
		},
		EnableAllProjectMcpServers: true,
	}

	data, err := os.ReadFile(settingsPath)
	if os.IsNotExist(err) {
		return settings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", settingsPath, err)
	}

	if err := json.Unmarshal(data, settings); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", settingsPath, err)
	}

	if settings.Permissions.Allow == nil {
		settings.Permissions.Allow = []string{}
	}
	if settings.Permissions.Deny == nil {
		settings.Permissions.Deny = []string{}
	}

	return settings, nil
}

// SaveSettings saves the settings to settings.json with pretty formatting
func SaveSettings(settingsPath string, settings *ClaudeSettings) error {
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Add trailing newline
	data = append(data, '\n')

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", settingsPath, err)
	}

	return nil
}

// AddPermissions adds permissions to the allow list
// Returns true if any permissions were added, false if all already present
func (s *ClaudeSettings) AddPermissions(permissions []string) bool {
	existingPerms := make(map[string]bool)
	for _, p := range s.Permissions.Allow {
		existingPerms[p] = true
	}

	added := false
	for _, perm := range permissions {
		if !existingPerms[perm] {
			s.Permissions.Allow = append(s.Permissions.Allow, perm)
			added = true
		}
	}

	return added
}

// RemovePermissions removes permissions from the allow list
// Returns true if any permissions were removed
func (s *ClaudeSettings) RemovePermissions(permissions []string) bool {
	permsToRemove := make(map[string]bool)
	for _, p := range permissions {
		permsToRemove[p] = true
	}

	newAllow := []string{}
	removed := false
	for _, p := range s.Permissions.Allow {
		if permsToRemove[p] {
			removed = true
		} else {
			newAllow = append(newAllow, p)
		}
	}

	s.Permissions.Allow = newAllow
	return removed
}

// HasPermission checks if a permission is in the allow list
func (s *ClaudeSettings) HasPermission(permission string) bool {
	for _, p := range s.Permissions.Allow {
		if p == permission {
			return true
		}
	}
	return false
}

// CountPermissionsWithPrefix counts permissions with a given prefix
func (s *ClaudeSettings) CountPermissionsWithPrefix(prefix string) int {
	count := 0
	for _, p := range s.Permissions.Allow {
		if len(p) >= len(prefix) && p[:len(prefix)] == prefix {
			count++
		}
	}
	return count
}

// GetSettingsPath returns the path to .claude/settings.json
func GetSettingsPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".claude", SettingsFile)
}

// EnsureClaudeDir creates the .claude directory if needed
func EnsureClaudeDir(projectRoot string) error {
	return os.MkdirAll(filepath.Join(projectRoot, ".claude"), 0755)
}
