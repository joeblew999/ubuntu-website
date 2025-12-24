package claude

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// InstallType represents the type of Claude CLI installation
type InstallType string

const (
	InstallTypeNative  InstallType = "native"  // Mach-O/ELF binary from Anthropic
	InstallTypeNPM     InstallType = "npm"     // Node.js script via npm
	InstallTypeBun     InstallType = "bun"     // Node.js script via bun
	InstallTypeUnknown InstallType = "unknown" // Unknown type
)

// Installation represents a detected Claude CLI installation
type Installation struct {
	Path        string      `json:"path"`
	Type        InstallType `json:"type"`
	Version     string      `json:"version"`
	VersionInfo *Version    `json:"version_info,omitempty"`
	IsSymlink   bool        `json:"is_symlink"`
	Target      string      `json:"target,omitempty"`
	Recommended bool        `json:"recommended"`
}

// Version represents a parsed semantic version
type Version struct {
	Major int    `json:"major"`
	Minor int    `json:"minor"`
	Patch int    `json:"patch"`
	Raw   string `json:"raw"`
}

// Compare compares two versions. Returns -1 if v < other, 0 if equal, 1 if v > other
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor != other.Minor {
		if v.Minor < other.Minor {
			return -1
		}
		return 1
	}
	if v.Patch != other.Patch {
		if v.Patch < other.Patch {
			return -1
		}
		return 1
	}
	return 0
}

// String returns the version as a string
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// ParseVersion parses a version string like "2.0.73" or "2.0.73 (Claude Code)"
func ParseVersion(s string) (*Version, error) {
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("invalid version format: %s", s)
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
		Raw:   s,
	}, nil
}

// KnownLocations lists all known Claude CLI installation locations
var KnownLocations = []struct {
	Path        string
	Type        InstallType
	Recommended bool
}{
	{"~/.local/bin/claude", InstallTypeNative, true},
	{"~/.bun/bin/claude", InstallTypeBun, false},
	{"/opt/homebrew/bin/claude", InstallTypeNPM, false},
	{"/usr/local/bin/claude", InstallTypeUnknown, false},
	{"/usr/bin/claude", InstallTypeUnknown, false},
}

// Detector handles Claude CLI detection
type Detector struct {
	home string
}

// NewDetector creates a new detector
func NewDetector() *Detector {
	home, _ := os.UserHomeDir()
	return &Detector{home: home}
}

// ExpandPath expands ~ to home directory
func (d *Detector) ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(d.home, path[2:])
	}
	return path
}

// DetectAll finds all Claude CLI installations
func (d *Detector) DetectAll() []*Installation {
	var installations []*Installation

	for _, loc := range KnownLocations {
		path := d.ExpandPath(loc.Path)
		if inst := d.detectAt(path, loc.Type, loc.Recommended); inst != nil {
			installations = append(installations, inst)
		}
	}

	return installations
}

// detectAt checks for a Claude CLI installation at a specific path
func (d *Detector) detectAt(path string, expectedType InstallType, recommended bool) *Installation {
	info, err := os.Lstat(path)
	if err != nil {
		return nil
	}

	inst := &Installation{
		Path:        path,
		Recommended: recommended,
	}

	// Check if symlink
	if info.Mode()&os.ModeSymlink != 0 {
		inst.IsSymlink = true
		target, err := os.Readlink(path)
		if err == nil {
			inst.Target = target
		}
	}

	// Determine type by examining the file
	inst.Type = d.detectType(path)

	// Get version
	inst.Version = d.getVersion(path)
	if inst.Version != "" {
		inst.VersionInfo, _ = ParseVersion(inst.Version)
	}

	return inst
}

