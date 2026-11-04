---
title: "通过 MCP 接入 Google Workspace：零配置 AI 集成"
meta_title: "通过 MCP 接入 Google Workspace：零配置 AI 集成 | Ubuntu Software"
description: "将任意 AI（云端、本地或您自己的）通过一条命令连接到 Gmail、Calendar、Drive、Sheets、Docs 和 Slides。MCP 意味着没有供应商锁定。只需提供您的电子邮箱。"
date: 2025-12-15T10:00:00Z
image: "/images/blog/google-mcp-integration.svg"
categories: ["Publish", "AI"]
author: "Gerard Webb"
tags: ["google", "mcp", "automation", "gmail", "calendar", "drive", "ai", "local-ai", "ollama", "vendor-lock-in", "spatial", "sensors"]
draft: true
---

如果把您的 AI 助手连接到整个 Google Workspace，只需要一条命令呢？

不是"配置这些凭据，启用这些 API，设置这些权限，然后重启"。只需提供您的电子邮件地址，一切就能正常运行。

这正是我们打造的。

## AI 集成的问题

每一个 AI 集成的故事都遵循同样令人疲惫的模式：

1. 在某个开发者控制台中创建 API 凭据
2. 手动启用十几个 API
3. 配置 OAuth 同意屏幕
4. 设置重定向 URI
5. 把密钥存放在某个安全的地方
6. 编写连接配置
7. 重启您的工具
8. 排查为什么不起作用
9. 发现您漏掉了第 3b 步
10. 从头再来

等您完成时，花在设置上的时间比接下来一个月能节省的时间还多。而且您仍然不确定它是否正常工作。

## WellKnown + MCP：一种不同的方法

[模型上下文协议（MCP）]({{< relref "/platform/foundation" >}}) 为 AI 助手访问外部工具和数据提供了一种标准方式。WellKnown 项目更进一步：通过自动化设置，在幕后处理所有复杂性。

以下是完整的设置流程：

```bash
mise run google-mcp:setup
```

就是这样。系统会：

1. **安装固定版本的 MCP 服务器**（如果尚未安装）
2. **将您的 Google OAuth 客户端存储在操作系统钥匙串中**，绝不存放在明文文件里
3. **在浏览器中打开引导式身份验证**
4. **通过一次真实的 Calendar 调用验证一切正常**

您的 AI 助手配置随项目一起提供（`.mcp.json`），因此无需任何连接配置。

无需手动启用 API。无需编辑 JSON 文件。无需反复重启。

## 任何 AI，同一套集成

MCP 与专有集成的根本区别在于：**它适用于任何 AI**。

| AI 类型 | 示例 | 相同的 MCP 集成 |
|---------|----------|---------------------|
| **云端 AI** | Claude, GPT-4, Gemini | ✓ |
| **本地 AI** | Ollama, LM Studio, llama.cpp | ✓ |
| **混合** | 云端推理 + 本地执行 | ✓ |
| **专用** | 空间推理、特定领域模型 | ✓ |
| **您自己的** | 自定义模型、微调部署 | ✓ |

MCP 服务器并不关心是哪个 AI 在调用它。云端 Claude 查询您的日历？可以。在笔记本电脑上运行的本地 Llama？同一套集成。在您数据中心运行的公司私有模型？设置完全相同。

这与供应商锁定恰恰相反。只需构建一次 Google Workspace 集成，它就能与您选择的任何 AI 配合使用——今天、明天，或五年后出现全新模型的时候。

**更换 AI 提供商，无需改动您的集成。** 当您尝试不同的模型、升级能力，或根据隐私要求在云端和本地处理之间切换时，MCP 层始终保持不变。

### 专用 AI：空间推理

我们的[空间平台]({{< relref "/platform/spatial" >}}) 包含具备通用云端模型所没有的能力的本地 AI：

- **空间推理** - 理解三维关系、几何约束、装配顺序
- **传感器集成** - 处理来自物联网设备、摄像头和工业设备的实时数据
- **领域知识** - CAD 格式、制造公差、施工规范

