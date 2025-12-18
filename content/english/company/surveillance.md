---
title: "Surveillance"
meta_title: "Real-Time Surveillance & Situational Awareness | Ubuntu Software"
description: "Turn firehose data from drones, satellites, and IoT sensors into actionable alerts. Self-sovereign, secure, with full provenance."
image: "/images/robotics.svg"
draft: false
---

## Drowning in Data

Modern operations generate more data than humans can process. Drones streaming video. Satellites capturing imagery. Sensors reporting continuously. The signal is there—buried in the noise.

Traditional approaches fail at scale:
- **Manual monitoring** doesn't scale beyond a few feeds
- **Simple thresholds** generate alert fatigue
- **Siloed systems** miss cross-source patterns
- **Cloud-only solutions** can't handle sensitive operations

You need a system that watches everything, understands context, and alerts you only when it matters.

---

## From Firehose to Insight

Our surveillance architecture processes data from multiple sources in real-time, applying AI at both edge and cloud, then filtering through operator-defined rules to surface actionable events.

```
┌─────────────────────────────────────────────────────────────────┐
│                         DATA SOURCES                            │
├─────────────────┬─────────────────┬─────────────────────────────┤
│     Drones      │   Satellites    │        IoT Sensors          │
│  (telemetry,    │  (imagery,      │  (environmental,            │
│   video, AI)    │   positioning)  │   industrial)               │
└────────┬────────┴────────┬────────┴────────────┬────────────────┘
         │                 │                      │
         ▼                 ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    FIREHOSE INGESTION                           │
│            NATS JetStream — persistent, ordered streams         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AI PROCESSING                              │
│   Edge (Jetson): Real-time inference, event detection           │
│   Cloud: Pattern recognition, correlation, anomaly detection    │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    OPERATOR RULES ENGINE                        │
│   "Alert me when..." — thresholds, patterns, conditions         │
│   Customizable per user, role, or organization                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                    ALERT + CONTEXT                              │
│   • What happened (event description)                           │
│   • Where (2D map, 3D spatial view)                             │
│   • Why it matters (AI reasoning)                               │
│   • Provenance (full data lineage, evidence chain)              │
└─────────────────────────────────────────────────────────────────┘
```

---

## Data Sources

We treat all data sources as first-class citizens—unified ingestion, consistent processing, correlated analysis.

| Source | Data Types | Examples |
|--------|------------|----------|
| **Drones** | Telemetry, video, sensor readings, AI inference results | Fleet position, thermal imagery, object detection |
| **Satellites** | Imagery, positioning, weather data | Change detection, coverage analysis, atmospheric conditions |
| **IoT Sensors** | Environmental, industrial, infrastructure | Temperature, vibration, water levels, air quality |

The same architecture that handles drone telemetry handles satellite imagery and ground sensors. One system. Complete picture.

---

## You Define What Matters

Operators configure what they care about. The system watches everything and alerts only on conditions you specify.

| Operator Says | System Does |
|---------------|-------------|
| "Alert me if any drone loses GPS" | Monitors all fleet telemetry, notifies on GPS loss |
| "Flag thermal anomalies on power lines" | Runs thermal analysis, alerts on hotspots |
| "Track vehicles entering restricted zones" | Geofence monitoring with instant alerts |
| "Notify me of crop stress above 20%" | Aggregates multispectral data, triggers threshold |
| "Correlate satellite change with ground sensors" | Cross-references multiple sources, detects patterns |

Rules can be simple thresholds or complex multi-source correlations. You control the sensitivity.

---

## Visualization with Provenance

*Roadmap capability*

When an alert fires, operators see the full picture:

- **2D Map View** — Location context with historical track
- **3D Spatial View** — Full environmental context, terrain, structures
- **Data Provenance** — Complete chain from raw sensor data to alert
- **Evidence Package** — Original data, processed results, AI confidence scores

**Every decision is auditable. Every alert is traceable.**

Operators can drill down from alert to raw data, understanding exactly why the system flagged an event. No black boxes.

---

## Security & Sovereignty

These systems handle sensitive operations. Security isn't optional—it's foundational.

| Principle | Implementation |
|-----------|----------------|
| **Self-Sovereign** | You own and control all data—we never access without explicit permission |
| **End-to-End Encryption** | Data encrypted in transit and at rest, keys you control |
| **Zero Trust** | Every request authenticated, every action authorized, no implicit trust |
| **Air-Gap Ready** | Full functionality without internet connectivity for sensitive deployments |
| **Complete Audit Trail** | Every access, every query, every alert logged with timestamps |
| **Data Residency** | Choose where your data lives—your infrastructure, your jurisdiction |

**Your data. Your infrastructure. Your control.**

Self-hosted deployments run entirely on your hardware. Managed deployments isolate your data with dedicated encryption keys. Either way, you maintain sovereignty.

---

## Use Cases

| Domain | Application |
|--------|-------------|
| **Border Security** | Multi-sensor perimeter monitoring, intrusion detection, response coordination |
| **Critical Infrastructure** | Pipeline and powerline surveillance, predictive maintenance, anomaly detection |
| **Environmental** | Wildlife tracking, pollution monitoring, wildfire detection, ecosystem health |
| **Agriculture** | Crop health monitoring, irrigation optimization, pest/disease early warning |
| **Maritime** | Vessel tracking, illegal fishing detection, port security, search and rescue |
| **Urban Operations** | Traffic flow analysis, crowd monitoring, emergency response coordination |

Each use case leverages the same core architecture—different sensors, different rules, same reliable platform.

---

## Technical Foundation

This surveillance capability is built on our production-grade fleet architecture:

- [Fleet Architecture]({{< relref "/fleet" >}}) — How data flows from edge to cloud
- [NATS Topology]({{< relref "/fleet/nats-topology" >}}) — The messaging backbone that handles the firehose
- [Streams & Events]({{< relref "/fleet/streams" >}}) — How we turn continuous streams into queryable state
- [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) — Edge processing and event detection
- [Authorization & Grants]({{< relref "/fleet/authorization" >}}) — How third parties get scoped, secure access

The same infrastructure that manages drone fleets powers surveillance operations. Battle-tested at scale.

---

## Get Started

Ready to turn your data firehose into actionable intelligence?

[Contact Us →](/contact) to discuss your surveillance requirements.

We'll help you understand:
- Which data sources to integrate
- What operator rules make sense for your domain
- Self-hosted vs. managed deployment options
- Security and compliance considerations

