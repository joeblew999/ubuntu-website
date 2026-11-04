---
title: "google"
import_path: www.ubuntusoftware.net/pkg/google
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Einheitliche Google Workspace CLI. Operationen für Gmail, Calendar, Drive, Docs, Sheets per API oder Browser-Automatisierung."
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
  command: mise run google:gmail server
  port: 8087
  health_path: /health
  disabled: true
  namespace: servers
---

## Funktionen

- **Gmail** - E-Mails per API oder Browser senden, auflisten und durchsuchen
- **Calendar** - Termine erstellen, Zeitplan auflisten, Kalender öffnen
- **Drive** - Dateien auflisten, durchsuchen, hochladen und herunterladen
- **Docs** - Dokumente lesen und erstellen
- **Sheets** - Tabellen lesen und aktualisieren
- **Slides** - Präsentationen erstellen

## CLI-Nutzung

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

## Taskfile-Nutzung

```bash
mise run google:gmail list
mise run google:gmail send --to x --subject y --body z
mise run google:gmail server

mise run google:calendar today
mise run google:calendar server
```
