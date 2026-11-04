---
title: "google-mcp-server"
import_path: www.ubuntusoftware.net/pkg/google-mcp-server
repo_url: https://github.com/joeblew999/google-mcp-server
description: "用于 Google Workspace 集成的 MCP 服务器。支持 Gmail、Calendar、Drive、Docs、Sheets 和 Slides。"
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

## 关于

[ngs/google-mcp-server](https://github.com/ngs/google-mcp-server) 的分支——一个 MCP（Model Context Protocol，模型上下文协议）服务器，为 Claude 和其他 AI 助手提供对 Google Workspace 服务的访问。

## 功能

- **Gmail** - 阅读、发送和管理电子邮件
- **Calendar** - 查看和创建日程
- **Drive** - 浏览、上传和下载文件
- **Docs** - 阅读和创建文档
- **Sheets** - 读取和更新电子表格
- **Slides** - 创建和管理演示文稿

## 安装

```bash
go install github.com/joeblew999/google-mcp-server@latest
```

## 配置

添加到您的 Claude Code MCP 设置中：

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

## 贡献

这是由 Ubuntu Software 维护的分支。我们欢迎通过拉取请求（pull request）提交贡献。