当这种专用本地 AI 通过 MCP 连接到您的 Google Workspace 时，您将获得强大的组合能力：

> "从 Drive 中提取传感器数据，分析热分布模式，并起草一封发给运营团队的维护警报"

> "根据 CAD 装配体创建一份演示文稿，并根据检测电子表格标注超出公差的部件"

与 Claude 或 GPT 配合使用的同一套 MCP 集成，也适用于具备空间感知能力的本地模型——只不过现在您的 AI 能理解几何，而不仅仅是文本。

### 您将获得什么

连接后，您的 AI 助手可以原生访问：

| 服务 | 功能 |
|---------|-------------|
| **Gmail** | 阅读、搜索、发送、起草电子邮件 |
| **Calendar** | 查看活动、创建会议、检查空闲时间 |
| **Drive** | 列出文件、搜索、下载、上传 |
| **Sheets** | 读取数据、更新单元格、查询区域 |
| **Docs** | 读取内容、创建文档、编辑文本 |
| **Slides** | 创建演示文稿、添加幻灯片、导出 PDF |

一切都通过自然对话完成：

> "给我看看上周关于项目提案的电子邮件"

> "我明天有哪些会议？"

> "创建一份 Google 文档，总结附件中的 PDF"

> "在预算电子表格中添加一行，填入这些数字"

## 多账户支持

对企业来说，有趣的地方在这里：**您可以连接多个 Google 账户**。

个人 Gmail。工作用的 Google Workspace。客户账户。每个账户独立认证，您的 AI 助手可以在正确的上下文中访问每一个账户。

```
mise run google-mcp:auth
# Authenticate first account

mise run google-mcp:auth
# Authenticate second account
# The system tracks both
```

然后，您的 AI 助手就可以跨账户操作：

> "检查我的工作日历和个人日历，看看下周二有没有冲突"

> "把那封客户邮件转发到我的工作账户"

> "列出我个人和企业 Drive 中的最近文件"

## 与发布的连接

这与 [WellKnown 的单一来源发布]({{< relref "/blog/one-source-every-screen" >}}) 直接集成。您的内容管道将 Google Workspace 作为一等发布目标：

**Markdown 源文件 → WellKnown 转换 → Google Docs**

用 Markdown 写作，发布到 Docs 时保留正确的格式。无需复制粘贴。无需手动重新排版。生成您网站的同一份源文件，也可以填充您的 Google Drive。

**电子表格数据 → 发布管道**

将 Sheets 中的数据直接拉取到您的内容中。产品目录、价格表、团队通讯录——始终最新，始终同步。

**Calendar → 自动化排程**

安排在未来日期发布的博客文章可以触发日历活动、提醒，甚至可以向相关方起草电子邮件。

## 自主掌控，而非依赖

与典型云端集成的关键区别在于：**您的凭据和数据保留在您自己的机器上**。

MCP 服务器在本地运行。身份验证令牌存储在您的主目录中，OAuth 客户端密钥则存储在操作系统的钥匙串中。没有任何中间云服务能看到您的 Google 数据。没有任何第三方能访问您的电子邮件或文件。

您使用的是 Google 的服务，但通过的是您自己掌控的基础设施。如果您决定断开连接，删除本地凭据即可。无需联系供应商。无需申请数据导出。

这与我们在[自主掌控的电子邮件基础设施]({{< relref "/blog/self-sovereign-email-ai-automation" >}}) 中讨论的原则一致：使用您需要的服务，但保持对集成层的掌控。

## AI 自动化级别的应用

还记得[自动化级别]({{< relref "/blog/self-sovereign-email-ai-automation" >}}) 吗？它们在这里同样适用：

**第 1 级 - 模式识别**
> "您有 3 封来自团队的未读邮件。其中一封提到了'紧急'。要我为您总结一下吗？"

**第 2 级 - 草稿生成**
> "根据 Sarah 关于提案截止日期的邮件，我起草了一封回复，确认周三交付。要审阅一下吗？"

