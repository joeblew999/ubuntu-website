---
title: "硬件栈"
meta_title: "无人机机队硬件规格 | Ubuntu Software"
description: "千架无人机机队的量产硬件规格：Holybro X500 机架、Pixhawk 6X 飞控，以及 Raspberry Pi 和 Jetson 机载计算机。"
image: "/images/robotics.svg"
draft: false
aliases:
  - "/fleet/hardware/"
---

{{< notice "info" >}}
**在找其他类型的载具？** 本页介绍无人机硬件。地面载具（汽车、卡车、AGV）请参阅[地面载具平台]({{< relref "/fleet/platforms/ground" >}})。如需完整概览，请参阅[支持的平台]({{< relref "/fleet/platforms" >}})。
{{< /notice >}}

## 从本地商店购买

组装无人机需要一定的学习过程。**请从能为您提供帮助的商店购买**——他们供应硬件，我们提供机队软件。

从合作商店购买时，您将获得：
- 经过预先测试、可直接连接的硬件
- 来自懂无人机的本地人员的支持
- 使用我们的托管机队基础设施

---

## 设有实体门店的商店

亲自前往门店，在组装过程中获得面对面的帮助：

{{< featured-shops >}}

**您是无人机商店吗？** [加入我们的合作伙伴计划]({{< relref "/partners/become-partner" >}})，在销售硬件的同时提供机队基础设施。

---

## 所需物品

| 套件 | 组件 | 常见来源 |
|-----|------------|-----------------|
| **整机无人机** | X500 机架、Pixhawk 6X、GPS、电机、电调、螺旋桨 | GetFPV, NewBeeDrone, Holybro |
| **机载计算机** | Raspberry Pi 4 + Jetson Orin Nano | Digi-Key, Mouser, SparkFun |
| **电源与电池** | 4S LiPo 6000mAh、电源模块 | GetFPV, Tattu, Amazon |
| **无线电** | ExpressLRS 接收机、SiK 数传 | GetFPV, HappyModel, Holybro |
| **网络** | 工业以太网交换机、5G 调制解调器 | PUSR, Quectel |

