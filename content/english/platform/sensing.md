---
title: "Sensing"
meta_title: "Sensing & Perception Platform | Ubuntu Software"
description: "Multi-sensor integration for spatial intelligence—LiDAR, cameras, and industrial sensors unified through a single edge agent with 5G/eSIM connectivity."
image: "/images/spatial.svg"
draft: false
---

## Spatial Sensing

Real-world perception for digital twins, robotics, and autonomous systems. Connect LiDAR, cameras, and industrial sensors to your spatial models.

---

## The Problem

Sensors produce data. Making sense of that data requires:

- **Context** — Where is the sensor? What's it looking at?
- **Fusion** — Combining multiple sensor streams into a coherent picture
- **Integration** — Connecting to design tools, not just dashboards

Most sensing solutions stop at data collection. We connect sensors to spatial models.

---

## Deployment Modes

Same edge agent. Same sensors. Different config.

| Mode | Platform | Use Case |
|------|----------|----------|
| **Aerial** | DJI enterprise drones | Surveying, inspection, mapping |
| **Ground** | Tripod, backpack | Interior scanning, construction |
| **Robot** | Viam RDK, ROS2 | Navigation, pick-and-place |
| **Fixed** | Permanent mount | Traffic, security, warehouse |

### GPS-Free Navigation

Visual-inertial odometry enables autonomous flight without GPS—critical for indoor environments, urban canyons, and GPS-denied zones.

{{< video src="videos/visual-inertial-nav-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*Real-time position tracking using only cameras and IMU. No external infrastructure required.*

{{< video src="videos/gps-free-vio-edited.mp4" width="100%" controls="true" loop="true" muted="true" autoplay="true" >}}

*Autonomous closed-loop flight using visual-inertial odometry.*

For detailed hardware specifications, see our [Drone Fleet Hardware →]({{< relref "/fleet/hardware" >}})

---

## Sensor Abstraction

**Hardware-agnostic by design.** Your code talks to our unified API, not individual sensor drivers.

### Supported Sensor Types

| Type | Examples |
|------|----------|
| **LiDAR** | Livox Mid-360, Avia |
| **RGB-D Cameras** | Intel RealSense, Luxonis OAK-D |
| **Position** | GPS/GNSS (u-blox RTK), IMU |
| **Industrial** | Modbus sensors, CAN bus |

Swap sensors without changing code. Configuration-driven, not code-driven.

---

## Edge Agent Architecture

Go binary that runs on your hardware—Raspberry Pi, Jetson, industrial Linux, or custom ARM.

| Capability | Description |
|------------|-------------|
| **Plugin system** | Add sensors via config, not code changes |
| **Local buffering** | Store-and-forward when offline |
| **Real-time streaming** | NATS JetStream to cloud |
| **Lightweight** | Single binary, no runtime dependencies |

---

## Connectivity

### 5G/LTE with eSIM OTA

No SIM swapping. No QR code scanning. Server-push provisioning.

- Modem ships with bootstrap profile
- Your platform triggers carrier profile download
- Switch carriers mid-deployment via API

Works for drones in the air, robots on the move, fixed installations in remote locations.

---

## Integration with Spatial

Sensors feed directly into your 3D models:

| Data Flow | Purpose |
|-----------|---------|
| Point clouds → Spatial model | Reality capture |
| GPS/IMU → Model positioning | Georeferencing |
| Environmental sensors → Twin | Live state updates |
| Industrial I/O → Automation | Closed-loop control |

Not just dashboards. Sensors in context.

---

## Built on Foundation

Sensing inherits all [Foundation](/platform/foundation/) capabilities automatically:

| Capability | What It Means |
|------------|---------------|
| **Offline-first** | Capture without internet, sync when connected |
| **Universal deployment** | Edge, mobile, desktop—same agent |
| **Self-sovereign** | Your sensors, your data, your servers |
| **Real-time sync** | Stream to multiple destinations simultaneously |

[Learn more about Foundation →](/platform/foundation/)

---

## Part of Something Bigger

Sensing is the perception layer of the Ubuntu Software platform.

For organizations that need 3D design and AI, our Spatial platform provides the design environment—with direct integration to your sensor data.

[Explore Spatial →](/platform/spatial/)

---

## Production Reference: 1,000-Drone Fleet

See how these sensing capabilities scale to production. Our reference architecture demonstrates a complete fleet deployment using the hardware and software stack described above—1,000 Holybro X500 drones with PX4, dual companion computers, and NATS JetStream for fleet-scale digital twinning.

[Explore Drone Fleet Architecture →]({{< relref "/fleet" >}})

---

## Build With Us

Deploying sensors? Building perception systems? Let's talk.

[Contact →](/contact)
