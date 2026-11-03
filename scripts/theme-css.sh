#!/usr/bin/env bash
# Generate assets/css/generated-theme.css from data/theme.json with Hugo Plate's own
# scripts/themeGenerator.js, kept unmodified (task theme:update syncs it from the pinned commit).
#
# The generator only knows exampleSite/ and themes/<name>/ layouts, not a site that imports Hugo
# Plate as a module, so it runs in a temporary exampleSite-shaped directory and the result is
# copied to assets/css/ (project assets override the theme module's default generated-theme.css).
#
# Used by: task theme:css (and hugo:* tasks), scripts/build.sh (Cloudflare Pages), CI, process-compose.
set -euo pipefail
root=$(cd "$(dirname "$0")/.." && pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

mkdir -p "$tmp/scripts" "$tmp/exampleSite/data"
cp "$root/scripts/themeGenerator.js" "$tmp/scripts/"
cp "$root/data/theme.json" "$tmp/exampleSite/data/"
: > "$tmp/exampleSite/hugo.toml"

node "$tmp/scripts/themeGenerator.js" >/dev/null
mkdir -p "$root/assets/css"
cp "$tmp/assets/css/generated-theme.css" "$root/assets/css/generated-theme.css"
echo "✓ assets/css/generated-theme.css (from data/theme.json)"
