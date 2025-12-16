# Email Infrastructure Planning

## Current State

**SMTP2GO** configured manually in Gmail "Send mail as" settings.
- Works but painful setup
- Free tier: 1,000 emails/month, 200/day

### SMTP2GO Credentials & URLs

| Item | URL/Value |
|------|-----------|
| Dashboard | https://app.smtp2go.com |
| SMTP Users | https://app-us.smtp2go.com/sending/smtp_users/ |
| Support | ticket@smtp2go.com |
| SMTP Host | `mail.smtp2go.com` |
| SMTP Port | `587` (TLS) or `465` (SSL) |
| Username | `ubuntusoftware.net` |
| Password | See `.env` → `SMTP2GO_PASSWORD` |

**To get credentials later:**
1. Login: https://app.smtp2go.com
2. Go to: Sending → SMTP Users
3. Click on `ubuntusoftware.net` to edit/view password

## Goals

1. **Better free tier** - more emails/month for us and users
2. **Gmail web UI support** - users can send from Gmail with custom domain
3. **Automated setup** - use Playwright to configure Gmail "Send as" for users
4. **Programmatic sending** - API for automated emails from our platform

## Service Comparison

| Service | Free/Month | Daily Limit | SMTP Relay | API | Notes |
|---------|-----------|-------------|------------|-----|-------|
| **Brevo** | 9,000 | 300/day | ✅ | ✅ | ⭐ Best free tier |
| **Mailjet** | 6,000 | 200/day | ✅ | ✅ | Good option |
| **Resend** | 3,000 | 100/day | ✅ | ✅ | Developer-friendly |
| **SMTP2GO** (current) | 1,000 | 200/day | ✅ | ✅ | Working now |
| **MailerSend** | 500 | None | ✅ | ✅ | Reduced from 3k (Oct 2025) |
| **SendGrid** | 100 | - | ✅ | ✅ | ❌ Expires after 60 days |

**Recommendation:** Switch to **Brevo** (9x more free emails than SMTP2GO)

## Proposed Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     USER ONBOARDING FLOW                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. User visits walkthrough page on ubuntusoftware.net          │
│                          ↓                                      │
│  2. Signs up for Brevo → gets SMTP credentials                  │
│                          ↓                                      │
│  3. Clicks "Auto-configure my Gmail" button                     │
│                          ↓                                      │
│  4. Playwright automation:                                      │
│     - Opens Gmail Settings → Accounts → Send mail as            │
│     - Clicks "Add another email address"                        │
│     - Fills name, email, SMTP settings                          │
│     - User already logged into Gmail                            │
│                          ↓                                      │
│  5. User receives verification email, clicks link               │
│                          ↓                                      │
│  6. Done! Can send from custom domain in Gmail UI               │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Implementation Tasks

### Phase 1: Switch to Brevo (for ourselves first)
- [ ] Create Brevo account
- [ ] Verify sender domain (ubuntusoftware.net)
- [ ] Get SMTP credentials
- [ ] Test manually in Gmail "Send as"
- [ ] Update CLAUDE.md with new settings

### Phase 2: Playwright Gmail Automation
- [ ] Create `cmd/gmail-setup/main.go` or extend `cmd/gmail/`
- [ ] Implement Gmail "Send as" configuration flow
- [ ] Handle edge cases (CAPTCHAs, security prompts)
- [ ] Add `task gmail:configure-sendas` command
- [ ] Test with different Gmail account states

### Phase 3: User Walkthrough Page
- [ ] Create `/content/english/guides/email-setup.md`
- [ ] Step-by-step Brevo signup instructions
- [ ] "Auto-configure" button (triggers Playwright via local tool?)
- [ ] Manual fallback instructions with screenshots
- [ ] Translate to de/zh/ja

### Phase 4: Programmatic API Integration
- [ ] Add Brevo API integration to `internal/gmail/` or new `internal/mailer/`
- [ ] Support both SMTP relay and direct API
- [ ] Taskfile commands for API sending

## Brevo SMTP Settings (to research/confirm)

```
Host: smtp-relay.brevo.com (or smtp-relay.sendinblue.com)
Port: 587 (TLS) or 465 (SSL)
Username: (from Brevo dashboard)
Password: (API key or SMTP key)
```

## Gmail "Send as" Automation Details

**Target URL:** `https://mail.google.com/mail/u/0/#settings/accounts`

**Form fields to fill:**
1. Name (display name for "From")
2. Email address (user's custom domain email)
3. SMTP Server hostname
4. Port (587 recommended)
5. Username
6. Password
7. Security (TLS/SSL)

**Verification:** Gmail sends a confirmation email to the address. User must click link or enter code.

## Questions to Resolve

1. **Brevo vs others?** - Brevo has best free tier, but verify SMTP works smoothly with Gmail
2. **Playwright delivery?** - How do users run Playwright automation? Local CLI tool? Browser extension?
3. **Security** - How to handle user's SMTP credentials safely?
4. **Multi-tenant** - Each user has their own Brevo account, or shared under our account?

## References

- Brevo SMTP docs: https://developers.brevo.com/docs/send-a-transactional-email
- Gmail API sendAs: https://developers.google.com/workspace/gmail/api/reference/rest/v1/users.settings.sendAs
- Current Playwright tool: `cmd/playwright/main.go`
- Current Gmail tool: `cmd/gmail/main.go`
