---
title: "Hardware Stack"
meta_title: "Drone Fleet Hardware Specifications | Ubuntu Software"
description: "Production hardware specifications for 1,000-drone fleet: Holybro X500 airframe, Pixhawk 6X flight controller, Raspberry Pi and Jetson companion computers."
image: "/images/robotics.svg"
draft: false
aliases:
  - "/fleet/hardware/"
---

{{< notice "info" >}}
**Looking for other vehicle types?** This page covers drone hardware. For ground vehicles (cars, trucks, AGVs), see [Ground Vehicle Platform]({{< relref "/fleet/platforms/ground" >}}). For a complete overview, see [Supported Platforms]({{< relref "/fleet/platforms" >}}).
{{< /notice >}}

## Buy from a Local Shop

Drone building has a learning curve. **Buy from a shop that can help you**—they stock the hardware, we provide the fleet software.

When you purchase from a partner shop, you get:
- Pre-tested hardware ready to connect
- Local support from people who know drones
- Access to our managed fleet infrastructure

---

## Shops with Physical Locations

Visit a shop in person for hands-on help with your build:

{{< featured-shops >}}

**Are you a drone shop?** [Join our partner program]({{< relref "/partners/become-partner" >}}) to offer fleet infrastructure with your hardware sales.

---

## What You Need

| Kit | Components | Typical Sources |
|-----|------------|-----------------|
| **Complete Drone** | X500 frame, Pixhawk 6X, GPS, motors, ESCs, props | GetFPV, NewBeeDrone, Holybro |
| **Companion Computers** | Raspberry Pi 4 + Jetson Orin Nano | Digi-Key, Mouser, SparkFun |
| **Power & Batteries** | 4S LiPo 6000mAh, power module | GetFPV, Tattu, Amazon |
| **Radios** | ExpressLRS receiver, SiK telemetry | GetFPV, HappyModel, Holybro |
| **Networking** | Industrial ethernet switch, 5G modem | PUSR, Quectel |

[View complete component list ↓](#full-bill-of-materials)

---

## Hardware for Fleet-Scale Operations

Fleet operations demand hardware that balances capability with maintainability. Every component in this stack was chosen for production reliability, parts availability, and serviceability in the field.

---

## Airframe

### Holybro X500 V2 ARF

The **Holybro X500 V2** provides an ideal platform for fleet deployment:

| Specification | Value |
|---------------|-------|
| **SKU** | 30125 |
| **Wheelbase** | 500mm |
| **Frame Weight** | ~410g (610g with motors/ESCs) |
| **Max Takeoff Weight** | ~2kg |
| **Flight Time** | 15-20 min (with standard payload) |
| **Motor Mount** | 16x16mm / 19x19mm |
| **Included** | Motors (2216 KV920), ESCs (20A BLHeli S), Props (1045) |

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
| **Pixhawk 6X** | SKU 11073 |
| **Baseboard** | SKU 18117 |
| **Accessories** | SKU 15011 (cables & case) |
| **Processor** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Triple redundant (ICM-42688-P, ICM-45686, BMI088) |
| **Barometer** | Dual (MS5611, ICP-20100) |
| **Magnetometer** | IST8310 |
| **Interfaces** | 3x CAN, 6x UART, SPI, I2C, PWM |

### Flight Controller Accessories

| Component | SKU | Purpose |
|-----------|-----|---------|
| **PM06 V2 Power Module** | 15019 | Power sensing + 5V/3A BEC |
| **Holybro M9N GPS** | 12027 | GPS + compass module |
| **GPS Mast** | 12033 | Folding mast for interference reduction |
| **Safety Switch & Buzzer** | 12007 | Arm safety + audio alerts |
| **Vibration Mount** | 12010 | Anti-vibration damping |
| **Cable Kit** | 12055 | JST-GH/CAN cables |

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
| **Model** | Raspberry Pi 4 Model B |
| **SKU** | RPI4-MODEL-B-8GB |
| **Processor** | Quad-core Cortex-A72 @ 1.8GHz |
| **RAM** | 8GB |
| **Storage** | Industrial microSD 128GB (SDSDQED-128G-GNSIN) |
| **Power** | 5V @ 3A typical |

**Responsibilities:**

- Environmental sensor drivers (temperature, humidity, air quality)
- GPS/GNSS data logging
- Camera capture (non-AI workloads)
- MAVLink routing to Jetson
- Store-and-forward when Jetson is offline

### NVIDIA Jetson (AI Companion)

Handles compute-intensive workloads:

| Model | SKU | GPU Cores | AI Performance | Power |
|-------|-----|-----------|----------------|-------|
| **Jetson Orin Nano 8GB** ★ | 900-13767-0000-000 | 1024 CUDA | 40 TOPS | 7-15W |
| **Jetson Orin Nano Super** | 945-13766-0000-000 | 1024 CUDA | 67 TOPS | 7-25W |
| **Jetson Orin NX** | 900-13767-0040-000 | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 900-13701-0000-000 | 2048 CUDA | 275 TOPS | 15-60W |

★ = Reference BOM selection

| Storage | SKU | Capacity |
|---------|-----|----------|
| **NVMe SSD** | WD-SN530-256G | 256GB |

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
| **Receiver** | ExpressLRS EP1 |
| **SKU** | HM-EP1-2400 |
| **Frequency** | 2.4GHz |
| **Range** | 10km+ (depending on power/antenna) |
| **Latency** | <5ms |
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
<summary><strong>View all sellers by component (57 items)</strong></summary>

{{< bom-sellers >}}

</details>

---

## Next

[Software Stack →]({{< relref "/fleet/software" >}})
