# MailerLite Integration Guide

This document explains how to connect Web3Forms submissions to MailerLite for automated software delivery emails.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         DEPLOYMENT OPTIONS                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  Option A: Direct Integration (Recommended for Production)               │
│  ─────────────────────────────────────────────────────────              │
│  Web3Forms → MailerLite (via Zapier/Make)                               │
│  - No server needed                                                      │
│  - Works with Web3Forms Pro                                             │
│                                                                          │
│  Option B: Cloudflare Worker (Serverless)                               │
│  ────────────────────────────────────────                               │
│  Web3Forms → Cloudflare Worker → MailerLite API                         │
│  - Free tier available                                                   │
│  - Global edge deployment                                                │
│                                                                          │
│  Option C: Self-Hosted Server                                            │
│  ────────────────────────────────────────                               │
│  Web3Forms → mailerlite server → MailerLite API                         │
│  - For development/testing                                               │
│  - Requires tunnel for production (ngrok, cloudflared)                  │
│                                                                          │
│  Option D: GitHub Actions (Periodic Sync)                               │
│  ────────────────────────────────────────                               │
│  Scheduled job that syncs Web3Forms submissions to MailerLite           │
│  - Batch processing                                                      │
│  - Good for low-volume                                                   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

## Option A: Direct Integration via Zapier (Recommended)

This is the simplest approach for production use.

### Prerequisites
- Web3Forms Pro account (for webhook feature)
- Zapier account (free tier works for low volume)
- MailerLite account with API key

### Setup Steps

1. **Create Zapier Zap:**
   - Trigger: Webhooks by Zapier → Catch Hook
   - Action: MailerLite → Create or Update Subscriber

2. **Configure Web3Forms:**
   - Go to Form Settings → Integrations
   - Add webhook URL from Zapier

3. **Map Fields in Zapier:**
   ```
   email    → Subscriber Email
   name     → Name field
   company  → Custom field "company"
   ```

## Option B: Cloudflare Worker (Serverless)

For self-hosted serverless deployment.

### Worker Code

Create `workers/mailerlite-webhook.js`:

```javascript
// Cloudflare Worker for Web3Forms → MailerLite
export default {
  async fetch(request, env) {
    if (request.method !== 'POST') {
      return new Response('Method not allowed', { status: 405 });
    }

    try {
      const formData = await request.formData();
      const email = formData.get('email');
      const name = formData.get('name') || '';
      const company = formData.get('company') || '';

      if (!email) {
        return new Response('Email required', { status: 400 });
      }

      // Add to MailerLite
      const response = await fetch('https://connect.mailerlite.com/api/subscribers', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${env.MAILERLITE_API_KEY}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          fields: { name, company },
          groups: [env.MAILERLITE_GROUP_ID], // "Get Started" group
        }),
      });

      if (!response.ok) {
        console.error('MailerLite error:', await response.text());
      }

      return new Response('OK', { status: 200 });
    } catch (error) {
      console.error('Worker error:', error);
      return new Response('OK', { status: 200 }); // Don't retry
    }
  },
};
```

### Deploy with Wrangler

```bash
# wrangler.toml
name = "mailerlite-webhook"
main = "workers/mailerlite-webhook.js"
compatibility_date = "2024-01-01"

[vars]
MAILERLITE_GROUP_ID = "173641063022462964"

# Add secret via CLI:
# wrangler secret put MAILERLITE_API_KEY
```

```bash
task wrangler:deploy
```

## Option C: Self-Hosted Server (Development)

Use the built-in mailerlite CLI server for development and testing.

### Local Development

```bash
# Start server
task mailerlite:server PORT=8086 GROUP_ID=173641063022462964

# Test with curl
curl -X POST http://localhost:8086/webhook \
  -d "name=Test&email=test@example.com&company=Test Co"
```

### Production with Tunnel

```bash
# Option 1: ngrok
ngrok http 8086

# Option 2: Cloudflare Tunnel
cloudflared tunnel --url http://localhost:8086

# Then configure Web3Forms webhook URL to the tunnel URL
```

## Option D: GitHub Actions (Batch Processing)

For periodic sync without real-time webhooks.

### Workflow: `.github/workflows/sync-mailerlite.yml`

```yaml
name: Sync MailerLite

on:
  schedule:
    - cron: '0 */6 * * *'  # Every 6 hours
  workflow_dispatch:

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'

      - name: Build mailerlite CLI
        run: go build -o mailerlite ./cmd/mailerlite

      - name: Sync subscribers
        env:
          MAILERLITE_API_KEY: ${{ secrets.MAILERLITE_API_KEY }}
        run: |
          # Check for new Web3Forms submissions via their API
          # Add to MailerLite
          ./mailerlite stats
```

## MailerLite Automation Setup

Once subscribers are in the "Get Started" group, set up an automation:

1. **Go to:** https://dashboard.mailerlite.com/automations
2. **Create automation:**
   - Trigger: "When subscriber joins group" → "Get Started"
   - Action: Send email with download links

### Email Template Variables

Use these in your MailerLite email template:

```
Download Links:
- macOS (Apple Silicon): https://github.com/joeblew999/ubuntu-website/releases/latest/download/software-darwin-arm64.tar.gz
- macOS (Intel): https://github.com/joeblew999/ubuntu-website/releases/latest/download/software-darwin-amd64.tar.gz
- Windows: https://github.com/joeblew999/ubuntu-website/releases/latest/download/software-windows-amd64.exe
- Linux: https://github.com/joeblew999/ubuntu-website/releases/latest/download/software-linux-amd64.tar.gz

Or visit: https://github.com/joeblew999/ubuntu-website/releases/latest
```

## CLI Commands Reference

```bash
# Subscriber management
task mailerlite:subscribers:list
task mailerlite:subscribers:add EMAIL=user@example.com NAME="John Doe"
task mailerlite:subscribers:delete EMAIL=user@example.com

# Group management
task mailerlite:groups:list
task mailerlite:groups:create NAME="Get Started"
task mailerlite:groups:assign GROUP_ID=xxx EMAIL=user@example.com

# Webhook server (development)
task mailerlite:server PORT=8086 GROUP_ID=xxx

# GitHub releases info
task mailerlite:releases:latest
task mailerlite:releases:urls
```

## Current Configuration

| Setting | Value |
|---------|-------|
| Get Started Group ID | `173641063022462964` |
| Newsletter Group ID | `173092547261891750` |
| Webhook Server Port | `8086` |
| Web3Forms Access Key | See `layouts/contact/list.html` |

## Sources

- [Web3Forms Webhooks](https://docs.web3forms.com/getting-started/pro-features/webhooks)
- [MailerLite Webhooks](https://www.mailerlite.com/help/webhooks)
- [Zapier MailerLite Integration](https://zapier.com/apps/mailerlite/integrations/webhook)
