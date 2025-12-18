---
title: "Drone Platform"
meta_title: "Drone Fleet Hardware & Protocols | Ubuntu Software"
description: "Production hardware specifications for drone fleets: Holybro X500 airframe, Pixhawk 6X flight controller, PX4, MAVLink protocol, and NVIDIA Jetson companion computers."
image: "/images/robotics.svg"
draft: false
---

## Hardware for Fleet-Scale Drone Operations

Fleet operations demand hardware that balances capability with maintainability. Every component in this stack was chosen for production reliability, parts availability, and serviceability in the field.

---

## Protocol: MAVLink

Drones in our fleet communicate using **MAVLink** (Micro Air Vehicle Link):

| Aspect | Details |
|--------|---------|
| **Version** | MAVLink 2.0 |
| **Transport** | Serial (UART), UDP, TCP |
| **Message Rate** | 1-50Hz depending on message type |
| **Encryption** | MAVLink 2.0 signing (optional) |

**MAVLink provides:**
- Standardized telemetry (position, attitude, battery, sensors)
- Command protocol (arm, takeoff, waypoints, RTL)
- Parameter management
- File transfer (logs, missions)

The [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) translates MAVLink messages to NATS subjects, enabling fleet-wide telemetry aggregation and command distribution.

---

## Airframe

### Holybro X500 V2 ARF

The **Holybro X500 V2** provides an ideal platform for fleet deployment:

| Specification | Value |
|---------------|-------|
| **Wheelbase** | 500mm |
| **Frame Weight** | ~410g |
| **Max Takeoff Weight** | ~2kg |
| **Flight Time** | 15-20 min (with standard payload) |
| **Motor Mount** | 16x16mm / 19x19mm |

**Why X500:**

- **Proven design** — Thousands deployed worldwide, extensive community knowledge
- **Parts availability** — Arms, landing gear, and hardware readily available
- **Maintenance-friendly** — Modular construction, field-serviceable
- **Payload flexibility** — Sufficient capacity for dual companion computers plus sensors

The ARF (Almost Ready to Fly) kit includes frame, motors, ESCs, and propellers—reducing assembly variability across a large fleet.

---

## Flight Controller

### Pixhawk 6X

The **Pixhawk 6X** running **PX4 v1.14 LTS** handles all flight-critical functions:

| Specification | Value |
|---------------|-------|
| **Processor** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Triple redundant (ICM-42688-P, ICM-45686, BMI088) |
| **Barometer** | Dual (MS5611, ICP-20100) |
| **Magnetometer** | IST8310 |
| **Interfaces** | 3x CAN, 6x UART, SPI, I2C, PWM |

**Why Pixhawk 6X:**

- **Redundancy** — Triple IMU, dual barometer for sensor voting
- **PX4 LTS support** — Long-term stability, security patches
- **FMUv6X standard** — Interchangeable with other compliant boards
- **Vibration isolation** — Built-in IMU dampening

The flight controller handles:

- Attitude estimation and control
- Position hold and navigation
- Failsafe behaviors (RTL, land, geofence)
- RC input processing
- MAVLink telemetry

---

## Companion Computers

Fleet vehicles run **dual companion computers** with distinct responsibilities:

### Raspberry Pi 4 / CM4 (Sensor Companion)

Handles lightweight sensor integration and data collection:

| Specification | Value |
|---------------|-------|
| **Processor** | Quad-core Cortex-A72 @ 1.8GHz |
| **RAM** | 4GB or 8GB |
| **Storage** | 32GB+ microSD or eMMC |
| **Power** | 5V @ 3A typical |

**Responsibilities:**

- Environmental sensor drivers (temperature, humidity, air quality)
- GPS/GNSS data logging
- Camera capture (non-AI workloads)
- MAVLink routing to Jetson
- Store-and-forward when Jetson is offline

### NVIDIA Jetson (AI Companion)

Handles compute-intensive workloads:

| Model | GPU Cores | AI Performance | Power |
|-------|-----------|----------------|-------|
| **Jetson Orin Nano** | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin NX** | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 2048 CUDA | 275 TOPS | 15-60W |
| **Jetson Xavier NX** | 384 CUDA | 21 TOPS | 10-20W |

**Responsibilities:**

- Computer vision (object detection, tracking)
- Visual-inertial odometry
- Path planning and obstacle avoidance
- Vehicle Gateway (NATS client)
- NATS leaf node

Fleet deployments typically use **Jetson Orin Nano** for cost-effective inference or **Orin NX** for demanding perception workloads.

