---
title: "mailerlite"
import_path: www.ubuntusoftware.net/pkg/mailerlite
repo_url: https://github.com/joeblew999/ubuntu-website
description: "MailerLite API用のGoクライアントライブラリとCLI。購読者、グループ、メールキャンペーンを管理します。"
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

## 機能

- **購読者管理** - 購読者の追加、更新、一覧表示、削除
- **グループ** - グループの作成とグループメンバーシップの管理
- **フォームとオートメーション** - フォームとオートメーションワークフローの一覧表示
- **Webhook** - Webhook連携の作成と管理
- **Web3Forms連携** - フォーム送信用のWebhookサーバー

## CLIの使い方

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

## ライブラリの使い方

```go
import "www.ubuntusoftware.net/pkg/mailerlite"

client := mailerlite.NewClient(apiKey)
subscriber, err := client.AddSubscriber(ctx, "user@example.com", nil)
```
