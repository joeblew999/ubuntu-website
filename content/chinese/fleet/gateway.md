---
title: "车辆网关"
meta_title: "自动驾驶车队的车载网关设计 | Ubuntu Software"
description: "车辆网关架构：协议转换（MAVLink、CAN、ROS 2）、状态降采样、事件提取、命令执行以及影子状态协调。"
image: "/images/robotics.svg"
draft: false
---

## 车辆协议与 NATS 之间的桥梁

车辆网关（Vehicle Gateway）是运行在每辆车边缘计算机上的服务，负责在车辆原生协议（MAVLink、CAN bus、ROS 2）与用于车队协调的 NATS 消息系统之间进行转换。

---

## 为什么需要网关？

车辆协议和车队消息传递的用途不同：

| 方面 | 车辆协议 | NATS |
|--------|------------------|------|
| **范围** | 单车 | 整个车队 |
| **协议** | 二进制，紧凑 | JSON，人类可读 |
| **速率** | 100+ Hz 遥测 | 为广域网降采样 |
| **持久化** | 无 | JetStream 流 |
| **模式** | 请求-响应 | 带持久化的发布/订阅 |

网关连接这两个世界，负责协议转换、速率限制和状态管理。

---

## 支持的协议

| 协议 | 车辆类型 | 集成方式 |
|----------|--------------|-------------|
| **MAVLink** | 无人机（PX4、ArduPilot） | UDP/串口，消息解析 |
| **CAN bus** | 轿车、卡车 | SocketCAN，DBC 解码 |
| **J1939** | 重型车辆 | PGN/SPN 映射 |
| **ROS 2** | 自主平台 | 话题订阅 |

每种协议都需要专用的协议适配器，但网关的核心逻辑（降采样、事件、命令、影子）是共享的。

---

## 核心职责

### 1. 协议接入

接收并解析车辆原生协议消息。以 MAVLink（无人机）为例：

```go
// Receive MAVLink frames from mavlink-router
conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 14550})

for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)

    frame, _ := mavlink.Parse(buf[:n])

    switch msg := frame.Message().(type) {
    case *mavlink.Heartbeat:
        handleHeartbeat(msg)
    case *mavlink.GlobalPositionInt:
        handlePosition(msg)
    case *mavlink.BatteryStatus:
        handleBattery(msg)
    // ... handle other message types
    }
}
```

**关键消息类型（MAVLink 示例）：**

| MAVLink 消息 | 内容 |
|-----------------|---------|
| `HEARTBEAT` | 模式、解锁状态、系统状态 |
| `GLOBAL_POSITION_INT` | 纬度、经度、高度、速度 |
| `ATTITUDE` | 横滚、俯仰、偏航 |
| `BATTERY_STATUS` | 电压、电流、剩余电量 |
| `GPS_RAW_INT` | GPS 定位、卫星数、HDOP |
| `SYS_STATUS` | CPU 负载、错误、健康状况 |

**CAN bus 对应项（地面车辆）：**

| CAN 信号 | 内容 |
|------------|---------|
| 发动机转速 | 当前发动机转速 |
| 车速 | 来自轮速传感器的对地速度 |
| GPS 位置 | 来自车载信息模块的经纬度 |
| 燃油液位 | 油箱填充百分比 |
| DTC 代码 | 诊断故障代码 |

### 2. 状态降采样

降低遥测速率以便通过广域网传输：

```go
// Position: downsample from 10Hz to 1Hz
positionSampler := NewDownsampler(100*time.Millisecond, func(msg *Position) {
    // Average positions over window
    return averagePosition(msg)
})

// Battery: only publish on change
batterySampler := NewChangeFilter(func(old, new *Battery) bool {
    return math.Abs(old.Remaining - new.Remaining) > 0.01
})
```

**降采样策略：**

| 策略 | 使用场景 | 示例 |
|----------|----------|---------|
| **基于时间** | 周期性状态 | 以 1Hz 发送位置 |
| **基于变化** | 离散值 | 模式变化 |
| **阈值** | 渐变数据 | 电量变化 Δ > 1% 时发送 |
| **聚合** | 高频数据 | 窗口内的最小/最大/平均值 |

全速率遥测保留在车上，供感知系统使用。只有降采样后的数据才会经过广域网。

### 3. 事件提取

根据状态转换生成事件：

