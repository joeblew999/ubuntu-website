# TODO

Can Hugo help with capturing clicks out of the site ?

Can our cfanalytcis.

We have many Now because of the partners page any otehrs that list drone stores, and i want my reports to tell me which are most popular, and also send a report to each of our Partners via email, to keep them in the loop. Needs to be polite because many are not signed up yet, but i am hoping that is we make it really NICE they will get curious and sign up or contact us !! 

---

Current ubuntu-website state: You have significant uncommitted changes - looks like you've been refactoring CLI tools (moving things into internal/ packages, deleting duplicates). These appear to be separate from the xplat package manager work. Optional enhancements we could do:
Clean up ubuntu-website - Commit the CLI refactoring work you have staged
Add more packages - sitecheck, translate, autotranslate, vanityimport could be added to registry
Create xplat v0.2.1 release - Include the local registry support changes
Document the workflow - Add usage examples to CLAUDE.md


---

All our internal code is wrapped by a cmd.

We need this to we easy to run using narun ? 

https://github.com/akhenakh/narun
https://github.com/akhenakh/nats2sse
https://github.com/akhenakh/nlock



---


https://github.com/line/line-bot-sdk-go
https://github.com/line/line-bot-sdk-go/releases/tag/v8.18.0

https://developers.line.biz/en/docs/basics/channel-access-token/

---

JUst like LINE has a QR system, we need them too

https://line.me/R/ti/p/@linedevth?from=page



---

Like the Google MCP, we need the same.

We have a LINE account so can test it.

Maybe someone already has a golang MPC that wraps this ? 

---

Claude got an:

MCP error -32602: Output validation error: Invalid structured content for tool directory_tree: [
  {
    "expected": "string",
    "code": "invalid_type",
    "path": [
      "content"
    ],
    "message": "Invalid input: expected string, received array"
  }
]

I wonder if our claude.json and our MCP server is messing things up ?

---

SCION web page regarding high security to avoid BGP crappiness ? 
https://www.scion.org
https://github.com/scionproto/scion



---

Finish the Email System ready for Google Verification !!!

Need this so we can access other peoples google drive, email, etc without htem having to do anything.


---

### GitHub Organization: ubuntusoftware-net

**Status:** Org created, repo transfer DEFERRED

**Created:** `github.com/ubuntusoftware-net` (Dec 2024)

**Why deferred:** People are currently viewing `joeblew999/ubuntu-website`. Don't want broken links during evaluation period.

**Transfer steps (when ready):**

1. **Notify contacts** - Let key people know about the migration
2. **Transfer repo:**
   ```bash
   # Via web: github.com/joeblew999/ubuntu-website/settings → Transfer
   # Select: ubuntusoftware-net
   ```
3. **Fork back to personal:**
   ```bash
   gh repo fork ubuntusoftware-net/ubuntu-website --clone=false
   ```
4. **Update local remote:**
   ```bash
   git remote set-url origin https://github.com/ubuntusoftware-net/ubuntu-website
   git remote add personal https://github.com/joeblew999/ubuntu-website
   ```
5. **Update package frontmatter** (repo_url in `content/english/pkg/*.md`):
   ```bash
   task pkg:update  # or manually edit repo_url fields
   ```
6. **Verify:**
   - GitHub Actions still run on org
   - Vanity imports work: `go get www.ubuntusoftware.net/pkg/mailerlite`
   - Old links redirect (GitHub does this automatically)

**Notes:**
- GitHub Actions free tier is same for orgs (2000 min/month private, unlimited public)
- Vanity URLs (`www.ubuntusoftware.net/pkg/*`) shield users from GitHub structure changes
- GitHub auto-redirects old URLs for a period after transfer

---

- https://github.com/romshark/toki might help use with Transaltion of pages ? 
- He is currently updating a branch to add Datastar real time Web gui, which is perfert for us.

It does not parse markdown yet, so we can fork and add that somehow, so that it will help with use translating markdown.

I have an issue to add Date/Time at: https://github.com/romshark/toki/issues/18


---

