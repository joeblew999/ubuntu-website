// Package claude provides unified Claude Code management.
//
// This includes:
// - Claude CLI binary detection and version management
// - MCP server management (list, restart, status)
// - Migration from npm/bun to native binary
package claude

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Actor represents who is running the CLI (affects behavior and messaging)
type Actor string

const (
	ActorUser Actor = "user" // Regular user who wants things to just work
	ActorDev  Actor = "dev"  // Developer who may be testing migrations
	ActorCI   Actor = "ci"   // Automated system (minimal output, machine-readable)
)

// CLI provides the unified command-line interface for Claude management
type CLI struct {
	version  string
	stdout   io.Writer
	stderr   io.Writer
	detector *Detector
	actor    Actor
}

// NewCLI creates a new CLI instance
func NewCLI(version string, stdout, stderr io.Writer) *CLI {
	return &CLI{
		version:  version,
		stdout:   stdout,
		stderr:   stderr,
		detector: NewDetector(),
		actor:    ActorUser,
	}
}

// Run executes a CLI command
func (c *CLI) Run(args []string) int {
	if len(args) < 1 {
		c.printUsage()
		return 0
	}

	// Check for global flags first
	args = c.parseGlobalFlags(args)

	if len(args) < 1 {
		c.printUsage()
		return 0
	}

	switch args[0] {
	// CLI management commands
	case "status":
		return c.cmdStatus(args[1:])
	case "check":
		return c.cmdCheck(args[1:])
	case "detect":
		return c.cmdDetect(args[1:])
	case "cleanup":
		return c.cmdCleanup(args[1:])
	case "version-check":
		return c.cmdVersionCheck(args[1:])
	case "ensure":
		return c.cmdEnsure(args[1:])

	// MCP management commands
	case "mcp":
		return c.handleMCP(args[1:])

	// Meta commands
	case "version", "--version", "-v":
		fmt.Fprintf(c.stdout, "claude %s\n", c.version)
		return 0
	case "help", "--help", "-h":
		c.printUsage()
		return 0
	default:
		fmt.Fprintf(c.stderr, "Unknown command: %s\n", args[0])
		c.printUsage()
		return 1
	}
}

func (c *CLI) parseGlobalFlags(args []string) []string {
	var remaining []string
	for _, arg := range args {
		switch arg {
		case "--dev":
			c.actor = ActorDev
		case "--ci":
			c.actor = ActorCI
		case "--user":
			c.actor = ActorUser
		default:
			remaining = append(remaining, arg)
		}
	}
	return remaining
}

func (c *CLI) printUsage() {
	fmt.Fprintln(c.stdout, `Usage: claude [--dev|--ci|--user] <command> [options]

CLI Management:
  status              Show Claude CLI installation status
  check               Check if native binary is installed (exit 0 if yes)
  detect              List all detected installations
  cleanup             Show what needs to be cleaned up
  version-check MIN   Check if installed version >= MIN (e.g., 2.0.0)
  ensure MIN          Ensure native CLI is installed and meets version

MCP Server Management:
  mcp list            List configured MCP servers
  mcp status          Show running MCP servers
  mcp restart [name]  Restart MCP server(s)
  mcp kill [name]     Kill MCP server(s)
  mcp refresh-tokens  Refresh Google tokens and restart server

Meta:
  version             Show version
  help                Show this help

Global Flags:
  --dev     Developer mode (verbose, shows test commands)
  --ci      CI mode (minimal output, machine-readable)
  --user    User mode (friendly messages, default)

Options:
  --json    Output as JSON where applicable

Examples:
  claude status                      # Show CLI status
  claude --dev status                # Developer mode with extra info
  claude --ci ensure 2.0.0           # CI mode, auto-fix
  claude mcp list                    # List MCP servers
  claude mcp restart google          # Restart Google MCP server`)
}

// ============================================================================
// CLI Management Commands
// ============================================================================

