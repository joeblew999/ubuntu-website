# Claude Assistant Notes

## CRITICAL REMINDERS

**Production Domain:** `www.ubuntusoftware.net`

Apex domain (`ubuntusoftware.net`) redirects to www via Cloudflare redirect rule.

USE TASKFILE - it makes conventions for development.

### Taskfile Conventions

Key points:
- Use `status:` for idempotent `check:deps` tasks
- Use `deps:` for declarative dependencies

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

### Page Images (banner, services, etc.)

Location: `assets/images/`
Format: SVG with explicit width/height attributes
Dimensions: banner 800x500, services 560x520, call-to-action 400x400
Style: Hugo Plate grayscale line-art (white, `#f5f5f5`, `#ccc`, `#999`, `#666`)

### Translation Workflow

Languages: de (German), zh (Chinese), ja (Japanese) - auto-loaded from `config/_default/languages.toml`.

**Architecture - Separation of Concerns:**
- `internal/translator/hugo.go` - ALL Hugo-specific code (language parsing, menu parsing)
- `internal/translator/checker.go` - Pure query functions (CheckStatus, CheckMissing, etc.)
- `internal/translator/mutator.go` - Side-effect functions (DoClean, DoDone, etc.)
- `internal/translator/presenter.go` - Terminal and Markdown output formatting
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


