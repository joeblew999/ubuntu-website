package claude

import (
	"runtime"
	"testing"
)

func TestGetTargetInfo(t *testing.T) {
	tests := []struct {
		target      Target
		wantName    string
		wantForDevs bool
		wantForUser bool
	}{
		{TargetVSCode, "Claude Code (VSCode)", true, false},
		{TargetProject, "Project MCP", true, false},
		{TargetClaude, "Claude Folder", true, false},
		{TargetDesktop, "Claude Desktop", false, true},
		{TargetUserGlobal, "User Global", true, true},
		{TargetCloud, "Claude Cloud", true, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.target), func(t *testing.T) {
			info := GetTargetInfo(tt.target)
			if info.Name != tt.wantName {
				t.Errorf("GetTargetInfo(%s).Name = %s, want %s", tt.target, info.Name, tt.wantName)
			}
			if info.ForDevs != tt.wantForDevs {
				t.Errorf("GetTargetInfo(%s).ForDevs = %v, want %v", tt.target, info.ForDevs, tt.wantForDevs)
			}
			if info.ForUsers != tt.wantForUser {
				t.Errorf("GetTargetInfo(%s).ForUsers = %v, want %v", tt.target, info.ForUsers, tt.wantForUser)
			}
		})
	}
}

func TestGetDesktopConfigPath(t *testing.T) {
	path, err := GetDesktopConfigPath()
	if err != nil {
		t.Fatalf("GetDesktopConfigPath() error = %v", err)
	}

	// Check path ends with expected filename
	wantSuffix := "claude_desktop_config.json"
	if len(path) < len(wantSuffix) {
		t.Fatalf("GetDesktopConfigPath() = %s, too short", path)
	}
	gotSuffix := path[len(path)-len(wantSuffix):]
	if gotSuffix != wantSuffix {
		t.Errorf("GetDesktopConfigPath() = %s, want suffix %s", path, wantSuffix)
	}

	// Check OS-specific path components
	switch runtime.GOOS {
	case "darwin":
		if !contains(path, "Library/Application Support/Claude") {
			t.Errorf("GetDesktopConfigPath() = %s, want to contain Library/Application Support/Claude", path)
		}
	case "linux":
		if !contains(path, ".config/Claude") {
			t.Errorf("GetDesktopConfigPath() = %s, want to contain .config/Claude", path)
		}
	}
}

func TestAllTargets(t *testing.T) {
	targets := AllTargets()
	if len(targets) != 6 {
		t.Errorf("AllTargets() returned %d targets, want 6", len(targets))
	}
}

func TestDevTargets(t *testing.T) {
	targets := DevTargets()
	for _, target := range targets {
		info := GetTargetInfo(target)
		if !info.ForDevs {
			t.Errorf("DevTargets() contains %s which is not for devs", target)
		}
	}
}

func TestUserTargets(t *testing.T) {
	targets := UserTargets()
	for _, target := range targets {
		info := GetTargetInfo(target)
		if !info.ForUsers {
			t.Errorf("UserTargets() contains %s which is not for users", target)
		}
	}
}

func TestGetTargetConfigPath(t *testing.T) {
	projectRoot := "/test/project"

	tests := []struct {
		target   Target
		wantPath string
		wantErr  bool
	}{
		{TargetVSCode, "/test/project/.vscode/mcp.json", false},
		{TargetProject, "/test/project/.mcp.json", false},
		{TargetClaude, "/test/project/.claude/mcp.json", false},
		{TargetCloud, "", false}, // No local config
	}

	for _, tt := range tests {
		t.Run(string(tt.target), func(t *testing.T) {
			path, err := GetTargetConfigPath(tt.target, projectRoot)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTargetConfigPath(%s) error = %v, wantErr %v", tt.target, err, tt.wantErr)
			}
			if path != tt.wantPath {
				t.Errorf("GetTargetConfigPath(%s) = %s, want %s", tt.target, path, tt.wantPath)
			}
		})
	}
}

// helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
