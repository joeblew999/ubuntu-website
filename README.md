# Ubuntu Software Website

Multi-language website built with Hugo Plate, featuring automated translation and deployment to Cloudflare Pages.

**Website:** https://www.ubuntusoftware.net
**Dashboard:** https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/ubuntusoftware.net

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





