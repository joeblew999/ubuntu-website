---
title: "软件栈"
meta_title: "无人机机队软件架构 | Ubuntu 软件"
description: "1000 架无人机机队的软件栈：PX4 飞控固件、Ubuntu Server、JetPack、ROS 2 以及基于 Go 的 Vehicle Gateway。"
image: "/images/robotics.svg"
draft: false
---

## 软件架构

每个软件组件都运行开源代码，并承担明确的职责。整个技术栈专为远程管理、空中（OTA）更新以及在 1,000 多台载具间保持一致行为而设计。

---

## 飞行控制器：Pixhawk

### PX4 v1.14 LTS

飞行控制器运行 **PX4 Autopilot v1.14 LTS**：

| 方面 | 详情 |
|--------|---------|
| **版本** | v1.14.x LTS（长期支持） |
| **支持周期** | 提供 2 年以上的安全与关键修复 |
| **配置** | 全机队统一的标准化参数 |
| **更新** | 通过 QGroundControl 或 MAVLink |

**为什么选择 LTS：**

- **稳定性** — 经过生产验证的代码库，变动极少
- **安全性** — 修补关键漏洞而不引入功能变更
- **可预测性** — 整个机队行为一致且可知
- **支持** — 可从 Dronecode 生态系统获得商业支持

**关键模块：**

- **EKF2** — 用于状态估计的扩展卡尔曼滤波器
- **Position Controller** — 多旋翼位置与速度控制
- **Navigator** — 任务执行、RTL（返航）、地理围栏
- **Commander** — 状态机、失效保护处理
- **MAVLink** — 遥测与指令协议

---

## 传感器协处理器：Raspberry Pi

### Ubuntu Server 22.04 LTS

Raspberry Pi 运行 **Ubuntu Server 22.04**：

| 方面 | 详情 |
|--------|---------|
| **发行版** | Ubuntu Server 22.04 LTS |
| **内核** | 带 Pi 补丁的主线 Linux |
| **初始化系统** | systemd |
| **支持周期** | 至 2027 年 4 月 |

**为什么选择 Ubuntu Server（而非 Raspberry Pi OS）：**

- **一致性** — Pi 与 Jetson 使用相同的操作系统
- **软件包** — 标准 apt 软件源，并支持 snap
- **LTS 支持** — 提供 5 年安全更新
- **工具链** — 运维团队熟悉的管理方式

### 核心服务

| 服务 | 用途 |
|---------|---------|
| **mavlink-router** | 在 FCU 与 Jetson 之间路由 MAVLink 消息 |
| **sensor-agent** | 采集环境传感器数据 |
| **gpsd** | 提供位置数据的 GPS/GNSS 守护进程 |
| **chrony** | NTP 时间同步 |

**mavlink-router** 至关重要——它通过 UDP 将 Pixhawk 串口上的 MAVLink 流量转发至 Jetson，使 Vehicle Gateway 能够处理遥测数据。

---

## AI 协处理器：Jetson

### JetPack 6.x

Jetson 运行 **NVIDIA JetPack 6.x**：

| 组件 | 版本 |
|-----------|---------|
| **L4T（Linux for Tegra）** | 基于 Ubuntu 22.04 |
| **CUDA** | 12.x |
| **cuDNN** | 8.x |
| **TensorRT** | 8.x |
| **VPI（视觉编程接口）** | 3.x |

**为什么选择 JetPack 6：**

- **最新 Orin 支持** — 在 Orin 系列上实现完整硬件加速
- **基于 Ubuntu 22.04** — 与 Pi 保持一致，便于统一运维
- **容器支持** — 使用 NVIDIA Container Runtime 隔离工作负载
- **OTA 更新** — 借助 NVIDIA 工具实现全机队固件更新

### CUDA + TensorRT

AI 推理通过 **TensorRT** 运行时执行：

- **模型优化** — 针对边缘部署进行 INT8/FP16 量化
- **批量推理** — 高效处理多帧图像
- **低延迟** — 直接访问 GPU 显存

典型工作负载：

- 目标检测（YOLO、SSD）
- 语义分割
- 光流
- 深度估计

### ROS 2 Humble

**ROS 2 Humble Hawksbill** 提供机器人中间件：

| 方面 | 详情 |
|--------|---------|
| **发行版** | Humble Hawksbill |
| **支持周期** | 至 2027 年 5 月 |
| **DDS** | CycloneDDS（默认）或 FastDDS |
| **构建系统** | colcon |

**为什么选择 ROS 2 Humble：**

- **LTS 版本** — 5 年支持周期与 Ubuntu/JetPack 一致
- **具备实时能力** — 可提供确定性执行路径
- **传感器集成** — 丰富的驱动程序包
- **仿真** — 集成 Gazebo 用于测试

**关键软件包：**

