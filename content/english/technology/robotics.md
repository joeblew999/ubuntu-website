---
title: "Robotics"
meta_title: "Robotics Stack | Ubuntu Software"
description: "Our robotics architecture built on Viam RDK - from design to deployment with hardware abstraction, computer vision, and industrial integration."
image: "/images/spatial.svg"
draft: false
---

## Robotics Stack

From design to deployment. Spatial provides the 3D design and simulation layer. Viam RDK provides the runtime to control and operate actual robots.

---

## Embedded Linux Runtime

Robotics runs on Linux. Our decades of Linux and embedded systems experience means we build for the real world—factory floors, outdoor environments, and resource-constrained edge devices.

### Viam RDK Runtime

**Why we chose Viam RDK:**

- **Open source** — Aligns with our open standards philosophy
- **Go-native** — Matches our technology stack
- **Modular** — Components can be swapped and extended
- **Cloud-optional** — Works offline-first, like our products
- **Hardware-agnostic** — Not locked to specific robot manufacturers
- **Runs on embedded Linux** — Raspberry Pi, NVIDIA Jetson, industrial controllers

Viam RDK is an open-source robotics development kit with language-agnostic SDKs (Go, Python, TypeScript) and a modular component architecture.

[Viam RDK Documentation →](https://docs.viam.com/)

### Embedded Platforms

| Platform | Use Case |
|----------|----------|
| **Raspberry Pi** | Prototyping, edge compute, kiosks |
| **NVIDIA Jetson** | GPU-accelerated vision and ML |
| **Industrial Linux** | Factory automation, harsh environments |
| **Custom ARM boards** | Application-specific deployments |

Same Go codebase. Same Viam RDK. Deploys from development workstation to embedded edge device.

---

## Capabilities

| Capability | Viam Service |
|------------|--------------|
| Motion planning with collision avoidance | motion service |
| Object detection and segmentation | vision service + ML models |
| Point cloud from depth sensors | Camera component |
| Hand-eye calibration | Frame system |
| SLAM with IMU | slam service |

---

## Vision & Perception

### Camera Abstraction

Viam provides a unified `rdk:component:camera` API that works across diverse camera types—webcams, IP cameras, LiDAR, and depth cameras like the Intel RealSense D435i.

**How it works:**

- **Standardized API** — Your code talks to the Viam camera interface, not the hardware
- **Built-in drivers** — The `viam-camera-realsense` module handles RealSense integration
- **Hardware agnostic** — Swap cameras without changing application logic
- **Configuration-driven** — Resolution, sensors, and streams configured via JSON, not code

Available methods: `GetImage()`, `GetImages()`, `GetPointCloud()`, plus intrinsic camera parameters.

### Intel RealSense D435i

RGB-D depth sensing for point cloud generation and spatial perception.

[librealsense on GitHub →](https://github.com/realsenseai/librealsense)

### ML Pipeline

- YOLOv8 object detection and segmentation
- Integration with Viam vision service
- Real-time inference on edge devices

---

## Hardware Abstraction

**The key value proposition: config change, not code change.**

- **Same software, bigger arm** — Swap xArm config for UR5e or KUKA, redeploy
- **xArm as teaching pendant** — Demonstrate motions on xArm, larger arm mimics
- **Multi-robot coordination** — xArm handles small parts, KUKA handles heavy lifting
- **Digital twin development** — Develop against xArm, simulate against KUKA kinematics

---

## Supported Arms

| Model | Payload | Reach | Use Case |
|-------|---------|-------|----------|
| xArm 6 | 5kg | 700mm | Development |
| UR5e | 5kg | 850mm | Production |
| KUKA KR6 | 6kg | 900mm | Small parts |
| KUKA KR10 | 10kg | 900-1100mm | Medium assembly |
| KUKA KR16 | 16kg | 1600mm | Welding, palletizing |
| KUKA KR30 | 30kg | 2000mm+ | Heavy handling |

[viam-kuka module →](https://github.com/viam-soleng/viam-kuka)

---

## Architecture

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

## Industrial Integration

**I/O Integration:**

Modbus → PLC for industrial automation integration. Control conveyors, sensors, actuators, and safety systems from the same runtime.

---

## How It Fits Together

| Layer | Technology | Purpose |
|-------|------------|---------|
| Design | Spatial | 3D work cell design and simulation |
| Runtime | Viam RDK | Robot control and operation |
| Vision | RealSense + YOLOv8 | Perception and object detection |
| Arms | xArm / UR / KUKA | Physical manipulation |
| I/O | Modbus / PLC | Industrial integration |

---

## Learn More

- [Technology Overview →](/technology/)
- [Spatial Platform →](/platform/spatial/)
- [Contact Us →](/contact/)
