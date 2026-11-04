# Ubuntu Software Website

Multi-language website built with Hugo Plate, featuring automated translation, served by a Cloudflare Worker (static assets, run by [mise-cloudflare](https://github.com/joeblew999/mise-cloudflare)).

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
mise trust && mise run root:setup   # pinned tools (mise) + npm packages (aube) + Go modules
mise run dev:up                     # Hugo + tools via process-compose (mise run dev:down to stop)
```

Pushing to `main` deploys the Worker (Cloudflare's Workers Builds); every other branch gets a preview (`mise run url:branch`). `mise run cloudflare` shows where the Worker is and what's next.

## 📋 Prerequisites

- [mise](https://mise.jdx.dev/getting-started.html): installs every tool pinned in `mise.toml` (Task, Hugo, Bun, Node, Wrangler, …). Go's version comes from `go.mod`.

## 🛠️ Development Tasks

```bash
mise run hugo:dev       # Hugo dev server
mise run root:build     # Production build
mise run hugo:preview   # Preview production build
mise run hugo:clean     # Clean artifacts
mise run upgrade:all    # Upgrade tools, Go modules, the theme, the actions (Hugo: mise run upgrade:hugo)
mise run root:ci        # Everything CI checks

mise run cloudflare   # The Worker: deployed, its login, what's next
```

## 🌐 Translation Workflow

1. Edit English content in `content/english/`
2. `mise run translate:content status`: what changed since the last translation
3. `mise run translate:content next`: next file to translate
4. Translate it into each language (see CLAUDE.md, Translation Workflow)
5. `mise run translate:content done`: move the checkpoint

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
├── mise.toml               # The control plane: tool pins, env, repo-wide tasks
├── .mise/tasks/            # A task file per Go binary (cmd/<binary>), checked by misecheck
└── wrangler.jsonc          # The Worker: public/ as assets, worker/index.js for /api/*
```

## 🔧 Configuration

Secrets are never committed:

- `.env` (gitignored): e.g. `CLOUDFLARE_API_TOKEN` (read by wrangler and the `cf:*` tasks), `DEEPL_API_KEY` (autotranslate)
- Google MCP OAuth client: macOS keychain via fnox (`fnox.toml`; `mise run google-mcp:credentials`)

## 🔗 References

- [Hugo Plate Template](https://github.com/zeon-studio/hugoplate)
- [Hugo Documentation](https://gohugo.io/documentation/)
- [Cloudflare Workers static assets](https://developers.cloudflare.com/workers/static-assets/)
- [Task Documentation](https://taskfile.dev/)

## 📄 License

MIT (Hugo Plate template)





