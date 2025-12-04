# TODO

## Completed

- [x] xplat v0.2.0 released with Windows support, mv command, unit tests
- [x] Remote taskfile includes work (validated in CI)
- [x] Cross-platform CI testing (Linux, macOS, Windows)

## In Progress

### Taskfile Registry Readiness

The taskfile system is designed to be a registry that devs can extend and users can consume via remote includes.

**Current structure:**
```
taskfiles/
├── Taskfile.tools.yml      → tools/Taskfile.xplat.yml (nested, reusable)
├── Taskfile.toolchain.yml  → toolchain/{golang,rust}.yml (nested, reusable)
├── Taskfile.ci.yml         (flat, mostly reusable)
├── Taskfile.git.yml        (flat, reusable)
└── Taskfile.lanip.yml      (flat, project-specific)
```

**What works:**
- `xplat` bootstraps via build-from-source OR GitHub release download
- Remote includes validated in `.github/workflows/remote-taskfile-test.yml`
- REQUIRES comment system for toolchain autodiscovery

**Next steps:**
- [ ] Test consumption from a fresh repo (real world validation)
- [ ] Document consumption pattern for external devs

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

Building toward:
1. **Via** - User IDE web GUI using Datastar
2. **Hugo extension** - VSCode extension for Hugo content editing
3. **Block editor** - Datastar-based editor for Hugo content/layouts

Users and VSCode devs call the same Go code, both modifying git repos.

### Split Taskfiles to .github Repo?

Could move reusable taskfiles to a dedicated repo so all repos can pull versioned taskfiles at runtime. Task supports this natively.

### narun

https://github.com/akhenakh/narun - could help generalize cmd/* tools

### Logo/Branding Updates

When logo changes, need to update: Bluesky, Gmail signature.
Could automate via Go code + task.
