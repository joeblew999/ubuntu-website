---
title: "JetStream 配置"
meta_title: "面向无人机机队数字孪生的 JetStream 流 | Ubuntu Software"
description: "用于车队数字孪生的 JetStream 流和 KV 存储配置：状态汇总、事件审计跟踪、命令队列和影子状态。"
image: "/images/robotics.svg"
draft: false
---

## 面向数字孪生的 JetStream

JetStream 为车队数字孪生提供持久化层。每种流类型都有特定用途，并配有量身定制的保留和复制设置。

---

## 流架构

五种 JetStream 资源支撑数字孪生模式：

| 资源 | 类型 | 用途 |
|----------|------|---------|
| **TWIN_STATE** | 流 | 带汇总（rollup）的遥测数据 |
| **TWIN_EVENTS** | 流 | 持久化的事件审计跟踪 |
| **TWIN_CMD** | 流 | 命令队列 |
| **TWIN_SHADOW** | KV 存储 | 期望/上报状态 |
| **TWIN_BLOBS** | 对象存储 | 大型制品 |

---

## TWIN_STATE：遥测流

使用**基于主题的汇总**存储持续的车辆状态：

```yaml
name: TWIN_STATE
subjects:
  - "fleet.prod.veh.*.state.>"
retention: limits
max_age: 1h
max_bytes: 10GB
storage: file
replicas: 3
discard: old
rollup_hdrs: true
```

### 配置说明

| 设置 | 值 | 理由 |
|---------|-------|-----------|
| **retention** | limits | 受时长和大小限制 |
| **max_age** | 1h | 保留近期状态以供重放 |
| **max_bytes** | 10GB | 限制每个枢纽的存储量 |
| **replicas** | 3 | 节点故障时仍可用 |
| **rollup_hdrs** | true | 启用按主题汇总 |

### 主题汇总

启用汇总后，消费者可以请求**每个主题仅最新的一条消息**：

```go
// Get latest position for all vehicles
sub, _ := js.Subscribe("fleet.prod.veh.*.state.position",
    nats.DeliverLastPerSubject())
```

这使流从数据洪流转变为可查询的状态存储。

### 消费者模式

| 模式 | 使用场景 |
|---------|----------|
| **DeliverAll** | 历史重放、分析 |
| **DeliverLast** | 当前状态快照 |
| **DeliverLastPerSubject** | 每辆车的最新状态 |
| **StartTime** | 从指定时间戳开始重放 |

---

## TWIN_EVENTS：审计跟踪

以**长期保留**方式存储离散事件，用于合规和调试：

```yaml
name: TWIN_EVENTS
subjects:
  - "fleet.prod.veh.*.evt.>"
retention: limits
max_age: 90d
max_bytes: 100GB
storage: file
replicas: 3
discard: old
```

### 配置说明

| 设置 | 值 | 理由 |
|---------|-------|-----------|
| **max_age** | 90d | 法规合规、事故分析 |
| **max_bytes** | 100GB | 事件比状态数据更小 |
| **storage** | file | 重启后仍持久保存 |

### 事件特性

事件与状态的区别：

| 特性 | 状态 | 事件 |
|----------|-------|--------|
| **频率** | 持续（Hz） | 离散（变化时） |
| **保留** | 短期（小时） | 长期（月） |
| **消费者** | 实时仪表板 | 审计、分析、告警 |
| **重放** | 重建当前状态 | 调查、合规 |

### 事件示例

```json
// Armed event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "previous_state": "disarmed",
    "armed_by": "operator-123",
    "reason": "mission_start"
  }
}

// Failsafe event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T11:45:00Z",
  "data": {
    "failsafe_type": "low_battery",
    "action_taken": "rtl",
    "battery_remaining": 0.15
  }
}
```

---

## TWIN_CMD：命令队列

以**工作队列保留**方式存储命令：

```yaml
name: TWIN_CMD
subjects:
  - "fleet.prod.veh.*.cmd.>"
retention: workqueue
max_age: 5m
storage: file
replicas: 3
```

### 配置说明

| 设置 | 值 | 理由 |
|---------|-------|-----------|
| **retention** | workqueue | 确认后删除消息 |
| **max_age** | 5m | 未确认的命令会过期 |
| **replicas** | 3 | 保证命令送达 |

### 工作队列语义

1. 命令发布到流
2. 车辆网关接收（精确一次投递）
3. 网关处理命令
4. 网关确认消息
5. 消息从流中移除

**未确认的命令会重新投递**，以确保断连期间命令不会丢失。

### 命令流程

```
Fleet Operator               TWIN_CMD Stream              Vehicle Gateway
     │                             │                             │
     │ Publish: cmd.takeoff        │                             │
     │ ─────────────────────────▶  │                             │
     │                             │ Deliver to consumer         │
     │                             │ ─────────────────────────▶  │
     │                             │                             │ Process
     │                             │                             │ command
     │                             │             Ack             │
     │                             │ ◀─────────────────────────  │
     │                             │ (message removed)           │
     │                             │                             │
```

---

