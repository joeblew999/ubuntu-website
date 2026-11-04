---
title: "mailerlite"
import_path: www.ubuntusoftware.net/pkg/mailerlite
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Thư viện client Go và CLI cho API MailerLite. Quản lý người đăng ký, nhóm và chiến dịch email."
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

## Tính năng

- **Quản lý người đăng ký** - Thêm, cập nhật, liệt kê và xóa người đăng ký
- **Nhóm** - Tạo nhóm và quản lý thành viên nhóm
- **Biểu mẫu & Tự động hóa** - Liệt kê biểu mẫu và quy trình tự động hóa
- **Webhook** - Tạo và quản lý các tích hợp webhook
- **Tích hợp Web3Forms** - Máy chủ webhook để nhận dữ liệu gửi từ biểu mẫu

## Sử dụng CLI

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

## Sử dụng thư viện

```go
import "www.ubuntusoftware.net/pkg/mailerlite"

client := mailerlite.NewClient(apiKey)
subscriber, err := client.AddSubscriber(ctx, "user@example.com", nil)
```
