# Claude Assistant Notes

## CRITICAL REMINDERS

**Production Domain:** `www.ubuntusoftware.net`

Apex domain (`ubuntusoftware.net`) redirects to www via Cloudflare redirect rule.

USE MISE: every task is `mise run <task>` (`mise tasks` lists them), and the root `mise.toml` is the control plane. There are no Taskfiles and no scripts.

### Tools: mise installs and runs

- Every third-party tool (hugo, node, aube, process-compose, gh, terraform, google-mcp-server, fnox, jq, git-lfs, go-mod-upgrade) is pinned in `mise.toml`; wrangler is pinned in `package.json` (aube installs it; mise-cloudflare requires it there). Go is pinned by the `toolchain` line in `go.mod` (its `go` line is the minimum language version). CI installs the same pins via `jdx/mise-action`; Cloudflare Pages installs only the build's (`root:pages`).
- Tasks call tools by name and never install them (no `brew install`, `curl | bash`, `go install` for third-party tools): mise installs every pin (`mise install`, or `mise run root:setup`).
- New tool → add it to `[tools]` in `mise.toml` with an exact version.

### mise tasks: the layout

- **`mise.toml` is the control plane**: tool pins, settings, env (`bin/` on PATH), tasks that act on the whole repo (`root:*`), and `[task_config] includes = [".mise/tasks"]`.
- **CI runs the control plane, nothing else** (SSOT): `.github/workflows/core-ci.yml` is `mise run root:ci`, and `root:ci` depends on every check (`root:check-mise`, `root:check-workflows` (actionlint), `root:check-actions` (pinact: actions pinned to SHAs), `root:test`, `root:test-worker`, `theme:check`, `root:build`, `root:check-taskfiles`). A new check is a task `root:ci` depends on; `mise run root:ci` locally does exactly what CI does.
- **One task file per Go binary: `.mise/tasks/<binary>.toml`.** Tasks are `<binary>:<command>`, each one line running `<binary> <command>`, i.e. `cmd/<binary>/<command>.go` (a file per command, holding the function named for it); `<binary>:build` builds it into `bin/`, and the others `depends` on it. Flags are declared in the task's `usage` and passed by name (`{% if usage.x %} --x{% endif %}`); a binary whose commands parse their own flags in `internal/` takes no `usage`, and mise passes the task's arguments through (`mise run translate:content missing --github-issue`; `internal/cli.Interspersed` lets a binary's global flags follow its command). No shell syntax; logic in Go, not in the TOML. A domain that only runs tools gets a task file of its own too (`.mise/tasks/hugo.toml`: `hugo:dev` = `hugo server …`; `dev.toml`: `dev:up` = process-compose), held to the rules every task is (mise-check ≥ v0.3.0; given no `GO_CMD`, a file of tool tasks passes on them).
- **misecheck holds it together** (`github.com/joeblew999/mise-check`; `mise_test.go`, `mise run root:check-mise`, in CI): `mise.toml` and every `.mise/tasks/*.toml`, each task file strictly against its own binary, both ways: every task's command is there, every command is run by a task, every flag on both sides, every task name the binary's Go prints is a task. A new task file is checked as soon as it exists.
- Examples: `.mise/tasks/google-mcp.toml` ↔ `cmd/google-mcp/`, `.mise/tasks/theme.toml` ↔ `cmd/theme/`, `.mise/tasks/sitecheck.toml` ↔ `cmd/sitecheck/` (its scheduled workflow, `monitor-sitecheck.yml`, runs `mise run sitecheck:issue`).
- **No scripts, no shell logic**: there is no `scripts/` folder. Where each piece goes: a step on the whole repo that only runs a tool → a `root:*` task in `mise.toml` (`root:build` is `hugo …`); anything with logic → Go, a command of its domain's binary (`cmd/<binary>/<command>.go`) with its task in that binary's `.mise/tasks/<binary>.toml` (the theme scripts became `cmd/theme/update.go`, `css.go`); a file from upstream that we run unmodified sits beside the Go that runs it (`cmd/theme/themeGenerator.js`). Taskfile shell moves the same way as each domain migrates.
- **The site is a Cloudflare Worker, run by mise-cloudflare** (moving off Pages, Cloudflare's dead end; TODO.md has the steps left): `wrangler.jsonc` serves the Hugo build (`public/`) as static assets and `worker/index.js` for `/api/*` (`/api/health`, `/api/webhook/mailerlite`, Web3Forms → MailerLite). mise-cloudflare (pinned in `.mise/config.toml`; `mise run cloudflare` says where it is and what's next; `docs/mise-cloudflare/README.md`; the `cloudflare` skill) owns its login (a wrangler profile, `ubuntusoftware-net`: no API token), its secret (`MAILERLITE_API_KEY`), deploys and previews; never run wrangler bare, run it through `mise run cloudflare:run-wrangler`. `package.json`'s `build` script is `mise run root:build`, which mise-cloudflare runs before serving and deploying. `static/.assetsignore` keeps the local airspace data and video cache out of the assets (25 MiB a file). Until the domain moves, `www.ubuntusoftware.net` is still Cloudflare Pages, whose build command (project setting) is `MISE_TRUSTED_CONFIG_PATHS=$PWD npx -y @jdxcode/mise@<min_version> run --skip-tools root:pages` (hugo-extended, node and aube only, then `root:build`). `min_version` in `mise.toml` is the mise version CI and Cloudflare run.
- **Cloudflare tokens**: none for the Worker (the repo's login). This repo's own tools that call Cloudflare's API (cfanalytics) read `UBUNTU_WEBSITE_CF_TOKEN` from `.env`, named so no Cloudflare tool reads it: wrangler takes `CLOUDFLARE_API_TOKEN` and `CF_API_TOKEN` over any login, and mise-cloudflare refuses either. The zone's cache and Email Routing are Cloudflare's own CLI, `cf` (`cf:*` tasks, on `cf:login`).

**Upgrading** (`.mise/tasks/upgrade.toml`, `cmd/upgrade`):

| What | Command |
|------|---------|
| **Everything except Hugo** (every row but Hugo's, in order) | `mise run upgrade:all` |
| Tools in `mise.toml` + Go (`toolchain` line in `go.mod`) | `mise run upgrade:tools` |
| Go module dependencies (`go-mod-upgrade`, incl. major versions; vet/test after) | `mise run upgrade:modules` (asks which) or `--force` |
| Hugo (the mise pin; Cloudflare builds with it too) | `mise run upgrade:hugo` → branch + Cloudflare preview → main |
| Theme: Hugo Plate + gethugothemes modules, and the npm packages they need | `mise run theme:update` |
| GitHub Actions in `.github/workflows` (pinned to SHAs; releases at least 7 days old, `.github/pinact.yaml`) | `mise run upgrade:actions` |

Each is safe to re-run: with nothing new upstream it changes nothing.
- New machine: `mise trust && mise run root:setup` (every pinned tool, the npm packages, the Go modules).
- **npm packages: aube, node the runtime.** aube (pinned in `mise.toml`) installs the site's npm packages from `aube-lock.yaml` (`mise run root:deps` = `aube install --frozen-lockfile`) and mise's `npm:` tools (wrangler) with its embedded aube; node is the one JavaScript runtime. No bun, no npm installs. `aube-workspace.yaml` also turns off aube's global virtual store: Hugo runs Tailwind under Node's permission model, which reads only inside the project. Dependency build scripts run only if aube's built-in list or `aube-workspace.yaml` allows them; `strictDepBuilds` fails the install on an unreviewed one (`aube ignored-builds`, `aube approve-builds`). Add or bump a package with `aube add -D <pkg>` and commit `package.json` + `aube-lock.yaml`.
- Not available through mise, so documented instead: tippecanoe (optional, `brew install tippecanoe`; without it airspace uses its Go tiler), gcloud (not needed: `mise run google:auth login` writes Terraform's credentials). The Claude CLI is your machine-wide tool (Anthropic's installer, `mise run claude-cli:ensure`), not a project pin. The project's own Go CLIs (`cmd/*`) are built by their `<binary>:build` tasks, into `bin/`.

### Code Structure (internal vs pkg)
- Keep project-only helpers in `internal/` (e.g., `internal/codecinstaller`), not API-stable.
- Put reusable surfaces in `pkg/` with small, clear interfaces; no `pkg` imports from `internal`.
- Share code across commands via `pkg/…`, not by duplicating under `cmd/`.

**Workflow Naming:** `{category}-{name}.yml`

| Category | Purpose | Examples |
|----------|---------|----------|
| `core-` | P0 - must pass for merge | `core-ci.yml` (`mise run root:ci`) |
| `monitor-` | Scheduled health checks | `monitor-analytics.yml`, `monitor-sitecheck.yml` |
| `syndication-` | Content distribution | `syndication-bluesky.yml` |
| `release-` | Build & release pipelines | `release-tools.yml` |

**Development Workflow:**

| Command | Purpose |
|---------|---------|
| `mise run dev:up` | Start Hugo + tools with TUI dashboard |
| `mise run dev:down` | Stop all processes |

### Branding Assets

**IMPORTANT:** Logo SVGs are generated from Go code, NOT edited directly!

- Source of truth: `cmd/genlogo/main.go`
- Regenerate: `mise run genlogo:all` (one asset: `genlogo:logo-svg`, `genlogo:favicon`, `genlogo:og`, …)
- Generated files: `assets/images/logo.svg`, `logo-darkmode.svg`, `favicon.png`, `og-image.png`
- DO NOT edit SVG files directly - changes will be overwritten

After regenerating, manually update: Bluesky, Gmail signature

### Blog Images

Location: `assets/images/blog/`
Format: SVG, 800x400 viewBox
Colors: `#58a6ff` (blue), `#121212` (dark), `#f8f9fa` (background)

### Bluesky Syndication

Blog posts auto-post to Bluesky via `.github/workflows/syndication-bluesky.yml`.

- Runs every 6 hours (or manual trigger)
- RSS feed: `https://www.ubuntusoftware.net/blog/index.xml`
- Account: `ubuntusoftware.net`
- Secret: `BLUESKY_APP_PASSWORD` (GitHub repo secret)

Social preview cards use `og-image.png` (site default) since SVGs don't work for social media previews.

### Contact Form

Uses Web3Forms (free, unlimited). Submissions go to `gerard.webb@ubuntusoftware.net`.

- Dashboard: https://web3forms.com
- Access key in `layouts/contact/list.html`
- Config: `config/_default/params.toml` → `contact_form_action`

**Custom Subject Lines:** Edit `data/contact_subjects.yaml` to add email subject mappings.

```yaml
security: "Security Vulnerability Report"
partnership: "Partnership Inquiry"
```

Link with `?subject=<key>` (e.g., `/contact/?subject=security`) to auto-fill.

### MailerLite Integration

CLI tool for subscriber management and email automation.

**CLI:** `cmd/mailerlite/main.go`
**Tasks:** `.mise/tasks/mailerlite.toml` (subcommands passed through)

**Key Commands:**
```bash
mise run mailerlite:subscribers list         # List subscribers
mise run mailerlite:subscribers add x@y.com  # Add subscriber
mise run mailerlite:groups list              # List groups
mise run mailerlite:stats                    # Account stats
mise run mailerlite:server                   # Start webhook server
mise run mailerlite:releases latest          # Show latest GitHub release
```

**Web3Forms → MailerLite Flow:**

```
User submits "Get Started" form
        ↓
Web3Forms sends webhook POST
        ↓
mailerlite server receives at /webhook
        ↓
Parses form data (name, email, company, platform, industry)
        ↓
Adds subscriber to MailerLite via API
        ↓
(Optional) Auto-assigns to group
```

**Running the webhook server:**
```bash
# Development (local)
mise run mailerlite:server PORT=8086

# With auto-group assignment
mise run mailerlite:server GROUP_ID=12345

# Production: expose via tunnel
ngrok http 8086
# Then configure Web3Forms webhook URL to the ngrok URL
```

**GitHub Releases Integration:**
```bash
mise run mailerlite:releases latest     # Show latest release info
mise run mailerlite:releases urls       # Get URLs for email templates
```

### Playwright Tool

Reusable browser automation for OAuth flows and testing.

**CLI:** `cmd/playwright/main.go`
**Tasks:** `.mise/tasks/playwright.toml` (arguments passed through)

**Commands:**
```bash
mise run playwright:install                # Install browsers
mise run playwright:oauth https://...      # OAuth flow with callback
mise run playwright:screenshot <url> <file> # Take screenshot
mise run playwright:open https://...       # Open URL in browser
```

**OAuth Output (JSON):**
```json
{
  "code": "authorization_code_here",
  "token": "access_token_if_present",
  "query": { "all": "query", "params": "here" }
}
```

Used by `cmd/google` (`google auth login`) for automated Google OAuth flows.



### Google MCP Setup (Prerequisites for Gmail/Calendar)

**Architecture:**
```
┌─────────────────────────────────────────────────────────────────┐
│                        PROJECT FILES                             │
├─────────────────────────────────────────────────────────────────┤
│ mise.toml               - Pins google-mcp-server, fnox, jq      │
│ .mise/tasks/google-mcp.toml - google-mcp:* setup tasks, each    │
│                           running cmd/google-mcp/<command>.go   │
│ fnox.toml               - OAuth client id/secret → keychain     │
│ .mcp.json               - MCP server config (checked in):       │
│                           mise exec -- fnox exec -- google-mcp- │
│                           server                                │
│ cmd/google/             - Unified Google CLI (gmail, calendar,  │
│                           drive, ... + auth login/open)         │
│ internal/googleauth/    - Shared token loading                  │
│ internal/browser/       - Shared browser automation             │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    USER-LEVEL FILES (outside project)           │
├─────────────────────────────────────────────────────────────────┤
│ ~/.google-mcp-accounts/ - OAuth tokens (MUST be outside git!)   │
│ macOS keychain          - GOOGLE_CLIENT_ID / _SECRET (via fnox) │
└─────────────────────────────────────────────────────────────────┘
```

The draft blog post (`content/*/blog/google-workspace-mcp-zero-config.md`) uses these `mise run google-mcp:*` names.
`google:*` is a different thing: the unified Google CLI, `cmd/google` (`.mise/tasks/google.toml`: `google:gmail`, `google:calendar`, `google:auth`, …).

**Why tokens are at `~/.google-mcp-accounts/`:**
- Security: Tokens contain access to Google account - NEVER commit to git
- Sharing: Same tokens work across multiple projects
- Standard: This is where `google-mcp-server` expects them

**One-Command Setup:**
```bash
mise run google-mcp:setup    # Full guided setup (interactive)
```

**Manual Setup Steps:**
```bash
# 1. Install the pinned MCP server, fnox, jq
mise install

# 2. Google Cloud OAuth screens + exact values to enter (opens browser)
mise run google-mcp:console

# 3. Store client id + secret in the keychain (hidden prompts)
mise run google-mcp:credentials

# 4. Authenticate (opens browser for Google sign-in), ends with a check
mise run google-mcp:auth
```

No step 5: `.mcp.json` is checked in. Restart Claude Code to pick up the server.

**Status Check:**
```bash
mise run google-mcp:status   # Binary, keychain secrets, accounts, MCP config
mise run google-mcp:check    # Token refresh at Google + a real Calendar call
mise run google-mcp:rotate   # New client secret (copy it in the console)
```

### Email Sending (SMTP2GO)

**Primary email relay service:** SMTP2GO (smtp2go.com)

Configured inside Gmail settings to send outbound email through SMTP2GO relay. This provides better deliverability than direct Gmail sending.

**Setup:** Gmail → Settings → Accounts → "Send mail as" → configured with SMTP2GO credentials
**From address:** `gerard.webb@ubuntusoftware.net`
**Dashboard:** https://app.smtp2go.com
**Support:** ticket@smtp2go.com

**⚠️ Setup was painful** - getting SMTP2GO working with Gmail's "Send mail as" required careful configuration on both sides. Don't casually reconfigure unless necessary.

**Why SMTP2GO:**
- Better deliverability (dedicated IP reputation)
- Free tier: 1,000 emails/month, 200/day (no expiry, no credit card)
- Delivery tracking and analytics
- No browser automation needed

**Future consideration:** Could add direct SMTP2GO API integration (`/email/send` endpoint) for programmatic sending without Gmail.

### Gmail CLI Tool

Unified email sending via the Gmail API, or via Gmail's compose window in your browser. **From address is always `gerard.webb@ubuntusoftware.net`** - hardcoded to prevent mistakes.

**CLI:** `cmd/google` (`google gmail …`), **Package:** `internal/google/gmail.go`, **Tasks:** `mise run google:gmail <subcommand>` (`.mise/tasks/google.toml`)

**Modes:**
| Mode | Command | Use Case |
|------|---------|----------|
| API | `gmail send --mode=api` | Headless, most reliable (default) |
| Browser | `gmail send --mode=browser` | Fallback when API unavailable: opens compose in your default browser, you click Send |
| Compose | `gmail compose` | Opens Gmail for user review before send |
| Server | `gmail server` | HTTP webhook endpoint for external triggers |

**Quick actions:**
```bash
mise run google:gmail check     # Verify API connection
mise run google:gmail list      # Recent messages
mise run google:gmail server    # Start webhook server (port 8087)
```

**Send Commands:**
```bash
# Send via API (recommended)
mise run google:gmail send --to user@example.com --subject "Hello" --body "Message"

# Open compose in your browser (fallback; you click Send)
mise run google:gmail send --mode=browser --to user@example.com --subject "Hello" --body "Message"

# Open compose for review (user clicks send)
mise run google:gmail compose --to user@example.com --subject "Review" --body "Please check"
```

**Server Mode (HTTP API):**
```bash
# Start server
mise run google:gmail server --port 8087

# Send via HTTP
curl -X POST http://localhost:8087/send \
  -H "Content-Type: application/json" \
  -d '{"to":"x@y.com","subject":"Test","body":"Hello"}'
```

**Process Compose:** Add `gmail-server` to dev workflow by enabling in `process-compose.yaml`.

**⚠️ NEVER use Playwright MCP directly for email!** Always use `mise run google:gmail …` - it guarantees the correct From address.

**Signature:** "Ubuntu Software Local AI" (not "Claude")

### Calendar CLI Tool

Unified Google Calendar management via API or browser automation.

**CLI:** `cmd/google` (`google calendar …`), **Package:** `internal/google/calendar.go`, **Tasks:** `mise run google:calendar <subcommand>` (`.mise/tasks/google.toml`)

**Modes:**
| Mode | Command | Use Case |
|------|---------|----------|
| API | `calendar create --mode=api` | Headless, most reliable (default) |
| Browser | `calendar create --mode=browser` | Fallback when API unavailable |
| Compose | `calendar compose` | Opens calendar for user review before save |
| Server | `calendar server` | HTTP webhook endpoint for external triggers |

**Quick actions:**
```bash
mise run google:calendar check     # Verify API connection
mise run google:calendar today     # List today's events
mise run google:calendar open      # Open Google Calendar in browser
mise run google:calendar server    # Start webhook server (port 8088)
```

**Create Commands:**
```bash
# Create via API (recommended)
mise run google:calendar create --title "Meeting" --start "tomorrow 2pm" --end "tomorrow 3pm"

# Create via browser automation (fallback)
mise run google:calendar create --mode=browser --title "Meeting" --start "tomorrow 2pm" --end "tomorrow 3pm"

# Open calendar for review (user saves)
mise run google:calendar compose --title "Review" --start "tomorrow 10am" --end "tomorrow 11am"
```

**List Commands:**
```bash
mise run google:calendar list --max 10   # List upcoming events
mise run google:calendar today          # List today's events
```

**Time Formats:**
- RFC3339: `2024-12-13T14:00:00+07:00`
- Relative: `"today 2pm"`, `"tomorrow 10am"`, `"+1h"`, `"+30m"`

**Server Mode (HTTP API):**
```bash
# Start server
mise run google:calendar server --port 8088

# Create event via HTTP
curl -X POST http://localhost:8088/create \
  -H "Content-Type: application/json" \
  -d '{"title":"Meeting","start":"2024-12-15T14:00:00+07:00","end":"2024-12-15T15:00:00+07:00"}'

# List today's events
curl http://localhost:8088/today
```

**Process Compose:** Add `calendar-server` to dev workflow by enabling in `process-compose.yaml`.

### Airspace Data Pipeline (BVLOS)

Interactive US airspace visualization for drone fleet operations. Displays FAA controlled airspace, special use airspace (MOAs, restricted areas), and LAANC ceiling altitudes.

**Demo:** `/airspace-demo/` (static HTML + Leaflet.js + PMTiles)
**Content Page:** `content/english/fleet/airspace-demo.md`
**CLI:** `cmd/airspace/main.go`

**Architecture:**
```
FAA ArcGIS APIs → Go sync → GeoJSON → tippecanoe → PMTiles → R2 → Frontend
                     ↓
              ETag tracking (idempotent)
                     ↓
              sync-history.json (cron tuning)
```

**Data Flow:**
1. `airspace sync` checks ETags against FAA sources
2. If changed: downloads GeoJSON (378K features for UAS = 422MB)
3. `airspace:tile` converts to PMTiles (123MB for UAS)
4. `r2:airspace:upload` pushes to Cloudflare R2
5. Frontend loads PMTiles progressively (manifest-driven)

**Key Commands:**
```bash
# Sync pipeline (idempotent - skips if no FAA changes)
mise run airspace:sync              # Check for changes, download if needed
mise run airspace:pipeline          # Sync + tile + manifest (no upload)
mise run airspace:upload            # Then R2

# Data inspection
mise run airspace:status            # Show local data files
mise run airspace:history           # Show sync history (for cron tuning)
go run ./cmd/airspace history   # Same, direct

# Manual operations
mise run airspace:sync --force      # Force download all
mise run airspace:tile              # Convert GeoJSON to PMTiles
task r2:airspace:upload         # Upload PMTiles to R2
```

**Files:**
| File | Purpose |
|------|---------|
| `data/airspace/etags.json` | ETag cache for change detection |
| `data/airspace/sync-result.json` | Last sync result (for pipeline idempotency) |
| `data/airspace/sync-history.json` | Rolling history of sync runs |
| `data/airspace/manifest.json` | Global manifest (regions) |
| `data/airspace/usa-manifest.json` | USA regional manifest (layers) |
| `static/airspace/*.geojson` | Downloaded FAA data (gitignored) |
| `static/airspace/tiles/*.pmtiles` | Generated tiles (gitignored) |

**Sync Timing:**
- No changes: ~7-10 seconds (HEAD requests only)
- UAS changed: ~5-10 minutes (422MB download + tile gen)
- All changed: ~10-15 minutes

**GitHub Actions:** `.github/workflows/sync-airspace.yml`
- Runs weekly (Thursday 06:00 UTC, aligned with AIRAC cycle)
- Caches ETags between runs
- Only tiles/uploads if changes detected

**Why PMTiles?**
- Single-file format for vector tiles (no tile server needed)
- Progressive loading (browser fetches only visible tiles)
- 422MB GeoJSON → 123MB PMTiles (compression + tiling)

**Limitations:**
- FAA only supports Query API (no sync/change tracking)
- Must download ALL features even for 1 change
- PMTiles regeneration is full rebuild (no incremental)

**R2 Bucket:** `ubuntu-website-assets`
**R2 Public URL:** `https://pub-97cfaeb734ae474c80c79c3e3cc6dbee.r2.dev/airspace/`

**Local Development:**
Hugo's dev server doesn't properly handle HTTP range requests, which PMTiles requires.
Use `pmtiles-server` alongside Hugo for local tile serving:

```bash
mise run pmtiles-server:serve    # Start PMTiles server on :8765
```

Or enable in process-compose (in `process-compose.yaml`, set `disabled: false` for `pmtiles-server`).

The airspace demo page auto-detects the local server and falls back to R2 if unavailable.

### Cloudflare R2 (Large Assets)

R2 object storage for static assets over the Worker's 25 MiB a file (the airspace tiles and data). CORS (`cloudflare/r2-cors.json`, `mise run r2:cors-set`) allows production, the Worker's workers.dev URL and every branch preview's (`https://*-ubuntusoftware-net.gedw99.workers.dev`), and localhost.

**Tasks:** `.mise/tasks/r2.toml` (wrangler's, as tool tasks)

**Key Commands:**
```bash
mise run r2:status              # The account's buckets
mise run r2:list                # List bucket contents
mise run r2:upload <key> <file> # Upload a file (then at the public URL/<key>)
mise run r2:cors-set            # CORS rules from r2-cors.json (browsers need them)
mise run url:cf-r2              # Open the R2 dashboard
```

**Configuration:**
- Bucket: `ubuntu-website-assets`
- Public URL: `https://pub-97cfaeb734ae474c80c79c3e3cc6dbee.r2.dev`
- Requires `CLOUDFLARE_API_TOKEN` with R2 permissions

### Page Images (banner, services, etc.)

Location: `assets/images/`
Format: SVG with explicit width/height attributes
Dimensions: banner 800x500, services 560x520, call-to-action 400x400
Style: Hugo Plate grayscale line-art (white, `#f5f5f5`, `#ccc`, `#999`, `#666`)

### Translation Workflow

Languages: de (German), zh (Chinese), ja (Japanese), vi (Vietnamese) - auto-loaded from `config/_default/languages.toml`.

**Architecture - Separation of Concerns:**
- `internal/translate/hugo.go` - ALL Hugo-specific code (language parsing, menu parsing)
- `internal/translate/checker.go` - Pure query functions (CheckStatus, CheckMissing, etc.)
- `internal/translate/mutator.go` - Side-effect functions (DoClean, DoDone, etc.)
- `internal/translate/presenter.go` - Terminal and Markdown output formatting
- `.mise/tasks/translate.toml` ↔ `cmd/translate/` - the tasks (`translate:content`, `translate:menu`, `translate:lang`, `translate:check`), subcommands passed through

**Tasks** (`mise run translate:<task> <subcommand>`):

| Task | Subcommands | Purpose |
|-----------|----------|---------|
| `translate:content` | status, diff, changed, next, done | Track English source changes |
| `translate:content` | missing, orphans, stale, clean | Find translation problems |
| `translate:menu` | check, sync | Manage navigation menus |
| `translate:lang` | list, add, remove, init, validate | Manage languages |

**Common Commands:**
- `mise run translate:content status` - what English files changed since last translation?
- `mise run translate:content missing` - what's missing in target languages?
- `mise run translate:content next` - which file should I translate next?
- `mise run translate:content done` - mark translations complete (moves and pushes the `last-translation` tag)
- `mise run translate:menu check` - validate menus for broken links
- `mise run translate:lang list` - show configured languages

**Lifecycle:**
1. Edit English file → `translate:content status` shows it
2. Run `translate:content next` → get next file to translate
3. Translate to all languages, commit translations
4. Run `translate:content done` → moves checkpoint (pushed), status is clean

**CI:** `monitor-translate.yml` runs weekly, creates GitHub Issue if missing translations.

### Auto-Translation (DeepL)

Automatic translation of Hugo markdown using DeepL API. Preserves front matter, shortcodes, code blocks.

**Setup:**
1. Get free API key at https://www.deepl.com/pro-api (500k chars/month free)
2. Add to `.env`: `DEEPL_API_KEY=your-key-here`

**CLI:** `cmd/autotranslate/main.go`
**Package:** `internal/autotranslate/`
**Tasks:** `.mise/tasks/autotranslate.toml` (arguments and flags passed through)

**Commands:**
```bash
mise run autotranslate:status                       # Check API config
mise run autotranslate:languages                    # List supported languages
mise run autotranslate:missing vi                   # Translate all missing Vietnamese
mise run autotranslate:missing vi --dry-run         # Preview what would translate
mise run autotranslate:missing vi --provider claude-cli   # With the Claude CLI
mise run autotranslate:file <path> vi               # Translate single file
```

**What it preserves (won't translate):**
- YAML front matter (except title, description, meta_title)
- Hugo shortcodes `{{< >}}` and `{{% %}}`
- Code blocks (fenced and inline)
- URLs in markdown links
- HTML tags
- Image references

**Provider Interface:** Designed to support multiple providers. Currently DeepL, can add Google/OpenAI later via `internal/autotranslate/provider.go`.

### Path Convention
**ALWAYS use `joeblew999` (with three 9s), NEVER `joeblew99` (with two 9s)**

Correct paths:
- `/Users/apple/workspace/go/src/github.com/joeblew999/ubuntu-website`
- `github.com/joeblew999/ubuntu-website`

### Hugo Version

**Current:** the `hugo-extended` pin in `mise.toml`, the only place: local, CI and Cloudflare Pages (`root:pages` installs it there) all use it. No Cloudflare `HUGO_VERSION`.
**Hugo Plate requires:** v0.165+ extended (its `theme.toml` `min_version`)

**Upgrade:** `mise run upgrade:hugo` (pins, proves a production build) → commit on a branch, push, check the Cloudflare preview build (it builds with the new pin) → merge to main.

Hugo >= 0.165 removed `tailwindcss` from the default `security.exec.allow`; `hugo.toml` re-adds it (`[security.exec]`, per gohugoio/hugo#15178, reported upstream as zeon-studio/hugoplate#250). Hugo then requires `tailwindcss` to be the npm package (`mise run root:deps`), not a standalone binary on PATH.

### Updating Hugo Plate (repeatable)

Hugo Plate doesn't tag releases ([issue #188](https://github.com/zeon-studio/hugoplate/issues/188)). It and the gethugothemes modules are Hugo Modules, vendored in `_vendor/`; the exact commits are in `_vendor/modules.txt` (`mise run theme:status`). Not in `go.mod`: that file belongs to the Go code and `go mod tidy` strips Hugo's entries.

```bash
mise run theme:update                 # latest upstream (re-run is a no-op when nothing changed)
mise run theme:update --ref <sha>     # pin Hugo Plate to a commit, or roll back
```

`theme update` (`cmd/theme/update.go`) does: `hugo mod get -u` + `hugo mod vendor` (go.mod restored) → syncs `cmd/theme/themeGenerator.js` from that commit (unmodified) → bumps/adds npm packages to match upstream's `package.json` → regenerates `assets/css/generated-theme.css` → 3-way merges upstream's changes into the templates our `layouts/` override (conflicts to `.theme-merge/`) → production build. Then review `git diff -- _vendor package.json cmd/theme layouts`, push a branch, check the Cloudflare preview build, merge to main.

**Don't use upstream's `themeUpdate.js`** (`curl …/tarball/main | tar`): unpinned, no checksums, and it copies templates without the npm packages/build steps they need (that's how `tailwind-bootstrap-grid` broke the build).

**Theme CSS:** upstream generates colors/fonts from `data/theme.json` with its `themeGenerator.js` (kept unmodified in `cmd/theme/`). Its paths assume `exampleSite/` or `themes/<name>/`, so `theme css` runs it in a temp dir shaped like that and writes `assets/css/generated-theme.css`, which is committed: builds use it as it is, and `theme:check` (in `root:ci`) fails if it isn't what `data/theme.json` makes. After editing `data/theme.json`: `mise run theme:css` and commit the CSS.

### Hugo Source code

When working on Hugo coding aspects use the source of the pinned version (`hugo-extended` in `mise.toml`):

- **Source code**: https://github.com/gohugoio/hugo (check out the tag matching the pin, e.g. `v0.169.0`)
- **Release notes** (breaking changes, deprecations): `gh release view v<version> -R gohugoio/hugo`

### Hugo Plate Source code

When working on Hugo Plate, which is a Hugo theme:

- **Source code (exact pinned commit)**: `_vendor/github.com/zeon-studio/hugoplate/` (commit in `_vendor/modules.txt`; full repo incl. `scripts/` in `$(go env GOMODCACHE)/github.com/zeon-studio/hugoplate@<version>/`)
- **Upstream**: https://github.com/zeon-studio/hugoplate

- Only use Hugo template things properly. I do not want to steer away from the standard Hugo Plate way of doing things !!

### URL Hygiene

**⚠️ CRITICAL: When renaming or moving ANY content file, ALWAYS add an alias for the old URL!**

This applies to:
- Renaming folders (e.g., `early-access/` → `get-started/`)
- Renaming files (e.g., `team/gerard-webb.md` → `founder/gerard-webb.md`)
- Moving pages between sections

```yaml
---
title: "Get Started"
aliases:
  - "/early-access/"
---
```

Hugo generates redirect HTML at old URLs with canonical tags (SEO-friendly).
Keep aliases for 6+ months to preserve bookmarks and search rankings.

**Internal Links: Use `relref` for build-time validation**

```markdown
# Don't do this (no validation, can break silently):
[Platform](/platform/)

# Do this (Hugo fails build if page doesn't exist):
[Platform]({{< relref "/platform" >}})
```

**Config**: `refLinksErrorLevel = "error"` in hugo.toml ensures broken relrefs fail the build.

### SEO & Structured Data

**SEO is automatic** - every page gets JSON-LD schema, Open Graph, and Twitter cards via `layouts/_partials/basic-seo.html`.

Schema types applied:
- **Homepage**: Organization (company info, founder, social links)
- **Blog posts**: Article (headline, author, dates, publisher)
- **Other pages**: WebPage (name, description, publisher)
- **All non-home pages**: BreadcrumbList (navigation hierarchy)

**When creating/editing content**, ensure these front matter fields are set:

```yaml
---
title: "Page Title"           # Required - used in breadcrumbs
meta_title: "SEO Title | Ubuntu Software"  # Optional - overrides title in <title> tag
description: "150-160 char description for search results and social sharing"
image: "images/page-image.png"  # Optional - for blog posts, used in Article schema
author: "Gerard Webb"           # For blog posts
---
```

**Validation tools:**
- Rich Results Test: https://search.google.com/test/rich-results
- Schema Validator: https://validator.schema.org/

**robots.txt**: Custom file at `static/robots.txt` - explicitly allows all AI crawlers.

### Theme Upgrade Policy

**DO NOT modify theme CSS, Tailwind config, or `data/theme.json`** - keep upgrades easy.

Allowed changes:
- Content files (`content/`)
- Custom layouts (`layouts/`)
- Config (`config/`)
- Images (`assets/images/`, `static/images/`)

Avoid:
- CSS/SCSS changes
- Tailwind plugin modifications
- Theme color overrides
- Any changes that would conflict with upstream Hugo Plate updates 

## Founder Social Presence

**NO LinkedIn** - Gerard is not on LinkedIn. Don't add LinkedIn links.

Social links to use:
- GitHub: https://github.com/joeblew999
- Bluesky: https://bsky.app/profile/ubuntusoftware.net

## Information about the company

This Project is a web site for my company. I have various info here about me and the company.

- **Source code**: `/Users/apple/Library/Mobile Documents/com~apple~CloudDocs/Thailand /`

### Company Registration Details (for App Store / Developer Programs)

**Ubuntu Software (Australia)**

| Field | Value |
|-------|-------|
| Legal Name | Ubuntu Software |
| ABN | 95 595 575 880 |
| D-U-N-S Number | 891770992 |
| Domain | ubuntusoftware.net |
| Support Email | support@ubuntusoftware.net |
| Developer Contact | gerard.webb@ubuntusoftware.net |
| Principal | Gerard Joseph Webb |

**Verification Links:**
- ABN Lookup: https://abr.business.gov.au/ABN/View?id=95595575880
- D-U-N-S Lookup (illion): https://express.illion.com.au/company (search by ABN)
- Company History: https://www.ubuntusoftware.net/pages/company-history/

**D-U-N-S Registered Program:** Applied at https://www.dunsregistered.com - awaiting confirmation email to `gerard.webb@ubuntusoftware.net`

**Bundle ID Pattern:** `net.ubuntusoftware.<appname>` (e.g., `net.ubuntusoftware.cad`)

See [notes/deeplink.md](notes/deeplink.md) for the App Store deployment plan. Product and business notes are in `notes/`.