func (c *CLI) cmdStatus(args []string) int {
	jsonOutput := hasFlag(args, "--json")

	minVersion, _ := ParseVersion("2.0.0")
	status := c.detector.GetStatus(minVersion)

	if jsonOutput {
		data, _ := json.MarshalIndent(status, "", "  ")
		fmt.Fprintln(c.stdout, string(data))
		return 0
	}

	// CI mode: minimal output
	if c.actor == ActorCI {
		if !status.Installed {
			fmt.Fprintln(c.stdout, "not-installed")
			return 1
		}
		if status.NativeInstalled && !status.NeedsMigration && !status.NeedsUpgrade {
			fmt.Fprintf(c.stdout, "ok:%s\n", status.ActiveVersion.String())
			return 0
		}
		if status.NeedsMigration {
			fmt.Fprintln(c.stdout, "needs-migration")
		}
		if status.NeedsUpgrade {
			fmt.Fprintln(c.stdout, "needs-upgrade")
		}
		if !status.NativeInstalled {
			fmt.Fprintln(c.stdout, "needs-install")
		}
		return 1
	}

	// Human-readable output
	fmt.Fprintln(c.stdout, "Claude CLI Status")
	fmt.Fprintln(c.stdout, "=================")
	if c.actor == ActorDev {
		fmt.Fprintf(c.stdout, "(mode: developer)\n")
	}
	fmt.Fprintln(c.stdout)

	if !status.Installed {
		fmt.Fprintln(c.stdout, "❌ Claude CLI is not installed")
		fmt.Fprintln(c.stdout)
		fmt.Fprintln(c.stdout, "Install with: task claude-cli:install")
		return 1
	}

	fmt.Fprintln(c.stdout, "Installations found:")
	for _, inst := range status.Installations {
		icon := "📦"
		note := ""

		switch inst.Type {
		case InstallTypeNative:
			icon = "✅"
			if inst.Recommended {
				note = " [RECOMMENDED]"
			}
		case InstallTypeNPM:
			icon = "⚠️"
			note = " [npm - should remove]"
		case InstallTypeBun:
			icon = "⚠️"
			note = " [bun - should remove]"
		}

		version := inst.Version
		if version == "" {
			version = "unknown version"
		}

		fmt.Fprintf(c.stdout, "  %s %s (%s)%s\n", icon, inst.Path, version, note)

		if inst.IsSymlink && inst.Target != "" {
			fmt.Fprintf(c.stdout, "      → %s\n", inst.Target)
		}
	}

	fmt.Fprintln(c.stdout)

	if status.ActivePath != "" {
		fmt.Fprintf(c.stdout, "Active (in PATH): %s\n", status.ActivePath)
		if status.ActiveVersion != nil {
			fmt.Fprintf(c.stdout, "Active version: %s\n", status.ActiveVersion.String())
		}
	}

	fmt.Fprintln(c.stdout)

	if status.NeedsMigration {
		fmt.Fprintln(c.stdout, "⚠️  Migration needed: Run 'task claude-cli:cleanup:all' to remove npm/bun versions")
	}
	if status.NeedsUpgrade {
		fmt.Fprintln(c.stdout, "⚠️  Upgrade needed: Run 'task claude-cli:upgrade' to update")
	}
	if !status.NativeInstalled {
		fmt.Fprintln(c.stdout, "⚠️  Native binary not installed: Run 'task claude-cli:install'")
	}
	if status.NativeInstalled && !status.NeedsMigration && !status.NeedsUpgrade {
		fmt.Fprintln(c.stdout, "✅ All good! Native binary installed and up to date.")
	}

	if c.actor == ActorDev {
		fmt.Fprintln(c.stdout)
		fmt.Fprintln(c.stdout, "Developer Commands:")
		fmt.Fprintln(c.stdout, "  task claude-cli:test:install-npm    # Install npm version for testing")
		fmt.Fprintln(c.stdout, "  task claude-cli:test:install-bun    # Install bun version for testing")
		fmt.Fprintln(c.stdout, "  task claude-cli:test:migration      # Full migration test")
		fmt.Fprintln(c.stdout, "  task claude-cli:test:detect-json    # JSON output for debugging")
	}

	return 0
}

