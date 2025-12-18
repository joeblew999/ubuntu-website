---
title: "NATS Topology"
meta_title: "NATS JetStream Topology for Drone Fleets | Ubuntu Software"
description: "Hierarchical NATS topology for 1,000-drone fleet: leaf nodes per vehicle, regional hub clusters, WAN isolation, and fleet-wide digital twin replay."
image: "/images/robotics.svg"
draft: false
---

## NATS Architecture for Fleet Scale

The messaging topology determines how telemetry flows, where state persists, and how the system behaves when connectivity degrades. This architecture scales from tens to thousands of vehicles while maintaining low latency and high reliability.

---

## The Challenge

Fleet messaging must solve competing requirements:

| Requirement | Challenge |
|-------------|-----------|
| **Low latency** | Vehicle-to-vehicle communication needs milliseconds, not seconds |
| **WAN resilience** | Cellular connections drop, latency spikes, bandwidth varies |
| **State persistence** | Digital twin data must survive restarts and reconnections |
| **Scale** | Thousands of concurrent publishers and subscribers |
| **Isolation** | Vehicle failures shouldn't cascade to the fleet |

Traditional architectures force tradeoffs. Centralized brokers add latency. Peer-to-peer meshes don't scale. Edge-only approaches lose fleet visibility.

---

## Hierarchical Topology

The solution is a **three-tier hierarchy**:

```
┌─────────────────────────────────────────────────────────────────┐
│                         GLOBAL (optional)                        │
│                    Cross-region replication                      │
│                    Fleet-wide aggregation                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  REGIONAL HUB   │ │  REGIONAL HUB   │ │  REGIONAL HUB   │
│   Americas      │ │   Europe        │ │   Asia-Pacific  │
│  3-node cluster │ │  3-node cluster │ │  3-node cluster │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
    ┌────┼────┐         ┌────┼────┐         ┌────┼────┐
    │    │    │         │    │    │         │    │    │
    ▼    ▼    ▼         ▼    ▼    ▼         ▼    ▼    ▼
┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐┌─────┐
│LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF ││LEAF │
│VID-1││VID-2││VID-3││VID-4││VID-5││VID-6││VID-7││VID-8││VID-9│
└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘└─────┘
```

---

## Tier 1: Vehicle Leaf Nodes

Each drone runs a **NATS leaf node** on the Jetson:

| Aspect | Configuration |
|--------|---------------|
| **Process** | nats-server in leaf mode |
| **Storage** | Local JetStream on SSD/eMMC |
| **Connection** | Single upstream to regional hub |
| **Resources** | ~50MB RAM, minimal CPU |

**What Leaf Nodes Provide:**

### Local Pub/Sub

On-vehicle services communicate through the local NATS instance:

```
Vehicle Gateway ──publish──▶ fleet.prod.veh.VID-001.state.position
                                       │
                                       ▼
                            ┌──────────────────┐
                            │   Local NATS     │
                            │   (Leaf Node)    │
                            └──────────────────┘
                                       │
                                       ▼
AI Perception ◀──subscribe── fleet.prod.veh.VID-001.state.position
```

Messages between on-vehicle services never leave the vehicle. Latency is sub-millisecond.

### Store-and-Forward

When the upstream connection drops, the leaf node:

1. Continues accepting publishes from local services
2. Stores messages in local JetStream
3. Queues outbound messages for the hub
4. Reconnects and replays when connectivity returns

This means **vehicles continue operating during connectivity loss**. State synchronizes when the link recovers.

### WAN Isolation

The leaf node acts as a buffer between the vehicle and WAN:

- **Backpressure** — If the hub can't keep up, the leaf buffers locally
- **Filtering** — Only subscribed subjects traverse the WAN link
- **Compression** — NATS protocol compresses efficiently
- **Reconnection** — Automatic reconnection with exponential backoff

---

## Tier 2: Regional Hub Clusters

Regional hubs aggregate vehicles by geographic area:

| Aspect | Configuration |
|--------|---------------|
| **Deployment** | 3-node NATS cluster (minimum) |
| **Location** | Cloud region or edge data center |
| **Storage** | JetStream on fast SSDs |
| **Capacity** | Hundreds of leaf connections per cluster |

**Why Regional:**

- **Latency** — Leaf nodes connect to nearby hubs
- **Bandwidth** — Telemetry stays regional by default
- **Compliance** — Data residency requirements
- **Failure isolation** — Regional outages don't affect other regions

**Hub Responsibilities:**

| Function | Description |
|----------|-------------|
| **Leaf aggregation** | Accept connections from vehicle leaf nodes |
| **Stream storage** | Persist digital twin data in JetStream |
| **Consumer support** | Serve fleet dashboards, APIs, analytics |
| **Cross-region routing** | Forward to global tier when needed |

### Cluster Configuration

A 3-node cluster provides:

- **High availability** — Survives single node failure
- **Data replication** — JetStream streams replicated across nodes
- **Load distribution** — Leaf connections balanced across nodes

```yaml
# Example hub cluster configuration
cluster:
  name: hub-americas
  routes:
    - nats://hub-americas-1:6222
    - nats://hub-americas-2:6222
    - nats://hub-americas-3:6222

jetstream:
  store_dir: /data/jetstream
  max_memory_store: 4GB
  max_file_store: 100GB

leafnodes:
  port: 7422
  authorization:
    users:
      - user: vehicle
        password: $VEHICLE_PASSWORD
        allowed_connection_types: ["LEAFNODE"]
```

---

## Tier 3: Global (Optional)

