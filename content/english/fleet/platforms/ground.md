---
title: "Ground Vehicle Platform"
meta_title: "Ground Vehicle Fleet Hardware & Protocols | Ubuntu Software"
description: "Fleet architecture for cars, trucks, and AGVs: CAN bus integration, J1939, ROS 2, industrial compute, and vehicle-specific safety models."
image: "/images/robotics.svg"
draft: false
---

## Hardware for Ground Vehicle Fleets

Ground vehicles—cars, trucks, AGVs, and autonomous platforms—use the same fleet architecture as drones but with vehicle-appropriate hardware, protocols, and safety models.

---

## Protocol Options

Ground vehicles communicate using industry-standard protocols:

### CAN Bus

| Aspect | Details |
|--------|---------|
| **Standard** | CAN 2.0B / CAN FD |
| **Data Rate** | 500kbps (CAN) / 5Mbps (CAN FD) |
| **Message Format** | DBC files define signal encoding |
| **Physical** | Differential pair, robust in noisy environments |

**CAN bus provides:**
- Direct access to vehicle systems (engine, transmission, brakes)
- Real-time telemetry (speed, RPM, fuel, temperatures)
- Diagnostic trouble codes (OBD-II / J1939)
- Aftermarket device integration

### J1939 (Heavy Vehicles)

For trucks, buses, and heavy equipment:

| Aspect | Details |
|--------|---------|
| **Standard** | SAE J1939 |
| **Transport** | CAN 2.0B at 250kbps |
| **Addressing** | 29-bit identifiers with PGN/SPN |
| **Applications** | Engine, transmission, brakes, trailer |

**J1939 provides:**
- Standardized parameter groups (PGNs) across manufacturers
- Fleet-wide telemetry consistency
- Integration with telematics platforms
- Compliance with regulatory requirements

### ROS 2 (Autonomous Platforms)

For purpose-built autonomous vehicles:

| Aspect | Details |
|--------|---------|
| **Middleware** | ROS 2 Humble / Iron |
| **Transport** | DDS (Cyclone DDS, Fast DDS) |
| **Message Types** | Standard (sensor_msgs, geometry_msgs) + custom |
| **QoS** | Configurable reliability, durability |

**ROS 2 provides:**
- Native sensor integration (LiDAR, cameras, IMUs)
- Standard message types for robotics
- Direct topic bridging to NATS
- Simulation compatibility (Gazebo, CARLA)

---

## Hardware Stack

### Compute Platform

Ground vehicles typically use ruggedized industrial compute:

| Platform | Use Case | Notes |
|----------|----------|-------|
| **NVIDIA Jetson AGX Orin** | Autonomous vehicles | High-performance perception |
| **NVIDIA Jetson Orin NX** | Fleet vehicles | Balance of performance/cost |
| **Industrial PC** | Telematics-only | Lower cost, no GPU needed |
| **Neousys** | Rugged automotive | Wide temp range, ignition control |

**Compute responsibilities:**
- Vehicle Gateway (CAN/ROS 2 to NATS)
- NATS leaf node
- Edge AI inference (optional)
- Store-and-forward during connectivity loss

### CAN Interface

| Device | Interface | Notes |
|--------|-----------|-------|
| **PEAK PCAN-USB** | USB | Popular, well-supported |
| **Kvaser Leaf Light** | USB | Industrial-grade |
| **SocketCAN devices** | USB/PCIe | Linux-native support |
| **Jetson native CAN** | Built-in | Available on AGX Orin |

### Connectivity

| Method | Use Case |
|--------|----------|
| **4G/5G Cellular** | Primary fleet connectivity |
| **WiFi** | Depot/facility operations |
| **Ethernet** | Wired dock connectivity |
| **Satellite** | Remote/rural operations |

---

## Vehicle Types

### Passenger Vehicles (Cars)

| Component | Typical Selection |
|-----------|-------------------|
| **Protocol** | CAN bus (OBD-II port or direct) |
| **Compute** | Jetson Orin NX or industrial PC |
| **Power** | 12V vehicle system |
| **Connectivity** | 4G/5G cellular |
| **Safety** | Driver override, e-stop |

**Use cases:** Ride-sharing fleets, rental cars, delivery vehicles

### Commercial Vehicles (Trucks)

| Component | Typical Selection |
|-----------|-------------------|
| **Protocol** | J1939 (CAN-based) |
| **Compute** | Ruggedized industrial PC or Jetson |
| **Power** | 12V/24V vehicle system |
| **Connectivity** | 4G/5G + satellite backup |
| **Safety** | Driver override, ELD compliance |

**Use cases:** Long-haul logistics, delivery fleets, service vehicles

### Autonomous Ground Vehicles (AGVs)

| Component | Typical Selection |
|-----------|-------------------|
| **Protocol** | ROS 2 native or proprietary |
| **Compute** | Jetson AGX Orin |
| **Power** | 24V/48V battery system |
| **Connectivity** | WiFi (facility) + cellular (outdoor) |
| **Safety** | E-stop, safety-rated sensors |

**Use cases:** Warehouse logistics, mining, agriculture, port operations

---

## Safety Model

Ground vehicle safety differs from drones but follows the same principles:

