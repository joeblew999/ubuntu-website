---
title: "Safety Model"
meta_title: "Autonomous Fleet Safety Architecture | Ubuntu Software"
description: "Safety model for autonomous vehicle fleets: manual override authority, vehicle failsafes, network isolation from control loops, and graceful degradation."
image: "/images/robotics.svg"
draft: false
---

## Safety by Design

Fleet-scale autonomous vehicle operations demand rigorous safety architecture. This design ensures that **network failures, software bugs, and system malfunctions never compromise vehicle safety**.

The core safety principles apply to **all vehicle types**—drones, cars, trucks, and AGVs. The specifics differ (RC vs steering wheel, RTL vs stop-in-place), but the architecture is consistent.

---

## Core Principle

**NATS is never in the real-time control loop.**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SAFETY HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   1. MANUAL OVERRIDE      ─────────────────────────▶  HIGHEST       │
│      Human control (RC, steering, e-stop), always available          │
│                                                                      │
│   2. VEHICLE FAILSAFES    ─────────────────────────▶  HIGH          │
│      Hardware-enforced, local to vehicle                             │
│                                                                      │
│   3. GATEWAY POLICY       ─────────────────────────▶  MEDIUM        │
│      Software enforcement, on-vehicle                                │
│                                                                      │
│   4. FLEET COMMANDS       ─────────────────────────▶  LOWEST        │
│      Network-dependent, advisory                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

Higher layers can always override lower layers. Network-dependent commands are the lowest authority.

---

## Manual Override is Primary Authority

The human operator has **unconditional override**—whether via RC transmitter, steering wheel, or e-stop button.

### Drone Override (RC)

| Aspect | Implementation |
|--------|----------------|
| **Protocol** | ExpressLRS direct to receiver |
| **Path** | TX → RX → Pixhawk (no software in loop) |
| **Mode switch** | Hardware failsafe to manual/stabilize |
| **Kill switch** | Immediate motor disarm |

**RC never depends on:**

- Jetson being online
- NATS connection
- Vehicle Gateway
- Any software component

If every computer fails, the pilot retains full manual control.

#### ExpressLRS Failsafe

When RC signal is lost:

```
1. RX detects signal loss (100ms)
2. RX sends failsafe values to Pixhawk
3. Pixhawk enters RC_LOSS failsafe mode
4. Configured action executes (RTL, land, etc.)
```

Failsafe behavior is configured in PX4, not dependent on any fleet software.

### Ground Vehicle Override

| Aspect | Cars/Trucks | AGVs |
|--------|-------------|------|
| **Primary** | Steering wheel + pedals | E-stop button |
| **Path** | Physical connection to actuators | Hardware interrupt to motors |
| **Backup** | Key-off / gear neutral | Wireless e-stop |
| **Network** | Never in control path | Never in control path |

**Ground vehicle override never depends on:**

- Onboard computer being online
- NATS connection
- Any autonomous software

If every computer fails, the driver/operator retains physical control or the vehicle stops safely.

---

## Vehicle Control System Enforces Failsafes

The vehicle's control system implements failsafes locally—no network required.

### Drone Failsafes (PX4)

The flight controller implements multiple failsafes:

#### RC Loss Failsafe

| Parameter | Value | Action |
|-----------|-------|--------|
| `COM_RC_LOSS_T` | 0.5s | Timeout before failsafe |
| `NAV_RCL_ACT` | 2 | RTL on RC loss |

#### Data Link Loss

| Parameter | Value | Action |
|-----------|-------|--------|
| `COM_DL_LOSS_T` | 10s | GCS connection timeout |
| `NAV_DLL_ACT` | 0 | Continue mission (link not critical) |

#### Battery Failsafes

| Level | Parameter | Value | Action |
|-------|-----------|-------|--------|
| Low | `BAT_LOW_THR` | 0.25 | Warning |
| Critical | `BAT_CRIT_THR` | 0.15 | RTL |
| Emergency | `BAT_EMERGEN_THR` | 0.07 | Land immediately |

#### Geofence

