#!/usr/bin/env bash
# Update the Hugo Plate theme, repeatably. Run via: task theme:update [REF=<hugoplate commit>]
#
# Hugo Modules pin the exact commits (recorded in _vendor/modules.txt) and `hugo mod vendor` turns
# every theme change into a reviewable diff. On top of that, this does what upstream's own
# themeUpdate.js (curl main | tar) does not:
#   - syncs scripts/themeGenerator.js from the pinned commit (kept unmodified)
#   - bumps/adds the npm packages upstream's package.json needs (e.g. tailwind-bootstrap-grid)
#   - regenerates assets/css/generated-theme.css and runs a production build
#   - lists changed theme templates that our layouts/ override (upstream fixes we don't get)
# go.mod/go.sum are restored afterwards: they belong to the Go code and `go mod tidy` strips
# Hugo's entries anyway; _vendor/modules.txt is the theme pin.
#
# Re-running with nothing new upstream changes nothing.
set -euo pipefail
cd "$(dirname "$0")/.."

module=github.com/zeon-studio/hugoplate
ref=${1:-}
# Upstream packages we deliberately don't use (process-compose runs the dev processes instead).
skip_deps=" concurrently "

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
cp go.mod go.sum "$tmp/"
restore_gomod() { cp "$tmp/go.mod" "$tmp/go.sum" .; }

pinned() { grep "^# $1 " _vendor/modules.txt | awk '{print $3}'; }
before=$(pinned "$module")

echo "==> Modules (now: hugoplate ${before##*-})"
if ! { hugo mod get -u && { [ -z "$ref" ] || hugo mod get "$module@$ref"; } && hugo mod vendor; }; then
  restore_gomod
  echo "✗ hugo mod get/vendor failed"
  exit 1
fi
restore_gomod
after=$(pinned "$module")
sha=${after##*-}
echo "    hugoplate ${before##*-} → $sha"

echo "==> scripts/themeGenerator.js"
src="$(go env GOMODCACHE)/$module@$after/scripts/themeGenerator.js"
if [ -f "$src" ]; then
  cp "$src" scripts/themeGenerator.js
else
  gh api -H "Accept: application/vnd.github.raw" "repos/zeon-studio/hugoplate/contents/scripts/themeGenerator.js?ref=$sha" > scripts/themeGenerator.js
fi
chmod 644 scripts/themeGenerator.js
git diff --quiet -- scripts/themeGenerator.js && echo "    unchanged" || echo "    updated from upstream @ $sha"

echo "==> npm packages (vs upstream package.json @ $sha)"
add=""
while read -r name range; do
  case "$skip_deps" in *" $name "*) continue ;; esac
  ours=$(jq -r --arg n "$name" '.devDependencies[$n] // .dependencies[$n] // ""' package.json)
  up_v=${range#[\^~]}
  our_v=${ours#[\^~]}
  if [ -z "$ours" ]; then
    add="$add $name@$range"
  elif [ "$up_v" != "$our_v" ] && [ "$(printf '%s\n%s\n' "$up_v" "$our_v" | sort -V | tail -1)" = "$up_v" ]; then
    add="$add $name@$range"
  fi
done < <(jq -r '((.devDependencies // {}) + (.dependencies // {})) | to_entries[] | "\(.key) \(.value)"' "_vendor/$module/package.json")
if [ -n "$add" ]; then
  echo "    bun add -d$add"
  # shellcheck disable=SC2086
  bun add -d $add
else
  echo "    up to date"
fi

echo "==> Theme CSS"
bash scripts/theme-css.sh

echo "==> Changed theme templates that our layouts/ override (upstream changes we don't get)"
overridden=0
while read -r f; do
  rel=${f#*/layouts/}
  for cand in "$rel" "${rel/#partials\//_partials/}" "${rel/#_partials\//partials/}"; do
    if [ -f "layouts/$cand" ]; then
      echo "    layouts/$cand"
      overridden=$((overridden + 1))
      break
    fi
  done
done < <(git diff --name-only HEAD -- '_vendor/*/layouts/*'; git ls-files --others --exclude-standard -- '_vendor/*/layouts/*')
[ "$overridden" -gt 0 ] || echo "    none"

echo "==> Production build"
if ! hugo -e production --gc --minify -d "$tmp/public" > "$tmp/build.log" 2>&1; then
  grep -E "^(ERROR|Error)" "$tmp/build.log" | head -5
  echo "✗ build failed: fix it before committing (full log above is from: hugo -e production)"
  exit 1
fi
grep -E "^WARN" "$tmp/build.log" | sed 's/^/    /' || true
grep -E "Pages |Total in" "$tmp/build.log" | sed 's/^/    /'

echo ""
if git diff --quiet -- _vendor package.json bun.lock scripts/themeGenerator.js assets/css/generated-theme.css \
  && [ -z "$(git ls-files --others --exclude-standard -- _vendor)" ]; then
  echo "✓ Already up to date (hugoplate $sha)"
else
  echo "✓ Updated. Review: git diff --stat -- _vendor package.json scripts assets/css"
  echo "  Ship: push a branch, check the Cloudflare preview build, then merge to main."
fi
