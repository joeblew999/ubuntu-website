# Claude Code & Google MCP Integration Plan

## Current Status: Ready for Testing

The `cmd/claude` CLI and Google MCP server have been fixed and are ready for testing.

## What Was Fixed (2025-12-19)

### 1. Fixed `.oauth_client.json` being loaded as account file
- **File**: `.src/google-mcp-server/auth/account_manager.go`
- **Issue**: The `.oauth_client.json` file (OAuth credentials) was being parsed as an account, causing "refresh token is not set" errors
- **Fix**: Skip hidden files (starting with `.`) and files without a valid token

### 2. Fixed browser opening on startup when accounts exist
- **File**: `.src/google-mcp-server/main.go`
- **Issue**: Even with valid multi-account tokens, the server tried to create a backward-compat OAuth client, opening a browser
- **Fix**: Skip backward-compat OAuth client creation when AccountManager already has accounts

### 3. Rebuilt google-mcp-server binary
- **Location**: `.build/google-mcp-server`
- **Timestamp**: 2025-12-19 16:26

## Architecture

```
cmd/claude/                     # CLI for Claude Code management
  main.go                       # Entry point

internal/claude/                # Core logic (moved from internal/mcp/)
  cli.go                        # CLI commands including mcp subcommands
  detect.go                     # Claude CLI binary detection
  google.go                     # Google MCP server definition
  mcpconfig.go                  # MCP config file management (.vscode/mcp.json)
  mcpserver.go                  # MCP server registry
  settings.go                   # .claude/settings.json management
  targets.go                    # Claude target definitions (VSCode, Desktop, etc.)

.src/google-mcp-server/         # Google MCP server source
  main.go                       # Server entry point
  auth/                         # OAuth and account management
    account_manager.go          # Multi-account token management
    oauth.go                    # OAuth client
  calendar/                     # Google Calendar integration
  drive/                        # Google Drive integration
  gmail/                        # Gmail integration
  sheets/                       # Google Sheets integration
  docs/                         # Google Docs integration
  slides/                       # Google Slides integration

.build/google-mcp-server        # Compiled binary
.vscode/mcp.json                # MCP server config for Claude Code
.claude/settings.json           # Permissions for MCP tools
~/.google-mcp-accounts/         # OAuth tokens (outside project)
```

## Next Steps

### Immediate: Test the Fix
1. **Restart Claude Code** (close and reopen VSCode)
2. Run: `go run ./cmd/claude mcp status` - should show google-mcp-server
3. Try using a Google MCP tool (e.g., `mcp__google__accounts_list`)

### If Token Expired
The current token expires at `2025-12-19T16:38:39+07:00`. If expired:
```bash
# Re-authenticate (will open browser)
task google-mcp:auth
```

### Future Enhancements

#### cmd/claude improvements
- [ ] Add `mcp auth google` command to trigger OAuth flow
- [ ] Add `mcp auth status` to check token expiry for all accounts
- [ ] Add `mcp auth refresh` to force token refresh via Google API
- [ ] Add `mcp setup google` for turnkey Google MCP setup

#### Google MCP Server improvements
- [ ] Automatic token refresh before expiry
- [ ] Better error messages for expired tokens
- [ ] Health check endpoint

#### Documentation
- [ ] Update CLAUDE.md with cmd/claude usage
- [ ] Add troubleshooting section for common MCP issues

## Quick Reference

### Check Claude CLI status
```bash
go run ./cmd/claude status
```

### Check MCP servers
```bash
go run ./cmd/claude mcp list      # Show configured servers
go run ./cmd/claude mcp status    # Show running servers
```

### Restart Google MCP server
```bash
go run ./cmd/claude mcp restart google
```

### Rebuild google-mcp-server
```bash
cd .src/google-mcp-server && go build -o ../../.build/google-mcp-server .
```

### Full cleanup and fresh start
```bash
# Remove all tokens
rm -rf ~/.google-mcp-accounts/*.json

# Re-authenticate
task google-mcp:auth

# Restart Claude Code
```

## Files Changed (uncommitted)

From git status:
- `internal/claude/` - New location for MCP code (moved from internal/mcp/)
- `internal/mcp/` - Deleted (moved to internal/claude/)
- `.src/google-mcp-server/auth/account_manager.go` - Fixed account loading
- `.src/google-mcp-server/main.go` - Fixed browser popup issue
