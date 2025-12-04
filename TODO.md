# TODO

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

## In Progress

(nothing currently)

---

## Broken/Incomplete

### cmd/env - DO NOT RUN

The env tool expects a specific .env structure and will corrupt your .env file.

Investigate https://github.com/helmfile/vals for proper secrets management.
vals supports multiple backends (Vault, AWS SSM, GCP Secrets, SOPS, etc.)

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

### narun

https://github.com/akhenakh/narun - could help generalize cmd/* tools

### gh CLI - Future Enhancements

Now that we have `tools:gh:*` tasks and migrated xplat release to use gh CLI:
- Consider migrating other workflows to use gh CLI
- Add `gh:workflow:run` to trigger CI from local (useful for manual releases)
- Explore `gh api` for more advanced GitHub integrations

### Logo/Branding Updates

When logo changes, need to update: Bluesky, Gmail signature.
Could automate via Go code + task.
