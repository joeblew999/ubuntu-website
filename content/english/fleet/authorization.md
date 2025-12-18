---
title: "Authorization & Grants"
meta_title: "Decentralized Authorization & Access Grants | Ubuntu Software"
description: "NATS JetStream's decentralized security model with grants enables third parties to securely participate in real-time systems without central coordination."
image: "/images/robotics.svg"
draft: false
---

## Security Without Bottlenecks

Traditional centralized authorization creates problems at scale:

| Problem | Impact |
|---------|--------|
| **Single point of failure** | Auth server down = entire system down |
| **Latency on every request** | Round-trip to auth server adds delay |
| **Can't work disconnected** | No network = no authorization |
| **Doesn't scale** | Thousands of devices overwhelm central server |

NATS takes a different approach: **decentralized authorization** with self-contained credentials.

---

## How NATS Authorization Works

```
┌─────────────────────────────────────────────────────────────────┐
│                    CREDENTIAL ISSUANCE                          │
│   Account Authority issues signed JWTs with embedded perms      │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌─────────┐     ┌─────────┐     ┌─────────┐
        │ Vehicle │     │ Operator│     │ Partner │
        │  Creds  │     │  Creds  │     │  Creds  │
        └────┬────┘     └────┬────┘     └────┬────┘
             │               │               │
             ▼               ▼               ▼
┌─────────────────────────────────────────────────────────────────┐
│                      NATS SERVER                                │
│   Validates JWT signature locally — no external call needed     │
│   Enforces permissions from embedded claims                     │
└─────────────────────────────────────────────────────────────────┘
```

**Key properties:**

- **Local validation** — Server validates credentials without calling any external service
- **Embedded permissions** — What you can do is encoded in your credential
- **Instant revocation** — Account-level revocation list propagates across cluster
- **Zero trust** — No credential = no access, period

This works during disconnection. No call-home for authorization. The credential is the proof.

---

## Permission Scopes

Credentials specify exactly what a party can do:

| Scope | Description | Example |
|-------|-------------|---------|
| **Publish** | Send messages to subjects | Vehicle publishes `fleet.prod.veh.V001.state.pos` |
| **Subscribe** | Receive messages from subjects | Operator subscribes to `fleet.prod.veh.*.state.*` |
| **Publish + Subscribe** | Both directions | Gateway handles command/response |
| **Deny** | Explicitly block patterns | Block vehicle from subscribing to other vehicles |

**Subject wildcards in permissions:**

```
fleet.prod.veh.V001.*        → Only this vehicle's subjects
fleet.prod.veh.*.state.*     → All vehicles' state (read-only)
fleet.prod.veh.*.cmd.*       → All vehicles' commands (operator)
fleet.prod.veh.*.evt.*       → All events (auditor)
```

This is how we enforce that vehicles can only publish their own data and can't impersonate other vehicles. See [Subject Naming]({{< relref "/fleet/subjects" >}}) for the full hierarchy.

---

## The Grants Model

Grants enable **fine-grained, delegatable permissions**. An operator can grant a third party access to specific data without giving them full system access.

```
Account Authority (root)
    │
    ├── Fleet Operator Account
    │       │
    │       ├── Grant: Emergency Services (incident-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.pos
    │       │       └── Subscribe: fleet.prod.veh.*.stream.video
    │       │       └── Expires: 4 hours
    │       │
    │       ├── Grant: Research Partner (data-scoped)
    │       │       └── Subscribe: fleet.prod.veh.*.state.env.*
    │       │       └── Deny: fleet.prod.veh.*.state.pos
    │       │       └── Expires: 30 days
    │       │
    │       └── Grant: Auditor (read-only, full)
    │               └── Subscribe: fleet.prod.veh.*.state.*
    │               └── Subscribe: fleet.prod.veh.*.evt.*
    │               └── Deny: fleet.prod.veh.*.cmd.*
    │               └── Expires: 7 days
```

**Grant properties:**

| Property | Benefit |
|----------|---------|
| **Scoped** | Only specific subjects/patterns—nothing more |
| **Time-bounded** | Automatic expiration, no forgotten access |
| **Revocable** | Instant revocation without touching grantee |
| **Delegatable** | Operators can grant sub-permissions to others |
| **Audited** | Every grant issuance logged |

---

## Third-Party Integration Scenarios

### Emergency Services Integration

```
Scenario: Wildfire response, fire department needs drone video

Grant issued:
  - Subject: fleet.prod.veh.*.stream.video
  - Subject: fleet.prod.veh.*.state.pos
  - Duration: 4 hours (incident duration)
  - Audit: All access logged with responder ID

Result: Fire department connects with their credential, receives
real-time video and position from all drones in area.
Grant expires automatically when incident closes.
```

**Why this works:**
- Fire department gets exactly what they need—live video and positions
- No access to commands—can't interfere with drone operations
- Time-bounded—access ends when incident ends
- Auditable—every frame they receive is logged

### Partner Data Sharing

