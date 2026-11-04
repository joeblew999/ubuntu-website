---
title: "AI 驱动的现场作业"
meta_title: "AI驱动的物流与现场作业 | Ubuntu Software"
description: "面向现场作业的实时 AI 路线优化。卡车司机、配送车队和服务团队可获得由 AI 规划的路线，并可通过语音或文字输入随时动态调整。"
image: "/images/robotics.svg"
draft: false
---

## 当三大平台融合

**Fleet + Publish + Foundation** 相结合，创造出全新的能力：AI 驱动的现场作业。在这个全面互联的系统中，人工操作员成为智能边缘节点。

手持 iPad 的卡车司机会收到由 AI 规划的收运路线。系统掌握每辆车的位置。操作员可以通过语音或文字更新路线，AI 会立即重新规划所有受影响的路线。

---

## 问题：动态世界中的静态路线

传统的现场作业在前一天晚上就已规划好：

- **路线规划**：调度员手动制定路线，或使用夜间批量优化
- **临时变更**：司机生病？街道堵塞？只能打电话、重新规划，一片混乱
- **缺乏实时可见性**：想知道“47 号卡车在哪里？”还得用无线电呼叫
- **纸质流程**：工单需要打印，完成状态需要手动重新录入
- **信息孤岛**：司机知道的情况，调度员不知道，客户更是一无所知

一旦实际情况偏离计划，整个系统就难以适应。

---

## 解决方案：AI 引导的现场作业

```
┌─────────────────────────────────────────────────────────────────────┐
│                    UBUNTU SOFTWARE UNIFIED STACK                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────────────┐  │
│  │   PUBLISH   │    │    FLEET    │    │      FOUNDATION         │  │
│  │             │    │             │    │                         │  │
│  │ Work orders │    │  Vehicles   │    │  Offline-first sync     │  │
│  │ Forms/PDFs  │◄──►│  Telemetry  │◄──►│  Automerge CRDT         │  │
│  │ Kiosks      │    │  Commands   │    │  NATS JetStream         │  │
│  │ Translation │    │  Digital    │    │  SQLite everywhere      │  │
│  │             │    │  Twins      │    │  WellKnown Gateway      │  │
│  └─────────────┘    └─────────────┘    └─────────────────────────┘  │
│         │                  │                       │                 │
│         └──────────────────┼───────────────────────┘                 │
│                            │                                         │
│                            ▼                                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                   FIELD OPERATOR (Human Edge)                │    │
│  │                                                               │    │
│  │   iPad/Phone with:                                           │    │
│  │   • Real-time route display (2D/3D)                          │    │
│  │   • AI-directed pickup notifications                         │    │
│  │   • SST/voice input for route changes                        │    │
│  │   • Offline capability (queues sync when connected)          │    │
│  │   • Position tracking (device GPS)                           │    │
│  │   • Work order completion/forms                              │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                            │                                         │
│                            ▼                                         │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                      VEHICLE (Truck)                         │    │
│  │                                                               │    │
│  │   • CAN bus telemetry (engine, fuel, diagnostics)            │    │
│  │   • GPS position (location tracking)                         │    │
│  │   • Gateway syncs to NATS JetStream                          │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

**核心洞察：** 人工操作员就是边缘计算节点。他们的设备提供定位、接收 AI 指令，并能随时更新系统，就像网络中的其他传感器一样。

---

## 示例：重新构想垃圾收运

以下具体示例展示了该架构的工作方式：

### 早晨：AI 规划一天的工作

1. 系统掌握所有收运地点、优先级和垃圾箱满溢程度（IoT 传感器）
2. AI 为每辆卡车计算最优路线：尽量缩短行驶距离，同时遵守时间窗口
3. 换班开始前，路线自动推送到司机的平板电脑
4. 每位司机都能在 2D/3D 视图中看到自己的专属路线，并获得逐向导航

### 作业期间：实时调整

**1 号司机报告：**“Main St 123 号因施工封路”

```
Driver speaks into tablet
    └─▶ SST converts to text
    └─▶ System marks location as blocked
    └─▶ AI refactors all affected routes
    └─▶ Updated routes push to affected drivers
    └─▶ Total time: seconds