## TWIN_SHADOW：KV 存储

JetStream KV 存储维护数字孪生的影子状态：

```yaml
bucket: TWIN_SHADOW
max_value_size: 1MB
history: 5
ttl: 0
replicas: 3
storage: file
```

### 配置说明

| 设置 | 值 | 理由 |
|---------|-------|-----------|
| **history** | 5 | 保留前 5 个版本 |
| **ttl** | 0 | 不过期（永久） |
| **replicas** | 3 | 高可用 |

### 键结构

```
fleet/prod/veh/VID-001/desired    → Desired state JSON
fleet/prod/veh/VID-001/reported   → Reported state JSON
fleet/prod/veh/VID-001/config     → Vehicle configuration
fleet/prod/veh/VID-001/meta       → Metadata (serial, type, etc.)
```

### 影子操作

```go
// Write desired state
kv.Put("fleet/prod/veh/VID-001/desired", desiredJSON)

// Read reported state
entry, _ := kv.Get("fleet/prod/veh/VID-001/reported")

// Watch for changes
watcher, _ := kv.Watch("fleet/prod/veh/VID-001/desired")
for entry := range watcher.Updates() {
    // Desired state changed, reconcile
}
```

### 用于调试的历史记录

启用历史记录后，可以查看先前的状态：

```go
// Get all historical values
history, _ := kv.History("fleet/prod/veh/VID-001/desired")
for _, entry := range history {
    fmt.Printf("Rev %d at %s: %s\n",
        entry.Revision(), entry.Created(), entry.Value())
}
```

---

## TWIN_BLOBS：对象存储

无法放入消息中的大型制品：

```yaml
bucket: TWIN_BLOBS
max_chunk_size: 128KB
storage: file
replicas: 3
```

### 使用场景

| 制品 | 大小 | 示例 |
|----------|------|---------|
| **任务文件** | 10-100KB | 航点定义 |
| **日志文件** | 1-10MB | 飞行日志、诊断信息 |
| **固件** | 10-50MB | PX4 固件镜像 |
| **地图** | 1-100MB | 离线地图瓦片 |
| **ML 模型** | 10-500MB | TensorRT 引擎文件 |

### 对象操作

```go
// Upload mission file
obj, _ := obs.PutFile("missions/mission-456.json", file)

// Download to vehicle
reader, _ := obs.Get("missions/mission-456.json")

// List available firmware
objects := obs.List(nats.ObjectSearchPrefix("firmware/"))
```

### 分块传输

对象存储通过以下方式处理大文件：

1. 拆分为分块（默认 128KB）
2. 将分块存储为流消息
3. 跟踪元数据（大小、哈希、分块）
4. 检索时重新组装

这使得在不可靠的连接上也能实现**可断点续传的传输**。

---

## 流容量规划

### 单车估算

| 流 | 消息大小 | 速率 | 每日数据量 |
|--------|--------------|------|--------------|
| **STATE** | 200 字节 | 10 Hz | 约 170 MB |
| **EVENTS** | 500 字节 | 10 次/天 | 约 5 KB |
| **CMD** | 300 字节 | 100 次/天 | 约 30 KB |

### 车队估算（1,000 辆车）

| 资源 | 1 小时 | 1 天 | 90 天 |
|----------|--------|-------|--------|
| **TWIN_STATE** | 7 GB | - | - |
| **TWIN_EVENTS** | - | 5 MB | 450 MB |
| **TWIN_SHADOW** | - | 1 GB | 1 GB |
| **TWIN_BLOBS** | - | - | 50 GB |

以上为基线估算。实际用量取决于消息速率和大小。

---

## 运维注意事项

### 流健康状况

监控以下指标：

| 指标 | 告警阈值 |
|--------|-----------------|
| **消费者延迟** | > 10,000 条消息 |
| **流字节数** | > 最大值的 80% |
| **副本同步** | 任一副本落后 |
| **待确认数** | 每个消费者 > 1,000 |

### 备份策略

| 资源 | 备份方式 | 频率 |
|----------|---------------|-----------|
| **流** | JetStream 快照 | 每日 |
| **KV 存储** | 键导出 | 每日 |
| **对象存储** | 文件级备份 | 每周 |

### 迁移

流可以在集群之间迁移：

1. 在目标集群上以相同配置创建流
2. 创建以原始流为源的镜像
3. 等待同步完成
4. 将生产者/消费者切换到目标集群
5. 删除原始流

---

## 总结

| 资源 | 类型 | 保留 | 用途 |
|----------|------|-----------|---------|
| **TWIN_STATE** | 流 | 1 小时 | 带汇总的遥测数据 |
| **TWIN_EVENTS** | 流 | 90 天 | 审计跟踪 |
| **TWIN_CMD** | 流 | 工作队列 | 命令投递 |
| **TWIN_SHADOW** | KV | 永久 | 期望/上报状态 |
| **TWIN_BLOBS** | 对象 | 永久 | 大型制品 |

---

## 下一步

[车辆网关 →]({{< relref "/fleet/gateway" >}})