```
Scenario: University research on environmental monitoring

Grant issued:
  - Subject: fleet.prod.veh.*.state.env.* (temp, humidity, air quality)
  - Deny: fleet.prod.veh.*.state.pos (no location data)
  - Duration: 30 days (research period)
  - Audit: Data access logged for compliance

Result: Researcher receives environmental sensor data without
knowing vehicle locations—privacy preserved.
```

**Why this works:**
- Researcher gets environmental data for their study
- Location data explicitly denied—privacy protected
- Clear audit trail for research compliance
- Automatic cleanup after research period

### Regulatory Compliance

```
Scenario: Aviation authority audit of flight operations

Grant issued:
  - Subject: fleet.prod.veh.*.state.* (all telemetry)
  - Subject: fleet.prod.veh.*.evt.* (all events)
  - Deny: fleet.prod.veh.*.cmd.* (no command capability)
  - Duration: 7 days (audit period)
  - Audit: Every query logged for audit trail

Result: Auditor gets complete read access to historical data
but cannot issue any commands to vehicles.
```

**Why this works:**
- Full transparency for auditor—nothing hidden
- Read-only—auditor cannot affect operations
- Complete audit trail of auditor's access
- Time-bounded—no permanent access hole

### Customer Integration

```
Scenario: Enterprise customer wants API access to their fleet data

Grant issued:
  - Subject: fleet.prod.veh.CUST-*.state.* (only their vehicles)
  - Subject: fleet.prod.veh.CUST-*.evt.*
  - Deny: fleet.prod.veh.OTHER-*.* (other customers' data)
  - Duration: Persistent (until revoked)
  - Audit: API access logged

Result: Customer builds their own dashboards using real-time
data from their vehicles—and only their vehicles.
```

### Subcontractor Access

```
Scenario: Maintenance team needs to run diagnostics on specific vehicles

Grant issued:
  - Subject: fleet.prod.veh.V001.cmd.diag.* (diagnostic commands only)
  - Subject: fleet.prod.veh.V001.state.health.*
  - Deny: fleet.prod.veh.V001.cmd.flight.* (no flight commands)
  - Duration: 8 hours (maintenance window)
  - Audit: All commands logged with technician ID

Result: Technician can run diagnostics but cannot arm or fly the vehicle.
Access expires after maintenance window.
```

---

## Security Properties

| Property | Implementation |
|----------|----------------|
| **Zero Trust** | No implicit trust—every connection authenticated |
| **Least Privilege** | Credentials contain minimum necessary permissions |
| **Defense in Depth** | Multiple layers: TLS + JWT + subject ACLs |
| **Offline Capable** | Authorization works without network connectivity |
| **Auditable** | Every access logged with credential identity |
| **Revocable** | Instant revocation via account revocation list |

---

## Operational Benefits

### For Fleet Operators

- **Add/remove third parties without system changes** — Issue or revoke credentials, no config files to edit
- **Grant temporary access without credential sharing** — Each party gets their own credential
- **Audit who accessed what and when** — Complete access log
- **Revoke access instantly if needed** — Takes effect on next connection

### For Third Parties

- **Self-service connection** — Use provided credentials, connect when ready
- **Clear scope** — Know exactly what you can and cannot access
- **No internal architecture knowledge needed** — Just connect and subscribe
- **Automatic expiration** — No cleanup required on your end

### For Compliance

- **Complete audit trail** — Every access traced to credential
- **Provable permission boundaries** — Demonstrate who had access to what
- **Time-bounded access for auditors** — No permanent access holes
- **No shared credentials** — Individual accountability

---

## Implementation Notes

### Credential Format

NATS uses **NKeys** and **JWTs**:

- **Ed25519 signatures** — Cryptographically secure, fast verification
- **Self-contained claims** — Permissions embedded in token
- **No external validation** — Server verifies signature locally

### Revocation

When you need to revoke access:

1. Add credential to account revocation list
2. List propagates to all servers in cluster
3. Takes effect on next connection attempt
4. No need to touch the grantee's systems

### Credential Rotation

For long-lived integrations:

1. Issue new credential before old one expires
2. Overlap period allows seamless transition
3. Revoke old credential after transition complete
4. No service interruption

---

## Why This Matters for Real Systems

Traditional authorization requires:
- Central auth server (single point of failure)
- Network call on every request (latency)
- Complex integration for third parties (months of work)
- Manual cleanup of expired access (forgotten credentials)

NATS decentralized authorization provides:
- **Resilience** — Works offline, no central bottleneck
- **Speed** — Local validation, no round-trip
- **Extensibility** — Add third parties in minutes, not months
- **Safety** — Automatic expiration, instant revocation

This is how you build systems that are both **highly secure** and **highly user-friendly**.

---

## Related Documentation

- [Subject Naming]({{< relref "/fleet/subjects" >}}) — How subject hierarchy enables fine-grained ACLs
- [NATS Topology]({{< relref "/fleet/nats-topology" >}}) — How the messaging infrastructure is structured
- [Safety Model]({{< relref "/fleet/safety" >}}) — How command authorization fits into safety hierarchy
- [Surveillance]({{< relref "/company/surveillance" >}}) — How grants enable third-party access for monitoring

---

## Next

[Streams & Events →]({{< relref "/fleet/streams" >}})

