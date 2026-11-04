---
title: "主题命名"
meta_title: "无人机机队的 NATS 主题命名 | Ubuntu Software"
description: "用于车队遥测、事件、命令及数字孪生状态管理的分层 NATS 主题命名约定。"
image: "/images/robotics.svg"
draft: false
---

## 主题层次结构设计

主题命名是消息系统的 API。设计良好的层次结构可以实现高效的路由、过滤和访问控制，同时让开发者易于理解。

---

## 设计原则

| 原则 | 理由 |
|-----------|-----------|
| **层次化** | 支持在任意层级进行通配符订阅 |
| **可预测** | 开发者无需查询即可构造主题 |
| **环境感知** | 将生产环境与预发布/开发环境分开 |
| **可过滤** | 可订阅单个载具、一个分组或整个机队 |
| **可扩展** | 无需重构即可添加新的消息类型 |

---

## 核心主题模式

所有机队主题都遵循以下模式：

```
fleet.<env>.veh.<vid>.<category>.<type>
```

| 段 | 描述 | 示例 |
|---------|-------------|----------|
| `fleet` | 根命名空间 | 始终为 "fleet" |
| `<env>` | 环境 | `prod`、`staging`、`dev` |
| `veh` | 资源类型 | 载具（未来：`gcs` 用于地面控制） |
| `<vid>` | 载具 ID | `VID-001`、`drone-alpha-7` |
| `<category>` | 消息类别 | `state`、`evt`、`cmd`、`cmdack` |
| `<type>` | 具体类型 | `position`、`battery`、`takeoff` |

---

## 消息类别

### 状态：`fleet.<env>.veh.<vid>.state.*`

载具的当前状态，持续发布：

| 主题 | 描述 | 频率 |
|---------|-------------|------|
| `state.position` | 纬度/经度/高度、速度 | 10 Hz |
| `state.attitude` | 横滚/俯仰/偏航 | 10 Hz |
| `state.battery` | 电压、电流、剩余电量 | 1 Hz |
| `state.mode` | 飞行模式（AUTO、MANUAL、RTL） | 变化时 |
| `state.health` | 系统健康摘要 | 1 Hz |

**发布示例：**

```go
subject := fmt.Sprintf("fleet.prod.veh.%s.state.position", vehicleID)
nc.Publish(subject, positionJSON)
```

### 事件：`fleet.<env>.veh.<vid>.evt.*`

离散事件，在发生时发布：

| 主题 | 描述 |
|---------|-------------|
| `evt.armed` | 载具已解锁 |
| `evt.disarmed` | 载具已上锁 |
| `evt.takeoff` | 已开始起飞 |
| `evt.landed` | 着陆完成 |
| `evt.waypoint` | 已到达航点 |
| `evt.failsafe` | 已触发失效保护 |
| `evt.geofence` | 越出地理围栏 |
| `evt.battery.low` | 低电量警告 |

**事件是状态转换**，而非连续数据。它们是持久化的——存储在 JetStream 中用于审计追踪。

### 命令：`fleet.<env>.veh.<vid>.cmd.*`

发送给载具的指令：

| 主题 | 描述 |
|---------|-------------|
| `cmd.arm` | 解锁电机 |
| `cmd.disarm` | 上锁电机 |
| `cmd.takeoff` | 开始起飞 |
| `cmd.land` | 开始着陆 |
| `cmd.rtl` | 返航 |
| `cmd.goto` | 导航至指定位置 |
| `cmd.mission.upload` | 上传任务 |
| `cmd.mission.start` | 开始任务 |
| `cmd.param.set` | 设置参数 |

**命令向下游流动**，从机队管理系统发往载具。

### 命令确认：`fleet.<env>.veh.<vid>.cmdack.*`

对命令的响应：

| 主题 | 描述 |
|---------|-------------|
| `cmdack.arm` | 解锁命令结果 |
| `cmdack.takeoff` | 起飞命令结果 |
| `cmdack.goto` | Goto 命令结果 |

**确认消息负载包括：**