For fleets spanning multiple regions, a global tier enables:

| Function | Description |
|----------|-------------|
| **Cross-region mirroring** | Replicate streams between regional hubs |
| **Global aggregation** | Fleet-wide dashboards and analytics |
| **Command routing** | Send commands to vehicles in any region |
| **Disaster recovery** | Failover between regions |

**Implementation Options:**

1. **NATS Supercluster** — Federate regional clusters via gateway connections
2. **JetStream Mirroring** — Replicate specific streams to central location
3. **Application-level** — Custom sync between regional APIs

Most fleets operate within a single region and don't need global tier initially.

---

## Benefits of This Topology

### Low Latency Local Pub/Sub

On-vehicle communication is **sub-millisecond**:

```
MAVLink message received
    └─▶ Vehicle Gateway publishes to local NATS
        └─▶ AI Perception subscribes from local NATS
            └─▶ Total latency: <1ms
```

No WAN round-trip for on-vehicle communication.

### WAN Isolation

Vehicle operations **continue during connectivity loss**:

- Perception systems keep running
- State accumulates locally
- Commands queue for execution
- Reconnection synchronizes automatically

### Fleet-Wide Digital Twin Replay

JetStream enables **temporal queries** across the fleet:

```
# Replay vehicle state from 10 minutes ago
nats stream get TWIN_STATE --start-time="2024-01-15T10:00:00Z"

# Subscribe to real-time state updates
nats sub "fleet.prod.veh.*.state.position"
```

Every vehicle's state history is queryable from the regional hub.

### Scalable Architecture

The topology scales horizontally:

| Scale Point | Approach |
|-------------|----------|
| **More vehicles** | Add leaf connections to existing hub |
| **Higher throughput** | Add nodes to hub cluster |
| **More regions** | Deploy additional regional hubs |
| **Global reach** | Connect hubs via supercluster |

---

## Connection Flow

When a vehicle powers on:

```
1. Jetson boots, NATS leaf node starts
2. Leaf node reads hub address from config
3. TLS connection established to regional hub
4. Authentication via credentials
5. Leaf node advertises local subjects
6. Hub routes relevant subscriptions to leaf
7. Bidirectional message flow begins
```

When connectivity drops:

```
1. Leaf node detects connection loss
2. Local pub/sub continues uninterrupted
3. Outbound messages queue locally
4. Reconnection attempts with backoff
5. On reconnect, queued messages replay
6. Stream state synchronizes
```

---

## Security

Every connection is authenticated and encrypted:

| Connection | Security |
|------------|----------|
| **Leaf → Hub** | TLS 1.3, credential authentication |
| **Hub cluster** | TLS, cluster routes authenticated |
| **Hub → Global** | TLS, gateway authentication |

Credentials are provisioned per-vehicle, enabling:

- Individual vehicle revocation
- Audit trails per vehicle
- Rate limiting per connection

---

## Summary

| Tier | Purpose | Deployment |
|------|---------|------------|
| **Leaf (Vehicle)** | Local pub/sub, store-and-forward | On Jetson |
| **Hub (Regional)** | Aggregation, persistence, consumers | Cloud/edge datacenter |
| **Global (Optional)** | Cross-region, fleet-wide views | Central cloud |

This topology ensures vehicles operate independently while maintaining fleet-wide visibility and control.

---

## Deployment Options

### Open Source Foundation

NATS JetStream is **100% open source** under the Apache 2.0 license:

- **Source code**: [github.com/nats-io/nats-server](https://github.com/nats-io/nats-server)
- **No commercial license required** — use freely in production
- **Full feature parity** — open source has all features
- **Active development** — backed by Synadia, used by thousands

You run the exact same code we run. No lock-in.

---

### Option 1: Self-Hosted

Deploy the entire stack on your infrastructure:

| Component | Your Responsibility |
|-----------|---------------------|
| Hub clusters | Provision, configure, maintain |
| Leaf nodes | Deploy on each vehicle |
| Monitoring | Set up observability |
| Updates | Manage upgrades |

**We provide:**
- Reference configurations (this documentation)
- Architecture consulting
- Implementation support

Best for: Organizations with ops teams, data residency requirements, or existing infrastructure.

---

### Option 2: Connect to Our Infrastructure

Tap into Ubuntu Software's managed NATS infrastructure:

| Component | Responsibility |
|-----------|----------------|
| Hub clusters | **We manage** |
| Leaf nodes | You deploy on vehicles |
| Monitoring | **We provide dashboards** |
| Scaling | **We handle** |

**How it works:**
1. We provision credentials for your fleet
2. Your leaf nodes connect to our regional hubs
3. Your data stays isolated (dedicated subjects)
4. You get fleet dashboards and API access

**Benefits:**
- No infrastructure to build or maintain
- Pre-tuned for drone fleet patterns
- Global coverage across regions
- Focus on your drones, not your messaging

**Pricing:**
- **Free tier** — Up to 10 vehicles, perfect for development and small deployments
- **Scale tier** — Per-vehicle pricing as you grow

[Contact Us →](/contact) to get started.

---

### Comparison

| Aspect | Self-Hosted | Managed |
|--------|-------------|---------|
| **Setup time** | Days to weeks | Hours |
| **Ops burden** | You maintain | We maintain |
| **Cost model** | Infrastructure + time | Per-vehicle |
| **Customization** | Full control | Standard config |
| **Data location** | Your choice | Our regions |
| **Best for** | Large ops teams | Fast deployment |

---

## Next

[Subject Naming →]({{< relref "/fleet/subjects" >}})
