---
title: "Vehicle Gateway"
meta_title: "Vehicle Gateway Design for Drone Fleets | Ubuntu Software"
description: "Go-based Vehicle Gateway architecture: MAVLink ingestion, state downsampling, event extraction, command execution, and shadow reconciliation."
image: "/images/robotics.svg"
draft: false
---

## The Bridge Between MAVLink and NATS

The Vehicle Gateway is a Go service running on each Jetson that translates between the MAVLink protocol spoken by PX4 and the NATS messaging system used for fleet coordination.

---

## Why a Gateway?

MAVLink and NATS serve different purposes:

| Aspect | MAVLink | NATS |
|--------|---------|------|
| **Scope** | Single vehicle | Fleet-wide |
| **Protocol** | Binary, compact | JSON, human-readable |
| **Rate** | 100+ Hz telemetry | Downsampled for WAN |
| **Persistence** | None | JetStream streams |
| **Pattern** | Request-response | Pub/sub with persistence |

The Gateway bridges these worlds, handling protocol translation, rate limiting, and state management.

---

## Core Responsibilities

### 1. MAVLink Ingest

Receive and parse MAVLink messages from the flight controller:

```go
// Receive MAVLink frames from mavlink-router
conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 14550})

for {
    buf := make([]byte, 1024)
    n, _, _ := conn.ReadFromUDP(buf)

    frame, _ := mavlink.Parse(buf[:n])

    switch msg := frame.Message().(type) {
    case *mavlink.Heartbeat:
        handleHeartbeat(msg)
    case *mavlink.GlobalPositionInt:
        handlePosition(msg)
    case *mavlink.BatteryStatus:
        handleBattery(msg)
    // ... handle other message types
    }
}
```

**Key message types:**

| MAVLink Message | Content |
|-----------------|---------|
| `HEARTBEAT` | Mode, armed state, system status |
| `GLOBAL_POSITION_INT` | Lat, lon, alt, velocity |
| `ATTITUDE` | Roll, pitch, yaw |
| `BATTERY_STATUS` | Voltage, current, remaining |
| `GPS_RAW_INT` | GPS fix, satellites, HDOP |
| `SYS_STATUS` | CPU load, errors, health |

### 2. State Downsampling

Reduce telemetry rate for WAN transmission:

```go
// Position: downsample from 10Hz to 1Hz
positionSampler := NewDownsampler(100*time.Millisecond, func(msg *Position) {
    // Average positions over window
    return averagePosition(msg)
})

// Battery: only publish on change
batterySampler := NewChangeFilter(func(old, new *Battery) bool {
    return math.Abs(old.Remaining - new.Remaining) > 0.01
})
```

**Downsampling strategies:**

| Strategy | Use Case | Example |
|----------|----------|---------|
| **Time-based** | Periodic state | Position at 1Hz |
| **Change-based** | Discrete values | Mode changes |
| **Threshold** | Gradual changes | Battery when Δ > 1% |
| **Aggregation** | High-frequency data | Min/max/avg over window |

Full-rate telemetry stays on-vehicle for perception systems. Only downsampled data crosses the WAN.

### 3. Event Extraction

Generate events from state transitions:

```go
type EventDetector struct {
    previousState *VehicleState
}

func (e *EventDetector) Process(current *VehicleState) []Event {
    var events []Event

    // Detect arm state change
    if current.Armed && !e.previousState.Armed {
        events = append(events, Event{
            Type:      "armed",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"reason": "manual"},
        })
    }

    // Detect mode change
    if current.Mode != e.previousState.Mode {
        events = append(events, Event{
            Type:      "mode_change",
            Timestamp: time.Now(),
            Data: map[string]interface{}{
                "from": e.previousState.Mode,
                "to":   current.Mode,
            },
        })
    }

    // Detect failsafe
    if current.FailsafeActive && !e.previousState.FailsafeActive {
        events = append(events, Event{
            Type:      "failsafe",
            Timestamp: time.Now(),
            Data:      map[string]interface{}{"type": current.FailsafeType},
        })
    }

    e.previousState = current
    return events
}
```

