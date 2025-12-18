---
title: "Supported Platforms"
meta_title: "Supported Vehicle Platforms | Ubuntu Software"
description: "Our fleet architecture supports drones, cars, trucks, and autonomous ground vehicles. Same core infrastructure, vehicle-specific integration."
image: "/images/robotics.svg"
draft: false
---

## One Architecture, Many Vehicles

The Ubuntu Software fleet architecture is **vehicle-agnostic by design**. The same NATS JetStream infrastructure, digital twin patterns, and authorization model that manages drone fleets also manages cars, trucks, and autonomous ground vehicles.

---

## What's Universal

These components work identically across all vehicle types:

| Component | Description |
|-----------|-------------|
| **NATS JetStream** | Messaging backbone, works for any telemetry |
| **Subject Hierarchy** | `fleet.{env}.veh.{id}.*` — vehicle ID is just an identifier |
| **Digital Twins** | State streams, event streams, shadow stores |
| **Authorization** | Decentralized security, grants, third-party access |
| **Streams & Events** | Telemetry rollup, audit trails, queryable state |

**The architecture doesn't care what's generating the telemetry.** A position update from a drone and a position update from a truck flow through the same infrastructure.

---

## What's Vehicle-Specific

Each platform has unique requirements:

| Aspect | Varies By Platform |
|--------|-------------------|
| **Hardware** | Airframes vs chassis, flight controllers vs ECUs |
| **Protocols** | MAVLink, CAN bus, ROS 2, J1939 |
| **Safety Model** | RC override vs e-stop, RTL vs stop-in-place |
| **Sensors** | Altitude vs ground clearance, airspeed vs wheel speed |
| **Regulations** | Aviation authorities vs road transport |

---

## Supported Platforms

### [Drones]({{< relref "/fleet/platforms/drones" >}})

Unmanned aerial vehicles running PX4 or ArduPilot:

- **Hardware**: Holybro X500, Pixhawk 6X, NVIDIA Jetson
- **Protocol**: MAVLink over serial/UDP
- **Safety**: RC override, Return-to-Launch, geofencing
- **Use Cases**: Inspection, mapping, surveillance, delivery

[Drone Platform Details →]({{< relref "/fleet/platforms/drones" >}})

---

### [Ground Vehicles]({{< relref "/fleet/platforms/ground" >}})

Cars, trucks, AGVs, and autonomous ground platforms:

- **Hardware**: Vehicle ECUs, industrial PCs, Jetson
- **Protocol**: CAN bus, J1939, ROS 2
- **Safety**: E-stop, stop-in-place, geofencing
- **Use Cases**: Logistics, mining, agriculture, last-mile delivery

[Ground Vehicle Details →]({{< relref "/fleet/platforms/ground" >}})

---

## Gateway Pattern

The [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) bridges vehicle-native protocols to NATS. Each platform has a gateway implementation:

```
┌─────────────────────────────────────────────────────────────────┐
│                     VEHICLE GATEWAY                              │
├─────────────────┬─────────────────┬─────────────────────────────┤
│   MAVLink GW    │    CAN GW       │      ROS 2 GW               │
│   (Drones)      │  (Cars/Trucks)  │   (ROS Robots)              │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS JETSTREAM                              │
│            Same subjects, same streams, same patterns            │
└─────────────────────────────────────────────────────────────────┘
```

The gateway:
- Translates protocol-specific messages to common telemetry format
- Publishes to standard subject hierarchy
- Enforces safety policies
- Handles store-and-forward during connectivity loss

---

## Adding New Platforms

The architecture supports any vehicle that can:

1. **Generate telemetry** — Position, velocity, health data
2. **Accept commands** — Start, stop, waypoint navigation
3. **Run a gateway** — Translate native protocol to NATS

Custom platforms integrate by implementing a gateway for their protocol. The rest of the infrastructure (NATS, streams, authorization) works unchanged.

---

## Learn More

| Platform | Documentation |
|----------|---------------|
| [Drones]({{< relref "/fleet/platforms/drones" >}}) | PX4, Pixhawk, MAVLink integration |
| [Ground Vehicles]({{< relref "/fleet/platforms/ground" >}}) | CAN bus, J1939, automotive integration |
| [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) | Protocol bridging architecture |
| [Safety Model]({{< relref "/fleet/safety" >}}) | Platform-specific safety requirements |

---

## Get Started

Building a mixed fleet? Deploying ground vehicles alongside drones? We can help with architecture and implementation.

[Contact Us →](/contact)