```

无需电话，无需调度员介入，系统会自动处理。

**与此同时：**
- 调度仪表板实时显示所有卡车
- 客户门户显示其收运的预计到达时间
- 每次完成的作业都会自动记录
- 每项变更都会被审计，并保留完整的来源记录

### 人机协同

操作员可以随时否决 AI 的决定：

| 输入方式 | 示例 | 系统响应 |
|--------------|---------|-----------------|
| **语音（SST）** | “跳过 Oak Ave 456 号，客户不在家” | 标记跳过，调整路线 |
| **触控** | 在地图上点击位置，选择“受阻” | 从所有路线中移除 |
| **表单** | 填写服务表单并添加备注 | 记录完成情况，触发计费 |
| **手动覆盖** | 拖动路线以更改顺序 | 接受更改，重新优化后续路线 |

AI 增强人的判断，而不是取代它。

---

## 人机协同：操作员即边缘智能

真正的突破不只是 AI，而是 **AI 与人的判断协同工作**。现场操作员不是被动接收指令的人，而是系统中的智能节点。

### 全面接入的设备

每位司机的 iPad 或手机既是传感器，也是执行器：

```
┌─────────────────────────────────────────────────────────────────┐
│                    DRIVER DEVICE (iPad/Phone)                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  INPUTS (Sensing)              OUTPUTS (Acting)                  │
│  ─────────────────             ────────────────                  │
│  • GPS position                • Route display (2D/3D)           │
│  • Voice commands (SST)        • Turn-by-turn directions         │
│  • Touch interactions          • Pickup notifications            │
│  • Form submissions            • Alert sounds/haptics            │
│  • Camera (proof of service)   • Customer ETA updates            │
│  • Accelerometer (driving)     • Route change confirmations      │
│                                                                  │
│  LOCAL INTELLIGENCE                                              │
│  ──────────────────                                              │
│  • SQLite database (offline state)                               │
│  • Automerge CRDT (conflict-free sync)                           │
│  • Background location tracking                                  │
│  • Push notification handling                                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

系统知道每位司机在哪里、正在做什么，并能即时联系到他们。但更关键的是，**司机也可以向系统反馈**。

### 语音驱动的路线重构

语音转文字（SST）让司机无需动手即可更新系统：

| 司机说 | 系统执行 |
|-------------|-------------|
| “Main 123 号因施工封路” | 将该地点标记为受阻，重构所有受影响的路线 |
| “跳过下一个收运点，家里没人” | 从路线中移除，通知客户，调整预计到达时间 |
| “Oak Street 的垃圾箱满溢了” | 标记为优先，可能改派给更大的卡车 |
| “休息 15 分钟” | 调整所有后续预计到达时间，通知受影响的客户 |
| “路线完成” | 关闭路线，更新仪表板，触发计费 |

无需给调度打电话，也没有无线电通话。司机说话，AI 倾听，系统随之调整。

### 决策权层级

人的判断优先于 AI 的建议：

