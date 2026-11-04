---
title: "授权与许可"
meta_title: "去中心化授权与访问授予 | Ubuntu Software"
description: "NATS JetStream 基于授权的去中心化安全模型，使第三方无需中央协调即可安全地参与实时系统。"
image: "/images/robotics.svg"
draft: false
---

## 无瓶颈的安全

传统的集中式授权在规模化时会带来问题：

| 问题 | 影响 |
|---------|--------|
| **单点故障** | 认证服务器宕机 = 整个系统宕机 |
| **每个请求都有延迟** | 往返认证服务器会增加延迟 |
| **无法在断网状态下工作** | 没有网络 = 无法授权 |
| **无法扩展** | 数千台设备会压垮中央服务器 |

NATS 采用了不同的方法：基于自包含凭证的**去中心化授权**。

---

## NATS 授权的工作原理

```
┌─────────────────────────────────────────────────────────────────┐
│                    CREDENTIAL ISSUANCE                          │
│   Account Authority issues signed JWTs with embedded perms      │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │ Vehicle │     │ Operator│     │ Partner │
        │  Creds  │     │  Creds  │     │  Creds  │
        └────┬────┘     └────┬────┘     └────┬────┘
             │               │               │
             ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS SERVER                                │
│   Validates JWT signature locally — no external call needed     │
│   Enforces permissions from embedded claims                     │
└─────────────────────────────────────────────────────────────────┘
```

**关键特性：**

- **本地验证** — 服务器无需调用任何外部服务即可验证凭证
- **内嵌权限** — 你能做什么已编码在你的凭证中
- **即时撤销** — 账户级撤销列表在整个集群中传播
- **零信任** — 没有凭证 = 没有访问权限，仅此而已

这在断网期间同样有效。授权无需回连中心。凭证本身就是证明。

---

## 权限范围

凭证精确规定了一方可以做什么：

| 范围 | 描述 | 示例 |
|-------|-------------|---------|
| **发布** | 向主题发送消息 | 车辆发布 `fleet.prod.veh.V001.state.pos` |
| **订阅** | 从主题接收消息 | 操作员订阅 `fleet.prod.veh.*.state.*` |
| **发布 + 订阅** | 双向 | 网关处理命令/响应 |
| **拒绝** | 显式阻止特定模式 | 阻止车辆订阅其他车辆 |

**权限中的主题通配符：**

```
fleet.prod.veh.V001.*        → Only this vehicle's subjects
fleet.prod.veh.*.state.*     → All vehicles' state (read-only)
fleet.prod.veh.*.cmd.*       → All vehicles' commands (operator)
fleet.prod.veh.*.evt.*       → All events (auditor)
```

我们正是通过这种方式确保车辆只能发布自己的数据，且无法冒充其他车辆。完整层级结构请参见[主题命名]({{< relref "/fleet/subjects" >}})。

---

## 授权模型（Grants）

授权（Grants）实现了**细粒度、可委托的权限**。操作员可以授予第三方访问特定数据的权限，而无需给予其完整的系统访问权限。

```
Account Authority (root)
    │
    ├── Fleet Operator Account
    │       │
    │       ├── Grant: Emergency Services (incident-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.pos
    │       │       └── Subscribe: fleet.prod.veh.*.stream.video
    │       │       └── Expires: 4 hours
    │       │
    │       ├── Grant: Research Partner (data-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.env.*
    │       │       └── Deny: fleet.prod.veh.*.state.pos
    │       │       └── Expires: 30 days
    │       │
    │       └── Grant: Auditor (read-only, full)
    │               └── Subscribe: fleet.prod.veh.*.state.*
    │               └── Subscribe: fleet.prod.veh.*.evt.*
    │               └── Deny: fleet.prod.veh.*.cmd.*
    │               └── Expires: 7 days
```

**授权特性：**

| 特性 | 优势 |
|----------|---------|
| **范围限定** | 仅限特定主题/模式——仅此而已 |
| **时间限定** | 自动过期，不会遗留被遗忘的访问权限 |
| **可撤销** | 即时撤销，无需触及被授权方 |
| **可委托** | 操作员可将子权限授予他人 |
| **可审计** | 每次授权签发均有记录 |

---

## 第三方集成场景

### 应急服务集成

```
Scenario: Wildfire response, fire department needs drone video

Grant issued:
  - Subject: fleet.prod.veh.*.stream.video
  - Subject: fleet.prod.veh.*.state.pos
  - Duration: 4 hours (incident duration)
  - Audit: All access logged with responder ID

Result: Fire department connects with their credential, receives
real-time video and position from all drones in area.
Grant expires automatically when incident closes.
```

**为什么可行：**
- 消防部门恰好获得所需的内容——实时视频和位置
- 无权访问命令——无法干扰无人机操作
- 时间限定——事件结束时访问权限随之终止
- 可审计——他们接收的每一帧都有记录

### 合作伙伴数据共享

```
Scenario: University research on environmental monitoring

Grant issued:
  - Subject: fleet.prod.veh.*.state.env.* (temp, humidity, air quality)
  - Deny: fleet.prod.veh.*.state.pos (no location data)
  - Duration: 30 days (research period)
  - Audit: Data access logged for compliance

Result: Researcher receives environmental sensor data without
knowing vehicle locations—privacy preserved.
```