**第 3 级 - 自主操作**
> "每周报告电子表格已更新。我已将摘要导出到相关方文档，并将分发邮件安排在周一上午 9 点发送。"

**第 4 级 - 主动智能**
> "我注意到您明天有一个客户会议，但还没有准备文档。要我根据最近三个邮件往来和 Drive 中的项目时间表创建一份吗？"

当 AI 能够纵览您的整个 Google Workspace，而不仅仅是单个服务时，每个级别都会变得更加强大。

## 技术架构

想了解底层原理的读者请看：

```
┌─────────────────────────────────────────────────────────────┐
│                    Your Local Machine                        │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐  │
│  │   Claude     │───▶│  MCP Server  │───▶│   Google     │  │
│  │   Code       │◀───│  (local)     │◀───│   APIs       │  │
│  └──────────────┘    └──────────────┘    └──────────────┘  │
│         │                   │                              │
│         │                   ▼                              │
│         │           ┌──────────────┐                       │
│         │           │   Tokens     │                       │
│         │           │   (~/.google-│                       │
│         │           │   mcp-...)   │                       │
│         │           └──────────────┘                       │
│         │                                                  │
│         ▼                                                  │
│  ┌──────────────────────────────────────────────────────┐ │
│  │              .mcp.json (checked into the project)     │ │
│  │   "google": {                                         │ │
│  │     "command": "mise",                                │ │
│  │     "args": ["exec", "--", "fnox", "exec", "--",      │ │
│  │              "google-mcp-server"]                     │ │
│  │   }                                                   │ │
│  └──────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

MCP 服务器充当您的 AI 助手与 Google API 之间的桥梁。它负责处理：

- **OAuth 令牌刷新**（自动、透明）
- **请求批处理**（高效使用 API）
- **响应格式化**（为 AI 处理而结构化）
- **多账户路由**（将请求定向到正确的账户）

全部在本地运行。全部由您掌控。

## 快速上手

### 快速开始（大多数用户）

```bash
# From your project directory
mise run google-mcp:setup
```

按照提示操作。重启您的 AI 助手。完成。

### 分步操作（如果您想要更多控制）

```bash
# 1. Install the MCP server (pinned release)
mise install

# 2. Store your Google OAuth client in the OS keychain
mise run google-mcp:credentials

# 3. Authenticate with Google (opens the browser, then checks)
mise run google-mcp:auth

# 4. Verify everything works
mise run google-mcp:status
```

### 检查已连接的内容

```bash
mise run google-mcp:accounts
# Shows all authenticated Google accounts

mise run google-mcp:status
# Full status: binary, keychain secrets, accounts, configuration
```

### 撤销访问权限

```bash
mise run google-mcp:reset --yes
# Removes the signed-in account tokens from this machine
# Revoke the app itself in your Google Account's third-party access settings
```

## 不止于 Google

同样的模式也适用于其他服务。MCP 提供协议；WellKnown 提供自动化设置。我们正在构建以下集成：

- **GitHub** - Issues、PR、代码搜索
- **Notion** - 页面、数据库、工作区
- **Slack** - 频道、消息、话题
- **Linear** - Issues、项目、路线图

每一项都遵循同样的原则：一条命令即可连接，完整的 AI 访问能力，本地掌控。

## 会带来哪些改变

当您的 AI 助手能够无缝访问您的 Google Workspace 时，工作流程会发生根本性的转变：

**之前**：您在各种工具之间来回切换。查看邮件。打开日历。搜索 Drive。把数据复制给 AI。获得回复。再粘贴回去。

**之后**：您描述您的需求。AI 在各个服务之间导航，收集上下文，执行操作，并向您汇报。

工具切换带来的认知负担消失了。您思考的是结果，而不是界面。

## 感兴趣吗？

我们目前正在为此集成招募早期合作伙伴。如果您正在使用 [Claude Code](https://docs.anthropic.com/en/docs/claude-code/overview) 或其他兼容 MCP 的助手，并希望零配置访问 Google Workspace，请[联系我们]({{< relref "/contact" >}})——我们将帮助您完成连接。