```
┌─────────────────────────────────────────────────────────────────┐
│                     DECISION AUTHORITY                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. DRIVER ON-SITE      ───────────────────────▶  HIGHEST      │
│      Sees reality, can override anything                         │
│                                                                  │
│   2. DISPATCHER          ───────────────────────▶  HIGH         │
│      Fleet-wide view, can reassign routes                        │
│                                                                  │
│   3. AI OPTIMIZER        ───────────────────────▶  MEDIUM       │
│      Suggests routes, adapts to changes                          │
│                                                                  │
│   4. SCHEDULED PLAN      ───────────────────────▶  LOWEST       │
│      Starting point, always subject to reality                   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

AI 固然强大，但最终决定权在亲临现场的人手中。

---

## 发布：单一来源，全渠道输出

[Publish 平台]({{< relref "/platform/publish" >}}) 从单一事实来源生成所有输出。

### 物流输出

| 源定义 | 生成的输出 |
|-------------------|-------------------|
| **收运记录** | 平板屏幕、PDF 清单、客户短信、信息亭显示 |
| **路线计划** | 司机逐向导航、调度仪表板、客户预计到达时间 |
| **服务表单** | 数字表单、可打印 PDF、OCR 扫描的纸质表单 |
| **每日汇总** | 电子邮件报告、仪表板小组件、合规导出 |

**一次定义，多处输出。** 修改源数据，所有渠道自动更新。

### 多语言现场团队

现场人员使用不同的语言？Publish 负责翻译：

- 1 号司机的平板：英语
- 2 号司机的平板：西班牙语
- 客户门户：客户的首选语言
- PDF 清单：公司标准语言

同一份数据，本地化呈现，无需重复管理内容。

### 可回流的表单

服务完成表单会回流到系统中：

```
Driver completes pickup form
    └─▶ Photo attached (proof of service)
    └─▶ Notes recorded ("Container damaged")
    └─▶ Time stamped automatically
    └─▶ GPS coordinates captured
    └─▶ Data syncs via Foundation
    └─▶ Billing system updated
    └─▶ Customer notified
    └─▶ Maintenance ticket created (if needed)
```

表单不仅用于采集数据，还会触发下游自动化。

### 数据主权

你的数据流向何处？[WellKnown Gateway]({{< relref "/platform/publish#wellknown-gateway-data-sovereignty" >}}) 确保即使在向 Salesforce、车队追踪系统或计费服务商等外部平台发布数据时，数据仍归你所有。

**核心原则：** 你向平台发布数据。数据归你所有，而不是归平台所有。

[完整了解 Publish 平台 →]({{< relref "/platform/publish" >}})

---

## Fleet 集成：车辆即数据源

**Fleet 架构**提供车辆层，让卡车成为数据丰富的数字孪生。

### 卡车掌握的信息

每辆车都通过 CAN 总线网关传输遥测数据：

| 数据点 | 来源 | 用途 |
|------------|--------|-----|
| **GPS 位置** | 车载信息模块 | 实时追踪、预计到达时间计算 |
| **发动机状态** | ECU | 健康监测、维护提醒 |
| **燃油量** | 燃油传感器 | 续航估算、加油调度 |
| **车速** | 轮速传感器 | 安全监测、路线计时 |
| **诊断信息** | OBD-II / J1939 | 预测性维护 |
| **车门状态** | 车身控制器 | 服务核验 |
| **载重** | 称重传感器 | 容量管理 |

卡车不只是运输工具，更是移动的传感器平台。

### 数字孪生同步

每辆卡车在系统中都有一个数字孪生：

```
Physical Truck                    Digital Twin
─────────────                    ────────────
Position: 34.05°, -118.24°   →   state.position
Fuel: 67%                    →   state.fuel
Engine: Running              →   state.engine
Speed: 25 mph                →   state.speed
Door: Closed                 →   state.door
Weight: 3,400 lbs            →   state.weight
```

调度员看到的是实际情况，AI 基于实际情况进行优化，客户获得准确的预计到达时间。

[Fleet 架构 →]({{< relref "/fleet" >}}) | [地面车辆 →]({{< relref "/fleet/platforms/ground" >}}) | [车辆网关 →]({{< relref "/fleet/gateway" >}})

---

## Foundation：一切以离线优先

**Foundation 平台**确保即使网络中断，系统也能正常运行。

### 为什么离线能力很重要

现场作业往往发生在网络失效的地方：

- 地下停车场
- 信号不稳定的农村地区
- 隧道和地下通道
- 活动期间人流密集的区域
- 国际边境（漫游盲区）

如果你的应用离线就卡死，你的业务也会随之停摆。这是不可接受的。

### Automerge CRDT：无冲突同步

离线时做出的更改会在网络恢复后干净地合并：

```
Driver A (online):    Marks pickup complete at 9:15am
Driver B (offline):   Marks same pickup skipped at 9:14am (tunnel)