**为什么可行：**
- 研究人员获得其研究所需的环境数据
- 位置数据被显式拒绝——隐私得到保护
- 清晰的审计追踪，满足研究合规要求
- 研究期结束后自动清理

### 监管合规

```
Scenario: Aviation authority audit of flight operations

Grant issued:
  - Subject: fleet.prod.veh.*.state.* (all telemetry)
  - Subject: fleet.prod.veh.*.evt.* (all events)
  - Deny: fleet.prod.veh.*.cmd.* (no command capability)
  - Duration: 7 days (audit period)
  - Audit: Every query logged for audit trail

Result: Auditor gets complete read access to historical data
but cannot issue any commands to vehicles.
```

**为什么可行：**
- 对审计员完全透明——毫无隐瞒
- 只读——审计员无法影响运营
- 审计员访问的完整审计追踪
- 时间限定——不存在永久性访问漏洞

### 客户集成

```
Scenario: Enterprise customer wants API access to their fleet data

Grant issued:
  - Subject: fleet.prod.veh.CUST-*.state.* (only their vehicles)
  - Subject: fleet.prod.veh.CUST-*.evt.*
  - Deny: fleet.prod.veh.OTHER-*.* (other customers' data)
  - Duration: Persistent (until revoked)
  - Audit: API access logged

Result: Customer builds their own dashboards using real-time
data from their vehicles—and only their vehicles.
```

### 分包商访问

```
Scenario: Maintenance team needs to run diagnostics on specific vehicles

Grant issued:
  - Subject: fleet.prod.veh.V001.cmd.diag.* (diagnostic commands only)
  - Subject: fleet.prod.veh.V001.state.health.*
  - Deny: fleet.prod.veh.V001.cmd.flight.* (no flight commands)
  - Duration: 8 hours (maintenance window)
  - Audit: All commands logged with technician ID

Result: Technician can run diagnostics but cannot arm or fly the vehicle.
Access expires after maintenance window.
```

---

## 安全特性

| 特性 | 实现方式 |
|----------|----------------|
| **零信任** | 无隐式信任——每个连接都经过认证 |
| **最小权限** | 凭证仅包含最低必要权限 |
| **纵深防御** | 多层防护：TLS + JWT + 主题 ACL |
| **离线可用** | 无需网络连接即可完成授权 |
| **可审计** | 每次访问均记录凭证身份 |
| **可撤销** | 通过账户撤销列表即时撤销 |

---

## 运营优势

### 对于车队运营方

- **无需更改系统即可添加/移除第三方** — 签发或撤销凭证，无需编辑配置文件
- **无需共享凭证即可授予临时访问权限** — 每一方都获得自己的凭证
- **审计谁在何时访问了什么** — 完整的访问日志
- **必要时即时撤销访问权限** — 在下次连接时生效

### 对于第三方

- **自助连接** — 使用提供的凭证，准备就绪即可连接
- **范围清晰** — 确切知道你能访问和不能访问什么
- **无需了解内部架构** — 只需连接并订阅
- **自动过期** — 你无需进行任何清理

### 对于合规

- **完整的审计追踪** — 每次访问都可追溯到凭证
- **可证明的权限边界** — 证明谁有权访问什么
- **审计员的时间限定访问** — 不存在永久性访问漏洞
- **无共享凭证** — 个人问责

---

## 实施说明

### 凭证格式

NATS 使用 **NKeys** 和 **JWT**：

- **Ed25519 签名** — 密码学安全，验证快速
- **自包含声明** — 权限内嵌于令牌中
- **无需外部验证** — 服务器在本地验证签名

### 撤销

当你需要撤销访问权限时：

1. 将凭证添加到账户撤销列表
2. 该列表传播到集群中的所有服务器
3. 在下次连接尝试时生效
4. 无需触及被授权方的系统

### 凭证轮换

对于长期集成：

1. 在旧凭证过期前签发新凭证
2. 重叠期可实现无缝过渡
3. 过渡完成后撤销旧凭证
4. 服务不中断

---

## 为什么这对实际系统很重要

传统授权需要：
- 中央认证服务器（单点故障）
- 每个请求都需网络调用（延迟）
- 第三方集成复杂（耗时数月）
- 手动清理过期访问权限（被遗忘的凭证）

NATS 去中心化授权提供：
- **弹性** — 离线可用，无中心瓶颈
- **速度** — 本地验证，无需往返
- **可扩展性** — 几分钟即可添加第三方，而非数月
- **安全性** — 自动过期，即时撤销

这就是构建既**高度安全**又**高度易用**的系统的方式。

---

## 相关文档

- [主题命名]({{< relref "/fleet/subjects" >}}) — 主题层级如何实现细粒度 ACL
- [NATS 拓扑]({{< relref "/fleet/nats-topology" >}}) — 消息基础设施的结构
- [安全模型]({{< relref "/fleet/safety" >}}) — 命令授权如何融入安全层级
- [监控]({{< relref "/company/surveillance" >}}) — 授权如何为监控实现第三方访问

---

## 下一步

[流与事件 →]({{< relref "/fleet/streams" >}})
