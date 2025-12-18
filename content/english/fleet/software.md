---
title: "Software Stack"
meta_title: "Drone Fleet Software Architecture | Ubuntu Software"
description: "Software stack for 1,000-drone fleet: PX4 flight firmware, Ubuntu Server, JetPack, ROS 2, and Go-based Vehicle Gateway."
image: "/images/robotics.svg"
draft: false
---

## Software Architecture

Every software component runs open-source code with clear responsibilities. The stack is designed for remote management, over-the-air updates, and consistent behavior across 1,000+ vehicles.

---

## Flight Controller: Pixhawk

### PX4 v1.14 LTS

The flight controller runs **PX4 Autopilot v1.14 LTS**:

| Aspect | Details |
|--------|---------|
| **Version** | v1.14.x LTS (Long Term Support) |
| **Support Window** | Security and critical fixes for 2+ years |
| **Configuration** | Standardized parameters across fleet |
| **Updates** | Via QGroundControl or MAVLink |

**Why LTS:**

- **Stability** — Production-tested codebase, minimal churn
- **Security** — Critical vulnerabilities patched without feature changes
- **Predictability** — Known behavior across entire fleet
- **Support** — Commercial support available from Dronecode ecosystem

**Key Modules:**

- **EKF2** — Extended Kalman Filter for state estimation
- **Position Controller** — Multicopter position and velocity control
- **Navigator** — Mission execution, RTL, geofence
- **Commander** — State machine, failsafe handling
- **MAVLink** — Telemetry and command protocol

---

## Sensor Companion: Raspberry Pi

### Ubuntu Server 22.04 LTS

The Raspberry Pi runs **Ubuntu Server 22.04**:

| Aspect | Details |
|--------|---------|
| **Distribution** | Ubuntu Server 22.04 LTS |
| **Kernel** | Mainline Linux with Pi patches |
| **Init System** | systemd |
| **Support Window** | Until April 2027 |

**Why Ubuntu Server (not Raspberry Pi OS):**

- **Consistency** — Same OS on Pi and Jetson
- **Packaging** — Standard apt repositories, snaps available
- **LTS Support** — Security updates for 5 years
- **Tooling** — Familiar administration for ops teams

### Core Services

| Service | Purpose |
|---------|---------|
| **mavlink-router** | MAVLink message routing between FCU and Jetson |
| **sensor-agent** | Environmental sensor data collection |
| **gpsd** | GPS/GNSS daemon for position data |
| **chrony** | NTP time synchronization |

**mavlink-router** is critical—it forwards MAVLink traffic from the Pixhawk's serial port to the Jetson over UDP, enabling the Vehicle Gateway to process telemetry.

---

## AI Companion: Jetson

### JetPack 6.x

The Jetson runs **NVIDIA JetPack 6.x**:

| Component | Version |
|-----------|---------|
| **L4T (Linux for Tegra)** | Ubuntu 22.04 based |
| **CUDA** | 12.x |
| **cuDNN** | 8.x |
| **TensorRT** | 8.x |
| **VPI (Vision Programming Interface)** | 3.x |

**Why JetPack 6:**

- **Latest Orin support** — Full hardware acceleration on Orin family
- **Ubuntu 22.04 base** — Matches Pi for consistent operations
- **Container support** — NVIDIA Container Runtime for isolated workloads
- **OTA updates** — NVIDIA tools for fleet-wide firmware updates

### CUDA + TensorRT

AI inference runs through the **TensorRT** runtime:

- **Model optimization** — INT8/FP16 quantization for edge deployment
- **Batch inference** — Process multiple frames efficiently
- **Low latency** — Direct GPU memory access

Typical workloads:

- Object detection (YOLO, SSD)
- Semantic segmentation
- Optical flow
- Depth estimation

### ROS 2 Humble

**ROS 2 Humble Hawksbill** provides the robotics middleware:

| Aspect | Details |
|--------|---------|
| **Distribution** | Humble Hawksbill |
| **Support Window** | Until May 2027 |
| **DDS** | CycloneDDS (default) or FastDDS |
| **Build System** | colcon |

