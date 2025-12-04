# Claude Assistant Notes

## CRITICAL REMINDERS

**Production Domain:** `www.ubuntusoftware.net`

Apex domain (`ubuntusoftware.net`) redirects to www via Cloudflare redirect rule.

USE TASKFILE - it makes conventions for development.

### Taskfile Conventions

**Environment Variables:**
- Taskfile loads `.env` automatically via `dotenv: ['.env']`
- All secrets/config should be in `.env` (gitignored)
- `.env.test` provides template with placeholder values

**Naming Convention (everything aligned):**

| Component | Pattern | Example |
|-----------|---------|---------|
| cmd/ tool | `cmd/<name>` | `cmd/analytics` |
| Taskfile tasks | `<name>:*` | `analytics:report` |
| CI task | `ci:<name>` | `ci:analytics` |

**Workflow Naming:** `{category}-{name}.yml`

| Category | Purpose | Examples |
|----------|---------|----------|
| `core-` | P0 - must pass for merge | `core-taskfile.yml`, `core-xplat.yml`, `core-tools.yml` |
| `monitor-` | Scheduled health checks | `monitor-analytics.yml`, `monitor-sitecheck.yml` |
| `syndication-` | Content distribution | `syndication-bluesky.yml` |
| `release-` | Build & release pipelines | `release-xplat.yml` |

**Task Suffixes:**
- `namespace:action` - Leaf task (`sitecheck:dns`)
- `namespace:all` - Calls other tasks in namespace (`sitecheck:all`)
- `namespace` (bare) - Default action (`sitecheck`)

**Task Dependencies** (see Taskfile header for full list)

**DRY Principle - GitHub Actions:**
- GitHub Actions call Taskfile tasks, not run commands directly
- This ensures local `task X` runs the same as CI
- Pattern: `run: task ci:<name>` in workflows

**CI Tasks:**
- `ci:*` namespace is the interface for GitHub Actions
- These tasks output markdown and use exit codes for workflow control
- Task descriptions include `[Category]` prefix matching the workflow

**Workflow → Task Mapping:**

| Workflow | Task | Purpose |
|----------|------|---------|
| `core-taskfile.yml` | `ci:taskfile` | Validate Taskfile across platforms |
| `monitor-analytics.yml` | `ci:analytics` | Weekly analytics check |
| `monitor-sitecheck.yml` | `ci:sitecheck` | Site reachability check |
| `monitor-health.yml` | `ci:health` | Health checks for external deps |
| `core-xplat.yml` | *(direct)* | xplat cross-platform build tests |
| `core-tools.yml` | *(direct)* | Binary tools build tests |
| `release-xplat.yml` | *(direct)* | xplat release automation |
| `syndication-bluesky.yml` | *(direct)* | Blog post syndication |

**Binary Pattern (using xplat binary:install):**

Binary tools use `xplat binary:install` for cross-platform installation. This command:
1. Checks if binary exists in PATH or install dir (skip if found)
2. Builds from local source if Go is available
3. Downloads from GitHub release as fallback

| Category | Tools | Installation |
|----------|-------|--------------|
| With releases | xplat, analytics, sitecheck, genlogo | `xplat binary:install` |
| Local-only | lanip, env, translate | `go run` (simple/broken) |

Version management:
- Versions defined in `versions.env` (single source of truth)
- Taskfiles reference `{{.TOOL_VERSION}}` from dotenv
- Release tag format: `<tool>-v<version>` (e.g., `analytics-v0.1.0`)
- Install location: `~/.local/bin/` (unix) or `~/bin/` (windows)

Example taskfile:
```yaml
vars:
  # ANALYTICS_VERSION comes from versions.env
  # XPLAT_BIN comes from root Taskfile.yml (handles .exe on Windows)
  ANALYTICS_REPO: joeblew999/ubuntu-website
  ANALYTICS_INSTALL_DIR: '{{if eq OS "windows"}}{{.HOME}}/bin{{else}}{{.HOME}}/.local/bin{{end}}'

tasks:
  check:deps:
    cmds:
      - '{{.XPLAT_BIN}} binary install analytics {{.ANALYTICS_VERSION}} {{.ANALYTICS_REPO}} --source {{.ROOT_DIR}}/cmd/analytics'

  report:
    deps: [check:deps]
    cmds:
      - '{{.ANALYTICS_INSTALL_DIR}}/analytics{{exeExt}}'
```

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

Tasks:
- `task translate:status` - what English files changed since last translation
- `task translate:missing` - which languages are missing content files
- `task translate:done` - mark translations complete (updates checkpoint)

Workflow:
1. `task translate:status` → translate changed files to all languages
2. `task translate:done` → update checkpoint

Languages: de (German), zh (Chinese), ja (Japanese)

Note: When adding/removing languages, also update Taskfile.yml

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