We need to ensure we release this via github and that we then finish the ability for anyone sending in an email from "Get Started" can have the software sent to them via email. I think we need to finsih the mailerlite stuff in order to do that . Will be a sprint ...

---

BTW i noticed that your broke an GOLDNE run in the wrangler taskfile. I see "bun" used in it. We agreed thats NOT ok. I dont know how this was missed, but its vital you dont so this !!!!   I am not sure how your going to validate this, but now that your have Archetypes and fmt and lint, i think it will be pretty easy, and you know what task includes what task  ? 

I need you to keep watch on where you see patterns or mistakes when you modify task file, and screw up, and then look at the xplat lint and fmt etc code to be updated. We dont want to go overboard and boil the ocean.

---

Currently honing our testing of task files, archtype classification and task files conformance system.

Getting single task file testing locally and in CI, will allow Me, other devs and claude to QUICKLY get things properly validated !!

--



xplat is nearing prime time.

We have "release" that does the release into github.

But we want to release via winget and homebrew. Dont bother with linux yet.

To make it such that we and any user can use xplat to release we want a smat way. We have the gh cli task file, and winget and homebew are github Repos, so maybe we can have a template inside xplat and push to these repos ?

I am SURE i have left out a few aspects though, so i need help.

---

Once oyu get the NON CGO FULLY working for the FULL lifecyle, Make sure that dummy cmd main.go can also require CGO. What i mean is that then you can test for both CGO and non CGO with the same binary for all the round tripping permutations. Its really easy. you just need to focus .  you could also make a 2nd main.go too of how ever you want to do it.   



---

## Follow-ups (2026-11-03: mise tooling, upgrades, Hugo 0.169, Hugo Plate update)

Upgrade commands are in CLAUDE.md ("Upgrading: always through Taskfile"). Tick items here when done.

### Decisions needed

