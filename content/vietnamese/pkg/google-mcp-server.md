---
title: "google-mcp-server"
import_path: www.ubuntusoftware.net/pkg/google-mcp-server
repo_url: https://github.com/joeblew999/google-mcp-server
description: "Máy chủ MCP để tích hợp Google Workspace. Gmail, Calendar, Drive, Docs, Sheets và Slides."
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

## Giới thiệu

Bản fork của [ngs/google-mcp-server](https://github.com/ngs/google-mcp-server) - một máy chủ MCP (Model Context Protocol) giúp Claude và các trợ lý AI khác truy cập các dịch vụ Google Workspace.

## Tính năng

- **Gmail** - Đọc, gửi và quản lý email
- **Calendar** - Xem và tạo sự kiện
- **Drive** - Duyệt, tải lên và tải xuống tệp
- **Docs** - Đọc và tạo tài liệu
- **Sheets** - Đọc và cập nhật bảng tính
- **Slides** - Tạo và quản lý bản trình bày

## Cài đặt

```bash
go install github.com/joeblew999/google-mcp-server@latest
```

## Cấu hình

Thêm vào phần cài đặt MCP của Claude Code:

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

## Đóng góp

Đây là bản fork do Ubuntu Software duy trì. Chúng tôi hoan nghênh mọi đóng góp qua pull request.
