# Claude Assistant Notes

## CRITICAL REMINDERS

**Production Domain:** `www.ubuntusoftware.net`

Apex domain (`ubuntusoftware.net`) redirects to www via Cloudflare redirect rule.

USE TASKFILE - it makes conventions for development.

### Taskfile Conventions

Key points:
- Use `status:` for idempotent `check:deps` tasks
- Use `deps:` for declarative dependencies

### Code Structure (internal vs pkg)
- Keep project-only helpers in `internal/` (e.g., `internal/codecinstaller`), not API-stable.
- Put reusable surfaces in `pkg/` with small, clear interfaces; no `pkg` imports from `internal`.
- Share code across commands via `pkg/…`, not by duplicating under `cmd/`.

**Workflow Naming:** `{category}-{name}.yml`

| Category | Purpose | Examples |
|----------|---------|----------|
| `core-` | P0 - must pass for merge | `core-taskfile.yml`, `core-tools.yml` |
| `monitor-` | Scheduled health checks | `monitor-analytics.yml`, `monitor-sitecheck.yml` |
| `syndication-` | Content distribution | `syndication-bluesky.yml` |
| `release-` | Build & release pipelines | `release-tools.yml` |

**Development Workflow:**

| Command | Purpose |
|---------|---------|
| `task dev:up` | Start Hugo + tools with TUI dashboard |
| `task dev:down` | Stop all processes |
| `task dev:ui` | Start Task-UI web dashboard |

### Branding Assets

**IMPORTANT:** Logo SVGs are generated from Go code, NOT edited directly!

- Source of truth: `cmd/genlogo/main.go`
- Regenerate: `task genlogo:all`
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
**Taskfile:** `taskfiles/Taskfile.mailerlite.yml`

**Key Commands:**
```bash
task mailerlite:subscribers:list          # List subscribers
task mailerlite:subscribers:add EMAIL=x   # Add subscriber
task mailerlite:groups:list               # List groups
task mailerlite:stats                     # Account stats
task mailerlite:server                    # Start webhook server
task mailerlite:releases:latest           # Show latest GitHub release
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
task mailerlite:server PORT=8086

# With auto-group assignment
task mailerlite:server GROUP_ID=12345

# Production: expose via tunnel
ngrok http 8086
# Then configure Web3Forms webhook URL to the ngrok URL
```

**GitHub Releases Integration:**
```bash
task mailerlite:releases:latest     # Show latest release info
task mailerlite:releases:urls       # Get URLs for email templates
```

### Playwright Tool

Reusable browser automation for OAuth flows and testing.

**CLI:** `cmd/playwright/main.go`
**Taskfile:** `taskfiles/tools/Taskfile.playwright.yml`

**Commands:**
```bash
task playwright:install                    # Install browsers
task playwright:oauth URL=https://...      # OAuth flow with callback
task playwright:screenshot URL=x FILE=y   # Take screenshot
task playwright:open URL=https://...       # Open URL in browser
```

**OAuth Output (JSON):**
```json
{
  "code": "authorization_code_here",
  "token": "access_token_if_present",
  "query": { "all": "query", "params": "here" }
}
```

Used by `cmd/google-auth` for automated Google OAuth flows.



### Google MCP Setup (Prerequisites for Gmail/Calendar)

**Architecture:**
```
┌─────────────────────────────────────────────────────────────────┐
│                        PROJECT FILES                             │
├─────────────────────────────────────────────────────────────────┤
│ cmd/google-auth/        - MCP config management (add/remove)    │
│ cmd/gmail/              - Gmail CLI (uses internal/gmail)       │
│ cmd/calendar/           - Calendar CLI (uses internal/calendar) │
│ internal/googleauth/    - Shared token loading                  │
│ internal/browser/       - Shared browser automation             │
│ .vscode/mcp.json        - MCP server config (gitignored: NO)    │
│ .env                    - OAuth credentials (gitignored: YES)   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    USER-LEVEL FILES (outside project)           │
├─────────────────────────────────────────────────────────────────┤
│ ~/.google-mcp-accounts/ - OAuth tokens (MUST be outside git!)   │
│ ~/go/bin/google-mcp-server - External MCP server binary         │
└─────────────────────────────────────────────────────────────────┘
```

**Why tokens are at `~/.google-mcp-accounts/`:**
- Security: Tokens contain access to Google account - NEVER commit to git
- Sharing: Same tokens work across multiple projects
- Standard: This is where `google-mcp-server` expects them

**One-Command Setup:**
```bash
task google-mcp:setup    # Full guided setup (interactive)
```

**Manual Setup Steps:**
```bash
# 1. Install external MCP server
task google-mcp:install

# 2. Create Google Cloud OAuth credentials (opens browser)
task google-mcp:guide

# 3. Save credentials
task google-mcp:credentials CLIENT_ID='xxx' CLIENT_SECRET='xxx'

# 4. Authenticate (opens browser for Google sign-in)
source .env && task google-mcp:auth

# 5. Add to Claude Code
task google-mcp:claude:add
```

