---
title: "技术"
meta_title: "技术栈 | Ubuntu 软件"
description: "基于 Go 和 NATS JetStream 构建——我们为性能、可靠性和简洁性而做出的技术选择。"
image: "/images/spatial.svg"
draft: false
---

## 我们的技术栈

我们选用的技术注重性能、可靠性和长期可维护性。

---

## Go

**为什么选择 Go：**

- **性能** — 编译型、静态类型，运行时开销极小
- **简洁** — 做事方式唯一，默认即具可读性
- **并发** — 语言内置 Goroutine 和 Channel
- **部署** — 单一二进制文件，无依赖，支持交叉编译
- **生态** — 强大的标准库，出色的工具链

Go 是我们贯穿整个技术栈的主要语言——从后端服务到 CLI 工具，再到机器人集成。

---

## NATS JetStream

**为什么选择 NATS JetStream：**

- **持久化** — 支持重放的持久消息存储
- **精确一次投递** — 可靠的消息处理
- **轻量级** — 单一二进制文件，资源占用极少
- **可扩展** — 内置集群和水平扩展能力
- **实时** — 发布/订阅延迟低于毫秒级

NATS JetStream 驱动着我们的事件驱动架构——连接设计工具、仿真、数字孪生和物理设备。

---

## 架构原则

| 原则 | 实现方式 |
|-----------|----------------|
| **离线优先** | 本地优先数据，联网时同步 |
| **事件驱动** | 所有服务间通信均采用 NATS JetStream |
| **开放标准** | STEP、IFC，不使用专有格式 |
| **硬件无关** | 通过抽象层实现可移植性 |
| **自主可控** | 随处部署——云端、本地、物理隔离环境 |

---

## 技术领域

### 机器人

我们的机器人技术栈基于 Viam RDK 构建——这是一个开源机器人开发套件，提供硬件抽象、运动规划和计算机视觉功能。

[机器人技术栈 →](/technology/robotics/)

### Linux 与跨平台

数十年的 Linux 专业经验。为 Windows、Mac、Linux、iOS 和 Android 开发跨平台应用。在多年使用 Qt、Flutter 和 Electron 之后，打造了自有的 GUI 框架。

[Linux 与跨平台 →](/technology/linux/)

### 安全与合规

自主可控的架构简化了 SOC 2、FedRAMP、HIPAA 和 ISO 27001 合规工作。支持物理隔离部署，无回传通信，完全掌控数据。

[安全与合规 →](/technology/security/)

---

## 部署平台

我们可部署到任何平台：

| 平台 | 用例 |
|----------|----------|
| **Linux** | 服务器、桌面、嵌入式 |
| **OpenBSD** | 安全关键型系统 |
| **Windows** | 企业桌面 |
| **macOS** | 设计与开发 |
| **iOS/Android** | 移动应用 |
| **Docker/Kubernetes** | 容器编排 |

---

## 开源基础

我们基于开源构建，并为开源做出贡献：

| 组件 | 技术 |
|-----------|------------|
| 语言 | Go |
| 消息传递 | NATS JetStream |
| 协作 | Automerge (CRDT) |
| 3D 格式 | STEP、IFC |
| 机器人 | Viam RDK |
| 视觉 | Intel RealSense、YOLOv8 |

---

## 了解更多

- [机器人技术栈 →](/technology/robotics/)
- [Linux 与跨平台 →](/technology/linux/)
- [安全与合规 →](/technology/security/)
- [空间平台 →](/platform/spatial/)
- [基金会 →](/platform/foundation/)
- [联系我们 →](/contact/)
