# Ubuntu Software Website

Multi-language website built with Hugo Plate, featuring automated translation and deployment to Cloudflare Pages.

http://localhost:1313

http://localhost:1313/fleet/airspace-demo/

https://www.ubuntusoftware.net

https://www.ubuntusoftware.net/sitemap.xml
https://www.ubuntusoftware.net/en/sitemap.xml

https://www.ubuntusoftware.net/fleet/airspace-demo/


## 📍 Touch Points (Edit Locations)

- **Gmail Signature**: [Gmail Settings → Signature](https://mail.google.com/mail/u/0/#settings/general)
- **Bluesky Profile**: [bsky.app/profile/ubuntusoftware.net](https://bsky.app/profile/ubuntusoftware.net) → Edit Profile
- **Cloudflare DNS**: [Cloudflare Dashboard](https://dash.cloudflare.com/) → ubuntusoftware.net → DNS
- **GitHub Profile**: [github.com/settings/profile](https://github.com/settings/profile)

## 🔗 Quick Links

For all URLs, dashboards, and project information, run:

```bash
task          # or: task info
```

## 🌍 Supported Languages

- 🇬🇧 English (en) - Source language
- 🇩🇪 German (de)
- 🇨🇳 Chinese Simplified (zh)
- 🇯🇵 Japanese (ja)
- 🇻🇳 Vietnamese (vi)


## 🚀 Quick Start

```bash
task setup      # mise install (pinned tools) + npm packages (aube) + Go modules
task dev:up     # Hugo + tools via process-compose (task dev:down to stop)
```

Pushing to `main` deploys to Cloudflare Pages; other branches get preview deployments.

## 📋 Prerequisites

- [mise](https://mise.jdx.dev/getting-started.html): installs every tool pinned in `mise.toml` (Task, Hugo, Bun, Node, Wrangler, …). Go's version comes from `go.mod`.

## 🛠️ Development Tasks

```bash
task dev           # Hugo dev server
task build         # Production build
task preview       # Preview production build
task clean         # Clean artifacts
task upgrade       # Upgrade tools, Go modules and the theme (Hugo: task hugo:upgrade)

task cf:status     # Cloudflare Pages deployment status
```

## 🌐 Translation Workflow

1. Edit English content in `content/english/`
2. `task translate:content:status`: what changed since the last translation
3. `task translate:content:next`: next file to translate
4. Translate it into each language (see CLAUDE.md, Translation Workflow)
5. `task translate:content:done`: move the checkpoint

## 📦 Project Structure

```
.
├── cmd/                    # Go CLIs (translate, autotranslate, google, mailerlite, …)
├── internal/               # Shared Go packages
├── content/
│   ├── english/            # Source (EN)
│   ├── german/             # DE
│   ├── chinese/            # ZH
│   ├── japanese/           # JA
│   └── vietnamese/         # VI
├── config/_default/
│   ├── languages.toml      # Language config
│   └── menus.*.toml        # Per-language menus
├── i18n/                   # Translation strings
├── _vendor/                # Hugo Plate + modules, pinned (task theme:update)
├── .github/workflows/      # CI (tools installed by mise)
├── mise.toml               # Pinned tools
├── Taskfile.yml            # Task definitions
└── wrangler.toml           # Cloudflare Pages config
```

## 🔧 Configuration

Secrets are never committed:

- `.env` (gitignored): e.g. `CLOUDFLARE_API_TOKEN` (read by wrangler and the `cf:*` tasks), `DEEPL_API_KEY` (autotranslate)
- Google MCP OAuth client: macOS keychain via fnox (`fnox.toml`; `task google-mcp:credentials`)

## 🔗 References

- [Hugo Plate Template](https://github.com/zeon-studio/hugoplate)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [Cloudflare Pages Docs](https://developers.cloudflare.com/pages/)
- [Task Documentation](https://taskfile.dev/)

## 📄 License

MIT (Hugo Plate template)