Driver B reconnects:
    └─▶ Automerge detects conflict
    └─▶ Applies "last write wins" or custom merge rule
    └─▶ Both changes preserved in history
    └─▶ Alert triggers for dispatcher review
    └─▶ No data loss, no app crash
```

传统数据库可能会损坏或丢失数据，而 CRDT 能从容应对。

### 无处不在的 SQLite

每台设备都有完整的本地数据库：

- **司机平板**：路线、客户、表单、历史记录
- **调度仪表板**：车队状态、警报、分析
- **客户门户**：客户的服务记录、即将进行的收运

没有网络？没问题。应用所需的一切都在本地。

[Foundation 平台 →]({{< relref "/platform/foundation" >}})

---

## 全面感知：2D/3D 可视化

操作员和调度员看到的是同一个现实：

### 司机视图（平板）
- 当前路线及逐向导航
- 下一个收运点的详情和客户备注
- 所在区域内的所有其他卡车（态势感知）
- 交通与路况叠加图层

### 调度视图（仪表板）
- 地图上所有车辆的位置
- 每位司机当前的路线进度
- 突出显示的警报和异常
- 历史轨迹和审计记录

### 客户视图（门户）
- 服务的预计到达时间
- 司机即将到达的通知
- 服务完成确认

所有视图都源自同一个实时数据流，并自动保持同步。

---

## 技术基础

### 离线韧性

网络中断时，系统仍能继续运行：

```
Connectivity lost
    └─▶ App continues with local SQLite
    └─▶ Changes queue in Automerge CRDT
    └─▶ Position caches locally

Connectivity returns
    └─▶ Automerge merges all changes
    └─▶ No conflicts, no data loss
    └─▶ System catches up automatically
```

现场作业不能因网络问题而停止。这套架构确保它们不会停止。

### 实时同步

联网时，更改会即时传播：

```
Driver completes pickup
    └─▶ Form data saves to local SQLite
    └─▶ Automerge syncs to NATS
    └─▶ JetStream persists event
    └─▶ Dashboard updates
    └─▶ Customer notified
    └─▶ Total latency: <1 second
```

### 通知架构

每台设备都订阅相关的 NATS 主题：

| 主题模式 | 接收方 |
|-----------------|--------------|
| `fleet.prod.veh.{id}.route.*` | 该司机的设备 |
| `fleet.prod.veh.*.state.pos` | 调度仪表板 |
| `fleet.prod.customer.{id}.*` | 客户门户 |

更改只需发布一次，即可自动送达所有相关方。

---

## 不止于垃圾收运

同样的模式适用于各类现场作业：

| 行业 | 应用 |
|----------|-------------|
| **废弃物管理** | 路线优化的收运、满溢程度传感器 |
| **配送与物流** | 最后一公里优化、交付证明 |
| **现场服务** | 技术人员调度、零部件库存 |
| **农业** | 农机协调、收获路线规划 |
| **建筑** | 设备与材料配送、工地物流 |
| **公用事业** | 抄表路线、服务上门 |
| **医疗保健** | 居家护理探访、样本采集 |

相同的架构，相同的模式，不同的领域数据。

---

## 安全与数据主权

你的运营数据始终归你所有：

| 原则 | 实现方式 |
|-----------|----------------|
| **自主主权** | 数据存储在你的基础设施或我们的基础设施上，由你选择 |
| **WellKnown Gateway** | 向平台发布数据，不被锁定 |
| **端到端加密** | 数据在传输和存储时均加密 |
| **零信任** | 每台设备都经过认证，每个请求都经过授权 |
| **审计记录** | 每项变更都有日志，并保留完整的来源记录 |
| **离线优先** | 数据存放在设备上，而不仅仅在云端 |

[WellKnown Gateway →]({{< relref "/platform/publish#wellknown-gateway-data-sovereignty" >}}) | [授权与许可 →]({{< relref "/fleet/authorization" >}})

---

## 立即开始

AI 驱动的现场作业彻底改变了工作方式。你的操作员将成为实时自适应系统中的智能边缘节点。

准备好了解 Fleet + Publish + Foundation 如何改变你的现场作业了吗？

[联系我们 →](/contact)