**Events detected:**

| Event | Trigger |
|-------|---------|
| `armed` | Armed state false → true |
| `disarmed` | Armed state true → false |
| `mode_change` | Flight mode transition |
| `takeoff` | In-air state false → true |
| `landed` | In-air state true → false |
| `failsafe` | Failsafe activated |
| `geofence` | Geofence breach detected |
| `battery.low` | Battery below threshold |
| `battery.critical` | Battery critical level |

### 4. Command Execution with Policy

Receive commands from NATS and execute via MAVLink:

```go
func (g *Gateway) handleCommand(cmd *Command) *CommandAck {
    // Validate command is allowed
    if !g.policy.Allows(cmd) {
        return &CommandAck{
            Status: "rejected",
            Error:  "command not allowed by policy",
        }
    }

    // Convert to MAVLink command
    mavCmd := cmd.ToMAVLink()

    // Send to flight controller
    if err := g.mavlink.Send(mavCmd); err != nil {
        return &CommandAck{
            Status: "failed",
            Error:  err.Error(),
        }
    }

    // Wait for MAVLink ACK
    mavAck, err := g.mavlink.WaitAck(mavCmd.CommandID, 5*time.Second)
    if err != nil {
        return &CommandAck{
            Status: "timeout",
            Error:  "no response from flight controller",
        }
    }

    return &CommandAck{
        Status: mavAck.Result.String(),
    }
}
```

**Policy enforcement:**

| Policy | Description |
|--------|-------------|
| **Geofence** | Reject goto commands outside boundary |
| **Altitude limit** | Cap maximum altitude commands |
| **Mode restrictions** | Disallow certain mode transitions |
| **Rate limiting** | Prevent command flooding |
| **Authentication** | Verify command source |

### 5. Shadow Reconciliation

Sync desired state with actual vehicle state:

```go
func (g *Gateway) reconcileLoop() {
    ticker := time.NewTicker(1 * time.Second)

    for range ticker.C {
        // Read desired state from KV
        desired, _ := g.kv.Get(g.desiredKey())

        // Compare with actual state
        actual := g.getCurrentState()

        // Generate commands to reconcile
        commands := g.reconcile(desired, actual)

        for _, cmd := range commands {
            g.executeCommand(cmd)
        }

        // Update reported state
        reported := g.buildReportedState(actual)
        g.kv.Put(g.reportedKey(), reported)
    }
}

func (g *Gateway) reconcile(desired, actual *State) []*Command {
    var commands []*Command

    // Mode reconciliation
    if desired.Mode != actual.Mode {
        commands = append(commands, &Command{
            Type: "set_mode",
            Data: map[string]interface{}{"mode": desired.Mode},
        })
    }

    // Geofence reconciliation
    if desired.GeofenceEnabled != actual.GeofenceEnabled {
        commands = append(commands, &Command{
            Type: "set_geofence",
            Data: map[string]interface{}{"enabled": desired.GeofenceEnabled},
        })
    }

    return commands
}
```

**Shadow state enables declarative management:**

- Fleet operator sets desired state
- Gateway detects differences
- Gateway issues commands to converge
- Gateway reports actual state
- Repeat continuously

---

## Architecture