func (c *CLI) cmdCheck(args []string) int {
	native := c.detector.GetRecommended()
	if native == nil {
		fmt.Fprintln(c.stderr, "native-not-installed")
		return 1
	}

	if native.Type != InstallTypeNative {
		fmt.Fprintln(c.stderr, "not-native")
		return 1
	}

	fmt.Fprintf(c.stdout, "native:%s:%s\n", native.Version, native.Path)
	return 0
}

func (c *CLI) cmdDetect(args []string) int {
	jsonOutput := hasFlag(args, "--json")

	installations := c.detector.DetectAll()

	if jsonOutput {
		data, _ := json.MarshalIndent(installations, "", "  ")
		fmt.Fprintln(c.stdout, string(data))
		return 0
	}

	if len(installations) == 0 {
		fmt.Fprintln(c.stdout, "No Claude CLI installations found")
		return 0
	}

	for _, inst := range installations {
		fmt.Fprintf(c.stdout, "%s|%s|%s|%v\n", inst.Path, inst.Type, inst.Version, inst.Recommended)
	}

	return 0
}

func (c *CLI) cmdCleanup(args []string) int {
	jsonOutput := hasFlag(args, "--json")

	type CleanupItem struct {
		Path   string `json:"path"`
		Type   string `json:"type"`
		Reason string `json:"reason"`
	}

	var items []CleanupItem

	for _, inst := range c.detector.DetectAll() {
		if inst.Type == InstallTypeNPM || inst.Type == InstallTypeBun {
			items = append(items, CleanupItem{
				Path:   inst.Path,
				Type:   string(inst.Type),
				Reason: "non-native installation",
			})
		}
	}

	for _, nm := range c.detector.FindNodeModules() {
		if nm.Exists {
			items = append(items, CleanupItem{
				Path:   nm.Path,
				Type:   "node_modules",
				Reason: "leftover npm/bun package",
			})
		}
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(items, "", "  ")
		fmt.Fprintln(c.stdout, string(data))
		return 0
	}

	if len(items) == 0 {
		fmt.Fprintln(c.stdout, "✅ No cleanup needed")
		return 0
	}

	fmt.Fprintln(c.stdout, "Items to clean up:")
	for _, item := range items {
		fmt.Fprintf(c.stdout, "  ⚠️  %s (%s)\n", item.Path, item.Reason)
	}

	return 1
}

func (c *CLI) cmdVersionCheck(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(c.stderr, "Usage: claude version-check <minimum-version>")
		return 1
	}

	minVersion, err := ParseVersion(args[0])
	if err != nil {
		fmt.Fprintf(c.stderr, "Invalid version format: %s\n", args[0])
		return 1
	}

	native := c.detector.GetRecommended()
	if native == nil || native.VersionInfo == nil {
		fmt.Fprintln(c.stderr, "not-installed")
		return 1
	}

	if native.VersionInfo.Compare(minVersion) < 0 {
		fmt.Fprintf(c.stdout, "upgrade-needed:%s:%s\n", native.Version, minVersion.String())
		return 1
	}

	fmt.Fprintf(c.stdout, "ok:%s\n", native.Version)
	return 0
}

