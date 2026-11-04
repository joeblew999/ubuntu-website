---
title: "谷歌"
import_path: www.ubuntusoftware.net/pkg/google
repo_url: https://github.com/joeblew999/ubuntu-website
description: "统一的 Google Workspace CLI。通过 API 或浏览器自动化执行 Gmail、Calendar、Drive、Docs、Sheets 操作。"
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

## 功能

- **Gmail** - 通过 API 或浏览器发送、列出、搜索电子邮件
- **Calendar** - 创建活动、列出日程、打开日历
- **Drive** - 列出、搜索、上传、下载文件
- **Docs** - 读取和创建文档
- **Sheets** - 读取和更新电子表格
- **Slides** - 创建演示文稿

## CLI 用法

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

## Taskfile 用法

```bash
task google:gmail:list
task google:gmail:send TO=x SUBJECT=y BODY=z
task google:gmail:server

task google:calendar:today
task google:calendar:server
```
