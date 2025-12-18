---
title: "Subject Naming"
meta_title: "NATS Subject Naming for Drone Fleets | Ubuntu Software"
description: "Hierarchical NATS subject naming conventions for fleet telemetry, events, commands, and digital twin state management."
image: "/images/robotics.svg"
draft: false
---

## Subject Hierarchy Design

Subject naming is the API of your messaging system. A well-designed hierarchy enables efficient routing, filtering, and access control while remaining intuitive for developers.

---

## Design Principles

| Principle | Rationale |
|-----------|-----------|
| **Hierarchical** | Enable wildcard subscriptions at any level |
| **Predictable** | Developers can construct subjects without lookup |
| **Environment-aware** | Separate production from staging/development |
| **Filterable** | Subscribe to one vehicle, group, or entire fleet |
| **Extensible** | Add new message types without restructuring |

---

## Core Subject Pattern

All fleet subjects follow this pattern:

```
fleet.<env>.veh.<vid>.<category>.<type>
```

| Segment | Description | Examples |
|---------|-------------|----------|
| `fleet` | Root namespace | Always "fleet" |
| `<env>` | Environment | `prod`, `staging`, `dev` |
| `veh` | Resource type | Vehicles (future: `gcs` for ground control) |
| `<vid>` | Vehicle ID | `VID-001`, `drone-alpha-7` |
| `<category>` | Message category | `state`, `evt`, `cmd`, `cmdack` |
| `<type>` | Specific type | `position`, `battery`, `takeoff` |

---

## Message Categories

### State: `fleet.<env>.veh.<vid>.state.*`

Current vehicle state, published continuously:

| Subject | Description | Rate |
|---------|-------------|------|
| `state.position` | Lat/lon/alt, velocity | 10 Hz |
| `state.attitude` | Roll/pitch/yaw | 10 Hz |
| `state.battery` | Voltage, current, remaining | 1 Hz |
| `state.mode` | Flight mode (AUTO, MANUAL, RTL) | On change |
| `state.health` | System health summary | 1 Hz |

**Example publish:**

```go
subject := fmt.Sprintf("fleet.prod.veh.%s.state.position", vehicleID)
nc.Publish(subject, positionJSON)
```

### Events: `fleet.<env>.veh.<vid>.evt.*`

Discrete occurrences, published when they happen:

| Subject | Description |
|---------|-------------|
| `evt.armed` | Vehicle armed |
| `evt.disarmed` | Vehicle disarmed |
| `evt.takeoff` | Takeoff initiated |
| `evt.landed` | Landing complete |
| `evt.waypoint` | Waypoint reached |
| `evt.failsafe` | Failsafe triggered |
| `evt.geofence` | Geofence breach |
| `evt.battery.low` | Low battery warning |

**Events are state transitions**, not continuous data. They're durable—stored in JetStream for audit trails.

### Commands: `fleet.<env>.veh.<vid>.cmd.*`

Instructions sent to vehicles:

| Subject | Description |
|---------|-------------|
| `cmd.arm` | Arm motors |
| `cmd.disarm` | Disarm motors |
| `cmd.takeoff` | Initiate takeoff |
| `cmd.land` | Initiate landing |
| `cmd.rtl` | Return to launch |
| `cmd.goto` | Navigate to position |
| `cmd.mission.upload` | Upload mission |
| `cmd.mission.start` | Start mission |
| `cmd.param.set` | Set parameter |

**Commands flow downstream** from fleet management to vehicles.

### Command Acknowledgments: `fleet.<env>.veh.<vid>.cmdack.*`

Responses to commands:

| Subject | Description |
|---------|-------------|
| `cmdack.arm` | Arm command result |
| `cmdack.takeoff` | Takeoff command result |
| `cmdack.goto` | Goto command result |

**Ack payload includes:**

