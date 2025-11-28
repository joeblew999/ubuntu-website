# Claude Assistant Notes

## CRITICAL REMINDERS

**Production Domain:** `www.ubuntusoftware.net`


USE TASKFILE - it makes conventions for development...

### Branding Assets

**IMPORTANT:** Logo SVGs are generated from Go code, NOT edited directly!

- Source of truth: `cmd/genlogo/main.go`
- Regenerate: `task generate:assets`
- Generated files: `assets/images/logo.svg`, `logo-darkmode.svg`, `favicon.png`, `og-image.png`
- DO NOT edit SVG files directly - changes will be overwritten

After regenerating, manually update: Bluesky, Gmail signature

### Blog Images

Location: `assets/images/blog/`
Format: SVG, 800x400 viewBox
Colors: `#58a6ff` (blue), `#121212` (dark), `#f8f9fa` (background)

### Bluesky Syndication

Blog posts auto-post to Bluesky via `.github/workflows/bluesky-syndication.yml`.

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

## Information about the company 

This Project is a web site for my company. I have various info here about me and the company.

- **Source code**: `/Users/apple/Library/Mobile Documents/com~apple~CloudDocs/Thailand /`