func (c *CLI) cmdEnsure(args []string) int {
	minVersionStr := "2.0.0"
	if len(args) > 0 {
		minVersionStr = args[0]
	}

	minVersion, err := ParseVersion(minVersionStr)
	if err != nil {
		fmt.Fprintf(c.stderr, "Invalid version format: %s\n", minVersionStr)
		return 1
	}

	status := c.detector.GetStatus(minVersion)

	if status.NativeInstalled && !status.NeedsMigration && !status.NeedsUpgrade {
		if c.actor != ActorCI {
			fmt.Fprintf(c.stdout, "✅ Claude CLI %s is ready\n", status.ActiveVersion.String())
		} else {
			fmt.Fprintf(c.stdout, "ok:%s\n", status.ActiveVersion.String())
		}
		return 0
	}

	if c.actor != ActorCI {
		fmt.Fprintln(c.stdout, "🔧 Claude CLI needs setup...")
		fmt.Fprintln(c.stdout)
	}

	if status.NeedsMigration {
		if c.actor != ActorCI {
			fmt.Fprintln(c.stdout, "Step 1: Cleaning up old installations...")
		}
		c.doCleanup()
	}

	if !status.NativeInstalled {
		if c.actor != ActorCI {
			fmt.Fprintln(c.stdout, "Step 2: Installing native binary...")
		}
		if err := c.doInstall(); err != nil {
			fmt.Fprintf(c.stderr, "Installation failed: %v\n", err)
			return 1
		}
	}

	if status.NeedsUpgrade {
		if c.actor != ActorCI {
			fmt.Fprintln(c.stdout, "Step 3: Upgrading to latest version...")
		}
		if err := c.doUpgrade(); err != nil {
			fmt.Fprintf(c.stderr, "Upgrade failed: %v\n", err)
			return 1
		}
	}

	newStatus := c.detector.GetStatus(minVersion)
	if newStatus.NativeInstalled && !newStatus.NeedsMigration && !newStatus.NeedsUpgrade {
		if c.actor != ActorCI {
			fmt.Fprintln(c.stdout)
			fmt.Fprintf(c.stdout, "✅ Claude CLI %s is ready\n", newStatus.ActiveVersion.String())
		} else {
			fmt.Fprintf(c.stdout, "ok:%s\n", newStatus.ActiveVersion.String())
		}
		return 0
	}

	fmt.Fprintln(c.stderr, "❌ Setup failed - please check manually")
	return 1
}

func (c *CLI) doCleanup() int {
	fmt.Fprintln(c.stdout, "🧹 Cleaning up old Claude CLI installations...")
	fmt.Fprintln(c.stdout)

	cleaned := 0

	fmt.Fprintln(c.stdout, "Removing via package managers...")

	if err := exec.Command("bun", "remove", "-g", "@anthropic-ai/claude-code").Run(); err == nil {
		fmt.Fprintln(c.stdout, "  ✓ Removed via bun")
		cleaned++
	}

	if err := exec.Command("npm", "uninstall", "-g", "@anthropic-ai/claude-code").Run(); err == nil {
		fmt.Fprintln(c.stdout, "  ✓ Removed via npm")
		cleaned++
	}

	binariesToRemove := []string{
		c.detector.ExpandPath("~/.bun/bin/claude"),
		"/opt/homebrew/bin/claude",
	}

	for _, path := range binariesToRemove {
		if info, err := os.Lstat(path); err == nil {
			if info.Mode()&os.ModeSymlink != 0 || info.Mode().IsRegular() {
				if err := os.Remove(path); err == nil {
					fmt.Fprintf(c.stdout, "  ✓ Removed %s\n", path)
					cleaned++
				}
			}
		}
	}

	for _, nm := range c.detector.FindNodeModules() {
		if nm.Exists {
			if err := os.RemoveAll(nm.Path); err == nil {
				fmt.Fprintf(c.stdout, "  ✓ Removed %s\n", nm.Path)
				cleaned++
			}
		}
	}

	usrLocalClaude := "/usr/local/bin/claude"
	if info, err := os.Lstat(usrLocalClaude); err == nil && info.Mode()&os.ModeSymlink != 0 {
		if target, err := os.Readlink(usrLocalClaude); err == nil {
			if strings.Contains(target, "node_modules") {
				fmt.Fprintf(c.stdout, "  ⚠️  %s → %s (run: sudo rm %s)\n", usrLocalClaude, target, usrLocalClaude)
			}
		}
	}

	fmt.Fprintln(c.stdout)
	if cleaned > 0 {
		fmt.Fprintf(c.stdout, "✅ Cleaned up %d item(s)\n", cleaned)
	} else {
		fmt.Fprintln(c.stdout, "✅ Nothing to clean up")
	}

	return 0
}