### Authority Hierarchy

```
Manual Override (highest authority)
    │
    ├── Steering wheel / brake pedal (driver present)
    ├── E-stop button (AGVs)
    └── Remote operator console
    │
    ▼
Autonomous System (software control)
    │
    ├── Perception → Planning → Control
    └── Geofencing, speed limits, operational domain
    │
    ▼
Fleet Management (lowest authority)
    │
    └── Mission assignment, routing, monitoring
```

### Failsafe Behaviors

| Trigger | Drone Equivalent | Ground Vehicle Response |
|---------|------------------|------------------------|
| Loss of comms | Return-to-Launch | Stop in place / pull over |
| Sensor failure | Land immediately | Stop safely, alert operator |
| Geofence breach | Return to boundary | Stop at boundary |
| E-stop activated | Motor cutoff | Immediate braking |
| Low battery | RTL / land | Return to depot / stop safely |

### Network Safety

**Same principle as drones:** NATS is never in the control loop.

| What NATS Does | What NATS Doesn't Do |
|----------------|---------------------|
| Telemetry collection | Real-time steering |
| Mission assignment | Brake commands |
| Fleet coordination | Throttle control |
| Monitoring & alerting | Safety-critical decisions |

The vehicle's onboard systems handle all safety-critical functions. Network loss means loss of monitoring, not loss of vehicle control.

---

## Gateway Implementation

The Vehicle Gateway bridges vehicle-native protocols to NATS:

```
┌─────────────────────────────────────────────────────────────────┐
│                    GROUND VEHICLE GATEWAY                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   CAN Bus ─────► DBC Parser ─────► Telemetry Subjects           │
│                                                                  │
│   J1939 ───────► PGN Decoder ────► Telemetry Subjects           │
│                                                                  │
│   ROS 2 ───────► Topic Bridge ───► Telemetry Subjects           │
│                                                                  │
│   Commands ◄───────────────────── Command Subjects              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS JETSTREAM                              │
│   fleet.prod.veh.{vehicle_id}.state.*                           │
│   fleet.prod.veh.{vehicle_id}.evt.*                             │
│   fleet.prod.veh.{vehicle_id}.cmd.*                             │
└─────────────────────────────────────────────────────────────────┘
```

**Gateway responsibilities:**
- Protocol translation (CAN/J1939/ROS 2 → common telemetry format)
- State downsampling (reduce high-frequency CAN to manageable rate)
- Event extraction (detect state changes, thresholds)
- Command validation (safety checks before execution)
- Store-and-forward (buffer during connectivity loss)

---

## Subject Mapping

Ground vehicles use the same subject hierarchy as drones:

| Subject | Ground Vehicle Content |
|---------|----------------------|
| `fleet.prod.veh.{id}.state.pos` | GPS position, heading |
| `fleet.prod.veh.{id}.state.vel` | Speed, acceleration |
| `fleet.prod.veh.{id}.state.health` | Engine status, diagnostics |
| `fleet.prod.veh.{id}.state.fuel` | Fuel level / battery SOC |
| `fleet.prod.veh.{id}.evt.trip` | Trip start/end events |
| `fleet.prod.veh.{id}.evt.geofence` | Zone entry/exit |
| `fleet.prod.veh.{id}.cmd.mission` | Route assignment |

The subject structure is vehicle-agnostic. A fleet mixing drones and trucks uses identical patterns.

---

## Integration Patterns

### ELD Compliance (US Trucks)

Electronic Logging Device requirements:
- Hours of Service tracking
- Driver identification
- Location logging
- Data transfer to authorities

The gateway captures J1939 engine data and integrates with ELD compliance systems.

### Telematics Integration

Existing telematics platforms can connect via:
- NATS grants (subscribe to specific vehicles)
- REST API (query historical data)
- Webhook events (real-time alerts)

### Fleet Management Systems

Integrate with enterprise TMS/FMS:
- Vehicle assignments
- Route optimization
- Fuel management
- Maintenance scheduling

---

## Summary

| Aspect | Drones | Ground Vehicles |
|--------|--------|-----------------|
| **Protocol** | MAVLink | CAN bus / J1939 / ROS 2 |
| **Compute** | Jetson Orin Nano | Jetson / Industrial PC |
| **Power** | 4S LiPo (14.8V) | 12V/24V/48V vehicle |
| **Safety Override** | RC transmitter | Steering wheel / E-stop |
| **Failsafe** | Return-to-Launch | Stop in place |
| **Regulations** | Aviation (FAA/EASA) | Road transport (DOT) |
| **Fleet Architecture** | **Identical** | **Identical** |

---

## Related Documentation

- [Supported Platforms]({{< relref "/fleet/platforms" >}}) — Overview of all vehicle types
- [Drone Platform]({{< relref "/fleet/platforms/drones" >}}) — UAV hardware and MAVLink
- [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) — Protocol bridging architecture
- [Safety Model]({{< relref "/fleet/safety" >}}) — Vehicle-agnostic safety principles
- [Subject Naming]({{< relref "/fleet/subjects" >}}) — Telemetry subject hierarchy

---

## Get Started

Deploying a ground vehicle fleet? We can help with architecture and integration.

[Contact Us →](/contact)
