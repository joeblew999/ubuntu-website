---
title: "google-mcp-server"
import_path: www.ubuntusoftware.net/pkg/google-mcp-server
repo_url: https://github.com/joeblew999/google-mcp-server
description: "Google Workspace連携用のMCPサーバー。Gmail、カレンダー、ドライブ、ドキュメント、スプレッドシート、スライドに対応。"
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

## 概要

[ngs/google-mcp-server](https://github.com/ngs/google-mcp-server) のフォークです。Claude やその他の AI アシスタントに Google Workspace サービスへのアクセスを提供する MCP（Model Context Protocol）サーバーです。

## 機能

- **Gmail** - メールの閲覧、送信、管理
- **Calendar** - 予定の表示と作成
- **Drive** - ファイルの閲覧、アップロード、ダウンロード
- **Docs** - ドキュメントの閲覧と作成
- **Sheets** - スプレッドシートの閲覧と更新
- **Slides** - プレゼンテーションの作成と管理

## インストール

```bash
go install github.com/joeblew999/google-mcp-server@latest
```

## 設定

Claude Code の MCP 設定に以下を追加してください：

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

## コントリビューション

これは Ubuntu Software が管理するフォークです。プルリクエストによるコントリビューションを歓迎します。
