---
title: "安全模型"
meta_title: "自主车队安全架构 | Ubuntu Software"
description: "自动驾驶车队安全模型：人工接管权限、车辆故障安全机制、控制回路网络隔离以及平稳降级。"
image: "/images/robotics.svg"
draft: false
---

## 安全始于设计

车队规模的自动驾驶车辆运营需要严谨的安全架构。本设计确保**网络故障、软件缺陷和系统失灵永远不会危及车辆安全**。

核心安全原则适用于**所有车辆类型**——无人机、汽车、卡车和 AGV。具体细节有所不同（遥控器 vs 方向盘，RTL vs 原地停车），但架构保持一致。

---

## 核心原则

**NATS 永远不在实时控制回路中。**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SAFETY HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   1. MANUAL OVERRIDE      ─────────────────────────▶  HIGHEST       │
│      Human control (RC, steering, e-stop), always available          │
│                                                                      │
│   2. VEHICLE FAILSAFES    ─────────────────────────▶  HIGH          │
│      Hardware-enforced, local to vehicle                             │
│                                                                      │
│   3. GATEWAY POLICY       ─────────────────────────▶  MEDIUM        │
│      Software enforcement, on-vehicle                                │
│                                                                      │
│   4. FLEET COMMANDS       ─────────────────────────▶  LOWEST        │
│      Network-dependent, advisory                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

上层始终可以覆盖下层。依赖网络的命令权限最低。

---

## 手动接管拥有首要权限

人类操作员拥有**无条件接管权**——无论是通过 RC 发射机、方向盘还是急停按钮。

### 无人机接管（RC）

| 方面 | 实现 |
|--------|----------------|
| **协议** | ExpressLRS 直连接收机 |
| **路径** | TX → RX → Pixhawk（回路中无软件） |
| **模式切换** | 硬件失效保护切换至手动/自稳模式 |
| **紧急停机开关** | 立即锁定电机 |

**RC 从不依赖：**

- Jetson 在线
- NATS 连接
- 车辆网关
- 任何软件组件

即使所有计算机都发生故障，飞手仍保留完全的手动控制权。

#### ExpressLRS 失效保护

当 RC 信号丢失时：

```
1. RX detects signal loss (100ms)
2. RX sends failsafe values to Pixhawk
3. Pixhawk enters RC_LOSS failsafe mode
4. Configured action executes (RTL, land, etc.)
```

失效保护行为在 PX4 中配置，不依赖任何车队软件。

### 地面车辆接管

| 方面 | 汽车/卡车 | AGV |
|--------|-------------|------|
| **主要方式** | 方向盘 + 踏板 | 急停按钮 |
| **路径** | 与执行器的物理连接 | 对电机的硬件中断 |
| **备用** | 关闭钥匙 / 挂空挡 | 无线急停 |
| **网络** | 永不在控制路径中 | 永不在控制路径中 |

**地面车辆接管从不依赖：**

- 车载计算机在线
- NATS 连接
- 任何自动驾驶软件

即使所有计算机都发生故障，驾驶员/操作员仍保留物理控制权，或车辆会安全停车。

---

## 车辆控制系统强制执行失效保护

车辆的控制系统在本地实现失效保护——无需网络。

### 无人机失效保护（PX4）

飞控实现了多种失效保护：

#### RC 丢失失效保护

| 参数 | 值 | 动作 |
|-----------|-------|--------|
| `COM_RC_LOSS_T` | 0.5s | 触发失效保护前的超时 |
| `NAV_RCL_ACT` | 2 | RC 丢失时 RTL |

#### 数据链路丢失

| 参数 | 值 | 动作 |
|-----------|-------|--------|
| `COM_DL_LOSS_T` | 10s | GCS 连接超时 |
| `NAV_DLL_ACT` | 0 | 继续任务（链路非关键） |

#### 电池失效保护

| 级别 | 参数 | 值 | 动作 |
|-------|-----------|-------|--------|
| 低 | `BAT_LOW_THR` | 0.25 | 警告 |
| 严重 | `BAT_CRIT_THR` | 0.15 | RTL |
| 紧急 | `BAT_EMERGEN_THR` | 0.07 | 立即降落 |

#### 地理围栏