- [ ] **Taskfile → mise tasks** (you prefer mise). Size: 36 Taskfiles, ~6,800 lines, 489 tasks; 6 CI workflows and process-compose call `task`; the published Google MCP post uses `task google-mcp:*`. Plan: migrate domain by domain, switch CI to `mise run`, delete each Taskfile once its CI passes, drop the `task` pin last.
- [x] **Global mise config**: `~/.config/mise/config.toml` got patch bumps at 2026-11-03 13:59 for ast-grep, cargo-binstall, git-lfs, pnpm, usage (during the project upgrade; cause unclear). Keep or revert. `task tools:upgrade` now passes an explicit project-only tool list. **Done:** kept (patch bumps of your own global tools; the project task can't touch them any more).

### Content

- [ ] Translation backlog: `task translate:content:status` lists many English pages changed since the last checkpoint (platform/*, pkg/*, …). Translate to de/zh/ja/vi, then `task translate:content:done`.
- [x] Google Workspace MCP draft post is dated 2024-12-15: probably the same year typo as the DJI post (fixed to 2025-12-24). Set the real date when publishing. **Done:** date fixed to 2025-12-15 (en + vi); still `draft: true`.
- [x] README top lists personal LAN URLs (`http://192.168.0.156:1313/`): keep or remove. **Done:** removed (`task hugo:dev:lan` prints the current LAN URL).

### Theme / Hugo

- [x] Review our 17 `layouts/` overrides of Hugo Plate templates (head, header, script, footer, style, blog/single, …): upstream fixes to those files don't reach us. `task theme:update` lists the ones upstream changed. **Done:** 3-way merged upstream's changes (base = previously vendored): 9 clean, 4 conflicts resolved by hand (header keeps our `xl:` nav + upstream's dropdown fix and external-link target; call-to-action takes `image_alt`; blog/single follows upstream, which removed Disqus), basic-seo kept ours. `task theme:update` now does this merge automatically (conflicts → `.theme-merge/`).
- [x] Upstream Hugo Plate imports optional modules we don't: open-remark, llms-txt, gethugothemes mermaid (we use hugomods/mermaid). Decide. **Done:** added llms-txt (`/llms.txt`, `/llms-full.txt`, `.md` per page); not open-remark (hosted comments needing a zeon.studio account); removed hugomods/mermaid (no page uses diagrams; upstream dropped the partial).
- [x] zeon-studio/hugoplate#250 (Hugo ≥ 0.165 needs `tailwindcss` in `security.exec.allow`): watch for a reply; add that `scripts/themeGenerator.js` has no Hugo Modules mode (we wrap it in `scripts/theme-css.sh`). **Done:** commented with the Hugo Modules gap (https://github.com/zeon-studio/hugoplate/issues/250#issuecomment-6598551147); replies arrive as GitHub notifications.
- [x] CLAUDE.md points at source checkouts that don't exist (`wellknown/.src/hugo`, `wellknown/.src/hugoplate`): clone them or drop the references. **Done:** now points at GitHub (Hugo, matching tag) and `_vendor/` + the Go module cache (Hugo Plate at the pinned commit).

### Tooling

- [ ] The "cmd/translate - DO NOT RUN" note below looks stale: `task translate:*` builds and uses cmd/translate and CLAUDE.md documents it. Verify the mutating commands, then remove the warning.
- [x] Two publish targets: `.github/workflows/release-pages.yml` deploys to GitHub Pages on every push to main, while production is Cloudflare Pages. Keep both or drop GitHub Pages. **Done:** dropped GitHub Pages (workflow removed, Pages disabled): it served a duplicate copy of the site; CI still builds Hugo in `taskfile.yml`.
- [x] `internal/browser.FindPlaywrightBinary` means two things: Gmail calls it expecting our `cmd/playwright` (`open --url`), but it prefers `node_modules/.bin/playwright` (npm CLI, different syntax) when present. Make Gmail call our tool explicitly. **Done:** Gmail's browser mode now opens compose in the default (signed-in) browser directly; the CLI shell-out never worked (`open --url` isn't our CLI's syntax) and is gone, as is the binary lookup.
- [x] `bun install` reports "Blocked 1 postinstall": check `bun pm untrusted` and trust it if it's needed. **Done:** it's `@parcel/watcher` compiling from source; the prebuilt binary is used, so it stays blocked.
- [x] Still installed outside mise: tippecanoe (not in the registry), gcloud (optional), task-ui (`go install …@latest` from xplat, unpinned), Claude CLI (curl/npm/bun routes in `taskfiles/tools/Taskfile.claude-cli.yml`). **Done:** task-ui removed (xplat is archived and dropped it; titpetric/task-ui is now a different product); Terraform flow no longer needs gcloud (credentials via `task google-mcp:tf:auth`); tippecanoe isn't in any mise backend (optional, `task airspace:install:tippecanoe`); Claude CLI installs one way (Anthropic's installer) and is a machine-wide tool, not a project pin.
- [x] Disk: old Playwright caches can go (`~/Library/Caches/ms-playwright` chromium-1228/1243, `ms-playwright-go/1.57.0`). **Done:** left alone: `~/Library/Caches/ms-playwright` is shared across projects (chromium-1243 is newer than ours, so another tool uses it).

### Also fixed along the way

- [x] Blog posts embedded Disqus from the theme vendor's demo account (`themefisher-template`); removed with upstream's change.
- [x] autotranslate left `[[NOTRANSLATE_n]]` placeholders in output at random (map-order restore of nested placeholders broke relref links); restores newest-first now, refuses to write leftovers, regression tests added.
- [x] `browser.GetProjectRoot` never matched `.git` (checked it as a file).
- [x] Hugo module `min` raised to 0.158.0 (what Hugo Plate needs).

### Decided (leave as is)

- Global `tailwindcss` mise shim stays: shadcn-templ-demo and go-htmx4 pin it.
- The youtube CLI keeps downloading its own ffmpeg/yt-dlp into `.build/`: the standalone binary needs it.
- `task google:auth:add` keeps writing the generic env-var MCP entry (standalone CLI) and never overwrites an existing entry.
- `go.mongodb.org/mongo-driver` v1 (deprecated) is indirect via playwright-go → golang-set: upstream's to fix.

---

## Completed

- [x] xplat v0.2.0 released with Windows support, mv command, unit tests
- [x] Remote taskfile includes work (validated in CI)
- [x] Cross-platform CI testing (Linux, macOS, Windows)
- [x] gh CLI taskfile added (tools:gh:*) - releases, workflows, PRs, issues
- [x] xplat release workflow migrated from softprops/action-gh-release to gh CLI
- [x] TOOL_BIN pattern established ({{exeExt}} for Windows .exe support)
- [x] Task version locked to 3.45.5 across all CI workflows
- [x] Task version check taskfile added (tools:task:check:deps)
- [x] Taskfile Registry validated via .github/workflows/remote-taskfile-test.yml
- [x] Deprecated deploy.yml.old workflow deleted
- [x] Centralized xplat.env - single source of truth for Task/Go versions
- [x] Go version standardized to 1.25 across all workflows (matches go.mod)
- [x] Version consistency check added (ci:check:versions)
- [x] DRY version management - xplat.env loaded via root Taskfile dotenv
  - Taskfile tasks read $GO_VERSION, $TASK_VERSION from xplat.env
  - Workflow YAML still has hardcoded values (GitHub Actions limitation)
  - ci:check:versions validates workflow values match xplat.env
- [x] Go toolchain version management (rustup-style parity)
  - Flexible mode (default): warns on version mismatch
  - Strict mode: fails on version mismatch (GO_VERSION_STRICT=true)
  - Bootstrap task: `task toolchain:golang:bootstrap` uses `golang.org/dl`
- [x] Taskfile refactor - extracted domain modules from root (900→355 lines)
  - New modules: hugo, cf, translate, seo, sitecheck, env, url
  - gh:secret/var tasks merged into tools:gh
  - All 146 tasks tested and working
- [x] Google MCP server integration with turnkey OAuth
  - Terraform creates GCP project with 6 Google APIs enabled
  - cmd/google-auth binary handles OAuth with PKCE
  - Three auth modes: manual (default), passkey (-assisted), auto (-auto with Playwright)
  - Account pre-selection via -account=EMAIL flag (uses login_hint)
  - Tasks: `google-mcp:tf:auth:passkey`, `google-mcp:tf:auth:auto`, `google-mcp:tf:plan/apply`
  - No brew dependencies - cross-platform compatible

## In Progress

(nothing currently)

---

## Broken/Incomplete

### cmd/translate - DO NOT RUN

Current workflow: Claude Code + Taskfile shell scripts (translate:status, translate:done, etc.)

The manual workflow is prototyping what cmd/translate should eventually do.
Once patterns are solid from manual use, codify them into cmd/translate.

---

## Future Ideas

### Bigger Picture

Via (Datastar web GUI) & Hugo extension work can leverage each other ? Maybe for collaboration aspects. real time, etc. Can redude what we are building to be simpler for inside vscode and outside in terms of code and compleyity. Def need to visit this !!

Building toward:
1. **Via** - User IDE web GUI using Datastar
2. **Hugo extension** - VSCode extension for Hugo content editing
3. **Block editor** - Datastar-based editor for Hugo content/layouts

Users and VSCode devs call the same Go code, both modifying git repos.

### Split Taskfiles to .github Repo?

Could move reusable taskfiles to a dedicated repo so all repos can pull versioned taskfiles at runtime. Task supports this natively.

### other weird ideas

https://github.com/akhenakh/narun - could help generalize cmd/* tools

https://github.com/akhenakh/nats2sse for feed off 

https://github.com/infogulch/xtemplate/tree/next


### gh CLI - Future Enhancements

Now that we have `tools:gh:*` tasks and migrated xplat release to use gh CLI:
- Consider migrating other workflows to use gh CLI
- Add `gh:workflow:run` to trigger CI from local (useful for manual releases)
- Explore `gh api` for more advanced GitHub integrations

### Logo/Branding Updates

When logo changes, need to update: Bluesky, Gmail signature.
Could automate via Go code + task.
