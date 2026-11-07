# Ubuntu Software Website

**The problem.** A company's website in five languages (English, German, Chinese, Japanese, Vietnamese) drifts a page at a time. An English page changes and nobody can say which of its four translations are now stale or missing; a menu links to a page one language doesn't have; a link breaks. And everything around the pages is a step someone has to remember, done a little differently on every machine: the build, the deploy and a preview of a change, the form that adds a sign-up to the mailing list, new blog posts on Bluesky, noticing that the site is down or that its traffic has changed.

**What ubuntu-website does.** It is [www.ubuntusoftware.net](https://www.ubuntusoftware.net): Hugo builds it on the Hugo Plate theme, from the English pages in `content/english/` and their translations, and a Cloudflare Worker serves the build, with `/api/*` for what needs code (a Web3Forms sign-up added to MailerLite). Every step is a mise task, every tool pinned. `mise run translate:content` says which English pages changed since the last translation and which to translate next, and every Monday a workflow opens an issue for what's missing or stale. `mise run root:ci`, the one thing GitHub's CI runs, builds the site as production does and checks that every link in it lands. A push to `main` deploys it through Cloudflare's Workers Builds, and every other branch gets a preview, run by [mise-cloudflare](https://github.com/joeblew999/mise-cloudflare). Scheduled workflows check the live site from around the world every 6 hours, report a big change in its traffic every Sunday, and post new blog posts to Bluesky.

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





