---
title: "JetStream Configuration"
meta_title: "JetStream Streams for Drone Fleet Digital Twins | Ubuntu Software"
description: "JetStream stream and KV store configuration for fleet digital twins: state rollup, event audit trails, command queues, and shadow state."
image: "/images/robotics.svg"
draft: false
---

## JetStream for Digital Twins

JetStream provides the persistence layer for fleet digital twins. Each stream type serves a specific purpose with tailored retention and replication settings.

---

## Stream Architecture

Five JetStream resources support the digital twin pattern:

| Resource | Type | Purpose |
|----------|------|---------|
| **TWIN_STATE** | Stream | Telemetry with rollup |
| **TWIN_EVENTS** | Stream | Durable event audit trail |
| **TWIN_CMD** | Stream | Command queue |
| **TWIN_SHADOW** | KV Store | Desired/reported state |
| **TWIN_BLOBS** | Object Store | Large artifacts |

---

## TWIN_STATE: Telemetry Stream

Stores continuous vehicle state with **subject-based rollup**:

```yaml
name: TWIN_STATE
subjects:
  - "fleet.prod.veh.*.state.>"
retention: limits
max_age: 1h
max_bytes: 10GB
storage: file
replicas: 3
discard: old
rollup_hdrs: true
```

### Configuration Explained

| Setting | Value | Rationale |
|---------|-------|-----------|
| **retention** | limits | Bounded by age and size |
| **max_age** | 1h | Keep recent state for replay |
| **max_bytes** | 10GB | Cap storage per hub |
| **replicas** | 3 | Survive node failures |
| **rollup_hdrs** | true | Enable per-subject rollup |

### Subject Rollup

With rollup enabled, consumers can request **only the latest message per subject**:

```go
// Get latest position for all vehicles
sub, _ := js.Subscribe("fleet.prod.veh.*.state.position",
    nats.DeliverLastPerSubject())
```

This transforms the stream from a firehose into a queryable state store.

### Consumer Patterns

| Pattern | Use Case |
|---------|----------|
| **DeliverAll** | Historical replay, analytics |
| **DeliverLast** | Current state snapshot |
| **DeliverLastPerSubject** | Latest state per vehicle |
| **StartTime** | Replay from specific timestamp |

---

## TWIN_EVENTS: Audit Trail

Stores discrete events with **long retention** for compliance and debugging:

```yaml
name: TWIN_EVENTS
subjects:
  - "fleet.prod.veh.*.evt.>"
retention: limits
max_age: 90d
max_bytes: 100GB
storage: file
replicas: 3
discard: old
```

### Configuration Explained

| Setting | Value | Rationale |
|---------|-------|-----------|
| **max_age** | 90d | Regulatory compliance, incident analysis |
| **max_bytes** | 100GB | Events are smaller than state |
| **storage** | file | Persistent across restarts |

### Event Properties

Events differ from state:

| Property | State | Events |
|----------|-------|--------|
| **Frequency** | Continuous (Hz) | Discrete (on change) |
| **Retention** | Short (hours) | Long (months) |
| **Consumers** | Real-time dashboards | Audit, analytics, alerting |
| **Replay** | Current state reconstruction | Investigation, compliance |

### Example Events

```json
// Armed event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T10:30:00Z",
  "data": {
    "previous_state": "disarmed",
    "armed_by": "operator-123",
    "reason": "mission_start"
  }
}

// Failsafe event
{
  "vid": "VID-001",
  "timestamp": "2024-01-15T11:45:00Z",
  "data": {
    "failsafe_type": "low_battery",
    "action_taken": "rtl",
    "battery_remaining": 0.15
  }
}
```

---

## TWIN_CMD: Command Queue

Stores commands with **workqueue retention**:

```yaml
name: TWIN_CMD
subjects:
  - "fleet.prod.veh.*.cmd.>"
retention: workqueue
max_age: 5m
storage: file
replicas: 3
```

### Configuration Explained

| Setting | Value | Rationale |
|---------|-------|-----------|
| **retention** | workqueue | Messages deleted after ack |
| **max_age** | 5m | Unacked commands expire |
| **replicas** | 3 | Command delivery guaranteed |

### Workqueue Semantics

1. Command published to stream
2. Vehicle Gateway receives (exactly-once delivery)
3. Gateway processes command
4. Gateway acks message
5. Message removed from stream

**Unacked commands re-deliver** to ensure commands aren't lost during disconnections.

### Command Flow

```
Fleet Operator               TWIN_CMD Stream              Vehicle Gateway
     │                             │                             │
     │ Publish: cmd.takeoff        │                             │
     │ ─────────────────────────▶  │                             │
     │                             │ Deliver to consumer         │
     │                             │ ─────────────────────────▶  │
     │                             │                             │ Process
     │                             │                             │ command
     │                             │             Ack             │
     │                             │ ◀─────────────────────────  │
     │                             │ (message removed)           │
     │                             │                             │
```