```json
{
  "cmd_id": "cmd-123456",
  "status": "accepted|rejected|completed|failed",
  "error": "optional error message",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## Wildcard Subscriptions

NATS wildcards enable flexible subscriptions:

| Wildcard | Meaning | Example |
|----------|---------|---------|
| `*` | Single token | `fleet.prod.veh.*.state.position` |
| `>` | One or more tokens | `fleet.prod.veh.VID-001.>` |

### Common Subscription Patterns

| Pattern | Use Case |
|---------|----------|
| `fleet.prod.veh.*.state.position` | All vehicle positions |
| `fleet.prod.veh.VID-001.>` | Everything from one vehicle |
| `fleet.prod.veh.*.evt.>` | All events from all vehicles |
| `fleet.prod.veh.*.state.battery` | All battery states |
| `fleet.*.veh.VID-001.>` | One vehicle across all environments |

---

## KV Stores: Digital Twin Shadow

JetStream KV stores maintain **desired** and **reported** state:

### Subject Pattern (KV)

```
fleet/<env>/veh/<vid>/desired
fleet/<env>/veh/<vid>/reported
```

Note: KV uses `/` separator (bucket path), while pub/sub uses `.` separator.

### Desired State

What the fleet management system **wants** the vehicle to do:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "geofence_enabled": true,
  "max_altitude": 120,
  "home_position": {
    "lat": 37.7749,
    "lon": -122.4194
  }
}
```

### Reported State

What the vehicle **actually** reports:

```json
{
  "mode": "AUTO",
  "mission_id": "mission-456",
  "mission_progress": 0.45,
  "geofence_enabled": true,
  "armed": true,
  "in_flight": true,
  "battery_remaining": 0.72
}
```

### Shadow Reconciliation

The Vehicle Gateway continuously:

1. Reads desired state from KV
2. Compares with actual vehicle state
3. Issues commands to align reality with desired
4. Updates reported state

This enables **declarative fleet management**—set the desired state, let vehicles converge.

---

## Access Control

NATS authorization maps users to allowed subjects:

```
# Vehicle can publish its own state/events
publish: fleet.prod.veh.VID-001.state.>
publish: fleet.prod.veh.VID-001.evt.>
publish: fleet.prod.veh.VID-001.cmdack.>

# Vehicle can subscribe to commands for itself
subscribe: fleet.prod.veh.VID-001.cmd.>

# Vehicle can read/write its own shadow
kv: fleet/prod/veh/VID-001/*
```

```
# Fleet operator can subscribe to all state/events
subscribe: fleet.prod.veh.*.state.>
subscribe: fleet.prod.veh.*.evt.>
subscribe: fleet.prod.veh.*.cmdack.>

# Fleet operator can publish commands
publish: fleet.prod.veh.*.cmd.>

# Fleet operator can manage shadows
kv: fleet/prod/veh/*/desired
```

This ensures **vehicles can't impersonate other vehicles** or publish unauthorized commands.

---

## Subject Design Tradeoffs

### Why not flat subjects?

Flat subjects like `vehicle-001-position` don't support:

- Wildcard subscriptions
- Hierarchical access control
- Logical grouping

### Why not deeper hierarchy?

Subjects like `fleet.prod.region.us-west.site.warehouse-a.veh.VID-001.state.position`:

- Harder to construct programmatically
- Wildcards become complex
- Changes ripple through the system

The four-level hierarchy balances flexibility with simplicity.

---

## Message Format

All messages use **JSON** with consistent structure:

```json
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00.123Z",
  "seq": 12345,
  "data": {
    // message-specific payload
  }
}
```

| Field | Description |
|-------|-------------|
| `vid` | Vehicle ID (redundant with subject, but useful for consumers) |
| `timestamp` | ISO 8601 with millisecond precision |
| `seq` | Monotonic sequence number for ordering |
| `data` | Message-specific payload |

---

## Summary

| Category | Subject Pattern | Direction |
|----------|-----------------|-----------|
| **State** | `fleet.<env>.veh.<vid>.state.<type>` | Vehicle → Hub |
| **Events** | `fleet.<env>.veh.<vid>.evt.<type>` | Vehicle → Hub |
| **Commands** | `fleet.<env>.veh.<vid>.cmd.<type>` | Hub → Vehicle |
| **Acks** | `fleet.<env>.veh.<vid>.cmdack.<type>` | Vehicle → Hub |
| **Desired** | `fleet/<env>/veh/<vid>/desired` (KV) | Hub → Vehicle |
| **Reported** | `fleet/<env>/veh/<vid>/reported` (KV) | Vehicle → Hub |

---

## Next

[Stream Configuration →]({{< relref "/fleet/streams" >}})