func (c *CLI) doInstall() error {
	cmd := exec.Command("bash", "-c", "curl -fsSL https://claude.ai/install.sh | bash")
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	return cmd.Run()
}

func (c *CLI) doUpgrade() error {
	cmd := exec.Command("bash", "-c", "curl -fsSL https://claude.ai/install.sh | bash -s latest")
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr
	return cmd.Run()
}

// ============================================================================
// MCP Management Commands
// ============================================================================

func (c *CLI) handleMCP(args []string) int {
	if len(args) < 1 {
		c.printMCPUsage()
		return 0
	}

	switch args[0] {
	case "list", "ls":
		return c.mcpList()
	case "status":
		return c.mcpStatus()
	case "add":
		if len(args) < 2 {
			fmt.Fprintf(c.stderr, "Usage: claude mcp add <server>\n")
			return 1
		}
		return c.mcpAdd(args[1])
	case "enable":
		if len(args) < 2 {
			fmt.Fprintf(c.stderr, "Usage: claude mcp enable <server>\n")
			return 1
		}
		return c.mcpEnable(args[1])
	case "restart":
		serverName := ""
		if len(args) > 1 {
			serverName = args[1]
		}
		return c.mcpRestart(serverName)
	case "refresh-tokens":
		return c.mcpRefreshTokens()
	case "kill":
		serverName := ""
		if len(args) > 1 {
			serverName = args[1]
		}
		return c.mcpKill(serverName)
	default:
		fmt.Fprintf(c.stderr, "Unknown mcp command: %s\n", args[0])
		c.printMCPUsage()
		return 1
	}
}

func (c *CLI) printMCPUsage() {
	fmt.Fprintln(c.stdout, `Usage: claude mcp <command> [arguments]

Commands:
  list                    List configured MCP servers
  status                  Show running MCP servers
  add <server>            Add MCP server to .vscode/mcp.json
  enable <server>         Enable MCP server in ~/.claude.json for this project
  restart [server]        Restart MCP server(s) by killing processes
  refresh-tokens          Refresh Google tokens and restart server
  kill [server]           Kill MCP server process(es)

Examples:
  claude mcp list
  claude mcp status
  claude mcp enable google
  claude mcp restart google
  claude mcp kill google
  claude mcp refresh-tokens`)
}

func (c *CLI) mcpList() int {
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to get working directory: %v\n", err)
		return 1
	}

	for name, loc := range Locations {
		var configPath string
		if loc.Dir != "" {
			configPath = filepath.Join(projectRoot, loc.Dir, loc.File)
		} else {
			configPath = filepath.Join(projectRoot, loc.File)
		}

		config, err := LoadConfig(configPath)
		if err != nil {
			continue
		}

		if config.IsEmpty() {
			continue
		}

		fmt.Fprintf(c.stdout, "\n%s (%s):\n", name, configPath)
		for serverName, server := range config.MCPServers {
			fmt.Fprintf(c.stdout, "  📦 %s\n", serverName)
			fmt.Fprintf(c.stdout, "      Command: %s\n", server.Command)
			if len(server.Args) > 0 {
				fmt.Fprintf(c.stdout, "      Args: %v\n", server.Args)
			}
			if len(server.Env) > 0 {
				fmt.Fprintf(c.stdout, "      Env: %d variables\n", len(server.Env))
			}
		}
	}

	return 0
}

