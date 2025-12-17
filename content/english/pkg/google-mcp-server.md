---
title: google-mcp-server
import_path: www.ubuntusoftware.net/pkg/google-mcp-server
repo_url: https://github.com/joeblew999/google-mcp-server
description: MCP server for Google Workspace integration. Gmail, Calendar, Drive, Docs, Sheets, and Slides.
version: v0.1.0
documentation_url: https://github.com/joeblew999/google-mcp-server#readme
license: MIT
author: Gerard Webb
created_at: 2024-12-16T00:00:00Z
updated_at: 2024-12-16T00:00:00Z
has_binary: true
binary_name: google-mcp-server
taskfile_path: taskfiles/tools/Taskfile.google-mcp.yml
---

## About

Fork of [ngs/google-mcp-server](https://github.com/ngs/google-mcp-server) - an MCP (Model Context Protocol) server that provides Claude and other AI assistants with access to Google Workspace services.

## Features

- **Gmail** - Read, send, and manage emails
- **Calendar** - View and create events
- **Drive** - Browse, upload, and download files
- **Docs** - Read and create documents
- **Sheets** - Read and update spreadsheets
- **Slides** - Create and manage presentations

## Installation

```bash
go install github.com/joeblew999/google-mcp-server@latest
```

## Configuration

Add to your Claude Code MCP settings:

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

## Contributing

This is a fork maintained by Ubuntu Software. We welcome contributions via pull requests.
