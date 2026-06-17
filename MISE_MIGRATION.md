# mise + nushell + fnox Migration

Aligns this repo with the standard toolchain used across our other repos:

| Tool | Role | Replaces / supplements |
|------|------|------------------------|
| [mise](https://mise.jdx.dev) | tool versions + env + task runner | go-task (Taskfile), ad-hoc installs |
| [nushell](https://www.nushell.sh) | task scripting | POSIX `sh` task bodies |
| [fnox](https://github.com/jdx/fnox) | encrypted/remote secrets | gitignored `.env` |

The migration is **incremental and non-disruptive**: `mise.toml` is additive and
coexists with the existing `Taskfile` (mise even keeps `task` as a managed tool),
so nothing breaks while we port module by module.

## Status

**Phase 0 — scaffold (this branch):**
- [x] `mise.toml` — pins tools (go, node, hugo, task, nushell, fnox) + project env
- [x] `mise-tasks/check-env` — pilot **nushell** task (`mise run check-env`)
- [x] Pilot `mise.toml` tasks wrapping existing Go commands (analytics, sitecheck, build)
- [ ] fnox wired for secrets (see [Secrets](#secrets-fnox) — needs `fnox init` + backend)

**Phase 1 — env & secrets:** migrate `.env` → fnox; enable the `fnox-env` plugin block in `mise.toml`.

**Phase 2 — port tasks:** convert `taskfiles/*.yml` modules to mise tasks (nushell where logic is non-trivial), one namespace at a time.

**Phase 3 — cut over:** once parity is reached, retire `Taskfile.yml` + `taskfiles/` and update `CLAUDE.md` ("USE TASKFILE" → "USE MISE").

## Quick start

```bash
mise install          # provision go, node, hugo, task, nushell, fnox
mise run check-env    # nushell pilot — verify secrets/env are present
mise tasks            # list mise tasks
mise run build        # hugo --gc --minify
```

## Taskfile → mise concept mapping

| Taskfile (go-task) | mise |
|--------------------|------|
| `tasks:` | `[tasks.*]` in `mise.toml`, or file tasks in `mise-tasks/` |
| `desc:` | `description` |
| `deps:` | `depends` |
| `cmds:` | `run` |
| `status:` (idempotency) | `sources` / `outputs` |
| `requires.vars` | assert in the task body (nushell `if ($x | is-empty) { exit 1 }`) |
| `includes:` | multiple config files / `[task_config] includes` |
| `vars:` + `{{.VAR}}` (Go templates) | `[env]` + `$VAR`; Tera `{{ }}` for templating |
| `dotenv: [.env]` | `[env._.file]` (`.env` auto-loaded) → later replaced by fnox |

### Nushell in tasks

go-task embeds a POSIX-sh interpreter (`mvdan/sh`), so nushell there means
wrapping calls as `nu -c "..."`. In mise, a **file task** with a `#!/usr/bin/env nu`
shebang in `mise-tasks/` is first-class (and gets IDE highlighting). See
`mise-tasks/check-env` for the reference pattern.

## Secrets (fnox)

Secrets must **not** live in `mise.toml` (it's committed). They move to fnox.
`fnox.toml` is intentionally **not** committed yet — generate it with `fnox init`
so the schema/backend matches our other repos, then add the keys below.

Secrets currently in `.env` to migrate:

| Key | Used by |
|-----|---------|
| `CLOUDFLARE_API_TOKEN` | cfanalytics, cfemail, r2, wrangler |
| `MAILERLITE_API_KEY` | mailerlite |
| `DEEPL_API_KEY` | autotranslate |
| `BLUESKY_APP_PASSWORD` | syndication (also a GitHub secret) |
| Google OAuth (`CLIENT_ID`/`CLIENT_SECRET`) | google-auth / gmail / calendar |

Setup outline:

```bash
mise use -g fnox                 # (already pinned in mise.toml)
fnox init                        # choose backend (1Password / Bitwarden / age / ...)
fnox set CLOUDFLARE_API_TOKEN    # repeat per key
fnox exec -- mise run check-env  # secrets injected for the command
```

Then enable the plugin block in `mise.toml` so secrets auto-load as env vars:

```toml
[plugins]
fnox-env = "https://github.com/jdx/mise-env-fnox"

[env]
_.fnox-env = { tools = true }
```

> If our other repos already have a committed `fnox.toml` convention, point me at
> one and I'll mirror its backend/profile structure instead of `fnox init` defaults.

## Notes

- **Hugo must be the extended build** (Hugo Plate requires v0.144+ extended). Confirm `hugo version` shows `+extended` after `mise install`.
- `node = "lts"` covers the Hugo Plate (tailwind) theme build.
- mise tasks deliberately reference only commands present on `main`; the `cfemail`
  tool lives on the `claude/mailerlite-integration-yyrdem` branch and gets its mise
  task once that merges.