// detectType determines the installation type by examining the binary
func (d *Detector) detectType(path string) InstallType {
	cmd := exec.Command("file", path)
	output, err := cmd.Output()
	if err != nil {
		return InstallTypeUnknown
	}

	outputStr := string(output)

	// Check for native binary signatures
	if strings.Contains(outputStr, "Mach-O") || strings.Contains(outputStr, "ELF") {
		return InstallTypeNative
	}

	// Check for script (npm/bun use Node.js scripts)
	if strings.Contains(outputStr, "script") || strings.Contains(outputStr, "text") {
		if strings.Contains(path, ".bun") {
			return InstallTypeBun
		}
		return InstallTypeNPM
	}

	return InstallTypeUnknown
}

// getVersion runs the CLI to get its version
func (d *Detector) getVersion(path string) string {
	cmd := exec.Command(path, "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetRecommended returns the recommended (native) installation, if any
func (d *Detector) GetRecommended() *Installation {
	for _, inst := range d.DetectAll() {
		if inst.Type == InstallTypeNative && inst.Recommended {
			return inst
		}
	}
	return nil
}

// GetActive returns the installation that would be used when running 'claude'
func (d *Detector) GetActive() *Installation {
	path, err := exec.LookPath("claude")
	if err != nil {
		return nil
	}

	for _, inst := range d.DetectAll() {
		if inst.Path == path {
			return inst
		}
	}

	return d.detectAt(path, InstallTypeUnknown, false)
}

// Status represents the overall Claude CLI status
type Status struct {
	Installed       bool            `json:"installed"`
	NativeInstalled bool            `json:"native_installed"`
	ActivePath      string          `json:"active_path"`
	ActiveType      InstallType     `json:"active_type"`
	ActiveVersion   *Version        `json:"active_version,omitempty"`
	NeedsMigration  bool            `json:"needs_migration"`
	NeedsUpgrade    bool            `json:"needs_upgrade"`
	Installations   []*Installation `json:"installations"`
}

// GetStatus returns the overall Claude CLI status
func (d *Detector) GetStatus(minVersion *Version) *Status {
	status := &Status{
		Installations: d.DetectAll(),
	}

	for _, inst := range status.Installations {
		status.Installed = true

		if inst.Type == InstallTypeNative {
			status.NativeInstalled = true
		}

		if inst.Type == InstallTypeNPM || inst.Type == InstallTypeBun {
			status.NeedsMigration = true
		}
	}

	if active := d.GetActive(); active != nil {
		status.ActivePath = active.Path
		status.ActiveType = active.Type
		status.ActiveVersion = active.VersionInfo

		if minVersion != nil && active.VersionInfo != nil {
			if active.VersionInfo.Compare(minVersion) < 0 {
				status.NeedsUpgrade = true
			}
		}
	}

	return status
}

// NodeModulesLocation represents a node_modules installation to clean up
type NodeModulesLocation struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

// FindNodeModules finds all @anthropic-ai/claude-code node_modules directories
func (d *Detector) FindNodeModules() []NodeModulesLocation {
	locations := []string{
		filepath.Join(d.home, ".bun/install/global/node_modules/@anthropic-ai/claude-code"),
		"/opt/homebrew/lib/node_modules/@anthropic-ai/claude-code",
		"/usr/local/lib/node_modules/@anthropic-ai/claude-code",
	}

	// Also check npm global prefix
	if cmd := exec.Command("npm", "prefix", "-g"); cmd != nil {
		if output, err := cmd.Output(); err == nil {
			prefix := strings.TrimSpace(string(output))
			npmPath := filepath.Join(prefix, "lib/node_modules/@anthropic-ai/claude-code")
			found := false
			for _, loc := range locations {
				if loc == npmPath {
					found = true
					break
				}
			}
			if !found {
				locations = append(locations, npmPath)
			}
		}
	}

	var results []NodeModulesLocation
	for _, loc := range locations {
		info, err := os.Stat(loc)
		results = append(results, NodeModulesLocation{
			Path:   loc,
			Exists: err == nil && info.IsDir(),
		})
	}

	return results
}