| 参数 | 值 | 动作 |
|-----------|-------|--------|
| `GF_ACTION` | 3 | 越界时 RTL |
| `GF_MAX_HOR_DIST` | 视情况而定 | 水平限制 |
| `GF_MAX_VER_DIST` | 120m | 高度限制 |

#### 位置丢失

| 参数 | 值 | 动作 |
|-----------|-------|--------|
| `COM_POS_FS_DELAY` | 1s | 位置超时 |
| `COM_POSCTL_NAVL` | 0 | 位置丢失时降落 |

**所有失效保护均在 Pixhawk 上执行**——无外部依赖。

### 地面车辆失效保护

地面车辆实现了同等的安全机制：

| 触发条件 | 无人机响应 | 地面车辆响应 |
|---------|---------------|------------------------|
| **通信丢失** | 返航（Return-to-Launch） | 原地停车 / 靠边停车 |
| **传感器故障** | 立即降落 | 安全停车，启用驻车 |
| **地理围栏越界** | 在边界处 RTL | 在边界处停车 |
| **急停激活** | 切断电机 | 立即制动 |
| **电量/燃油不足** | RTL 后降落 | 返回场站 / 安全停车 |
| **检测到障碍物** | 避障或悬停 | 停车等待 |

**地面车辆失效保护通常包括：**

- **激光雷达/雷达障碍物检测** — 路径受阻时紧急停车
- **限速** — 根据环境设定最高速度
- **区域限制** — 在本地强制执行禁入区域
- **看门狗定时器** — 自动驾驶软件挂起时停车
- **制动冗余** — 多套独立制动系统

**所有失效保护均在本地执行**——不依赖网络。

---

## NATS 仅用于监控

网络层负责：

| 功能 | 描述 |
|----------|-------------|
| **遥测采集** | 用于仪表盘的车辆状态 |
| **事件日志** | JetStream 中的审计跟踪 |
| **任务协调** | 车队级规划 |
| **影子状态** | 期望/上报状态同步 |

网络层**不负责：**

| 功能 | 原因 |
|----------|---------|
| **姿态控制** | 延迟不可接受，属于安全关键功能 |
| **电机命令** | 必须在 FCU 本地执行 |
| **传感器融合** | 实时性要求 |
| **失效保护决策** | 必须在无网络时也能工作 |

### 命令验证

当命令通过 NATS 到达时：

```go
func (g *Gateway) validateCommand(cmd *Command) error {
    // Check geofence
    if cmd.Type == "goto" {
        if !g.geofence.Contains(cmd.Target) {
            return errors.New("target outside geofence")
        }
    }

    // Check altitude
    if cmd.Altitude > g.policy.MaxAltitude {
        return errors.New("altitude exceeds limit")
    }

    // Check vehicle state
    if cmd.RequiresArmed && !g.state.Armed {
        return errors.New("vehicle not armed")
    }

    // Check rate limiting
    if g.rateLimiter.Exceeded(cmd.Type) {
        return errors.New("command rate exceeded")
    }

    return nil
}
```

命令在**执行前会经过验证**——但即使错误命令被执行，PX4 失效保护也会提供最后一道防护。

---

## 优雅降级

当组件发生故障时，系统会安全降级：

### Jetson 故障

```
Jetson offline
    └─▶ No perception, no AI inference
    └─▶ Pi continues MAVLink routing
    └─▶ RC control fully operational
    └─▶ PX4 failsafes active
    └─▶ Vehicle lands safely or continues mission
```

### NATS 连接丢失

```
WAN disconnected
    └─▶ Leaf node buffers locally
    └─▶ Gateway continues operation
    └─▶ Vehicle follows current mission
    └─▶ State syncs when connection returns
```

### 车辆网关崩溃

```
Gateway process dies
    └─▶ systemd restarts within 5 seconds
    └─▶ PX4 continues autonomous flight
    └─▶ RC override always available
    └─▶ No fleet commands until restart
```

### Hub 集群故障

```
Regional hub down
    └─▶ Leaf nodes queue messages
    └─▶ Vehicles continue operating
    └─▶ RC and PX4 failsafes unaffected
    └─▶ Fleet visibility lost temporarily
```

### 级联故障预防