func (c *CLI) mcpAdd(serverName string) int {
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to get working directory: %v\n", err)
		return 1
	}

	// Add server to VSCode config
	result, err := AddMCPServer(serverName, TargetVSCode, projectRoot)
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: %v\n", err)
		return 1
	}

	if result.ServerAdded {
		fmt.Fprintf(c.stdout, "✅ Added %s to %s\n", serverName, result.ConfigPath)
	} else {
		fmt.Fprintf(c.stdout, "ℹ️  %s already configured in %s\n", serverName, result.ConfigPath)
	}

	if result.PermissionsSet {
		fmt.Fprintf(c.stdout, "✅ Added permissions to %s\n", result.SettingsPath)
	}

	if result.BackupPath != "" {
		fmt.Fprintf(c.stdout, "📦 Backup saved to %s\n", result.BackupPath)
	}

	return 0
}

func (c *CLI) mcpEnable(serverName string) int {
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to get working directory: %v\n", err)
		return 1
	}

	// Enable server in ~/.claude.json by adding to enabledMcpjsonServers
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to get home directory: %v\n", err)
		return 1
	}

	claudeConfigPath := filepath.Join(home, ".claude.json")

	// Read current config
	data, err := os.ReadFile(claudeConfigPath)
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to read %s: %v\n", claudeConfigPath, err)
		return 1
	}

	var claudeConfig map[string]interface{}
	if err := json.Unmarshal(data, &claudeConfig); err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to parse %s: %v\n", claudeConfigPath, err)
		return 1
	}

	// Get or create projects map
	projects, ok := claudeConfig["projects"].(map[string]interface{})
	if !ok {
		projects = make(map[string]interface{})
		claudeConfig["projects"] = projects
	}

	// Get or create project entry
	projectConfig, ok := projects[projectRoot].(map[string]interface{})
	if !ok {
		projectConfig = make(map[string]interface{})
		projects[projectRoot] = projectConfig
	}

	// Get or create enabledMcpjsonServers array
	var enabledServers []string
	if existing, ok := projectConfig["enabledMcpjsonServers"].([]interface{}); ok {
		for _, s := range existing {
			if str, ok := s.(string); ok {
				enabledServers = append(enabledServers, str)
			}
		}
	}

	// Check if already enabled
	alreadyEnabled := false
	for _, s := range enabledServers {
		if s == serverName {
			alreadyEnabled = true
			break
		}
	}

	if alreadyEnabled {
		fmt.Fprintf(c.stdout, "ℹ️  %s already enabled for this project\n", serverName)
		return 0
	}

	// Add server to enabled list
	enabledServers = append(enabledServers, serverName)
	projectConfig["enabledMcpjsonServers"] = enabledServers

	// Write back
	newData, err := json.MarshalIndent(claudeConfig, "", "  ")
	if err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to marshal config: %v\n", err)
		return 1
	}

	if err := os.WriteFile(claudeConfigPath, newData, 0644); err != nil {
		fmt.Fprintf(c.stderr, "Error: Failed to write %s: %v\n", claudeConfigPath, err)
		return 1
	}

	fmt.Fprintf(c.stdout, "✅ Enabled %s for project %s\n", serverName, projectRoot)
	fmt.Fprintln(c.stdout, "\n💡 Restart Claude Code to load the server")

	return 0
}

type processInfo struct {
	pid  int
	name string
}

func (c *CLI) findMCPProcesses() []processInfo {
	cmd := exec.Command("pgrep", "-fl", "mcp-server")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var processes []processInfo
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) >= 2 {
			var pid int
			fmt.Sscanf(parts[0], "%d", &pid)
			processes = append(processes, processInfo{pid: pid, name: parts[1]})
		}
	}

	return processes
}

func (c *CLI) mcpStatus() int {
	processes := c.findMCPProcesses()

	if len(processes) == 0 {
		fmt.Fprintln(c.stdout, "No MCP server processes found")
		return 0
	}

	fmt.Fprintln(c.stdout, "Running MCP servers:")
	for _, p := range processes {
		fmt.Fprintf(c.stdout, "  🟢 PID %d: %s\n", p.pid, p.name)
	}

	return 0
}

