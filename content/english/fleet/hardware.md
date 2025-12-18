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

**Purchase:** [Holybro Store](https://holybro.com/products/x500-v2-kits) | [GetFPV](https://www.getfpv.com/holybro-x500-v2-arf-kit.html) | [NewBeeDrone](https://newbeedrone.com/products/holybro-x500-v2-kits)

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
| **SKU** | 11070 (Pixhawk 6X Pro) |
| **Processor** | STM32H753 (480MHz Cortex-M7) |
| **IMU** | Triple redundant (ICM-42688-P, ICM-45686, BMI088) |
| **Barometer** | Dual (MS5611, ICP-20100) |
| **Magnetometer** | IST8310 |
| **Interfaces** | 3x CAN, 6x UART, SPI, I2C, PWM |
| **Cable Set SKU** | 18119 (Standard Baseboard V2) |

**Purchase:** [Holybro Store](https://holybro.com/collections/pixhawk-6x-series) | [Documentation](https://docs.holybro.com/autopilot/pixhawk-6x)

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
| **Model** | Compute Module 4 (CM4) |
| **SKU** | CM4104032 (4GB RAM, 32GB eMMC, WiFi) |
| **Processor** | Quad-core Cortex-A72 @ 1.8GHz |
| **RAM** | 4GB or 8GB |
| **Storage** | 32GB+ eMMC (or microSD for Pi 4) |
| **Power** | 5V @ 3A typical |

**Purchase:** [Raspberry Pi](https://www.raspberrypi.com/products/compute-module-4/) | [Digi-Key](https://www.digikey.com/en/products/filter/embedded-system-on-module/660?s=N4IgTCBcDaIIwFYwA4C0SA7ALvAxgYwEsAnAcw1wEMAXfAIxAF0BfIA) | [Mouser](https://www.mouser.com/c/embedded-solutions/computing/single-board-computers-sbcs/?m=Raspberry%20Pi)

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
| **Jetson Orin Nano Super** | 945-13766-0000-000 | 1024 CUDA | 67 TOPS | 7-25W |
| **Jetson Orin NX** | 900-13767-0000-000 | 1024 CUDA | 100 TOPS | 10-25W |
| **Jetson AGX Orin** | 900-13701-0000-000 | 2048 CUDA | 275 TOPS | 15-60W |
| **Jetson Xavier NX** | 900-83668-0000-000 | 384 CUDA | 21 TOPS | 10-20W |

**Purchase:** [NVIDIA Store](https://www.nvidia.com/en-us/autonomous-machines/embedded-systems/jetson-orin/) | [SparkFun](https://www.sparkfun.com/nvidia-jetson-orin-nano-developer-kit.html) | [Seeed Studio](https://www.seeedstudio.com/NVIDIAr-Jetson-Orintm-Nano-Developer-Kit-p-5617.html)

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
| **Protocol** | Open-source (CRSF) |
| **Receiver** | BetaFPV ELRS Nano (~0.7g) |

**Purchase:** [BetaFPV](https://betafpv.com/products/elrs-nano-receiver) | [GetFPV](https://www.getfpv.com/betafpv-expresslrs-nano-915mhz-receiver.html) | [Pyrodrone](https://pyrodrone.com/products/betafpv-elrs-nano-receiver-915mhz)

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

**Purchase:** [Holybro Store](https://holybro.com/products/sik-telemetry-radio-v3) | [SparkFun](https://www.sparkfun.com/sik-telemetry-radio-v3-915mhz-100mw.html) | [GetFPV](https://www.getfpv.com/holybro-sik-telemetry-radio-v3-500mw-915mhz-2pcs.html)

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
| **Modem** | Quectel RM520N-GL |
| **SKU** | RM520NGLAA-M20-SGASA |
| **Bands** | 4G LTE + 5G NR Sub-6GHz |
| **SIM** | eSIM with OTA provisioning |
| **Interface** | M.2 (USB 3.0 / PCIe) |
| **Form Factor** | M.2 NGFF |

**Purchase:** [Quectel](https://www.quectel.com/product/5g-rm520n-gl) | [Digi-Key](https://www.digikey.com/en/products/filter/rf-transceiver-modules/872?s=N4IgTCBcDaICwFYEFoDMBGBAGAHNqBOIAugL5A) | [Mouser](https://www.mouser.com/c/rf-wireless/rf-transceiver-modules/?m=Quectel)

**eSIM Benefits:**

- No physical SIM swapping across large fleet
- Remote carrier provisioning
- Carrier switching for coverage optimization
- Centralized subscription management

See [Sensing platform](/platform/sensing/) for details on 5G/LTE connectivity architecture.

---

## Summary

| Component | Selection | SKU | Rationale |
|-----------|-----------|-----|-----------|
| **Frame** | Holybro X500 V2 ARF | 30125 | Proven, available, maintainable |
| **FCU** | Pixhawk 6X Pro | 11070 | Redundancy, PX4 LTS support |
| **Sensor Companion** | Raspberry Pi CM4 | CM4104032 | Cost-effective, GPIO-rich |
| **AI Companion** | Jetson Orin Nano Super | 945-13766-0000-000 | GPU inference, CUDA ecosystem |
| **RC Receiver** | BetaFPV ELRS Nano 915MHz | — | Open-source, long range, 0.7g |
| **Telemetry** | Holybro SiK Radio V3 | 17012/17013 | Backup link, debugging |
| **Power** | 4S LiPo 5000mAh | — | Fleet standardization |
| **Cellular** | Quectel RM520N-GL | RM520NGLAA-M20-SGASA | OTA provisioning, 5G coverage |

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

## Next

[Software Stack →]({{< relref "/fleet/software" >}})
