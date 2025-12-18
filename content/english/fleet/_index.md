---
title: "Autonomous Vehicle Fleet Architecture"
meta_title: "Fleet Architecture for Drones, Cars & Trucks | Ubuntu Software"
description: "Production-grade reference architecture for managing autonomous vehicle fleets—drones, cars, trucks, AGVs—using NATS JetStream for fleet-scale digital twinning."
image: "/images/robotics.svg"
draft: false
---

## Production-Grade Fleet Management at Scale

This reference architecture demonstrates how to deploy and manage **fleets of autonomous vehicles**—drones, cars, trucks, AGVs—using open-source components and proven patterns. Every design decision targets real-world production requirements: reliability, observability, and graceful degradation.

The core architecture is **vehicle-agnostic**. NATS JetStream, digital twin patterns, and the authorization model work identically whether you're managing drones or delivery trucks. Only the edge hardware and protocols change.

---

## The Challenge

Fleet-scale operations require solving problems that don't exist at small scale:

- **Thousands of concurrent telemetry streams** — Every vehicle reporting position, velocity, health, and sensor data
- **Command and control at scale** — Sending instructions to specific vehicles or groups without flooding the network
- **Digital twin synchronization** — Maintaining accurate state representation for every vehicle in real-time
- **Offline resilience** — Vehicles that continue operating when connectivity drops
- **Safety guarantees** — Ensuring network failures never compromise vehicle safety

Traditional approaches—direct protocol connections, centralized databases, polling architectures—collapse under these requirements.

---

## Our Approach

The architecture combines three proven technologies:

| Layer | Technology | Role |
|-------|------------|------|
| **Vehicle Control** | Platform-specific (PX4, ECU, ROS 2) | Autonomous operation, failsafes, manual override |
| **Edge Computing** | Jetson + Industrial PC | Sensor processing, AI inference, local decisions |
| **Fleet Messaging** | NATS JetStream | Pub/sub, persistence, digital twin state |

Each layer operates independently. Network failures degrade gracefully—vehicles continue operating, edge computers continue processing, and state synchronizes when connectivity returns.

---

## Architecture Components

### Supported Platforms

The architecture supports multiple vehicle types with platform-specific hardware and protocols:

| Platform | Protocol | Control System | Use Cases |
|----------|----------|----------------|-----------|
| **Drones** | MAVLink | PX4/ArduPilot | Inspection, mapping, delivery |
| **Ground Vehicles** | CAN bus / J1939 | ECU / Autoware | Logistics, mining, agriculture |
| **ROS 2 Robots** | ROS 2 topics | Custom | Warehouse, research, AGVs |

[Supported Platforms →]({{< relref "/fleet/platforms" >}})

### Hardware Stack (Drones)

Standard hardware choices that balance capability, availability, and maintainability:

- **Airframe**: Holybro X500 V2 ARF — proven platform, excellent parts availability
- **Flight Controller**: Pixhawk 6X running PX4 v1.14 LTS
- **Sensor Companion**: Raspberry Pi 4 / CM4 for lightweight sensor integration
- **AI Companion**: NVIDIA Jetson (Orin/Xavier/Nano) for computer vision and inference

[Drone Hardware Details →]({{< relref "/fleet/platforms/drones" >}}) | [Ground Vehicle Details →]({{< relref "/fleet/platforms/ground" >}})

### Software Stack

Open-source software across every layer:

- **PX4 v1.14 LTS** — Flight control with proven stability
- **Ubuntu Server 22.04** — Consistent Linux environment on all companions
- **ROS 2 Humble** — Robot middleware for sensor integration
- **Go** — Vehicle Gateway implementation

[Software Details →]({{< relref "/fleet/software" >}})

### NATS Architecture

Hierarchical messaging topology designed for WAN deployment:

- **Leaf nodes** on each vehicle — local pub/sub, store-and-forward
- **Regional hub clusters** — aggregate vehicles by geographic region
- **Global mirroring** — cross-region replication when required

[NATS Topology →]({{< relref "/fleet/nats-topology" >}})

### Digital Twin Design

JetStream streams and KV stores maintain fleet state:

- **Subject hierarchy** for routing and filtering
- **State streams** for telemetry rollup
- **Event streams** for audit trails
- **Shadow stores** for desired/reported state reconciliation

[Subject Naming →]({{< relref "/fleet/subjects" >}}) | [Stream Configuration →]({{< relref "/fleet/streams" >}})

### Vehicle Gateway

Service running on each vehicle's edge computer that bridges native protocols to NATS:

- Protocol translation (MAVLink, CAN bus, ROS 2)
- State downsampling and aggregation
- Event extraction from telemetry
- Command execution with policy enforcement
- Shadow state reconciliation

[Gateway Design →]({{< relref "/fleet/gateway" >}})

### Safety Model

Network connectivity is never trusted for vehicle safety:

- **Manual override is primary authority** — Operator always has control (RC for drones, steering/e-stop for ground)
- **Vehicle control system enforces failsafes** — RTL/stop-in-place, geofence
- **NATS is never in the control loop** — Monitoring and coordination only
- **Graceful degradation** — Loss of AI or network triggers safe modes

[Safety Details →]({{< relref "/fleet/safety" >}})

---

## Why This Architecture

### Proven at Scale

Every component has been deployed in production environments:

- **NATS** powers Synadia's global messaging infrastructure
- **PX4** flies on thousands of commercial drones worldwide
- **Jetson** runs inference in autonomous vehicles and robots

### Open Standards

No vendor lock-in:

- MAVLink is an open protocol with multiple implementations
- NATS is open-source with commercial support available
- PX4 runs on multiple flight controller hardware platforms

### Your Infrastructure or Ours

NATS JetStream is **100% open source** (Apache 2.0). Run it yourself or connect to our managed infrastructure—same protocol, same code, your choice.

- **Self-hosted** — Deploy on your infrastructure with our reference configs
- **Managed** — Connect your leaf nodes to our regional hubs

Free for small fleets. Scales with you.

[Deployment Options →]({{< relref "/fleet/nats-topology#deployment-options" >}})

### Operational Reality

Designed for the real world:

- Vehicles can be serviced without specialized tools
- Software updates deploy over-the-air
- Telemetry data feeds standard observability stacks
- Fleet state exports to existing enterprise systems

---

## Learn More

Explore each component in detail:

| Section | Description |
|---------|-------------|
| [Supported Platforms]({{< relref "/fleet/platforms" >}}) | Drones, ground vehicles, AGVs |
| [Drone Platform]({{< relref "/fleet/platforms/drones" >}}) | Airframe, flight controller, MAVLink |
| [Ground Vehicles]({{< relref "/fleet/platforms/ground" >}}) | Cars, trucks, CAN bus, J1939 |
| [Software]({{< relref "/fleet/software" >}}) | Operating systems, middleware, applications |
| [NATS Topology]({{< relref "/fleet/nats-topology" >}}) | Leaf nodes, hubs, WAN connectivity |
| [Subject Naming]({{< relref "/fleet/subjects" >}}) | Hierarchical subject structure |
| [Stream Configuration]({{< relref "/fleet/streams" >}}) | JetStream setup for digital twins |
| [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) | Protocol-to-NATS bridge |
| [Authorization & Grants]({{< relref "/fleet/authorization" >}}) | Decentralized security, third-party access |
| [Safety Model]({{< relref "/fleet/safety" >}}) | Failsafes and graceful degradation |

---

## Get Started

Building an autonomous vehicle fleet? Deploying drones, trucks, or AGVs at scale? We can help with architecture, implementation, and operations.

[Contact Us →](/contact)