[查看完整组件列表 ↓](#full-bill-of-materials)

---

## 机队规模运营所需硬件

机队运营要求硬件在性能与可维护性之间取得平衡。该技术栈中的每个组件都是基于生产可靠性、零件可获得性以及现场可维修性而选定的。

---

## 机架

### Holybro X500 V2 ARF

**Holybro X500 V2** 是机队部署的理想平台：

| 规格 | 数值 |
|---------------|-------|
| **SKU** | 30125 |
| **轴距** | 500mm |
| **机架重量** | 约 410g（含电机/电调为 610g） |
| **最大起飞重量** | 约 2kg |
| **续航时间** | 15-20 分钟（标准载荷） |
| **电机座** | 16x16mm / 19x19mm |
| **包含** | 电机（2216 KV920）、电调（20A BLHeli S）、螺旋桨（1045） |

**为什么选择 X500：**

- **成熟设计** — 全球已部署数千台，社区知识丰富
- **零件易得** — 机臂、起落架和五金件随时可购
- **便于维护** — 模块化结构，可在现场维修
- **载荷灵活** — 足以搭载双机载计算机及传感器

ARF（Almost Ready to Fly，近乎可飞）套件包含机架、电机、电调和螺旋桨——减少了大型机队中的组装差异。

---

## 飞行控制器

### Pixhawk 6X

运行 **PX4 v1.14 LTS** 的 **Pixhawk 6X** 负责所有飞行关键功能：

| 规格 | 数值 |
|---------------|-------|
| **Pixhawk 6X** | SKU 11073 |
| **底板** | SKU 18117 |
| **配件** | SKU 15011（线缆和外壳） |
| **处理器** | STM32H753（480MHz Cortex-M7） |
| **IMU** | 三重冗余（ICM-42688-P、ICM-45686、BMI088） |
| **气压计** | 双气压计（MS5611、ICP-20100） |
| **磁力计** | IST8310 |
| **接口** | 3x CAN、6x UART、SPI、I2C、PWM |

### 飞行控制器配件

| 组件 | SKU | 用途 |
|-----------|-----|---------|
| **PM06 V2 电源模块** | 15019 | 电源检测 + 5V/3A BEC |
| **Holybro M9N GPS** | 12027 | GPS + 罗盘模块 |
| **GPS 支架** | 12033 | 可折叠支架，用于减少干扰 |
| **安全开关和蜂鸣器** | 12007 | 解锁安全 + 声音警报 |
| **减震座** | 12010 | 防振阻尼 |
| **线缆套件** | 12055 | JST-GH/CAN 线缆 |

**为什么选择 Pixhawk 6X：**

- **冗余** — 三重 IMU、双气压计，支持传感器投票
- **PX4 LTS 支持** — 长期稳定性、安全补丁
- **FMUv6X 标准** — 可与其他符合标准的飞控板互换
- **振动隔离** — 内置 IMU 减震

飞行控制器负责：

- 姿态估计与控制
- 位置保持与导航
- 失效保护行为（RTL、降落、地理围栏）
- RC 输入处理
- MAVLink 遥测

---

## 机载计算机

机队载具运行**双机载计算机**，各自承担不同职责：

### Raspberry Pi 4 / CM4（传感器机载计算机）

负责轻量级传感器集成和数据采集：

| 规格 | 数值 |
|---------------|-------|
| **型号** | Raspberry Pi 4 Model B |
| **SKU** | RPI4-MODEL-B-8GB |
| **处理器** | 四核 Cortex-A72 @ 1.8GHz |
| **内存** | 8GB |
| **存储** | 工业级 microSD 128GB（SDSDQED-128G-GNSIN） |
| **功耗** | 典型值 5V @ 3A |

**职责：**

- 环境传感器驱动（温度、湿度、空气质量）
- GPS/GNSS 数据记录
- 相机采集（非 AI 工作负载）
- 向 Jetson 路由 MAVLink
- Jetson 离线时进行存储转发

### NVIDIA Jetson（AI 机载计算机）

负责计算密集型工作负载：

| 型号 | SKU | GPU 核心 | AI 性能 | 功耗 |
|-------|-----|-----------|----------------|-------|
| **Jetson Orin Nano 8GB** ★ | 900-13767-0000-000 | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin Nano Super** | 945-13766-0000-000 | 1024 CUDA | 67 TOPS | 7-25W |
| **Jetson Orin NX** | 900-13767-0040-000 | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 900-13701-0000-000 | 2048 CUDA | 275 TOPS | 15-60W |

★ = 参考 BOM 选型

| 存储 | SKU | 容量 |
|---------|-----|----------|
| **NVMe SSD** | WD-SN530-256G | 256GB |

**职责：**

- 计算机视觉（目标检测、跟踪）
- 视觉惯性里程计
- 路径规划与避障
- 载具网关（NATS 客户端）
- NATS 叶节点

机队部署通常使用 **Jetson Orin Nano** 进行高性价比推理，或使用 **Orin NX** 处理要求较高的感知工作负载。

---

## 无线电系统

### RC 控制：ExpressLRS

**ExpressLRS** 提供飞手接管链路：

| 规格 | 数值 |
|---------------|-------|
| **接收机** | ExpressLRS EP1 |
| **SKU** | HM-EP1-2400 |
| **频率** | 2.4GHz |
| **范围** | 10km+（取决于功率/天线） |
| **延迟** | <5ms |
| **Protocol** | Open-source (CRSF) |

**Why ExpressLRS:**

- **Open-source** — No vendor lock-in, community-driven development
- **Range** — Reliable link for visual line-of-sight operations
- **Low latency** — Responsive manual control when needed
- **Cost** — Affordable receivers for fleet scale

### Telemetry: SiK Radio

**SiK 915 MHz** radios provide ground station telemetry:

| Specification | Value |
|---------------|-------|
| **SKU** | 17012 (100mW) / 17013 (500mW) |
| **Frequency** | 915MHz (US) / 433MHz (EU) |
| **Range** | 300m-2km (depending on power) |
| **Data Rate** | Up to 250kbps |
| **Interface** | Serial (UART) |
| **Protocol** | MAVLink framing, FHSS |

**Role in Fleet:**

- Backup telemetry when cellular unavailable
- Ground control station connectivity
- Local testing and debugging

For production fleet operations, telemetry primarily flows through the Jetson's cellular connection to NATS. SiK radios serve as backup and for field diagnostics.

---

## Power System

### Battery: 4S LiPo

| Specification | Value |
|---------------|-------|
| **SKU** | TAA60004S35 |
| **Configuration** | 4S (14.8V nominal) |
| **Capacity** | 6000mAh |
| **Discharge Rate** | 35C |
| **Connector** | XT60 |

**Fleet Considerations:**

- Standardize on single battery configuration for logistics
- Use smart batteries with telemetry when available
- Implement battery rotation and health tracking
- Plan for 3:1 battery-to-vehicle ratio for continuous operations

### Power Distribution

| Component | SKU | Purpose |
|-----------|-----|---------|
| **PM06 V2 Power Module** | 15019 | Power sensing + 5V/3A BEC |

- **BEC** — 5V regulated supply for Pi, servos
- **Jetson power** — Direct from battery through regulator (12V typical)

### Networking

| Component | SKU | Purpose |
|-----------|-----|---------|
| **Industrial Ethernet Switch** | USR-ES105 | 3-port GbE between Pi, Jetson, payload |

Enables high-bandwidth communication between companion computers and connected payloads (cameras, sensors).

---

## Connectivity

### Cellular: 4G/5G + eSIM

Each Jetson connects via cellular modem:

| Specification | Value |
|---------------|-------|
| **Modem** | Quectel RM520N-GL |
| **SKU** | RM520NGLAA-M20-SGASA |
| **Bands** | 4G LTE + 5G NR Sub-6GHz |
| **SIM** | eSIM with OTA provisioning |
| **Interface** | M.2 (USB 3.0 / PCIe) |
| **Form Factor** | M.2 NGFF |

**eSIM Benefits:**

- No physical SIM swapping across large fleet
- Remote carrier provisioning
- Carrier switching for coverage optimization
- Centralized subscription management

See [Sensing platform](/platform/sensing/) for details on 5G/LTE connectivity architecture.

---

## Full Bill of Materials

| Category | Component | SKU |
|----------|-----------|-----|
| **Airframe** | Holybro X500 V2 ARF Kit | SKU30125 |
| **Flight Control** | Pixhawk 6X | SKU11073 |
| **Flight Control** | Pixhawk 6X Baseboard | SKU18117 |
| **Flight Control** | Pixhawk 6X Accessories | SKU15011 |
| **Power** | PM06 V2 Power Module | SKU15019 |
| **Navigation** | Holybro M9N GPS | SKU12027 |
| **Navigation** | GPS Mast | SKU12033 |
| **Safety** | Safety Switch & Buzzer | SKU12007 |
| **Companion** | Raspberry Pi 4 8GB | RPI4-MODEL-B-8GB |
| **Storage** | Industrial microSD 128GB | SDSDQED-128G-GNSIN |
| **Companion** | Jetson Orin Nano 8GB | 900-13767-0000-000 |
| **Storage** | NVMe SSD 256GB | WD-SN530-256G |
| **RC** | ExpressLRS EP1 | HM-EP1-2400 |
| **Telemetry** | SiK Telemetry Radio V3 | SKU17012 |
| **Networking** | Industrial Ethernet Switch | USR-ES105 |
| **Power** | LiPo Battery 4S 6000mAh | TAA60004S35 |
| **Mounting** | Pixhawk Vibration Mount | SKU12010 |
| **Cabling** | Holybro Cable Kit | SKU12055 |
| **Connectivity** | Quectel RM520N-GL 5G Modem | RM520NGLAA-M20-SGASA |

---

## Component Sellers

<details>
<summary><strong>按组件查看所有卖家（57 项）</strong></summary>

{{< bom-sellers >}}

</details>

---

## 下一步

[软件栈 →]({{< relref "/fleet/software" >}})
