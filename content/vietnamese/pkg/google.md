---
title: "google"
import_path: www.ubuntusoftware.net/pkg/google
repo_url: https://github.com/joeblew999/ubuntu-website
description: "CLI Google Workspace hợp nhất. Các thao tác Gmail, Calendar, Drive, Docs, Sheets qua API hoặc tự động hóa trình duyệt."
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/google
license: MIT
author: Gerard Webb
created_at: 2025-12-17T00:00:00Z
updated_at: 2025-12-17T00:00:00Z
has_binary: true
binary_name: google
tasks_path: .mise/tasks/google.toml
process:
  command: mise run google:gmail server
  port: 8087
  health_path: /health
  disabled: true
  namespace: servers
---

## Tính năng

- **Gmail** - Gửi, liệt kê, tìm kiếm email qua API hoặc trình duyệt
- **Calendar** - Tạo sự kiện, liệt kê lịch trình, mở lịch
- **Drive** - Liệt kê, tìm kiếm, tải lên, tải xuống tệp
- **Docs** - Đọc và tạo tài liệu
- **Sheets** - Đọc và cập nhật bảng tính
- **Slides** - Tạo bài thuyết trình

## Sử dụng CLI

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

## Sử dụng Taskfile

```bash
mise run google:gmail list
mise run google:gmail send --to x --subject y --body z
mise run google:gmail server

mise run google:calendar today
mise run google:calendar server
```