---

## Radio Systems

### RC Control: ExpressLRS

**ExpressLRS** provides the pilot override link:

| Specification | Value |
|---------------|-------|
| **Frequency** | 915MHz (US) / 868MHz (EU) |
| **Range** | 10km+ (depending on power/antenna) |
| **Latency** | <5ms |
| **Protocol** | Open-source |

**Why ExpressLRS:**

- **Open-source** — No vendor lock-in, community-driven development
- **Range** — Reliable link for visual line-of-sight operations
- **Low latency** — Responsive manual control when needed
- **Cost** — Affordable receivers for fleet scale

### Telemetry: SiK Radio

**SiK 915 MHz** radios provide ground station telemetry:

| Specification | Value |
|---------------|-------|
| **Frequency** | 915MHz (US) / 433MHz (EU) |
| **Range** | 1-2km typical |
| **Data Rate** | Up to 250kbps |
| **Interface** | Serial (UART) |

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
| **Configuration** | 4S (14.8V nominal) |
| **Capacity** | 5000-6000mAh typical |
| **Discharge Rate** | 30C+ |
| **Connector** | XT60 |

**Fleet Considerations:**

- Standardize on single battery configuration for logistics
- Use smart batteries with telemetry when available
- Implement battery rotation and health tracking
- Plan for 3:1 battery-to-vehicle ratio for continuous operations

### Power Distribution

- **PDB/PMS** — Power management board for companion computers
- **BEC** — 5V regulated supply for Pi, servos
- **Jetson power** — Direct from battery through regulator (12V typical)

---

## Connectivity

### Cellular: 4G/5G + eSIM

Each Jetson connects via cellular modem:

| Specification | Value |
|---------------|-------|
| **Modem** | Quectel RM520N-GL or similar |
| **Bands** | 4G LTE + 5G NR |
| **SIM** | eSIM with OTA provisioning |
| **Interface** | USB 3.0 or M.2 |

**eSIM Benefits:**

- No physical SIM swapping across large fleet
- Remote carrier provisioning
- Carrier switching for coverage optimization
- Centralized subscription management

---

## Safety Model

Drone-specific safety features are enforced at multiple levels:

| Level | Mechanism | Behavior |
|-------|-----------|----------|
| **RC Override** | ExpressLRS | Pilot always has manual control |
| **Flight Controller** | PX4 failsafes | RTL, land, geofence enforcement |
| **Companion** | Vehicle Gateway | Command validation, safety checks |
| **Network** | NATS never in control loop | Monitoring only, not flight-critical |

See [Safety Model]({{< relref "/fleet/safety" >}}) for complete safety architecture.

---

## Summary

| Component | Selection | Rationale |
|-----------|-----------|-----------|
| **Frame** | Holybro X500 V2 ARF | Proven, available, maintainable |
| **FCU** | Pixhawk 6X | Redundancy, PX4 LTS support |
| **Sensor Companion** | Raspberry Pi 4/CM4 | Cost-effective, GPIO-rich |
| **AI Companion** | Jetson Orin/Xavier | GPU inference, CUDA ecosystem |
| **RC** | ExpressLRS | Open-source, long range |
| **Telemetry** | SiK 915MHz | Backup link, debugging |
| **Power** | 4S LiPo 5000mAh | Fleet standardization |
| **Cellular** | 4G/5G + eSIM | OTA provisioning, coverage |
| **Protocol** | MAVLink 2.0 | Industry standard, PX4 native |

---

## Where to Buy

We don't sell hardware—our partners do. **Certified drone shops** can supply complete kits with fleet infrastructure included.

**For drone buyers:**
Purchase from a [certified partner]({{< relref "/partners" >}}) and get:
- Pre-configured hardware ready to connect
- Credentials for our managed infrastructure
- Local support from your shop

**For drone shops:**
[Join our partner program]({{< relref "/partners/become-partner" >}}) to offer fleet infrastructure with your hardware sales.

---

## Related Documentation

- [Supported Platforms]({{< relref "/fleet/platforms" >}}) — Overview of all vehicle types
- [Ground Vehicles]({{< relref "/fleet/platforms/ground" >}}) — Cars, trucks, and AGVs
- [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) — MAVLink-to-NATS bridge
- [Safety Model]({{< relref "/fleet/safety" >}}) — Failsafe architecture
- [Software Stack]({{< relref "/fleet/software" >}}) — Operating systems and middleware
