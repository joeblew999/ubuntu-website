---
title: "Safety Model"
meta_title: "Drone Fleet Safety Architecture | Ubuntu Software"
description: "Safety model for 1,000-drone fleet: RC authority, PX4 failsafes, network isolation from control loops, and graceful degradation."
image: "/images/robotics.svg"
draft: false
---

## Safety by Design

Fleet-scale drone operations demand rigorous safety architecture. This design ensures that **network failures, software bugs, and system malfunctions never compromise flight safety**.

---

## Core Principle

**NATS is never in the real-time control loop.**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SAFETY HIERARCHY                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   1. RC OVERRIDE          ─────────────────────────▶  HIGHEST       │
│      Manual pilot control, always available                          │
│                                                                      │
│   2. PX4 FAILSAFES        ─────────────────────────▶  HIGH          │
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

## RC is Primary Authority

The pilot's RC transmitter has **unconditional override**:

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

### ExpressLRS Failsafe

When RC signal is lost:

```
1. RX detects signal loss (100ms)
2. RX sends failsafe values to Pixhawk
3. Pixhawk enters RC_LOSS failsafe mode
4. Configured action executes (RTL, land, etc.)
```

Failsafe behavior is configured in PX4, not dependent on any fleet software.

---

## PX4 Enforces Failsafes

The flight controller implements multiple failsafes:

### RC Loss Failsafe

| Parameter | Value | Action |
|-----------|-------|--------|
| `COM_RC_LOSS_T` | 0.5s | Timeout before failsafe |
| `NAV_RCL_ACT` | 2 | RTL on RC loss |

### Data Link Loss

| Parameter | Value | Action |
|-----------|-------|--------|
| `COM_DL_LOSS_T` | 10s | GCS connection timeout |
| `NAV_DLL_ACT` | 0 | Continue mission (link not critical) |

### Battery Failsafes

| Level | Parameter | Value | Action |
|-------|-----------|-------|--------|
| Low | `BAT_LOW_THR` | 0.25 | Warning |
| Critical | `BAT_CRIT_THR` | 0.15 | RTL |
| Emergency | `BAT_EMERGEN_THR` | 0.07 | Land immediately |

### Geofence

| Parameter | Value | Action |
|-----------|-------|--------|
| `GF_ACTION` | 3 | RTL on breach |
| `GF_MAX_HOR_DIST` | varies | Horizontal limit |
| `GF_MAX_VER_DIST` | 120m | Altitude limit |

### Position Loss

| Parameter | Value | Action |
|-----------|-------|--------|
| `COM_POS_FS_DELAY` | 1s | Position timeout |
| `COM_POSCTL_NAVL` | 0 | Land on position loss |

**All failsafes execute on the Pixhawk**—no external dependency.

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

| Layer | Authority | Failure Mode |
|-------|-----------|--------------|
| **RC** | Highest | Direct pilot control |
| **PX4** | High | Autonomous failsafe |
| **Gateway** | Medium | Policy enforcement |
| **Fleet** | Lowest | Coordination only |

The safety model ensures:

- **Network is never safety-critical**
- **Pilot always has override authority**
- **Hardware failsafes execute locally**
- **Failures degrade gracefully**
- **No single point causes fleet-wide impact**

---

## Further Reading

- [PX4 Safety Configuration](https://docs.px4.io/main/en/config/safety.html)
- [MAVLink Common Message Set](https://mavlink.io/en/messages/common.html)
- [NATS Security Model](https://docs.nats.io/running-a-nats-service/configuration/securing_nats)

---

## Back to Overview

[← Drone Fleet Architecture]({{< relref "/fleet" >}})
