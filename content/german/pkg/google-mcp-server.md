---
title: "google-mcp-server"
import_path: www.ubuntusoftware.net/pkg/google-mcp-server
repo_url: https://github.com/joeblew999/google-mcp-server
description: "MCP-Server für die Google-Workspace-Integration. Gmail, Kalender, Drive, Docs, Sheets und Slides."
version: v0.1.0
documentation_url: https://github.com/joeblew999/google-mcp-server#readme
license: MIT
author: Gerard Webb
created_at: 2024-12-16T00:00:00Z
updated_at: 2024-12-16T00:00:00Z
has_binary: true
binary_name: google-mcp-server
tasks_path: .mise/tasks/google-mcp.toml
---

## Über

Fork von [ngs/google-mcp-server](https://github.com/ngs/google-mcp-server) – ein MCP-Server (Model Context Protocol), der Claude und anderen KI-Assistenten Zugriff auf Google-Workspace-Dienste bietet.

## Funktionen

- **Gmail** – E-Mails lesen, senden und verwalten
- **Calendar** – Termine anzeigen und erstellen
- **Drive** – Dateien durchsuchen, hochladen und herunterladen
- **Docs** – Dokumente lesen und erstellen
- **Sheets** – Tabellen lesen und aktualisieren
- **Slides** – Präsentationen erstellen und verwalten

## Installation

```bash
go install github.com/joeblew999/google-mcp-server@latest
```

## Konfiguration

Fügen Sie Folgendes zu Ihren Claude Code MCP-Einstellungen hinzu:

```json
{
  "mcpServers": {
    "google": {
      "command": "google-mcp-server",
      "args": ["--accounts-dir", "~/.google-mcp-accounts"]
    }
  }
}
```

## Mitwirken

Dies ist ein von Ubuntu Software gepflegter Fork. Beiträge über Pull Requests sind willkommen.
