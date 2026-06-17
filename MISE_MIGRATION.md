# mise + nushell + fnox Migration

Aligns this repo with the house toolchain (see `joeblew999/.github`, the shared
mise task library, and `joeblew999/.github-example`, the canonical consumer):

| Tool | Role |
|------|------|
| [mise](https://mise.jdx.dev) | tool versions + env + task runner |
| [nushell](https://www.nushell.sh) | task scripting (inline in `tasks/*.toml`) |
| [fnox](https://github.com/jdx/fnox) | encrypted/remote secrets (replaces `.env`) |

## Architecture (house convention)

**Two layers, consumed by-reference — identical locally and in CI:**

1. **Shared library** (`joeblew999/.github`) — generic tooling included by git ref,
   never copied:
   ```toml
   [task_config]
   includes = [
     "git::https://github.com/joeblew999/.github.git//tasks/tool-cf.toml?ref=main",
     # ...pin ?ref= to a release tag in production; main is the rolling ref
   ]
   ```
   Provides: `cf:*`, `gh:*`, `wrangler:*`, `fnox:*`, `secrets:*`, `mise:*`, `ci`, `release`.

2. **Project tasks** (local `tasks/*.toml`) — this repo's modules, same format:
   inline nushell via `run = '''#!/usr/bin/env nu …'''`, named `namespace:task`.

Tasks are **not** separate `.nu` files in `mise-tasks/` — they are inline nushell
inside `tasks/*.toml` (matches the shared library exactly).

## Status

**Phase 0 — foundation (done, verified with `mise tasks ls`):**
- [x] `mise.toml`: tools (`go`, `node`, `hugo-extended`, `github:nushell/nushell`, `fnox`, `gh`, `jq`, `usage`, `task`), project env, `.env` auto-load
- [x] Shared library included by-reference (remote git includes resolve)
- [x] Pilot modules: `tasks/hugo.toml`, `tasks/analytics.toml`, `tasks/sitecheck.toml`

**Phase 1 — port project modules** (in progress — see checklist below).

**Phase 2 — secrets:** `.env` → fnox; enable the `fnox-env` plugin block in `mise.toml`.

**Phase 3 — cut over:** once parity is verified, delete `Taskfile.yml` + `taskfiles/`,
update CI workflows from `task …` to `mise run …`, flip `CLAUDE.md` ("USE TASKFILE" → "USE MISE").

## Module port checklist

Each row = one `taskfiles/*.yml` → local `tasks/*.toml`, **unless** covered by the
shared library (then no local port — just call the shared task).

| Source Taskfile | Target | Notes |
|-----------------|--------|-------|
| hugo | ✅ `tasks/hugo.toml` | done |
| cfanalytics | ✅ `tasks/analytics.toml` | done |
| sitecheck | ✅ `tasks/sitecheck.toml` | done |
| mailerlite | `tasks/mailerlite.toml` | wraps cmd/mailerlite |
| translate | `tasks/translate.toml` | wraps cmd/translate |
| autotranslate | `tasks/autotranslate.toml` | wraps cmd/autotranslate |
| airspace | `tasks/airspace.toml` | wraps cmd/airspace + tippecanoe |
| genlogo | `tasks/genlogo.toml` | wraps cmd/genlogo |
| youtube | `tasks/youtube.toml` | wraps cmd/youtube |
| google | `tasks/google.toml` | wraps cmd/google |
| google-mcp | `tasks/google-mcp.toml` | MCP server setup |
| vanityimport | `tasks/vanityimport.toml` | wraps cmd/vanityimport |
| url | `tasks/url.toml` | browser-open shortcuts |
| env | `tasks/env.toml` | wraps cmd/env |
| cf-r2 | `tasks/r2.toml` | project R2 bucket/endpoints (uses shared `cf:r2-*`) |
| workers | `tasks/workers.toml` | Cloudflare Workers (uses shared `wrangler:*`) |
| lanip | `tasks/lanip.toml` | wraps cmd/lanip |
| toki | `tasks/toki.toml` | i18n toolkit |
| playwright | `tasks/playwright.toml` | browser automation |
| dev / pc / task-ui | `tasks/dev.toml` | process-compose + Task-UI orchestration |
| cf | **shared** `cf:*` + local project-deploy task | Pages deploy is project-specific |
| gh | **shared** `gh:*` | no local port |
| wrangler | **shared** `wrangler:*` | no local port |
| release | **shared** `release` + per-tool `release:build` | |
| git | keep on `task` (dev-only) | low value to port |
| coder / coder-fly / coder-templates | TBD | confirm still in use before porting |
| claude-cli | TBD | dev-only |

## Quick start

```bash
mise install            # provision tools (needs GitHub auth for aqua/github backends)
mise tasks ls           # list all tasks (shared + local)
mise run hugo:dev       # start Hugo
mise run sitecheck:all  # run all site checks
```

## Secrets (fnox)

Secrets must **not** live in `mise.toml` (committed) or `.env` long-term. They move
to fnox; the shared `secrets:*` / `fnox:*` tasks manage them. `fnox.toml` is generated
by `fnox init` so the backend/profile structure matches our other repos.

Secrets to migrate: `CLOUDFLARE_API_TOKEN`, `MAILERLITE_API_KEY`, `DEEPL_API_KEY`,
`BLUESKY_APP_PASSWORD`, Google OAuth (`CLIENT_ID`/`CLIENT_SECRET`).

Then enable in `mise.toml`:
```toml
[plugins]
fnox-env = "https://github.com/jdx/mise-env-fnox"
[env]
_.fnox-env = { tools = true }
```

## Notes

- **Hugo extended** is enforced via the `hugo-extended` tool pin.
- In this sandbox, aqua/github tool installs need a GitHub token (unauthenticated
  rate limit); remote *task includes* (git clone) work fine.
- `cfemail` (on branch `claude/mailerlite-integration-yyrdem`) gets `tasks/cfemail.toml` once merged.