```go
type EventDetector struct {
    previousState *VehicleState
}

func (e *EventDetector) Process(current *VehicleState) []Event {
    var events []Event

    // Detect arm state change
    if current.Armed && !e.previousState.Armed {
        events = append(events, Event{
            Type:      "armed",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"reason": "manual"},
        })
    }

    // Detect mode change
    if current.Mode != e.previousState.Mode {
        events = append(events, Event{
            Type:      "mode_change",
            Timestamp: time.Now(),
            Data: map[string]interface{}{
                "from": e.previousState.Mode,
                "to":   current.Mode,
            },
        })
    }

    // Detect failsafe
    if current.FailsafeActive && !e.previousState.FailsafeActive {
        events = append(events, Event{
            Type:      "failsafe",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"type": current.FailsafeType},
        })
    }

    e.previousState = current
    return events
}
```

**检测的事件：**

| 事件 | 触发条件 |
|-------|---------|
| `armed` | 解锁状态 false → true |
| `disarmed` | 解锁状态 true → false |
| `mode_change` | 飞行模式切换 |
| `takeoff` | 空中状态 false → true |
| `landed` | 空中状态 true → false |
| `failsafe` | 失效保护已激活 |
| `geofence` | 检测到越出地理围栏 |
| `battery.low` | 电量低于阈值 |
| `battery.critical` | 电量处于临界水平 |

### 4. 带策略的命令执行

从 NATS 接收命令并通过 MAVLink 执行：

```go
func (g *Gateway) handleCommand(cmd *Command) *CommandAck {
    // Validate command is allowed
    if !g.policy.Allows(cmd) {
        return &CommandAck{
            Status: "rejected",
            Error:  "command not allowed by policy",
        }
    }

    // Convert to MAVLink command
    mavCmd := cmd.ToMAVLink()

    // Send to flight controller
    if err := g.mavlink.Send(mavCmd); err != nil {
        return &CommandAck{
            Status: "failed",
            Error:  err.Error(),
        }
    }

    // Wait for MAVLink ACK
    mavAck, err := g.mavlink.WaitAck(mavCmd.CommandID, 5*time.Second)
    if err != nil {
        return &CommandAck{
            Status: "timeout",
            Error:  "no response from flight controller",
        }
    }

    return &CommandAck{
        Status: mavAck.Result.String(),
    }
}
```

**策略执行：**

| 策略 | 描述 |
|--------|-------------|
| **地理围栏** | 拒绝目标超出边界的 goto 命令 |
| **高度限制** | 限制命令中的最大高度 |
| **模式限制** | 禁止某些模式切换 |
| **速率限制** | 防止命令泛滥 |
| **身份验证** | 验证命令来源 |

### 5. 影子状态协调

将期望状态与车辆实际状态同步：

```go
func (g *Gateway) reconcileLoop() {
    ticker := time.NewTicker(1 * time.Second)

    for range ticker.C {
        // Read desired state from KV
        desired, _ := g.kv.Get(g.desiredKey())

        // Compare with actual state
        actual := g.getCurrentState()

        // Generate commands to reconcile
        commands := g.reconcile(desired, actual)

        for _, cmd := range commands {
            g.executeCommand(cmd)
        }

        // Update reported state
        reported := g.buildReportedState(actual)
        g.kv.Put(g.reportedKey(), reported)
    }
}

func (g *Gateway) reconcile(desired, actual *State) []*Command {
    var commands []*Command

    // Mode reconciliation
    if desired.Mode != actual.Mode {
        commands = append(commands, &Command{
            Type: "set_mode",
            Data: map[string]interface{}{"mode": desired.Mode},
        })
    }

    // Geofence reconciliation
    if desired.GeofenceEnabled != actual.GeofenceEnabled {
        commands = append(commands, &Command{
            Type: "set_geofence",
            Data: map[string]interface{}{"enabled": desired.GeofenceEnabled},
        })
    }

    return commands
}
```

**影子状态支持声明式管理：**

- 车队运营者设置期望状态
- 网关检测差异
- 网关发出命令使状态收敛
- 网关报告实际状态
- 持续循环

---

## 架构