**Status Check:**
```bash
task google-mcp:status   # Shows binary, accounts, config status
task google-mcp:check    # Shows next step if incomplete
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

Unified email sending via API or browser automation. **From address is always `gerard.webb@ubuntusoftware.net`** - hardcoded to prevent mistakes.

**CLI:** `cmd/gmail/main.go`
**Package:** `internal/gmail/`
**Taskfile:** `taskfiles/Taskfile.gmail.yml`

**Modes:**
| Mode | Command | Use Case |
|------|---------|----------|
| API | `gmail send --mode=api` | Headless, most reliable (default) |
| Browser | `gmail send --mode=browser` | Fallback when API unavailable |
| Compose | `gmail compose` | Opens Gmail for user review before send |
| Server | `gmail server` | HTTP webhook endpoint for external triggers |

**TaskUI Quick Actions (no variables needed):**
```bash
task gmail:check        # Verify API connection
task gmail:open         # Open Gmail inbox
task gmail:server       # Start webhook server (port 8087)
```

**Send Commands:**
```bash
# Send via API (recommended)
task gmail:send TO=user@example.com SUBJECT="Hello" BODY="Message"

# Send via browser automation (fallback)
task gmail:send:browser TO=user@example.com SUBJECT="Hello" BODY="Message"

# Open compose for review (user clicks send)
task gmail:compose TO=user@example.com SUBJECT="Review" BODY="Please check"
```

**Templates:**
```bash
# Pre-defined blog update email
task gmail:templates:blog-update TO=contact@example.com
```

**Server Mode (HTTP API):**
```bash
# Start server
task gmail:server PORT=8087

# Send via HTTP
curl -X POST http://localhost:8087/send \
  -H "Content-Type: application/json" \
  -d '{"to":"x@y.com","subject":"Test","body":"Hello"}'
```

**Process Compose:** Add `gmail-server` to dev workflow by enabling in `process-compose.yaml`.

**⚠️ NEVER use Playwright MCP directly for email!** Always use `task gmail:*` commands - they guarantee correct From address.

**Signature:** "Ubuntu Software Local AI" (not "Claude")

### Calendar CLI Tool

Unified Google Calendar management via API or browser automation.

**CLI:** `cmd/calendar/main.go`
**Package:** `internal/calendar/`
**Taskfile:** `taskfiles/Taskfile.calendar.yml`

**Modes:**
| Mode | Command | Use Case |
|------|---------|----------|
| API | `calendar create --mode=api` | Headless, most reliable (default) |
| Browser | `calendar create --mode=browser` | Fallback when API unavailable |
| Compose | `calendar compose` | Opens calendar for user review before save |
| Server | `calendar server` | HTTP webhook endpoint for external triggers |

**TaskUI Quick Actions (no variables needed):**
```bash
task calendar:check        # Verify API connection
task calendar:today        # List today's events
task calendar:open         # Open Google Calendar in browser
task calendar:server       # Start webhook server (port 8088)
```

**Create Commands:**
```bash
# Create via API (recommended)
task calendar:create TITLE="Meeting" START="tomorrow 2pm" END="tomorrow 3pm"

# Create via browser automation (fallback)
task calendar:create:browser TITLE="Meeting" START="tomorrow 2pm" END="tomorrow 3pm"

# Open calendar for review (user saves)
task calendar:compose TITLE="Review" START="tomorrow 10am" END="tomorrow 11am"
```

**List Commands:**
```bash
task calendar:list              # List upcoming events (default 10)
task calendar:list:week         # List this week's events
task calendar:today             # List today's events
```

**Time Formats:**
- RFC3339: `2024-12-13T14:00:00+07:00`
- Relative: `"today 2pm"`, `"tomorrow 10am"`, `"+1h"`, `"+30m"`

**Server Mode (HTTP API):**
```bash
# Start server
task calendar:server PORT=8088

# Create event via HTTP
curl -X POST http://localhost:8088/create \
  -H "Content-Type: application/json" \
  -d '{"title":"Meeting","start":"2024-12-15T14:00:00+07:00","end":"2024-12-15T15:00:00+07:00"}'

