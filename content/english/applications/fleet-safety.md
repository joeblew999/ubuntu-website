---
title: "Public Safety & Emergency"
meta_title: "IoT Fleet for Public Safety | Ubuntu Software"
description: "Public safety with IoT fleets: drones for aerial surveillance, emergency vehicles with real-time dispatch, and building sensors for threat detection. Unified architecture for coordinated response."
image: "/images/robotics.svg"
draft: false
---

## When Minutes Matter

Emergency response depends on information. IoT fleets provide situational awareness when responders need it most—faster deployment, better coverage, safer operations.

---

## Fleet Types in Public Safety

| Fleet Type | Devices | Role |
|------------|---------|------|
| **Aerial** | Surveillance drones, search drones | Search and rescue, incident overwatch, fire mapping |
| **Ground Vehicles** | Emergency vehicles, patrol cars | Real-time dispatch, telematics, asset tracking |
| **Fixed IoT** | Building alarms, traffic sensors, gunshot detectors | Threat detection, situational awareness, evidence |

All device types share the same architecture: edge AI, NATS JetStream messaging, and unified fleet management.

---

## The Challenge

Emergency situations suffer from information gaps:

- **Search and rescue** — Vast areas to cover, limited personnel
- **Disaster assessment** — Infrastructure damaged, access blocked
- **Active incidents** — Need eyes on scene before personnel arrive
- **Coordination** — Multiple agencies, fragmented communication
- **Documentation** — Evidence preservation, incident reconstruction

Ground-based assessment is slow and puts responders at risk. Helicopters are expensive and limited. The first minutes often determine outcomes.

---

## How IoT Fleets Solve This

### Aerial: Rapid Deployment

Pre-positioned drones launch immediately:

- **Automated dispatch** — Drones en route while crews mobilize
- **First on scene** — Visual assessment before personnel arrive
- **Live streaming** — Command sees what the drone sees
- **Night capability** — Thermal imaging in total darkness

### Ground: Vehicle Intelligence

Emergency vehicles with real-time coordination:

- **Automatic dispatch** — Closest available unit based on live positions
- **Route optimization** — Real-time traffic, fastest path to scene
- **Telematics** — Vehicle status, fuel, maintenance alerts
- **Resource tracking** — Who has what equipment, where

### Fixed IoT: Threat Detection

Sensors providing continuous awareness:

- **Gunshot detection** — Acoustic sensors pinpoint location
- **Traffic sensors** — Congestion, accidents, signal priority
- **Building alarms** — Fire, intrusion, elevator emergency
- **Weather stations** — Local conditions affecting response

### Area Coverage

Systematic search patterns cover ground fast:

- **Grid searches** — Automated flight paths, no gaps
- **Multi-drone coordination** — Cover more area simultaneously
- **Persistent surveillance** — Hours of coverage per battery swap
- **All-weather operation** — Wind, rain, cold—within limits

---

## Onboard AI (Edge Processing)

Real-time inference on NVIDIA Jetson for immediate detection:

| Capability | Public Safety Application |
|------------|--------------------------|
| **Person Detection** | Locate subjects in wilderness, disaster debris, crowds |
| **Vehicle Recognition** | Track vehicles of interest, identify abandoned cars |
| **Thermal Imaging** | Find heat signatures through smoke, at night, in foliage |
| **Fire Detection** | Identify active flames, hotspots, fire spread direction |
| **Structural Assessment** | Detect building damage, collapse risk, access routes |

**Why Edge Matters:**
- Alert responders immediately upon detection
- Work in areas with no cellular coverage
- Process thermal and visual simultaneously
- Continue operating during network congestion

---

## Cloud AI (Fleet Analytics)

Via NATS JetStream, coordinate response across agencies:

| Capability | Public Safety Application |
|------------|--------------------------|
| **Multi-Agency Coordination** | Share feeds across fire, police, EMS, SAR |
| **Resource Dispatch** | Optimal routing based on real-time conditions |
| **Pattern Analysis** | Track incident progression, predict spread |
| **Evidence Management** | Chain of custody, timestamp verification |
| **After-Action Review** | Incident reconstruction, training material |

**Why Cloud Matters:**
- Unified view for command centers
- Coordination across jurisdictions
- Historical data for similar incidents
- Integration with CAD/dispatch systems

---

## Hardware Configuration

For public safety deployments, we recommend:

| Component | Selection | Rationale |
|-----------|-----------|-----------|
| **Airframe** | Holybro X500 or ruggedized variant | Reliability in harsh conditions |
| **Flight Controller** | Pixhawk 6X | Proven autonomy, fail-safe modes |
| **Sensor Companion** | Raspberry Pi CM4 | Telemetry, basic sensors |
| **AI Companion** | Jetson Orin NX | High performance for thermal + visual |
| **Payload** | Dual RGB + thermal camera | Day and night capability |
| **Connectivity** | LTE + mesh radio | Resilience when towers are down |

[Full hardware specifications →]({{< relref "/fleet/hardware" >}})

---

## Use Cases

### Search and Rescue
Missing person searches. Wilderness SAR. Water rescue support.

### Fire Response
Structure fire assessment. Wildfire mapping. Hot spot detection.

### Law Enforcement
Active incident overwatch. Event security. Traffic incident management.

### Disaster Response
Damage assessment. Survivor location. Resource staging guidance.

### Border and Maritime
Coastal patrol. Port security. Border surveillance.

---

## Deployment Models

### Pre-Positioned Fleets
Drones staged at fire stations, police precincts, SAR bases. Launch on dispatch.

### Mobile Command
Vehicle-mounted systems. Deploy anywhere. Rapid setup.

### Persistent Coverage
Long-duration missions. Battery swap or charging stations. Continuous operation.

---

## Technical Deep Dive

This application is built on our production-grade fleet architecture:

| Component | Role in Public Safety |
|-----------|----------------------|
| [NATS Topology]({{< relref "/fleet/nats-topology" >}}) | Resilient messaging during infrastructure damage |
| [Vehicle Gateway]({{< relref "/fleet/gateway" >}}) | Real-time detection with immediate alerting |
| [JetStream Streams]({{< relref "/fleet/streams" >}}) | Evidence-grade video storage with timestamps |
| [Safety Model]({{< relref "/fleet/safety" >}}) | Operation near people, structures, aircraft |

[Explore Full Architecture →]({{< relref "/fleet" >}})

---

## Compliance

Public safety operations require:

- **FAA Part 107** waiver support for night, BVLOS operations
- **CJIS compliance** for law enforcement data handling
- **Evidence chain of custody** with cryptographic verification
- **Multi-agency data sharing** with access controls

Our architecture supports these requirements by design.

---

## Get Started

When seconds count, infrastructure matters. Deploy drone fleets that respond faster and see more.

[Contact Us →](/contact)