```
┌────────────────────────────────────────────────────────────────────┐
│                         Vehicle Gateway                             │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐     │
│  │   Protocol   │    │    State     │    │    NATS Client   │     │
│  │   Adapter    │───▶│   Machine    │───▶│    Publisher     │     │
│  └──────────────┘    └──────────────┘    └──────────────────┘     │
│         │                   │                      │               │
│         │            ┌──────▼──────┐               │               │
│         │            │   Event     │               │               │
│         │            │  Detector   │───────────────┤               │
│         │            └─────────────┘               │               │
│         │                                          │               │
│         │            ┌─────────────┐               │               │
│         │            │   Shadow    │◀──────────────┤               │
│         │            │ Reconciler  │               │               │
│         │            └──────┬──────┘               │               │
│         │                   │                      │               │
│         │            ┌──────▼──────┐    ┌─────────▼────────┐      │
│         │            │  Command    │    │   NATS Client    │      │
│         │◀───────────│  Executor   │◀───│   Subscriber     │      │
│         │            └─────────────┘    └──────────────────┘      │
│  ┌──────▼──────┐                                                   │
│  │   Protocol  │                                                   │
│  │   Sender    │                                                   │
│  └─────────────┘                                                   │
└────────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────┐                       ┌─────────────────┐
│  Vehicle Control│                       │   NATS Leaf     │
│ (Pixhawk/ECU/   │                       │   (localhost)   │
│  ROS 2 node)    │                       │                 │
└─────────────────┘                       └─────────────────┘
```

**协议适配器：**
- **MAVLink 适配器** — 用于无人机（PX4、ArduPilot）
- **CAN 适配器** — 用于轿车/卡车（SocketCAN + DBC）
- **ROS 2 适配器** — 用于基于 ROS 的平台

---

## 实现细节

### 依赖

```go
import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/bluenviern/go-mavlink/v2"
)
```

**核心库：**

| 库 | 用途 |
|---------|---------|
| **nats.go** | NATS 客户端、JetStream、KV |
| **go-mavlink** | MAVLink 协议实现 |
| **slog** | 结构化日志 |

### 配置

```yaml
# gateway.yaml
vehicle_id: VID-001
environment: prod

mavlink:
  local_addr: ":14550"
  system_id: 1
  component_id: 1

nats:
  url: "nats://localhost:4222"
  credentials: "/etc/gateway/vehicle.creds"

sampling:
  position_hz: 1
  attitude_hz: 1
  battery_change_threshold: 0.01

policy:
  max_altitude: 120
  geofence_file: "/etc/gateway/geofence.json"
```

### 部署

```yaml
# systemd service
[Unit]
Description=Vehicle Gateway
After=network.target nats.service

[Service]
Type=simple
ExecStart=/usr/local/bin/vehicle-gateway --config /etc/gateway/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## 错误处理

### MAVLink 错误

| 错误 | 响应 |
|-------|----------|
| **连接丢失** | 以退避策略重新连接 |
| **解析错误** | 记录日志并跳过该帧 |
| **超时** | 重试或使命令失败 |

### NATS 错误

| 错误 | 响应 |
|-------|----------|
| **连接丢失** | 本地 NATS 继续运行，重新连接到中心节点 |
| **发布失败** | 本地缓冲，重试 |
| **流错误** | 记录日志、告警、继续运行 |

### 命令错误

| 错误 | 响应 |
|-------|----------|
| **违反策略** | 立即拒绝 |
| **MAVLink 拒绝** | 返回失败 ACK |
| **超时** | 返回超时 ACK |

---

## 指标

网关暴露 Prometheus 指标：

```
# MAVLink
gateway_mavlink_messages_received_total{type="heartbeat"}
gateway_mavlink_messages_sent_total{type="command_long"}
gateway_mavlink_parse_errors_total

# NATS
gateway_nats_messages_published_total{subject="state.position"}
gateway_nats_messages_received_total{subject="cmd.takeoff"}
gateway_nats_publish_errors_total

# Commands
gateway_commands_received_total{type="takeoff"}
gateway_commands_executed_total{type="takeoff",result="success"}
gateway_command_latency_seconds{type="takeoff"}

# Shadow
gateway_shadow_reconciliations_total
gateway_shadow_commands_issued_total
```

---

## 总结

| 职责 | 输入 | 输出 |
|----------------|-------|--------|
| **协议接入** | 车辆消息 | 解析后的遥测 |
| **状态降采样** | 100Hz 遥测 | 1Hz 状态 |
| **事件提取** | 状态转换 | 离散事件 |
| **命令执行** | NATS 命令 | 车辆命令 |
| **影子状态协调** | 期望状态 | 收敛命令 |

车辆网关是关键组件，它使车队规模的运营成为可能，同时保留底层车辆控制系统的安全保障。

---

## 相关文档

- [支持的平台]({{< relref "/fleet/platforms" >}}) — 所有车辆类型概览
- [无人机平台]({{< relref "/fleet/platforms/drones" >}}) — MAVLink 集成详情
- [地面车辆]({{< relref "/fleet/platforms/ground" >}}) — CAN bus 和 J1939 集成

---

## 下一步

[安全模型 →]({{< relref "/fleet/safety" >}})
