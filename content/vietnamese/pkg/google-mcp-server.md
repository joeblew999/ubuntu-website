---
title: máy chủ google-mcp
import_path: www.ubuntusoftware.net/pkg/máy chủ google-mcp
repo_url: https://github.com/joeblew999/máy chủ google-mcp
description: MCP server cho tích hợp Google Workspace. Gmail, Lịch, Drive, Docs, Sheets và Slides.
version: v0.1.0
documentation_url: https://github.com/joeblew999/máy chủ google-mcp#readme
license: MIT
author: Gerard Webb
created_at: 2024-12-16T00:00:00Z
updated_at: 2024-12-16T00:00:00Z
---

## Về

Fork of [ngs/google-mcp-server](https://github.com/ngs/google-mcp-server) - an MCP (Model Context Protocol) server that provides Claude and other AI assistants with access to Google Workspace services.

## Tính năng

- **Gmail** - Đọc, gửi và quản lý email
- **Lịch** - Xem và tạo sự kiện
- **Drive** - Duyệt, tải lên và tải xuống tệp tin
- **Tài liệu** - Đọc và tạo tài liệu
- **Bảng tính** - Đọc và cập nhật bảng tính
- **Slides** - Tạo và quản lý bài thuyết trình

## Cài đặt

```bash
go install github.com/joeblew999/google-mcp-server@latest
```

## Cấu hình

Thêm vào cài đặt Claude Code MCP của bạn:

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

Đây là một nhánh được duy trì bởi Ubuntu Software. Chúng tôi hoan nghênh các đóng góp thông qua yêu cầu kéo (pull requests).