```json
{
  "cmd_id": "cmd-123456",
  "status": "accepted|rejected|completed|failed",
  "error": "optional error message",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## 通配符订阅

NATS 通配符支持灵活的订阅方式：

| 通配符 | 含义 | 示例 |
|----------|---------|---------|
| `*` | 单个令牌 | `fleet.prod.veh.*.state.position` |
| `>` | 一个或多个令牌 | `fleet.prod.veh.VID-001.>` |

### 常见订阅模式

| 模式 | 用例 |
|---------|----------|
| `fleet.prod.veh.*.state.position` | 所有载具的位置 |
| `fleet.prod.veh.VID-001.>` | 单个载具的所有消息 |
| `fleet.prod.veh.*.evt.>` | 所有载具的所有事件 |
| `fleet.prod.veh.*.state.battery` | 所有电池状态 |
| `fleet.*.veh.VID-001.>` | 单个载具在所有环境中的消息 |

---

## KV 存储：数字孪生影子

JetStream KV 存储维护**期望**状态和**上报**状态：

### 主题模式（KV）

```
fleet/<env>/veh/<vid>/desired
fleet/<env>/veh/<vid>/reported
```

注意：KV 使用 `/` 分隔符（存储桶路径），而发布/订阅使用 `.` 分隔符。

### 期望状态

机队管理系统**希望**载具执行的操作：

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "geofence_enabled": true,
  "max_altitude": 120,
  "home_position": {
    "lat": 37.7749,
    "lon": -122.4194
  }
}
```

### 上报状态

载具**实际**上报的内容：

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "mission_progress": 0.45,
  "geofence_enabled": true,
  "armed": true,
  "in_flight": true,
  "battery_remaining": 0.72
}
```

### 影子协调

载具网关持续执行以下操作：

1. 从 KV 读取期望状态
2. 与载具的实际状态进行比较
3. 发出命令，使实际状态与期望状态保持一致
4. 更新上报状态

这实现了**声明式机队管理**——设置期望状态，让载具自行收敛。

---

## 访问控制

NATS 授权将用户映射到其被允许的主题：

```
# Vehicle can publish its own state/events
publish: fleet.prod.veh.VID-001.state.>
publish: fleet.prod.veh.VID-001.evt.>
publish: fleet.prod.veh.VID-001.cmdack.>

# Vehicle can subscribe to commands for itself
subscribe: fleet.prod.veh.VID-001.cmd.>

# Vehicle can read/write its own shadow
kv: fleet/prod/veh/VID-001/*
```

```
# Fleet operator can subscribe to all state/events
subscribe: fleet.prod.veh.*.state.>
subscribe: fleet.prod.veh.*.evt.>
subscribe: fleet.prod.veh.*.cmdack.>

# Fleet operator can publish commands
publish: fleet.prod.veh.*.cmd.>

# Fleet operator can manage shadows
kv: fleet/prod/veh/*/desired
```

这确保了**载具无法冒充其他载具**，也无法发布未经授权的命令。

有关授权模型的详细介绍，包括第三方授权和凭证管理，请参阅 [授权与授予]({{< relref "/fleet/authorization" >}})。

---

## 主题设计权衡

### 为什么不使用扁平主题？

像 `vehicle-001-position` 这样的扁平主题不支持：

- 通配符订阅
- 层次化访问控制
- 逻辑分组

### 为什么不使用更深的层次结构？

像 `fleet.prod.region.us-west.site.warehouse-a.veh.VID-001.state.position` 这样的主题：

- 更难以通过编程方式构造
- 通配符会变得复杂
- 变更会波及整个系统

四级层次结构在灵活性与简洁性之间取得了平衡。

---

## 消息格式

所有消息都使用结构一致的 **JSON**：

```json
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00.123Z",
  "seq": 12345,
  "data": {
    // message-specific payload
  }
}
```

| 字段 | 描述 |
|-------|-------------|
| `vid` | 载具 ID（与主题重复，但对消费者很有用） |
| `timestamp` | 毫秒精度的 ISO 8601 时间戳 |
| `seq` | 用于排序的单调递增序列号 |
| `data` | 特定于消息的负载 |

---

## 总结

| 类别 | 主题模式 | 方向 |
|----------|-----------------|-----------|
| **状态** | `fleet.<env>.veh.<vid>.state.<type>` | 载具 → 中心 |
| **事件** | `fleet.<env>.veh.<vid>.evt.<type>` | 载具 → 中心 |
| **命令** | `fleet.<env>.veh.<vid>.cmd.<type>` | 中心 → 载具 |
| **确认** | `fleet.<env>.veh.<vid>.cmdack.<type>` | 载具 → 中心 |
| **期望** | `fleet/<env>/veh/<vid>/desired`（KV） | 中心 → 载具 |
| **上报** | `fleet/<env>/veh/<vid>/reported`（KV） | 载具 → 中心 |

---

## 下一步

[流配置 →]({{< relref "/fleet/streams" >}})
