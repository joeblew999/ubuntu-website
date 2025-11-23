# Ubuntu Software Website

Multi-language website built with Hugo Plate, featuring automated translation and deployment to Cloudflare Pages.

## 🔗 Quick Links

For all URLs, dashboards, and project information, run:

```bash
task          # or: task info
```

This will show:
- Development server URL (local)
- Production and preview URLs
- Cloudflare dashboards
- Quick commands

You can also open any URL directly in your browser:

```bash
task url:dev      # Open local dev server
task url:prod     # Open production site
task url:preview  # Open preview site
task cf:open      # Open Cloudflare dashboard
```





## 🌍 Supported Languages

- 🇬🇧 English (en) - Source language
- 🇩🇪 German (de)
- 🇸🇪 Swedish (sv)
- 🇨🇳 Chinese Simplified (zh)
- 🇯🇵 Japanese (ja)
- 🇹🇭 Thai (th)

## 🚀 Quick Start

```bash
# 1. Environment Setup (first time only - or anytime to sync)
task env:all                      # Complete env workflow (recommended)
task setup                        # Install Hugo, Bun, deps

# 2. Develop
task dev                          # Start dev server

# 3. Translate
task translate:all                # Translate changed content

# 4. Deploy
task cf:deploy                    # Deploy to Cloudflare Pages
```

## 🔑 Environment Setup

The recommended way to manage environment configuration:

```bash
task env:all  # Complete workflow (idempotent - safe to re-run anytime)
```

This runs the complete unidirectional flow: `setup → list → push → verify`

Individual commands (for advanced use):

```bash
task env:local:setup  # Setup local .env (interactive wizard)
task env:local:list   # List local .env configuration
task env:gh:list      # List GitHub secrets
task env:gh:push      # Push to GitHub secrets for CI/CD
```

All commands are idempotent - safe to run multiple times without side effects.

## 📋 Prerequisites

- **Go** 1.24+ (for Hugo and translation tool)
- **Bun** latest (JavaScript runtime)
- **Task** ([install](https://taskfile.dev/installation/))

## 🛠️ Development Tasks

```bash
task dev           # Start dev server
task build         # Build production
task preview       # Preview production
task clean         # Clean artifacts

task translate:check     # Check changed files
task translate:all       # Translate all
task translate:lang LANG=de  # Translate to German

task cf:deploy     # Deploy to Cloudflare
task cf:status     # Check status
```

## 🔒 HTTPS Development

The project uses **Caddy + mkcert** for local HTTPS development:

- **Hugo**: `https://localhost/` or `https://192.168.x.x/`
- **Via GUI**: `https://localhost/admin/` or `https://192.168.x.x/admin/`

### How It Works

1. **Automatic Setup**: Caddy starts automatically when you run the dev server or Via GUI
2. **Smart Certificates**: mkcert generates certificates for localhost + *.local + your current LAN IP
3. **Idempotent**: Safe to run multiple times - certificates regenerate only when your LAN IP changes
4. **Mobile Testing**: Works on iOS Safari without manual CA installation (just accept the certificate prompt)

### Certificate Management

Certificates are stored in `.caddy/certs/` (gitignored) and automatically regenerated when:
- Certificates don't exist
- Your LAN IP address changes (switching networks)

### Mobile Device Setup (iOS/Android)

**iOS Safari**: No installation needed! Safari will show a trust prompt the first time you visit the HTTPS site.

**iOS/Android (for apps or other browsers)**:
1. Install mkcert CA certificate on your mobile device
2. Find the CA at: `~/Library/Application Support/mkcert/rootCA.pem`
3. Transfer to your device via AirDrop, email, or file sharing
4. Install the certificate and trust it in Settings

### Architecture

```
┌─────────────────────────────────────────────┐
│  Caddy (HTTPS on port 443)                  │
│  - localhost, *.local, 192.168.x.x         │
│  - mkcert certificates                      │
└──────────┬──────────────────────┬───────────┘
           │                      │
           ▼                      ▼
  ┌────────────────┐    ┌─────────────────┐
  │ Hugo           │    │ Via GUI         │
  │ (HTTP :1313)   │    │ (HTTP :3000)    │
  │ /              │    │ /admin/*        │
  └────────────────┘    └─────────────────┘
```

### CLI Commands

```bash
# Start the system
go run cmd/env/main.go web-gui          # Starts Caddy + Via GUI
go run cmd/env/main.go build            # Starts Caddy + Hugo

# Manage Caddy manually (if needed)
go run cmd/env/main.go caddy-start      # Start Caddy
go run cmd/env/main.go caddy-stop       # Stop Caddy
go run cmd/env/main.go caddy-status     # Check if running

# Cleanup
# Press Ctrl+C in the web-gui terminal - auto-cleanup of Caddy and Hugo
```

## 🌐 Translation Workflow

1. Edit English content in `content/english/`
2. Commit changes: `git commit -m "Update: page"`
3. Check: `task translate:check`
4. Translate: `task translate:all`
5. Deploy: `git push` (auto-deploys via CI/CD)

The translation tool:
- Uses Git to track changes
- Only translates modified files
- Preserves Hugo shortcodes & code blocks
- Calls Claude API for translations
- Creates checkpoints automatically

## 📦 Project Structure

```
.
├── cmd/env/                # Environment setup CLI
├── cmd/translate/          # Translation CLI
├── internal/env/           # Environment management
├── internal/translator/    # Translation logic
├── content/
│   ├── english/           # Source (EN)
│   ├── german/            # DE translations
│   ├── swedish/           # SV translations
│   ├── chinese/           # ZH translations
│   ├── japanese/          # JA translations
│   └── thai/              # TH translations
├── config/_default/
│   ├── languages.toml     # Language config
│   └── menus.*.toml       # Per-language menus
├── i18n/                  # Translation strings
├── .github/workflows/     # CI/CD
├── Taskfile.yml           # Task definitions
└── wrangler.toml          # Cloudflare config
```

## 🔧 Configuration

The `.env` file contains your API credentials:

```bash
# Cloudflare credentials (for deployment)
CLOUDFLARE_API_TOKEN=your-token-here
CLOUDFLARE_ACCOUNT_ID=your-account-id
CLOUDFLARE_PROJECT_NAME=your-project-name

# Claude API key (for translation)
CLAUDE_API_KEY=your-api-key-here
```

**Setup:** Run `task env:local:setup` for interactive wizard, then `task env:gh:push` to sync to GitHub for CI/CD.

## 🔗 References

- [Hugo Plate Template](https://github.com/zeon-studio/hugoplate)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [Cloudflare Pages Docs](https://developers.cloudflare.com/pages/)
- [Task Documentation](https://taskfile.dev/)

## 📄 License

MIT (Hugo Plate template)





