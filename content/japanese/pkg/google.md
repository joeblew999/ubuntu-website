---
title: "Google"
import_path: www.ubuntusoftware.net/pkg/google
repo_url: https://github.com/joeblew999/ubuntu-website
description: "Google Workspace統合CLI。API またはブラウザ自動化によるGmail、Calendar、Drive、Docs、Sheetsの操作。"
version: v0.1.0
documentation_url: https://pkg.go.dev/www.ubuntusoftware.net/pkg/google
license: MIT
author: Gerard Webb
created_at: 2025-12-17T00:00:00Z
updated_at: 2025-12-17T00:00:00Z
has_binary: true
binary_name: google
taskfile_path: taskfiles/Taskfile.google.yml
process:
  command: task google:gmail:server
  port: 8087
  health_path: /health
  disabled: true
  namespace: servers
---

## 機能

- **Gmail** - API またはブラウザ経由でメールの送信、一覧表示、検索
- **Calendar** - イベントの作成、スケジュールの一覧表示、カレンダーを開く
- **Drive** - ファイルの一覧表示、検索、アップロード、ダウンロード
- **Docs** - ドキュメントの読み取りと作成
- **Sheets** - スプレッドシートの読み取りと更新
- **Slides** - プレゼンテーションの作成

## CLI の使い方

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

## Taskfile の使い方

```bash
task google:gmail:list
task google:gmail:send TO=x SUBJECT=y BODY=z
task google:gmail:server

task google:calendar:today
task google:calendar:server
```
