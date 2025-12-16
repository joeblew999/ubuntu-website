---
title: mailerLite
import_path: www.ubuntusoftware.net/pkg/mailerLite
repo_url: https://github.com/joeblew999/ubuntu-website
description: Thư viện khách hàng và giao diện dòng lệnh (CLI) cho API MailerLite. Quản lý người đăng ký, nhóm và chiến dịch email.
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/mailerLite
license: MIT
author: Gerard Webb
created_at: 2024-12-16T00:00:00Z
updated_at: 2025-12-16T10:23:36.694709+07:00
has_binary: true
---

## Tính năng

- **Quản lý người đăng ký** - Thêm, cập nhật, liệt kê và xóa người đăng ký
- **Nhóm** - Tạo nhóm và quản lý thành viên trong nhóm
- **Biểu mẫu & Tự động hóa** - Danh sách các biểu mẫu và quy trình làm việc tự động hóa
- **Webhooks** - Tạo và quản lý tích hợp webhook
- **Tích hợp Web3Forms** - Máy chủ webhook cho việc gửi biểu mẫu

## Cách sử dụng CLI

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
