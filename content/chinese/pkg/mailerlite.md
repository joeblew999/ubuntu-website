---
title: "mailerlite"
import_path: www.ubuntusoftware.net/pkg/mailerlite
repo_url: https://github.com/joeblew999/ubuntu-website
description: "MailerLite API 的 Go 客户端库和 CLI。管理订阅者、群组和电子邮件营销活动。"
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

## 功能特性

- **订阅者管理** - 添加、更新、列出和删除订阅者
- **群组** - 创建群组并管理群组成员
- **表单与自动化** - 列出表单和自动化工作流
- **Webhooks** - 创建和管理 Webhook 集成
- **Web3Forms 集成** - 用于处理表单提交的 Webhook 服务器

## CLI 用法

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

## 库用法

```go
import "www.ubuntusoftware.net/pkg/mailerlite"

client := mailerlite.NewClient(apiKey)
subscriber, err := client.AddSubscriber(ctx, "user@example.com", nil)
```
