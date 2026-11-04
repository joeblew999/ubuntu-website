---
title: "感知"
meta_title: "传感与感知平台 | Ubuntu 软件"
description: "面向空间智能的多传感器集成——激光雷达、摄像头与工业传感器通过单一边缘代理统一接入，并支持 5G/eSIM 连接。"
image: "/images/spatial.svg"
draft: false
---

## 空间感知

为数字孪生、机器人和自主系统提供真实世界感知能力。将 LiDAR、摄像头和工业传感器连接到您的空间模型。

---

## 问题所在

传感器会产生数据。要理解这些数据，需要：

- **上下文**：传感器在哪里？它在观测什么？
- **融合**：将多路传感器数据流整合为一幅连贯的图景
- **集成**：连接到设计工具，而不仅仅是仪表板

大多数感知解决方案止步于数据采集。我们则将传感器连接到空间模型。

---

## 部署模式

同一个边缘代理。同样的传感器。不同的配置。

| 模式 | 平台 | 用例 |
|------|----------|----------|
| **空中** | DJI 企业级无人机 | 测绘、巡检、制图 |
| **地面** | 三脚架、背包 | 室内扫描、施工 |
| **机器人** | Viam RDK、ROS2 | 导航、抓取与放置 |
| **固定** | 永久安装 | 交通、安防、仓储 |

### 无 GPS 导航

视觉惯性里程计可在没有 GPS 的情况下实现自主飞行，这对室内环境、城市峡谷和 GPS 拒止区域至关重要。

{{< video src="videos/visual-inertial-nav-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*仅使用摄像头和 IMU 进行实时位置跟踪。无需外部基础设施。*

{{< video src="videos/gps-free-vio-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*基于视觉惯性里程计的自主闭环飞行。*

如需详细的硬件规格，请参阅我们的[无人机机队硬件 →]({{< relref "/fleet/hardware" >}})

---

## 传感器抽象

**从设计上与硬件无关。** 您的代码对接的是我们的统一 API，而不是各个传感器驱动程序。

### 支持的传感器类型

| 类型 | 示例 |
|------|----------|
| **LiDAR** | [Livox Mid-360](https://www.livoxtech.com/mid-360)（[SDK](https://github.com/Livox-SDK/Livox-SDK2)）、[Avia](https://www.livoxtech.com/avia)（[演示](https://terra-1-g.djicdn.com/65c028cd298f4669a7f0e40e50ba1131/Download/Avia/02_Intro%20Demo.mp4)） |
| **RGB-D 摄像头** | [Intel RealSense](https://www.intelrealsense.com/)、[Luxonis OAK-D](https://shop.luxonis.com/products/oak-d) |
| **定位** | [u-blox RTK](https://www.u-blox.com/en/technologies/rtk-real-time-kinematic) GPS/GNSS、IMU |
| **工业** | Modbus 传感器、CAN 总线 |

更换传感器无需修改代码。由配置驱动，而非由代码驱动。

---

## 边缘代理架构

在您的硬件上运行的 Go 二进制程序，支持 Raspberry Pi、Jetson、工业 Linux 或定制 ARM 设备。

| 能力 | 描述 |
|------------|-------------|
| **插件系统** | 通过配置添加传感器，无需修改代码 |
| **本地缓冲** | 离线时存储转发 |
| **实时流传输** | 通过 NATS JetStream 传输到云端 |
| **轻量级** | 单一二进制文件，无运行时依赖 |

---

## 连接

### 支持 eSIM OTA 的 5G/LTE

无需更换 SIM 卡。无需扫描二维码。由服务器推送完成配置。

- 调制解调器出厂时预置引导配置文件
- 由您的平台触发运营商配置文件下载
- 部署期间可通过 API 切换运营商

适用于空中飞行的无人机、移动中的机器人以及偏远地区的固定设施。

---

## 与 Spatial 集成

传感器数据直接输入您的 3D 模型：

| 数据流 | 用途 |
|-----------|---------|
| 点云 → 空间模型 | 实景捕捉 |
| GPS/IMU → 模型定位 | 地理配准 |
| 环境传感器 → 数字孪生 | 实时状态更新 |
| 工业 I/O → 自动化 | 闭环控制 |

不仅仅是仪表板。而是置于上下文中的传感器。

---

## 基于 Foundation 构建

Sensing 自动继承 [Foundation](/platform/foundation/) 的全部能力：

| 能力 | 含义 |
|------------|---------------|
| **离线优先** | 无网络时也能采集，联网后同步 |
| **通用部署** | 边缘、移动、桌面，使用同一个代理 |
| **自主可控** | 您的传感器、您的数据、您的服务器 |
| **实时同步** | 同时向多个目的地推送数据流 |

[进一步了解 Foundation →](/platform/foundation/)

---

## 更宏大体系的一部分

Sensing 是 Ubuntu Software 平台的感知层。

对于需要 3D 设计和 AI 的组织，我们的 Spatial 平台提供设计环境，并可直接集成您的传感器数据。

[探索 Spatial →](/platform/spatial/)

---

## 生产参考：1,000 架无人机机队

了解这些感知能力如何扩展到生产环境。我们的参考架构展示了一套完整的机队部署方案，采用上述硬件和软件栈：1,000 架搭载 PX4 的 Holybro X500 无人机，配备双伴随计算机，并使用 NATS JetStream 实现机队规模的数字孪生。

[探索无人机机队架构 →]({{< relref "/fleet" >}})

---

## 与我们共建

正在部署传感器？正在构建感知系统？我们来聊聊。

[联系我们 →](/contact)