| Parameter | Value | Action |
|-----------|-------|--------|
| `GF_ACTION` | 3 | RTL on breach |
| `GF_MAX_HOR_DIST` | varies | Horizontal limit |
| `GF_MAX_VER_DIST` | 120m | Altitude limit |

#### Position Loss

| Parameter | Value | Action |
|-----------|-------|--------|
| `COM_POS_FS_DELAY` | 1s | Position timeout |
| `COM_POSCTL_NAVL` | 0 | Land on position loss |

**All failsafes execute on the Pixhawk**—no external dependency.

### Ground Vehicle Failsafes

Ground vehicles implement equivalent safety mechanisms:

| Trigger | Drone Response | Ground Vehicle Response |
|---------|---------------|------------------------|
| **Comms loss** | Return-to-Launch | Stop in place / pull over |
| **Sensor failure** | Land immediately | Stop safely, engage parking |
| **Geofence breach** | RTL at boundary | Stop at boundary |
| **E-stop activated** | Motor cutoff | Immediate braking |
| **Low battery/fuel** | RTL then land | Return to depot / stop safely |
| **Obstacle detected** | Avoid or hover | Stop and wait |

**Ground vehicle failsafes typically include:**

- **Lidar/radar obstacles** — Emergency stop if path blocked
- **Speed limiting** — Max speed based on environment
- **Zone restrictions** — No-go areas enforced locally
- **Watchdog timers** — Stop if autonomy software hangs
- **Brake redundancy** — Multiple independent brake systems

**All failsafes execute locally**—no network dependency.

---

## NATS is Monitoring Only

The network layer handles:

| Function | Description |
|----------|-------------|
| **Telemetry collection** | Vehicle state for dashboards |
| **Event logging** | Audit trail in JetStream |
| **Mission coordination** | Fleet-level planning |
| **Shadow state** | Desired/reported synchronization |

The network layer **does not handle:**

| Function | Why Not |
|----------|---------|
| **Attitude control** | Latency unacceptable, safety-critical |
| **Motor commands** | Must be local to FCU |
| **Sensor fusion** | Real-time requirements |
| **Failsafe decisions** | Must work without network |

### Command Validation

When commands arrive via NATS:

```go
func (g *Gateway) validateCommand(cmd *Command) error {
    // Check geofence
    if cmd.Type == "goto" {
        if !g.geofence.Contains(cmd.Target) {
            return errors.New("target outside geofence")
        }
    }

    // Check altitude
    if cmd.Altitude > g.policy.MaxAltitude {
        return errors.New("altitude exceeds limit")
    }

    // Check vehicle state
    if cmd.RequiresArmed && !g.state.Armed {
        return errors.New("vehicle not armed")
    }

    // Check rate limiting
    if g.rateLimiter.Exceeded(cmd.Type) {
        return errors.New("command rate exceeded")
    }

    return nil
}
```

Commands are **validated before execution**—but even if a bad command executes, PX4 failsafes provide the final protection.

---

## Graceful Degradation

The system degrades safely when components fail:

### Jetson Failure

```
Jetson offline
    └─▶ No perception, no AI inference
    └─▶ Pi continues MAVLink routing
    └─▶ RC control fully operational
    └─▶ PX4 failsafes active
    └─▶ Vehicle lands safely or continues mission
```

### NATS Connection Lost

```
WAN disconnected
    └─▶ Leaf node buffers locally
    └─▶ Gateway continues operation
    └─▶ Vehicle follows current mission
    └─▶ State syncs when connection returns
```

### Vehicle Gateway Crash

```
Gateway process dies
    └─▶ systemd restarts within 5 seconds
    └─▶ PX4 continues autonomous flight
    └─▶ RC override always available
    └─▶ No fleet commands until restart
```

### Hub Cluster Failure

```
Regional hub down
    └─▶ Leaf nodes queue messages
    └─▶ Vehicles continue operating
    └─▶ RC and PX4 failsafes unaffected
    └─▶ Fleet visibility lost temporarily
```

### Cascading Failure Prevention

| Failure | Blast Radius | Mitigation |
|---------|--------------|------------|
| Single vehicle | That vehicle only | Failsafe to land |
| Hub node | No impact (replicated) | Cluster self-heals |
| Hub cluster | Regional fleet coordination | Vehicles continue autonomously |
| Global tier | Cross-region visibility | Regional operations continue |

