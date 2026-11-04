---
title: "机器人技术"
meta_title: "机器人技术栈 | Ubuntu 软件"
description: "我们基于 Viam RDK 构建的机器人架构——从设计到部署，涵盖硬件抽象、计算机视觉和工业集成。"
image: "/images/spatial.svg"
draft: false
---

## 机器人技术栈

从设计到部署。Spatial 提供 3D 设计与仿真层。Viam RDK 提供用于控制和操作实体机器人的运行时。

---

## 嵌入式 Linux 运行时

机器人运行在 Linux 上。凭借数十年的 Linux 和嵌入式系统经验，我们为真实世界而构建——工厂车间、户外环境以及资源受限的边缘设备。

### Viam RDK 运行时

**我们为何选择 Viam RDK：**

- **开源** — 契合我们的开放标准理念
- **Go 原生** — 与我们的技术栈一致
- **模块化** — 组件可替换、可扩展
- **云端可选** — 离线优先，与我们的产品理念相同
- **硬件无关** — 不受限于特定机器人制造商
- **运行于嵌入式 Linux** — Raspberry Pi、NVIDIA Jetson、工业控制器

Viam RDK 是一个开源机器人开发套件，提供与语言无关的 SDK（Go、Python、TypeScript）以及模块化组件架构。

[Viam RDK 文档 →](https://docs.viam.com/)

### 嵌入式平台

| 平台 | 使用场景 |
|----------|----------|
| **Raspberry Pi** | 原型开发、边缘计算、自助终端 |
| **NVIDIA Jetson** | GPU 加速的视觉与机器学习 |
| **工业 Linux** | 工厂自动化、恶劣环境 |
| **定制 ARM 板** | 特定应用部署 |

同一套 Go 代码库。同一个 Viam RDK。从开发工作站部署到嵌入式边缘设备。

---

## 功能

| 功能 | Viam 服务 |
|------------|--------------|
| 带碰撞规避的运动规划 | motion service |
| 目标检测与分割 | vision service + 机器学习模型 |
| 基于深度传感器的点云 | Camera 组件 |
| 手眼标定 | Frame system |
| 结合 IMU 的 SLAM | slam service |

---

## 视觉与感知

### 相机抽象

Viam 提供统一的 `rdk:component:camera` API，适用于各种类型的相机——网络摄像头、IP 摄像头、LiDAR，以及 Intel RealSense D435i 等深度相机。

**工作原理：**

- **标准化 API** — 你的代码与 Viam 相机接口交互，而不是直接与硬件交互
- **内置驱动** — `viam-camera-realsense` 模块负责 RealSense 集成
- **硬件无关** — 更换相机无需修改应用逻辑
- **配置驱动** — 分辨率、传感器和数据流通过 JSON 配置，而非代码

可用方法：`GetImage()`、`GetImages()`、`GetPointCloud()`，以及相机内参。

### Intel RealSense D435i

RGB-D 深度感知，用于点云生成与空间感知。

[GitHub 上的 librealsense →](https://github.com/realsenseai/librealsense)

### 机器学习流水线

- YOLOv8 目标检测与分割
- 与 Viam vision service 集成
- 在边缘设备上进行实时推理

---

## 硬件抽象

**核心价值主张：改配置，而不是改代码。**

- **同样的软件，更大的机械臂** — 将 xArm 配置替换为 UR5e 或 KUKA，重新部署即可
- **将 xArm 用作示教器** — 在 xArm 上演示动作，更大的机械臂随之模仿
- **多机器人协同** — xArm 处理小型零件，KUKA 负责重物搬运
- **数字孪生开发** — 基于 xArm 开发，基于 KUKA 运动学进行仿真

---

## 支持的机械臂

| 型号 | 负载 | 工作半径 | 使用场景 |
|-------|---------|-------|----------|
| xArm 6 | 5kg | 700mm | 开发 |
| UR5e | 5kg | 850mm | 生产 |
| KUKA KR6 | 6kg | 900mm | 小型零件 |
| KUKA KR10 | 10kg | 900-1100mm | 中型装配 |
| KUKA KR16 | 16kg | 1600mm | 焊接、码垛 |
| KUKA KR30 | 30kg | 2000mm+ | 重型搬运 |

[viam-kuka 模块 →](https://github.com/viam-soleng/viam-kuka)

---

## 架构

```
SOFTWARE STACK
┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐
│ Vision   │  │ Motion   │  │   ML     │  │ Business │
│ Pipeline │  │ Planning │  │ Models   │  │  Logic   │
└────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘
     └─────────────┴─────────────┴─────────────┘
                         │
                  Viam Arm API
              (rdk:component:arm)
                         │
         ┌───────────────┼───────────────┐
         ▼               ▼               ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │  xArm 6  │   │  UR5e    │   │  KUKA    │
    │  (dev)   │   │  (prod)  │   │  (heavy) │
    │ 5kg/700mm│   │ 5kg/850mm│   │ 30kg+    │
    └──────────┘   └──────────┘   └──────────┘
```

---

## 工业集成

**I/O 集成：**

Modbus → PLC，用于工业自动化集成。在同一运行时中控制输送机、传感器、执行器和安全系统。

---

## 整体协作方式

| 层级 | 技术 | 用途 |
|-------|------------|---------|
| 设计 | Spatial | 3D 工作单元设计与仿真 |
| 运行时 | Viam RDK | 机器人控制与操作 |
| 视觉 | RealSense + YOLOv8 | 感知与目标检测 |
| 机械臂 | xArm / UR / KUKA | 物理操作 |
| I/O | Modbus / PLC | 工业集成 |

---

## 了解更多

- [技术概览 →](/technology/)
- [Spatial 平台 →](/platform/spatial/)
- [联系我们 →](/contact/)