# List today's events
curl http://localhost:8088/today
```

**Process Compose:** Add `calendar-server` to dev workflow by enabling in `process-compose.yaml`.

### Airspace Demo (BVLOS)

Interactive US airspace visualization for drone fleet operations. Displays FAA controlled airspace, special use airspace (MOAs, restricted areas), and LAANC ceiling altitudes.

**Demo:** `/airspace-demo/` (static HTML + Leaflet.js)
**Content Page:** `content/english/fleet/airspace-demo.md`
**Data:** FAA UDDS GeoJSON files stored in Cloudflare R2

**Architecture:**
- Development: Loads from local `static/airspace/*.geojson` (if present)
- Production: Loads from R2 (`https://pub-97cfaeb734ae474c80c79c3e3cc6dbee.r2.dev/airspace/`)
- Auto-detects environment via `window.location.hostname`

**CLI:** `cmd/airspace/main.go`

**Taskfile Commands:**
```bash
task airspace:demo              # Start standalone server (port 9091)
task airspace:status            # Show data file status and age
task airspace:download          # Refresh all FAA data
task airspace:download:uas      # Download only UAS Facility Map (LAANC)
task airspace:download:boundary # Download only Airspace Boundary
task airspace:download:sua      # Download only Special Use Airspace
task r2:airspace:upload         # Sync data to R2
task r2:endpoints               # List all R2 asset URLs
task r2:endpoints:test          # Verify R2 endpoints accessible
```

**Data Files (44MB total, gitignored):**
- `faa_airspace_boundary.geojson` (14MB) - Class B/C/D/E
- `faa_special_use_airspace.geojson` (28MB) - MOAs, Restricted
- `faa_uas_facility_map.geojson` (2.2MB) - LAANC grid

**R2 Bucket:** `ubuntu-website-assets`
**R2 Dashboard:** `task r2:open:bucket`

### Cloudflare R2 (Large Assets)

R2 object storage for static assets that exceed Cloudflare Pages' 25MB file limit.

**Taskfile:** `taskfiles/Taskfile.cf-r2.yml`

**Key Commands:**
```bash
task r2:status          # Show bucket status and dashboard URLs
task r2:list            # List bucket contents
task r2:endpoints       # List all asset endpoints
task r2:endpoints:test  # Test endpoint accessibility
task r2:upload FILE=x   # Upload single file
task r2:sync DIR=x      # Sync directory to R2
task r2:open            # Open R2 dashboard
task r2:open:bucket     # Open bucket browser
task r2:open:settings   # Open bucket settings
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
- `taskfiles/Taskfile.translate.yml` - CLI interface (calls Go binary)

**Taskfile Commands (all namespaced):**

| Namespace | Commands | Purpose |
|-----------|----------|---------|
| `content:` | status, diff, changed, next, done | Track English source changes |
| `content:` | missing, orphans, stale, clean | Find translation problems |
| `menu:` | check, sync | Manage navigation menus |
| `lang:` | list, add, remove, init, validate | Manage languages |

**Common Commands:**
- `task translate:content:status` - what English files changed since last translation?
- `task translate:content:missing` - what's missing in target languages?
- `task translate:content:next` - which file should I translate next?
- `task translate:content:done` - mark translations complete
- `task translate:menu:check` - validate menus for broken links
- `task translate:lang:list` - show configured languages

**Lifecycle:**
1. Edit English file → `translate:content:status` shows it
2. Run `translate:content:next` → get next file to translate
3. Translate to all languages, commit translations
4. Run `translate:content:done` → moves checkpoint, status is clean

**CI:** `monitor-translate.yml` runs weekly, creates GitHub Issue if missing translations.

### Auto-Translation (DeepL)

Automatic translation of Hugo markdown using DeepL API. Preserves front matter, shortcodes, code blocks.

**Setup:**
1. Get free API key at https://www.deepl.com/pro-api (500k chars/month free)
2. Add to `.env`: `DEEPL_API_KEY=your-key-here`

**CLI:** `cmd/autotranslate/main.go`
**Package:** `internal/autotranslate/`
**Taskfile:** `taskfiles/Taskfile.autotranslate.yml`

**Commands:**
```bash
task autotranslate:status                    # Check API config
task autotranslate:languages                 # List supported languages
task autotranslate:missing LANG=vi           # Translate all missing Vietnamese
task autotranslate:missing LANG=vi DRY_RUN=true  # Preview what would translate
task autotranslate:file FILE=path LANG=vi    # Translate single file
task autotranslate:vi                        # Shortcut: all missing Vietnamese
task autotranslate:all DRY_RUN=true          # Preview all languages
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

### Via Framework

When working on Via coding aspects:

- **Source code**: `/Users/apple/workspace/go/src/github.com/joeblew999/wellknown/.src/via`

- use a closure variable pattern instead of trying to access the signal's value.
- use PicoCSS

### Hugo Version

**Current:** 0.152.2 (extended)
**Hugo Plate requires:** v0.144+ extended

Three places must stay aligned:
1. **Local**: `hugo version` (install via `brew install hugo`)
2. **Cloudflare**: `HUGO_VERSION` secret (check build logs via dashboard)
3. **This doc**: Update version above when changing

Hugo Plate doesn't tag releases (see [issue #188](https://github.com/zeon-studio/hugoplate/issues/188)), so we track by commit hash in `go.mod`.

**Updating modules:** `hugo mod get -u` updates all modules together. The gethugothemes modules share a monorepo, so they stay in sync automatically.

### Hugo Source code

When working on Hugo coding aspects use the source at:

- **Source code**: `/Users/apple/workspace/go/src/github.com/joeblew999/wellknown/.src/hugo`

### Hugo Plate Source code

When working on Hugo Plate, which is a Hugo theme:

- **Source code**: `/Users/apple/workspace/go/src/github.com/joeblew999/wellknown/.src/hugoplate`

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

See [DEEPLINK.md](DEEPLINK.md) for full App Store deployment plan.