---

## Operational Safety Procedures

### Pre-Flight Checklist

Before arming any fleet vehicle:

1. **RC link verified** — Pilot confirms control authority
2. **Failsafe test** — Verify RTL/land behavior
3. **Geofence loaded** — Boundaries confirmed in PX4
4. **Battery threshold** — Failsafe levels appropriate
5. **Weather check** — Wind within limits
6. **Airspace clear** — No conflicts

### Emergency Procedures

| Emergency | Response |
|-----------|----------|
| **Loss of visual** | Activate RTL via RC |
| **Flyaway** | RC kill switch or geofence triggers RTL |
| **Multi-vehicle conflict** | Fleet-wide land command |
| **Complete comms loss** | PX4 RC_LOSS failsafe |
| **Fire/crash** | RC disarm, emergency services |

### Kill Chain

Multiple methods to stop a vehicle:

```
1. RC kill switch      → Immediate motor stop
2. RC mode switch      → Land mode
3. GCS command         → RTL or land
4. Fleet command       → RTL or land
5. Geofence breach     → Auto-RTL
6. Battery critical    → Auto-land
```

At least three methods are always available (RC, PX4 failsafe, local timeout).

---

## Certification Considerations

This architecture supports regulatory compliance:

| Requirement | Implementation |
|-------------|----------------|
| **Pilot in command** | RC override at all times |
| **Failsafe behaviors** | PX4 certified failsafes |
| **Geofencing** | Hardware-enforced boundaries |
| **Logging** | JetStream audit trail |
| **Remote ID** | Broadcast via separate module |

**Note:** Regulatory requirements vary by jurisdiction. This architecture provides the technical foundation—specific compliance requires additional measures per local regulations.

---

## Testing Safety Systems

### Unit Tests

```go
func TestGeofenceRejection(t *testing.T) {
    gw := NewGateway(config)
    gw.SetGeofence(testGeofence)

    cmd := &Command{
        Type:   "goto",
        Target: Position{Lat: 0, Lon: 0}, // Outside geofence
    }

    err := gw.validateCommand(cmd)
    assert.ErrorContains(t, err, "outside geofence")
}
```

### Integration Tests

```
1. Simulate RC loss → Verify RTL activates
2. Simulate battery drain → Verify land at threshold
3. Simulate geofence breach → Verify return behavior
4. Simulate Jetson crash → Verify continued flight
5. Simulate NATS disconnect → Verify local operation
```

### Flight Tests

Before fleet deployment:

1. **Single vehicle failsafe validation**
2. **RC loss response verification**
3. **Geofence breach testing**
4. **Battery failsafe validation**
5. **Multi-vehicle coordination testing**

---

## Summary

| Layer | Authority | Drones | Ground Vehicles |
|-------|-----------|--------|-----------------|
| **Manual** | Highest | RC transmitter | Steering / E-stop |
| **Vehicle** | High | PX4 failsafes | ECU / safety controller |
| **Gateway** | Medium | Policy enforcement | Policy enforcement |
| **Fleet** | Lowest | Coordination only | Coordination only |

The safety model ensures:

- **Network is never safety-critical**
- **Operator always has override authority**
- **Hardware failsafes execute locally**
- **Failures degrade gracefully**
- **No single point causes fleet-wide impact**

---

## Related Documentation

- [Supported Platforms]({{< relref "/fleet/platforms" >}}) — Overview of all vehicle types
- [Drone Platform]({{< relref "/fleet/platforms/drones" >}}) — Drone-specific safety details
- [Ground Vehicles]({{< relref "/fleet/platforms/ground" >}}) — Ground vehicle safety details

---

## Further Reading

- [PX4 Safety Configuration](https://docs.px4.io/main/en/config/safety.html)
- [MAVLink Common Message Set](https://mavlink.io/en/messages/common.html)
- [NATS Security Model](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

---

## Back to Overview

[← Autonomous Vehicle Fleet Architecture]({{< relref "/fleet" >}})
