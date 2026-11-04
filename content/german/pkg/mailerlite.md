---
title: "mailerlite"
import_path: www.ubuntusoftware.net/pkg/mailerlite
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Go-Clientbibliothek und CLI für die MailerLite-API. Verwalten Sie Abonnenten, Gruppen und E-Mail-Kampagnen."
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/mailerlite
license: MIT
author: Gerard Webb
created_at: 2024-12-16T00:00:00Z
updated_at: 2025-12-16T10:23:36.694709+07:00
has_binary: true
binary_name: mailerlite
tasks_path: .mise/tasks/mailerlite.toml
process:
  command: mise run mailerlite:server
  port: 8086
  health_path: /health
  disabled: true
  namespace: servers
---

## Funktionen

- **Abonnentenverwaltung** - Abonnenten hinzufügen, aktualisieren, auflisten und löschen
- **Gruppen** - Gruppen erstellen und Gruppenmitgliedschaften verwalten
- **Formulare & Automatisierungen** - Formulare und Automatisierungs-Workflows auflisten
- **Webhooks** - Webhook-Integrationen erstellen und verwalten
- **Web3Forms-Integration** - Webhook-Server für Formularübermittlungen

## CLI-Nutzung

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

## Bibliotheksnutzung

```go
import "www.ubuntusoftware.net/pkg/mailerlite"

client := mailerlite.NewClient(apiKey)
subscriber, err := client.AddSubscriber(ctx, "user@example.com", nil)
```
