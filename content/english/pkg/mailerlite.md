---
title: mailerlite
import_path: www.ubuntusoftware.net/pkg/mailerlite
repo_url: https://github.com/joeblew999/ubuntu-website
description: Go client library and CLI for the MailerLite API. Manage subscribers, groups, and email campaigns.
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/mailerlite
license: MIT
author: Gerard Webb
created_at: 2024-12-16T00:00:00Z
updated_at: 2025-12-16T10:23:36.694709+07:00
has_binary: true
---

## Features

- **Subscriber Management** - Add, update, list, and delete subscribers
- **Groups** - Create groups and manage group membership
- **Forms & Automations** - List forms and automation workflows
- **Webhooks** - Create and manage webhook integrations
- **Web3Forms Integration** - Webhook server for form submissions

## CLI Usage

```bash
# List subscribers
mailerlite subscribers list

# Add a subscriber
mailerlite subscribers add user@example.com "John Doe"

# Show account stats
mailerlite stats

# Start webhook server for Web3Forms
mailerlite server
```

## Library Usage

```go
import "www.ubuntusoftware.net/pkg/mailerlite"

client := mailerlite.NewClient(apiKey)
subscriber, err := client.AddSubscriber(ctx, "user@example.com", nil)
```
