package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// ReleaseCmd is the parent command for release operations
var ReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Release build orchestration",
	Long: `Commands for orchestrating cross-platform release builds.

Reads tool configuration from Taskfiles to determine build strategy:
- CGO=0: Can cross-compile all platforms from any OS
- CGO=1: Requires native build on each target platform

Examples:
  xplat release matrix tui          # Output platform matrix as JSON
  xplat release build tui           # Build for all platforms
  xplat release build tui --current # Build for current platform only`,
}

// Platform represents a target platform for builds
type Platform struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Runner   string `json:"runner,omitempty"`   // GitHub Actions runner
	CrossCompile bool `json:"cross_compile"`    // Can be built from Linux
}

// BuildMatrix represents the full build configuration
type BuildMatrix struct {
	Tool      string     `json:"tool"`
	BinName   string     `json:"bin_name"`  // Actual binary name (e.g., task-ui vs tui)
	Lang      string     `json:"lang"`
	CGO       bool       `json:"cgo"`
	Platforms []Platform `json:"platforms"`
}

// Standard platforms we support
var allPlatforms = []Platform{
	{OS: "linux", Arch: "amd64", Runner: "ubuntu-latest", CrossCompile: true},
	{OS: "linux", Arch: "arm64", Runner: "ubuntu-latest", CrossCompile: true},
	{OS: "darwin", Arch: "amd64", Runner: "macos-latest", CrossCompile: true},
	{OS: "darwin", Arch: "arm64", Runner: "macos-latest", CrossCompile: true},
	{OS: "windows", Arch: "amd64", Runner: "windows-latest", CrossCompile: true},
	{OS: "windows", Arch: "arm64", Runner: "windows-latest", CrossCompile: true},
}

// ReleaseMatrixCmd outputs the build matrix for a tool
var ReleaseMatrixCmd = &cobra.Command{
	Use:   "matrix <tool>",
	Short: "Output platform build matrix for a tool",
	Long: `Reads the tool's Taskfile to determine CGO requirements and outputs
a JSON matrix suitable for GitHub Actions or other CI systems.

For CGO=0 tools, all platforms can be cross-compiled from Linux.
For CGO=1 tools, each platform needs a native runner.

Examples:
  xplat release matrix tui
  xplat release matrix tui --format github`,
	Args: cobra.ExactArgs(1),
	RunE: runReleaseMatrix,
}

// ReleaseBuildCmd builds a tool for all (or current) platform
var ReleaseBuildCmd = &cobra.Command{
	Use:   "build <tool>",
	Short: "Build a tool for release",
	Long: `Builds a tool for all target platforms by calling the tool's
Taskfile release:build task with appropriate GOOS/GOARCH.

For CGO=0 tools, builds all platforms from current machine.
For CGO=1 tools, only builds for current platform (use CI for others).

Examples:
  xplat release build tui              # Build all platforms
  xplat release build tui --current    # Build current platform only
  xplat release build tui --platform linux/amd64`,
	Args: cobra.ExactArgs(1),
	RunE: runReleaseBuild,
}

var (
	matrixFormat   string
	buildCurrent   bool
	buildPlatform  string
)

func init() {
	ReleaseMatrixCmd.Flags().StringVar(&matrixFormat, "format", "json", "Output format: json, github")
	ReleaseBuildCmd.Flags().BoolVar(&buildCurrent, "current", false, "Only build for current platform")
	ReleaseBuildCmd.Flags().StringVar(&buildPlatform, "platform", "", "Build for specific platform (e.g., linux/amd64)")

	ReleaseCmd.AddCommand(ReleaseMatrixCmd)
	ReleaseCmd.AddCommand(ReleaseBuildCmd)
}

// TaskfileVars represents the vars section of a Taskfile
type TaskfileVars struct {
	Vars map[string]interface{} `yaml:"vars"`
}