| 软件包 | 用途 |
|---------|---------|
| **px4_msgs** | 面向 ROS 2 的 MAVLink 消息定义 |
| **px4_ros_com** | PX4-ROS 2 桥接（uXRCE-DDS） |
| **image_transport** | 相机数据处理 |
| **tf2** | 坐标系变换 |

---

## Vehicle Gateway

### Go 实现

**Vehicle Gateway** 是运行在每台 Jetson 上的自研 Go 服务：

| 方面 | 详情 |
|--------|---------|
| **语言** | Go 1.21+ |
| **依赖** | nats.go、mavlink 库 |
| **部署** | systemd 服务或容器 |
| **配置** | 环境变量 + 配置文件 |

**为什么选择 Go：**

- **单一二进制文件** — 无运行时依赖，部署简单
- **性能** — 低延迟，内存使用高效
- **并发** — 使用 Goroutine 并行处理 MAVLink/NATS
- **NATS 生态** — 一流的 nats.go 客户端库

**职责：**

1. **MAVLink 接入** — 接收并解析来自 mavlink-router 的 MAVLink 消息
2. **状态降采样** — 降低遥测速率以适应广域网传输
3. **事件提取** — 根据状态转换生成事件
4. **指令执行** — 接收 NATS 指令，校验后发送至 FCU
5. **影子状态协调** — 与机队后端同步期望状态/上报状态

[网关详情 →]({{< relref "/fleet/gateway" >}})

---

## 控制平面

### NATS + JetStream

机队协调运行在 **NATS** 之上，并由 **JetStream** 提供持久化：

| 组件 | 角色 |
|-----------|------|
| **NATS Server** | 核心发布/订阅消息传递 |
| **JetStream** | 流持久化、KV 存储 |
| **Leaf Nodes** | 载具本地的 NATS 实例 |
| **Hub Clusters** | 区域汇聚节点 |

**为什么选择 NATS：**

- **性能** — 单节点每秒可处理数百万条消息
- **简洁** — 单一二进制文件，配置极少
- **持久化** — JetStream 提供持久化流和 KV
- **拓扑** — 叶节点、集群、超级集群，适应任意规模

**机队拓扑：**

```
┌─────────────────────────────────────────────────────┐
│                   Global Hub (optional)             │
│                   NATS Supercluster                 │
└─────────────────┬───────────────────┬───────────────┘
                  │                   │
      ┌───────────▼───────┐ ┌─────────▼───────────┐
      │   Regional Hub A  │ │   Regional Hub B    │
      │   3-node cluster  │ │   3-node cluster    │
      └─────────┬─────────┘ └─────────┬───────────┘
                │                     │
    ┌───────────┼───────────┐         │
    │           │           │         │
┌───▼───┐ ┌─────▼───┐ ┌─────▼───┐ ┌───▼───┐
│ Leaf  │ │  Leaf   │ │  Leaf   │ │ Leaf  │
│ VID-1 │ │ VID-2   │ │ VID-3   │ │VID-N  │
└───────┘ └─────────┘ └─────────┘ └───────┘
```

每台载具运行一个连接到其区域枢纽的**叶节点**。叶节点负责本地发布/订阅，并在断开连接时执行存储转发。

[NATS 拓扑详情 →]({{< relref "/fleet/nats-topology" >}})

---

## 软件更新策略

### 全机队 OTA

| 组件 | 更新机制 |
|-----------|------------------|
| **PX4** | 通过 QGC 或自研工具以 MAVLink 上传固件 |
| **Pi/Ubuntu** | apt + unattended-upgrades，或自定义镜像部署 |
| **Jetson/JetPack** | NVIDIA OTA 工具、容器更新 |
| **Vehicle Gateway** | 通过 systemd 替换二进制文件，或更新容器 |
| **NATS Leaf** | 更新二进制文件并平滑重启 |

**原则：**

- **分阶段发布** — 先更新一部分，验证后再扩大范围
- **回滚能力** — 始终保留上一版本
- **健康检查** — 更新后自动验证
- **最小停机时间** — 在维护窗口内进行更新

---

## 总结

| 层级 | 软件 | 版本 |
|-------|----------|---------|
| **飞行控制器** | PX4 | v1.14 LTS |
| **传感器协处理器** | Ubuntu Server | 22.04 LTS |
| **传感器协处理器** | mavlink-router | 最新版 |
| **AI 协处理器** | JetPack | 6.x |
| **AI 协处理器** | CUDA/TensorRT | 12.x/8.x |
| **AI 协处理器** | ROS 2 | Humble |
| **AI 协处理器** | Vehicle Gateway | Go 1.21+ |
| **控制平面** | NATS + JetStream | 最新版 |

---

## 下一步

[NATS 拓扑 →]({{< relref "/fleet/nats-topology" >}})
