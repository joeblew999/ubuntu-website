#!/bin/bash
# Cloudflare Pages build script
# Uses production environment for main branch, default environment for preview branches

set -e

bash scripts/theme-css.sh

if [ "$CF_PAGES_BRANCH" = "main" ]; then
  echo "🚀 Building production environment for main branch..."
  hugo -e production --gc --minify
else
  echo "🔍 Building preview environment for branch: $CF_PAGES_BRANCH"
  hugo --gc --minify
fi