// findTaskfile locates the Taskfile for a given tool
func findTaskfile(tool string) (string, error) {
	// Map short names to taskfile locations
	// Convention: taskfiles/tools/Taskfile.<tool>.yml or taskfiles/Taskfile.<tool>.yml

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Try multiple locations
	candidates := []string{
		filepath.Join(cwd, "taskfiles", "tools", fmt.Sprintf("Taskfile.%s.yml", tool)),
		filepath.Join(cwd, "taskfiles", fmt.Sprintf("Taskfile.%s.yml", tool)),
	}

	// Also try with common name mappings
	nameMap := map[string]string{
		"tui": "task-ui",
		"pc":  "process-compose",
	}
	if fullName, ok := nameMap[tool]; ok {
		candidates = append(candidates,
			filepath.Join(cwd, "taskfiles", "tools", fmt.Sprintf("Taskfile.%s.yml", fullName)),
			filepath.Join(cwd, "taskfiles", fmt.Sprintf("Taskfile.%s.yml", fullName)),
		)
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("taskfile not found for tool: %s\nSearched: %v", tool, candidates)
}

// parseTaskfileVars extracts vars from a Taskfile
func parseTaskfileVars(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var tf TaskfileVars
	if err := yaml.Unmarshal(data, &tf); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for k, v := range tf.Vars {
		// Handle different value types
		switch val := v.(type) {
		case string:
			result[k] = val
		case int:
			result[k] = fmt.Sprintf("%d", val)
		case bool:
			result[k] = fmt.Sprintf("%t", val)
		}
	}

	return result, nil
}

// getToolConfig extracts build configuration from Taskfile vars
func getToolConfig(tool string) (*BuildMatrix, error) {
	taskfile, err := findTaskfile(tool)
	if err != nil {
		return nil, err
	}

	vars, err := parseTaskfileVars(taskfile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse taskfile: %w", err)
	}

	// Find CGO var and BIN var (convention: <TOOL>_CGO, <TOOL>_BIN)
	cgo := false
	binName := tool // Default to tool name
	for k, v := range vars {
		upperK := strings.ToUpper(k)
		if strings.HasSuffix(upperK, "_CGO") {
			cgo = v == "1" || strings.ToLower(v) == "true"
		}
		if strings.HasSuffix(upperK, "_BIN") {
			// Extract binary name, removing {{exeExt}} template
			binName = strings.TrimSuffix(v, "{{exeExt}}")
			binName = strings.TrimSuffix(binName, ".exe")
		}
	}

	// Determine platforms based on CGO
	platforms := make([]Platform, len(allPlatforms))
	copy(platforms, allPlatforms)

	if cgo {
		// CGO=1: Cannot cross-compile, need native runners
		for i := range platforms {
			platforms[i].CrossCompile = false
		}
	}

	return &BuildMatrix{
		Tool:      tool,
		BinName:   binName,
		Lang:      "go", // For now, assume Go. Could detect from source later
		CGO:       cgo,
		Platforms: platforms,
	}, nil
}

func runReleaseMatrix(cmd *cobra.Command, args []string) error {
	tool := args[0]

	matrix, err := getToolConfig(tool)
	if err != nil {
		return err
	}

	switch matrixFormat {
	case "github":
		// Output format suitable for GitHub Actions matrix
		type GHInclude struct {
			OS     string `json:"os"`
			GOOS   string `json:"goos"`
			GOARCH string `json:"goarch"`
			Runner string `json:"runner"`
		}
		includes := make([]GHInclude, 0, len(matrix.Platforms))
		for _, p := range matrix.Platforms {
			includes = append(includes, GHInclude{
				OS:     p.OS,
				GOOS:   p.OS,
				GOARCH: p.Arch,
				Runner: p.Runner,
			})
		}
		output := map[string]interface{}{
			"include":        includes,
			"bin_name":       matrix.BinName,
			"cgo":            matrix.CGO,
			"cross_compile":  !matrix.CGO,
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))

	default:
		// Default JSON output
		data, _ := json.MarshalIndent(matrix, "", "  ")
		fmt.Println(string(data))
	}

	return nil
}

func runReleaseBuild(cmd *cobra.Command, args []string) error {
	tool := args[0]

	matrix, err := getToolConfig(tool)
	if err != nil {
		return err
	}

	// Determine which platforms to build
	var targetPlatforms []Platform

	if buildCurrent {
		// Only current platform
		targetPlatforms = []Platform{{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		}}
	} else if buildPlatform != "" {
		// Specific platform
		parts := strings.Split(buildPlatform, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid platform format: %s (expected os/arch)", buildPlatform)
		}
		targetPlatforms = []Platform{{
			OS:   parts[0],
			Arch: parts[1],
		}}
	} else {
		// All platforms
		if matrix.CGO {
			// CGO=1: Can only build current platform
			fmt.Printf("Warning: %s requires CGO, can only build for current platform\n", tool)
			targetPlatforms = []Platform{{
				OS:   runtime.GOOS,
				Arch: runtime.GOARCH,
			}}
		} else {
			targetPlatforms = matrix.Platforms
		}
	}

	// Build each platform
	for _, p := range targetPlatforms {
		fmt.Printf("Building %s for %s/%s...\n", tool, p.OS, p.Arch)

		taskCmd := exec.Command("task", fmt.Sprintf("%s:release:build", tool))
		taskCmd.Env = append(os.Environ(),
			fmt.Sprintf("GOOS=%s", p.OS),
			fmt.Sprintf("GOARCH=%s", p.Arch),
		)
		taskCmd.Stdout = os.Stdout
		taskCmd.Stderr = os.Stderr

		if err := taskCmd.Run(); err != nil {
			return fmt.Errorf("build failed for %s/%s: %w", p.OS, p.Arch, err)
		}
	}

	fmt.Printf("\nOK: Built %s for %d platform(s)\n", tool, len(targetPlatforms))
	return nil
}