| 故障 | 影响范围 | 缓解措施 |
|---------|--------------|------------|
| 单台车辆 | 仅该车辆 | 失效保护降落 |
| Hub 节点 | 无影响（已复制） | 集群自愈 |
| Hub 集群 | 区域车队协调 | 车辆继续自主运行 |
| 全局层 | 跨区域可见性 | 区域运营继续 |

---

## 运行安全规程

### 飞行前检查清单

在解锁任何车队车辆之前：

1. **RC 链路已验证** — 飞手确认控制权
2. **失效保护测试** — 验证 RTL/降落行为
3. **地理围栏已加载** — 在 PX4 中确认边界
4. **电池阈值** — 失效保护级别设置合理
5. **天气检查** — 风速在限制范围内
6. **空域畅通** — 无冲突

### 应急程序

| 紧急情况 | 响应 |
|-----------|----------|
| **失去目视** | 通过 RC 激活 RTL |
| **飞丢** | RC 紧急停机开关或地理围栏触发 RTL |
| **多机冲突** | 全车队降落命令 |
| **通信完全丢失** | PX4 RC_LOSS 失效保护 |
| **起火/坠毁** | RC 上锁，联系紧急服务 |

### 停机链

停止车辆的多种方法：

```
1. RC kill switch      → Immediate motor stop
2. RC mode switch      → Land mode
3. GCS command         → RTL or land
4. Fleet command       → RTL or land
5. Geofence breach     → Auto-RTL
6. Battery critical    → Auto-land
```

始终至少有三种方法可用（RC、PX4 失效保护、本地超时）。

---

## 认证考量

该架构支持法规合规：

| 要求 | 实现 |
|-------------|----------------|
| **机长指挥权** | 随时可通过 RC 接管 |
| **失效保护行为** | PX4 认证的失效保护 |
| **地理围栏** | 硬件强制执行的边界 |
| **日志记录** | JetStream 审计跟踪 |
| **Remote ID** | 通过独立模块广播 |

**注意：** 法规要求因司法管辖区而异。该架构提供技术基础——具体合规需要根据当地法规采取额外措施。

---

## 测试安全系统

### 单元测试

```go
func TestGeofenceRejection(t *testing.T) {
    gw := NewGateway(config)
    gw.SetGeofence(testGeofence)

    cmd := &Command{
        Type:   "goto",
        Target: Position{Lat: 0, Lon: 0}, // Outside geofence
    }

    err := gw.validateCommand(cmd)
    assert.ErrorContains(t, err, "outside geofence")
}
```

### 集成测试

```
1. Simulate RC loss → Verify RTL activates
2. Simulate battery drain → Verify land at threshold
3. Simulate geofence breach → Verify return behavior
4. Simulate Jetson crash → Verify continued flight
5. Simulate NATS disconnect → Verify local operation
```

### 飞行测试

在车队部署之前：

1. **单车失效保护验证**
2. **RC 丢失响应验证**
3. **地理围栏越界测试**
4. **电池失效保护验证**
5. **多车协调测试**

---

## 总结

| 层级 | 权限 | 无人机 | 地面车辆 |
|-------|-----------|--------|-----------------|
| **手动** | 最高 | RC 发射机 | 方向盘 / 急停 |
| **车辆** | 高 | PX4 失效保护 | ECU / 安全控制器 |
| **网关** | 中 | 策略执行 | 策略执行 |
| **车队** | 最低 | 仅协调 | 仅协调 |

该安全模型确保：

- **网络永远不是安全关键环节**
- **操作员始终拥有接管权**
- **硬件失效保护在本地执行**
- **故障时优雅降级**
- **不存在导致全车队受影响的单点故障**

---

## 相关文档

- [支持的平台]({{< relref "/fleet/platforms" >}}) — 所有车辆类型概览
- [无人机平台]({{< relref "/fleet/platforms/drones" >}}) — 无人机专属安全细节
- [地面车辆]({{< relref "/fleet/platforms/ground" >}}) — 地面车辆安全细节

---

## 延伸阅读

- [PX4 安全配置](https://docs.px4.io/main/en/config/safety.html)
- [MAVLink 通用消息集](https://mavlink.io/en/messages/common.html)
- [NATS 安全模型](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

---

## 返回概览

[← 自动驾驶车队架构]({{< relref "/fleet" >}})