func (c *CLI) mcpRestart(serverName string) int {
	fmt.Fprintln(c.stdout, "🔄 Restarting MCP servers...")

	killed := c.killMCPProcesses(serverName)

	if killed == 0 {
		fmt.Fprintln(c.stdout, "No MCP server processes found to restart")
		fmt.Fprintln(c.stdout, "\n💡 Tip: The MCP server will be restarted automatically by Claude Code on next use")
	} else {
		fmt.Fprintf(c.stdout, "✅ Killed %d MCP server process(es)\n", killed)
		fmt.Fprintln(c.stdout, "\n💡 Claude Code will restart the server automatically on next MCP tool call")
	}

	return 0
}

func (c *CLI) mcpKill(serverName string) int {
	killed := c.killMCPProcesses(serverName)

	if killed == 0 {
		fmt.Fprintln(c.stdout, "No matching MCP server processes found")
	} else {
		fmt.Fprintf(c.stdout, "✅ Killed %d MCP server process(es)\n", killed)
	}

	return 0
}

func (c *CLI) killMCPProcesses(serverName string) int {
	var pattern string
	if serverName != "" {
		pattern = serverName + "-mcp-server"
	} else {
		pattern = "mcp-server"
	}

	cmd := exec.Command("pkill", "-f", pattern)
	if err := cmd.Run(); err != nil {
		return 0
	}

	time.Sleep(100 * time.Millisecond)

	remaining := c.findMCPProcesses()
	if len(remaining) == 0 {
		return 1
	}

	killedCount := 0
	for _, p := range remaining {
		if serverName == "" || strings.Contains(p.name, serverName) {
			killedCount++
		}
	}

	return killedCount
}

func (c *CLI) mcpRefreshTokens() int {
	fmt.Fprintln(c.stdout, "🔄 Refreshing Google tokens...")

	home, _ := os.UserHomeDir()
	tokenDir := filepath.Join(home, ".google-mcp-accounts")

	files, err := filepath.Glob(filepath.Join(tokenDir, "*.json"))
	if err != nil || len(files) == 0 {
		fmt.Fprintln(c.stderr, "Error: No Google accounts found. Run 'task google-mcp:auth' first")
		return 1
	}

	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}

		var account struct {
			Email string `json:"email"`
			Token struct {
				Expiry string `json:"expiry"`
			} `json:"token"`
		}

		if err := json.Unmarshal(data, &account); err != nil {
			continue
		}

		fmt.Fprintf(c.stdout, "  📧 %s (expires: %s)\n", account.Email, account.Token.Expiry)
	}

	for _, f := range files {
		now := time.Now()
		os.Chtimes(f, now, now)
	}

	fmt.Fprintln(c.stdout, "\n✅ Token files touched")
	fmt.Fprintln(c.stdout, "\n🔄 Restarting google-mcp-server...")

	killed := c.killGoogleMCPServer()
	if killed {
		fmt.Fprintln(c.stdout, "✅ Google MCP server killed")
		fmt.Fprintln(c.stdout, "\n💡 Claude Code will restart it automatically on next Google tool call")
	} else {
		fmt.Fprintln(c.stdout, "ℹ️  No google-mcp-server process found")
	}

	return 0
}

func (c *CLI) killGoogleMCPServer() bool {
	cmd := exec.Command("pgrep", "-f", "google-mcp-server")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	pids := strings.Fields(strings.TrimSpace(string(output)))
	if len(pids) == 0 {
		return false
	}

	for _, pidStr := range pids {
		var pid int
		fmt.Sscanf(pidStr, "%d", &pid)
		if pid > 0 {
			syscall.Kill(pid, syscall.SIGTERM)
		}
	}

	return true
}

// ============================================================================
// Helpers
// ============================================================================

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
