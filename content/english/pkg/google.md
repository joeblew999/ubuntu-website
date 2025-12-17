---
title: google
import_path: www.ubuntusoftware.net/pkg/google
repo_url: https://github.com/joeblew999/ubuntu-website
description: Unified Google Workspace CLI. Gmail, Calendar, Drive, Docs, Sheets operations via API or browser automation.
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/google
license: MIT
author: Gerard Webb
created_at: 2025-12-17T00:00:00Z
updated_at: 2025-12-17T00:00:00Z
has_binary: true
binary_name: google
taskfile_path: taskfiles/Taskfile.google.yml
process:
  command: task google:gmail:server
  port: 8087
  health_path: /health
  disabled: true
  namespace: servers
---

## Features

- **Gmail** - Send, list, search emails via API or browser
- **Calendar** - Create events, list schedule, open calendar
- **Drive** - List, search, upload, download files
- **Docs** - Read and create documents
- **Sheets** - Read and update spreadsheets
- **Slides** - Create presentations

## CLI Usage

```bash
# Gmail
google gmail list                    # List recent emails
google gmail send --to=x --subject=y --body=z
google gmail server                  # Start webhook server (port 8087)

# Calendar
google calendar today                # Show today's events
google calendar create --title="Meeting" --start="tomorrow 2pm"
google calendar server               # Start webhook server (port 8088)

# Drive
google drive list                    # List files
google drive search "report"         # Search files
```

## Taskfile Usage

```bash
task google:gmail:list
task google:gmail:send TO=x SUBJECT=y BODY=z
task google:gmail:server

task google:calendar:today
task google:calendar:server
```