---

## TWIN_SHADOW: KV Store

JetStream KV stores maintain digital twin shadow state:

```yaml
bucket: TWIN_SHADOW
max_value_size: 1MB
history: 5
ttl: 0
replicas: 3
storage: file
```

### Configuration Explained

| Setting | Value | Rationale |
|---------|-------|-----------|
| **history** | 5 | Keep 5 previous versions |
| **ttl** | 0 | No expiration (permanent) |
| **replicas** | 3 | High availability |

### Key Structure

```
fleet/prod/veh/VID-001/desired    → Desired state JSON
fleet/prod/veh/VID-001/reported   → Reported state JSON
fleet/prod/veh/VID-001/config     → Vehicle configuration
fleet/prod/veh/VID-001/meta       → Metadata (serial, type, etc.)
```

### Shadow Operations

```go
// Write desired state
kv.Put("fleet/prod/veh/VID-001/desired", desiredJSON)

// Read reported state
entry, _ := kv.Get("fleet/prod/veh/VID-001/reported")

// Watch for changes
watcher, _ := kv.Watch("fleet/prod/veh/VID-001/desired")
for entry := range watcher.Updates() {
    // Desired state changed, reconcile
}
```

### History for Debugging

With history enabled, you can inspect previous states:

```go
// Get all historical values
history, _ := kv.History("fleet/prod/veh/VID-001/desired")
for _, entry := range history {
    fmt.Printf("Rev %d at %s: %s\n",
        entry.Revision(), entry.Created(), entry.Value())
}
```

---

## TWIN_BLOBS: Object Store

Large artifacts that don't fit in messages:

```yaml
bucket: TWIN_BLOBS
max_chunk_size: 128KB
storage: file
replicas: 3
```

### Use Cases

| Artifact | Size | Example |
|----------|------|---------|
| **Mission files** | 10-100KB | Waypoint definitions |
| **Log files** | 1-10MB | Flight logs, diagnostics |
| **Firmware** | 10-50MB | PX4 firmware images |
| **Maps** | 1-100MB | Offline map tiles |
| **ML models** | 10-500MB | TensorRT engine files |

### Object Operations

```go
// Upload mission file
obj, _ := obs.PutFile("missions/mission-456.json", file)

// Download to vehicle
reader, _ := obs.Get("missions/mission-456.json")

// List available firmware
objects := obs.List(nats.ObjectSearchPrefix("firmware/"))
```

### Chunked Transfer

Object Store handles large files by:

1. Splitting into chunks (128KB default)
2. Storing chunks as stream messages
3. Tracking metadata (size, hash, chunks)
4. Reassembling on retrieval

This enables **resumable transfers** over unreliable connections.

---

## Stream Sizing

### Per-Vehicle Estimates

| Stream | Message Size | Rate | Daily Volume |
|--------|--------------|------|--------------|
| **STATE** | 200 bytes | 10 Hz | ~170 MB |
| **EVENTS** | 500 bytes | 10/day | ~5 KB |
| **CMD** | 300 bytes | 100/day | ~30 KB |

### Fleet Estimates (1,000 vehicles)

| Resource | 1-Hour | 1-Day | 90-Day |
|----------|--------|-------|--------|
| **TWIN_STATE** | 7 GB | - | - |
| **TWIN_EVENTS** | - | 5 MB | 450 MB |
| **TWIN_SHADOW** | - | 1 GB | 1 GB |
| **TWIN_BLOBS** | - | - | 50 GB |

These are baseline estimates. Actual usage depends on message rates and sizes.

---

## Operational Considerations

### Stream Health

Monitor these metrics:

| Metric | Alert Threshold |
|--------|-----------------|
| **Consumer lag** | > 10,000 messages |
| **Stream bytes** | > 80% of max |
| **Replica sync** | Any replica behind |
| **Ack pending** | > 1,000 per consumer |

### Backup Strategy

| Resource | Backup Method | Frequency |
|----------|---------------|-----------|
| **Streams** | JetStream snapshot | Daily |
| **KV Store** | Key export | Daily |
| **Object Store** | File-level backup | Weekly |

### Migration

Streams can be migrated between clusters:

1. Create stream on target with same configuration
2. Create mirror sourcing from original
3. Wait for sync
4. Switch producers/consumers to target
5. Remove original stream

---

## Summary

| Resource | Type | Retention | Purpose |
|----------|------|-----------|---------|
| **TWIN_STATE** | Stream | 1 hour | Telemetry with rollup |
| **TWIN_EVENTS** | Stream | 90 days | Audit trail |
| **TWIN_CMD** | Stream | Workqueue | Command delivery |
| **TWIN_SHADOW** | KV | Permanent | Desired/reported state |
| **TWIN_BLOBS** | Object | Permanent | Large artifacts |

---

## Next

[Vehicle Gateway →]({{< relref "/fleet/gateway" >}})