```
┌────────────────────────────────────────────────────────────────────┐
│                         Vehicle Gateway                             │
│                                                                     │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐     │
│  │   MAVLink    │    │    State     │    │    NATS Client   │     │
│  │   Receiver   │───▶│   Machine    │───▶│    Publisher     │     │
│  └──────────────┘    └──────────────┘    └──────────────────┘     │
│         │                   │                      │               │
│         │            ┌──────▼──────┐               │               │
│         │            │   Event     │               │               │
│         │            │  Detector   │───────────────┤               │
│         │            └─────────────┘               │               │
│         │                                          │               │
│         │            ┌─────────────┐               │               │
│         │            │   Shadow    │◀──────────────┤               │
│         │            │ Reconciler  │               │               │
│         │            └──────┬──────┘               │               │
│         │                   │                      │               │
│         │            ┌──────▼──────┐    ┌─────────▼────────┐      │
│         │            │  Command    │    │   NATS Client    │      │
│         │◀───────────│  Executor   │◀───│   Subscriber     │      │
│         │            └─────────────┘    └──────────────────┘      │
│  ┌──────▼──────┐                                                   │
│  │   MAVLink   │                                                   │
│  │   Sender    │                                                   │
│  └─────────────┘                                                   │
└────────────────────────────────────────────────────────────────────┘
         │                                           │
         ▼                                           ▼
┌─────────────────┐                       ┌─────────────────┐
│   Pixhawk 6X    │                       │   NATS Leaf     │
│   (via serial)  │                       │   (localhost)   │
└─────────────────┘                       └─────────────────┘
```

---

## Implementation Details

### Dependencies

```go
import (
    "github.com/nats-io/nats.go"
    "github.com/nats-io/nats.go/jetstream"
    "github.com/bluenviern/go-mavlink/v2"
)
```

**Core libraries:**

| Library | Purpose |
|---------|---------|
| **nats.go** | NATS client, JetStream, KV |
| **go-mavlink** | MAVLink protocol implementation |
| **slog** | Structured logging |

### Configuration

```yaml
# gateway.yaml
vehicle_id: VID-001
environment: prod

mavlink:
  local_addr: ":14550"
  system_id: 1
  component_id: 1

nats:
  url: "nats://localhost:4222"
  credentials: "/etc/gateway/vehicle.creds"

sampling:
  position_hz: 1
  attitude_hz: 1
  battery_change_threshold: 0.01

policy:
  max_altitude: 120
  geofence_file: "/etc/gateway/geofence.json"
```

### Deployment

```yaml
# systemd service
[Unit]
Description=Vehicle Gateway
After=network.target nats.service

[Service]
Type=simple
ExecStart=/usr/local/bin/vehicle-gateway --config /etc/gateway/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

---

## Error Handling

### MAVLink Errors

| Error | Response |
|-------|----------|
| **Connection lost** | Reconnect with backoff |
| **Parse error** | Log and skip frame |
| **Timeout** | Retry or fail command |

### NATS Errors

| Error | Response |
|-------|----------|
| **Connection lost** | Local NATS continues, reconnect to hub |
| **Publish failed** | Buffer locally, retry |
| **Stream error** | Log, alert, continue |

### Command Errors

| Error | Response |
|-------|----------|
| **Policy violation** | Reject immediately |
| **MAVLink rejection** | Return failure ACK |
| **Timeout** | Return timeout ACK |

---

## Metrics

The Gateway exposes Prometheus metrics:

```
# MAVLink
gateway_mavlink_messages_received_total{type="heartbeat"}
gateway_mavlink_messages_sent_total{type="command_long"}
gateway_mavlink_parse_errors_total

# NATS
gateway_nats_messages_published_total{subject="state.position"}
gateway_nats_messages_received_total{subject="cmd.takeoff"}
gateway_nats_publish_errors_total

# Commands
gateway_commands_received_total{type="takeoff"}
gateway_commands_executed_total{type="takeoff",result="success"}
gateway_command_latency_seconds{type="takeoff"}

# Shadow
gateway_shadow_reconciliations_total
gateway_shadow_commands_issued_total
```

---

## Summary

| Responsibility | Input | Output |
|----------------|-------|--------|
| **MAVLink Ingest** | UDP packets | Parsed messages |
| **State Downsampling** | 100Hz telemetry | 1Hz state |
| **Event Extraction** | State transitions | Discrete events |
| **Command Execution** | NATS commands | MAVLink commands |
| **Shadow Reconciliation** | Desired state | Convergence commands |

The Vehicle Gateway is the critical component that makes fleet-scale operations possible while preserving the safety guarantees of the underlying PX4 autopilot.

---

## Next

[Safety Model →]({{< relref "/fleet/safety" >}})
