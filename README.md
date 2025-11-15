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
# 1. Setup
cp .env.example .env              # Add your API keys
task setup                        # Install Hugo, Bun, deps

# 2. Develop
task dev                          # Start dev server

# 3. Translate
task translate:all                # Translate changed content

# 4. Deploy
task cf:deploy                    # Deploy to Cloudflare Pages
```

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
├── cmd/translate/          # Go translation CLI
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

### Environment Variables (.env)

```bash
CLOUDFLARE_API_TOKEN=xxx    # From dash.cloudflare.com
CLOUDFLARE_ACCOUNT_ID=xxx   # In dashboard URL
CLAUDE_API_KEY=xxx          # From console.anthropic.com
```

### GitHub Secrets

Add to repo Settings → Secrets:
- `CLOUDFLARE_API_TOKEN`
- `CLOUDFLARE_ACCOUNT_ID`
- `CLAUDE_API_KEY`

## 🔗 References

- [Hugo Plate Template](https://github.com/zeon-studio/hugoplate)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [Cloudflare Pages Docs](https://developers.cloudflare.com/pages/)
- [Task Documentation](https://taskfile.dev/)

## 📄 License

MIT (Hugo Plate template)





