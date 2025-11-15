# Claude Assistant Notes

## CRITICAL REMINDERS

### Path Convention
**ALWAYS use `joeblew999` (with three 9s), NEVER `joeblew99` (with two 9s)**

Correct paths:
- `/Users/apple/workspace/go/src/github.com/joeblew999/ubuntu-website`
- `github.com/joeblew999/ubuntu-website`

### Project Structure

```
ubuntu-website/
├── cmd/translate/          # Go translation CLI tool
│   └── main.go            # Entry point, version: 0.1.0
├── internal/translator/    # Translation logic
│   ├── translator.go      # Main translator
│   ├── claude.go          # Claude API client
│   ├── git.go             # Git tracking
│   └── markdown.go        # MD parsing
├── content/
│   ├── english/           # Source (EN)
│   ├── german/            # DE translations
│   ├── swedish/           # SV translations
│   ├── chinese/           # ZH translations
│   ├── japanese/          # JA translations
│   └── thai/              # TH translations
├── config/_default/
│   ├── languages.toml     # 6 languages configured
│   └── menus.*.toml       # Per-language menus
├── i18n/                  # Translation strings (YAML)
├── .github/workflows/
│   └── deploy.yml         # CI/CD pipeline
├── Taskfile.yml           # Task runner (25 tasks)
├── wrangler.toml          # Cloudflare Pages config
├── .env.example           # Environment template
└── go.mod                 # Go 1.24, yaml.v3
```

## Technology Stack

- **Hugo**: v0.152.2+extended (installed via Go)
- **Go**: v1.25.4 (requires 1.24+)
- **Bun**: v1.3.2 (JavaScript runtime, NOT npm)
- **Task**: v3.45.5 (task runner)
- **Wrangler**: v4.22.0 (Cloudflare CLI)
- **Git**: Configured as joeblew999

## Key Variables (Taskfile.yml)

```yaml
PROJECT_NAME: ubuntusoftware-net
DOMAIN: www.ubuntusoftware.net
TRANSLATE_CMD: cmd/translate/main.go
DIR_PUBLIC: public
DIR_RESOURCES: resources
DIR_NODE_MODULES: node_modules
DIR_BIN: bin
HUGO_BUILD_FLAGS: --gc --minify
HUGO_DEV_FLAGS: --buildDrafts --buildFuture
GO_VERSION: "1.24"
CGO_ENABLED: "1"
```

## Available Tasks

### Setup & Development
- `task setup` - Complete setup (Hugo, Bun, deps)
- `task dev` - Start Hugo dev server
- `task build` - Build production site
- `task preview` - Preview production build
- `task clean` - Clean build artifacts

### Translation
- `task translate:check` - Check which files need translation
- `task translate:all` - Translate all changed content
- `task translate:lang LANG=de` - Translate to specific language
- `task translate:i18n` - Translate i18n files

### Cloudflare
- `task cf:login` - Login to Cloudflare
- `task cf:init` - Initialize Cloudflare Pages project
- `task cf:deploy` - Build and deploy to Cloudflare Pages
- `task cf:status` - Check deployment status
- `task cf:delete` - Delete Cloudflare Pages project

### CI/CD
- `task ci:setup` - Setup CI environment
- `task ci:build` - Build in CI
- `task ci:deploy` - Deploy from CI

### Code Quality
- `task fmt` - Format Go code
- `task lint` - Lint Go code
- `task test` - Run Go tests

### Utilities
- `task generate:gitignore` - Generate .gitignore from vars
- `task help` - Show all tasks

## Environment Variables (.env)

Required for deployment and translation:
- `CLOUDFLARE_API_TOKEN` - From https://dash.cloudflare.com/profile/api-tokens
- `CLOUDFLARE_ACCOUNT_ID` - 7384af54e33b8a54ff240371ea368440
- `CLOUDFLARE_PROJECT_NAME` - ubuntusoftware-net
- `CLAUDE_API_KEY` - From https://console.anthropic.com/settings/keys

Optional:
- `GO_VERSION` - 1.24
- `BUN_VERSION` - latest
- `HUGO_VERSION` - latest

## Language Configuration

6 languages supported (config/\_default/languages.toml):
1. English (en-us) - Source, weight 1
2. German (de-de) - weight 2
3. Swedish (sv) - weight 3
4. Chinese (zh-cn) - weight 4, hasCJKLanguage: true
5. Japanese (ja) - weight 5, hasCJKLanguage: true
6. Thai (th) - weight 6

## Translation Workflow

1. Edit English content in `content/english/`
2. Commit changes: `git commit -m "Update: page"`
3. Check: `task translate:check`
4. Translate: `task translate:all`
5. Deploy: `git push` (auto-deploys via CI/CD)

Translation tool:
- Uses Git to track changes since last checkpoint
- Only translates modified files
- Preserves Hugo shortcodes & code blocks
- Calls Claude API for translations
- Creates checkpoints automatically with git tags

## Taskfile Philosophy

- Simple, clean structure (no emojis or fancy UX)
- Variable-driven configuration
- Hugo installed via Go (`go install -tags extended`)
- Uses Bun (not npm)
- Single source of truth for all operations
- Tasks use `deps:` for dependencies, not `task:` in cmds

## Common Issues to Avoid

1. **Path mistakes**: Always use `joeblew999` (three 9s)
2. **YAML syntax**: Use `deps: [build]` not `task: build` in cmds
3. **String formatting**: Avoid backticks in multi-line strings for Go code
4. **Hugo modules**: Uses Hugo modules, not themes
5. **CJK languages**: Chinese and Japanese need `hasCJKLanguage: true`

## Git Configuration

- User: joeblew999
- Email: joeblew999@users.noreply.github.com
- Current branch: main
- Main branch for PRs: main

## Deployment

- **Target**: Cloudflare Pages
- **Domain**: www.ubuntusoftware.net
- **Dashboard**: https://dash.cloudflare.com/7384af54e33b8a54ff240371ea368440/ubuntusoftware.net
- **Build command**: `task ci:setup && task ci:build`
- **Output dir**: public
- **Production branch**: main

## Testing Status

All systems tested and verified:
- ✓ Hugo build successful (52 EN pages + 9 per language)
- ✓ Translation tool compiles and runs
- ✓ All 25 Taskfile tasks functional
- ✓ Multi-language URLs working
- ✓ All CF tasks validated

## Missing/TODO

- `.env` file needs to be created from `.env.example`
- API keys need to be added to `.env`:
  - CLOUDFLARE_API_TOKEN
  - CLAUDE_API_KEY
