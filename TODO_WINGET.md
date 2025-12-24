TODO: Winget via Static Manifests (No REST Source)
===============================================

Intent
- For us: publish winget manifests through the existing Hugo/Cloudflare Pages pipeline—no extra server to run or maintain. CI writes manifests into `static/winget/**` and Pages makes them public.
- For users: one-click install/update driven by our CLI/Gio app; they never add sources or trust certs. Under the hood we call `winget install --manifest <latest.yaml>` from our own hosted URL.

Goal: ship winget installs/updates using static manifests hosted on our Hugo/Cloudflare Pages/GitHub Pages pipeline. No rewinged server. Our CLI/Gio app drives winget for the user.

What to build
- Hugo publishing: emit single-file winget manifests into `static/winget/<pkg>/<version>.yaml` and copy/symlink to `static/winget/<pkg>/latest.yaml`. Hugo publishes these URLs.
- Package metadata: extend `internal/vanityimport/package.go` to carry winget fields (e.g., `WingetID`, `InstallerURL`, `InstallerSHA256`, `InstallerType`, `InstallerArch`, `InstallerLocale`). Reuse `HasBinary`.
- Manifest generator: add a vanityimport subcommand (e.g., `vanityimport winget <package>`) that:
  - Pulls version/asset URL from GitHub (reuse gh integration).
  - Downloads the release asset to compute SHA256.
  - Writes a singleton manifest (ManifestVersion ≥1.6) to `static/winget/...`.
  - Updates `latest.yaml` pointer.
- Docs snippet: update CLI docs builder (`internal/cli/docs.go`) to include install text: `winget install --manifest https://www.ubuntusoftware.net/winget/<pkg>/latest.yaml`.
- App/CLI flow: add a helper in `internal/cli` or app code:
  - Ensure winget/App Installer is present; prompt if missing.
  - Fetch `latest.yaml` URL, download to temp, verify SHA256 matches manifest.
  - Shell out: `winget install --manifest <tempfile>` (or `winget upgrade --manifest` for updates).
  - For upgrades, compare manifest Version vs installed version; if newer, run install.

Manifest shape (singleton example)
```yaml
ManifestType: singleton
ManifestVersion: 1.6.0
Id: YourCo.YourApp           # from metadata
Name: YourApp
Publisher: YourCo
Version: 1.2.3
InstallerType: msix          # or exe/msi/zip
InstallerLocale: en-US
Installers:
  - Architecture: x64
    InstallerUrl: https://github.com/yourorg/yourrepo/releases/download/v1.2.3/YourApp_1.2.3_x64.msix
    InstallerSha256: <sha256>
```

Pipelines
- Generator runs in CI (GitHub Actions/Task): regenerate manifests on release, commit to repo so Hugo publishes.
- Hugo/Pages already publishes `static/**`; ensure the winget paths are included.
- (Optional) publish a small JSON index at `winget/<pkg>/index.json` with `{version, manifest}` for easy “latest” lookup.

Notes
- No winget source = no `winget search/upgrade` auto-discovery; our app/CLI owns the update check by pulling `latest.yaml`.
- Using GitHub URLs means trusted HTTPS; no cert install required.
- If we later want REST-source features, rewinged has useful patterns (manifest parsing, auto-internalization); worth mining for helper code/config ideas even if we stay static-only.