**Why ROS 2 Humble:**

- **LTS release** — 5-year support matches Ubuntu/JetPack
- **Real-time capable** — Deterministic execution paths available
- **Sensor integration** — Extensive driver packages
- **Simulation** — Gazebo integration for testing

**Key Packages:**

| Package | Purpose |
|---------|---------|
| **px4_msgs** | MAVLink message definitions for ROS 2 |
| **px4_ros_com** | PX4-ROS 2 bridge (uXRCE-DDS) |
| **image_transport** | Camera data handling |
| **tf2** | Coordinate frame transforms |

---

## Vehicle Gateway

### Go Implementation

The **Vehicle Gateway** is a custom Go service running on each Jetson:

| Aspect | Details |
|--------|---------|
| **Language** | Go 1.21+ |
| **Dependencies** | nats.go, mavlink library |
| **Deployment** | systemd service or container |
| **Configuration** | Environment variables + config file |

**Why Go:**

- **Single binary** — No runtime dependencies, simple deployment
- **Performance** — Low latency, efficient memory usage
- **Concurrency** — Goroutines for parallel MAVLink/NATS handling
- **NATS ecosystem** — First-class nats.go client library

**Responsibilities:**

1. **MAVLink Ingest** — Receive and parse MAVLink messages from mavlink-router
2. **State Downsampling** — Reduce telemetry rate for WAN transmission
3. **Event Extraction** — Generate events from state transitions
4. **Command Execution** — Receive NATS commands, validate, send to FCU
5. **Shadow Reconciliation** — Sync desired/reported state with fleet backend

[Gateway Details →]({{< relref "/fleet/gateway" >}})

---

## Control Plane

### NATS + JetStream

Fleet coordination runs on **NATS** with **JetStream** persistence:

| Component | Role |
|-----------|------|
| **NATS Server** | Core pub/sub messaging |
| **JetStream** | Stream persistence, KV stores |
| **Leaf Nodes** | Vehicle-local NATS instances |
| **Hub Clusters** | Regional aggregation points |

**Why NATS:**

- **Performance** — Millions of messages per second per node
- **Simplicity** — Single binary, minimal configuration
- **Persistence** — JetStream for durable streams and KV
- **Topology** — Leaf nodes, clusters, superclusters for any scale

**Fleet Topology:**

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

Each vehicle runs a **leaf node** that connects to its regional hub. The leaf node handles local pub/sub and store-and-forward when disconnected.

[NATS Topology Details →]({{< relref "/fleet/nats-topology" >}})

---

## Software Update Strategy

### Fleet-Wide OTA

| Component | Update Mechanism |
|-----------|------------------|
| **PX4** | MAVLink firmware upload via QGC or custom tooling |
| **Pi/Ubuntu** | apt + unattended-upgrades, or custom image deployment |
| **Jetson/JetPack** | NVIDIA OTA tools, container updates |
| **Vehicle Gateway** | Binary replacement via systemd, or container update |
| **NATS Leaf** | Binary update with graceful restart |

**Principles:**

- **Staged rollout** — Update subset, verify, expand
- **Rollback capability** — Previous version always available
- **Health checks** — Automated verification post-update
- **Minimal downtime** — Updates during maintenance windows

---

## Summary

| Layer | Software | Version |
|-------|----------|---------|
| **Flight Controller** | PX4 | v1.14 LTS |
| **Sensor Companion** | Ubuntu Server | 22.04 LTS |
| **Sensor Companion** | mavlink-router | Latest |
| **AI Companion** | JetPack | 6.x |
| **AI Companion** | CUDA/TensorRT | 12.x/8.x |
| **AI Companion** | ROS 2 | Humble |
| **AI Companion** | Vehicle Gateway | Go 1.21+ |
| **Control Plane** | NATS + JetStream | Latest |

---

## Next

[NATS Topology →]({{< relref "/fleet/nats-topology" >}})
