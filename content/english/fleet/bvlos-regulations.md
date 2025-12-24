---
title: "BVLOS Regulations"
meta_title: "Beyond Visual Line of Sight Drone Regulations | Ubuntu Software"
description: "Global BVLOS (Beyond Visual Line of Sight) drone regulations by country. Understanding regulatory requirements for commercial fleet operations."
image: "/images/robotics.svg"
draft: false
---

## Why BVLOS Matters for Fleet Operations

Beyond Visual Line of Sight (BVLOS) regulations fundamentally determine what commercial drone fleets can accomplish. Understanding the regulatory landscape is essential for fleet platform design and market strategy.

| Challenge | Impact on Fleet Platform |
|-----------|-------------------------|
| **Operational Limits** | VLOS restrictions (400-500m range) limit fleet utility. BVLOS enables infrastructure inspection, delivery, and agriculture at scale. |
| **Compliance Integration** | Fleet software must enforce regulatory boundaries—Remote ID, detect-and-avoid, and operational authorization requirements inform telemetry and safety features. |
| **Market Timing** | Countries are at different stages of BVLOS enablement. Understanding where regulations are opening helps prioritize target markets. |
| **Partner Readiness** | Fleet operators need BVLOS approvals for real missions. Tracking who has approvals identifies potential customers and partners. |

---

## Global Trends

Across all major aviation authorities, four trends are emerging:

| Trend | Description |
|-------|-------------|
| **SORA Adoption** | FAA, EASA, CASA, and Transport Canada using Specific Operations Risk Assessment methodology [[1]](#ref-1) |
| **Remote ID Mandates** | Becoming mandatory across all jurisdictions |
| **Waiver → Rules** | Transitioning from case-by-case waivers to standardized performance-based rules |
| **UTM Integration** | Countries developing UAS Traffic Management for scalable BVLOS |

---

## Americas

### United States (FAA)

**Status:** Waiver-based, moving to rule-based (Part 108) [[2]](#ref-2)

**Timeline** [[3]](#ref-3):

- 190 BVLOS waivers issued as of Oct 2024 (up from 6 in 2020)
- Part 108 proposed rule expected Jan 2026
- FAA Reauthorization Act 2024 extended BEYOND program to 2029 [[4]](#ref-4)

**Operators with BVLOS Approvals:**

| Operator | Use Case | Approval Date | Ref |
|----------|----------|---------------|-----|
| DroneUp | Medical deliveries | Jan 2024 | [[5]](#ref-5) |
| Amazon Prime Air | Consumer delivery | May 2024 | [[6]](#ref-6) |
| Asylon | Security/surveillance | Apr 2024 | [[7]](#ref-7) |
| Multiple (DFW area) | First multi-operator airspace | Jul 2024 | [[8]](#ref-8) |

---

### Canada (Transport Canada)

**Status:** SFOC required, moving to routine BVLOS [[9]](#ref-9)

**Framework:**

- Apr 2025: Routine BVLOS permitted without SFOC in low-risk conditions (≤150kg, uncontrolled airspace, sparse population)
- Requires Advanced Pilot Certificate
- National UTM framework via Nav Canada by 2030

**Operators:** TBD

---

## Europe

### European Union (EASA)

**Status:** Specific category with SORA/PDRA pathways [[10]](#ref-10)

**Framework** [[1]](#ref-1):

| Aspect | Details |
|--------|---------|
| Standard Scenario | STS-02: BVLOS with airspace observers |
| Remote ID | Mandatory since Jan 2024 [[11]](#ref-11) |
| Approval Timeline | 4-9 months (medical), 6-12 months (consumer) [[12]](#ref-12) |
| Fees | €2,000-€15,000 |

**Operators:** Various operators via national sandboxes and fast-track member state programs.

---

### United Kingdom (CAA)

**Status:** Operational Authorization (OA) via UK SORA

**Framework:**

- BVLOS requires CAA Operational Authorization
- UK SORA process (since Apr 2025)
- Detect-and-avoid technology typically required (ADS-B, FLARM, ground radar)

**Operators:**

| Operator | Achievement | Ref |
|----------|-------------|-----|
| sees.ai | First UK routine BVLOS (2021), National Grid inspections (2023) | [[13]](#ref-13) |

---

### Ireland (IAA)

**Status:** PDRA G-03 pathway

**Framework:** Opens door to streamlined EU authorizations via EASA.

**Operators:**

| Operator | Achievement | Ref |
|----------|-------------|-----|
| sees.ai | First BVLOS permission in Ireland (2025) | [[13]](#ref-13) |

---

## Asia-Pacific

### Australia (CASA)

**Status:** SORA-based assessment under Part 101 [[14]](#ref-14)

**Framework** [[15]](#ref-15):

- Adopting SORA 2.5 with quantitative risk modeling
- Broad Area BVLOS Self-Assessment framework available
- Streamlined process attracting global companies [[16]](#ref-16)

**Operators:**

| Operator | Achievement | Ref |
|----------|-------------|-----|
| Sphere Drones | Broad Area BVLOS approval (May 2024) | [[17]](#ref-17) |

---

### Japan (MLIT/JCAB)

**Status:** Level 4 BVLOS over populated areas permitted [[18]](#ref-18)

**Framework** [[19]](#ref-19):

| Level | Description | License Required |
|-------|-------------|------------------|
| Level 1-2 | VLOS operations | Second Class |
| Level 3 | BVLOS over uninhabited areas | Second Class |
| Level 4 | BVLOS over urban/populated areas | First Class |

- Licenses valid 3 years, require registered drone school + flight test
- 422,879 UAVs registered as of Oct 2024

**Operators:**

- Only 1 UAV type has obtained First Class (Level 4) certification as of Oct 2024
- 5 UAV types have Second Class certification

---

### Singapore (CAAS)

**Status:** BVLOS prohibited without Operator + Class 1 Activity Permit [[20]](#ref-20)

**Framework** [[21]](#ref-21):

- Strict authorization required for any BVLOS
- Drones >25kg or BVLOS require full Operator Permit + Class 1 permit
- B-RID (Broadcast Remote ID) mandatory from Dec 2025 for drones >250g
- B-RID enables foundation for future scaled BVLOS/automated ops
- Regulatory sandbox programs available for fast-track approvals

**Operators:** TBD (sandbox trials ongoing)

---

### China (CAAC)

**Status:** Special authorization required for BVLOS [[22]](#ref-22)

**Framework** [[23]](#ref-23):

- Interim Regulations on Management of Unmanned Aircraft Flights (Jan 2024)
- BVLOS requires Advanced Operations Certificate + Special Flight Operations Certificate (SFOC)
- 120m altitude limit nationwide; above requires authorization
- Simplified digital application processes enabling more BVLOS flights [[16]](#ref-16)

**Operators:** TBD (domestic operators under CAAC framework)

---

### India (DGCA)

**Status:** No formal BVLOS framework yet; trial phase only [[24]](#ref-24)

**Framework** [[25]](#ref-25):

- Operating under Drone Rules, 2021 with project-based exemptions
- Apr 2025: Separate BVLOS pilot license category introduced for urban/semi-urban
- DGCA "in advance stage of finalising BVLOS rules" [[26]](#ref-26)
- Experimental corridors active (e.g., "Medicine from the Sky" in Telangana)

**Operators:**

| Operator | Achievement | Ref |
|----------|-------------|-----|
| Skye Air Mobility | World's longest BVLOS medicine delivery (104km, Baruipur-Medinipur) | [[27]](#ref-27) |
| Flipkart Health | Partner on 104km delivery trial | — |

---

## Compliance Requirements

Across all jurisdictions, BVLOS regulations converge on common technical requirements. This section distills what our fleet platform must implement.

### 1. Remote ID (Mandatory Everywhere)

All major regulators require drones to broadcast identification and location data.

| Jurisdiction | Requirement | Deadline | Ref |
|--------------|-------------|----------|-----|
| USA (FAA) | Standard Remote ID or broadcast module | Sep 2023 (in effect) | [[2]](#ref-2) |
| EU (EASA) | Remote ID for Specific category + Class 1+ | Jan 2024 (in effect) | [[11]](#ref-11) |
| Singapore | B-RID for drones >250g | Dec 2025 | [[21]](#ref-21) |
| Australia | Remote ID under consideration | TBD | [[14]](#ref-14) |

**Platform Requirements:**
- Broadcast module hardware integration (OpenDroneID compatible)
- Network Remote ID capability via NATS for fleet-wide tracking
- Logging of all Remote ID transmissions to JetStream for audit

---

### 2. Detect-and-Avoid (DAA)

BVLOS operations require ability to detect and avoid other aircraft and obstacles.

| Jurisdiction | Requirement | Notes |
|--------------|-------------|-------|
| USA (FAA) | "Well clear" standard; DAA system required for Part 108 | Performance-based, not prescriptive |
| UK (CAA) | DAA technology typically required (ADS-B, FLARM, ground radar) | Case-by-case in OA process |
| Japan | Required for Level 4 urban operations | [[18]](#ref-18) |

**Platform Requirements:**
- ADS-B In receiver integration (traffic awareness)
- Sensor fusion on Jetson (camera, radar, lidar options)
- Autonomous avoidance maneuvers or operator alerting
- Logging of all DAA events to JetStream

---

### 3. Geofencing

Hardware-enforced boundaries preventing flight into restricted areas.

| Jurisdiction | Requirement | Notes |
|--------------|-------------|-------|
| USA (FAA) | Geofence part of operational approval | Defined in waiver/Part 108 |
| EU (EASA) | Geo-awareness required in Specific category | Per SORA ground risk assessment |
| All | Altitude limits (typically 120m/400ft AGL) | Fundamental to all frameworks |

**Platform Requirements:**
- Geofence enforcement on flight controller (PX4 `GF_*` parameters)
- Dynamic geofence updates via Gateway (temporary restrictions)
- Altitude ceiling enforcement
- Breach logging and alerting to fleet

---

### 4. Command & Control (C2) Link

Reliable, redundant communication between operator and aircraft.

| Jurisdiction | Requirement | Notes |
|--------------|-------------|-------|
| USA (FAA) | C2 link performance standards in Part 108 | Latency, availability, security |
| EU (EASA) | C2 link reliability per SORA | Risk-dependent requirements |
| All | Failsafe behavior on C2 loss | Return-to-home, land, or loiter |

**Platform Requirements:**
- Primary C2 via LTE/5G or dedicated radio
- Backup C2 path (e.g., satellite for long-range)
- Automatic failsafe on link loss (RTL configured in PX4)
- C2 link status telemetry to fleet
- Link quality logging for audit

---

### 5. Operational Authorization Documentation

Regulators require documented risk assessment and operational procedures.

| Framework | Used By | Key Documents |
|-----------|---------|---------------|
| SORA | EASA, UK CAA, CASA, Transport Canada | ConOps, Ground Risk, Air Risk, OSO mitigations |
| Part 107 Waiver | FAA (current) | Safety case, operational limits |
| Part 108 | FAA (coming) | Performance-based compliance |

**Platform Requirements:**
- Export flight logs in regulator-accepted formats
- Generate ConOps-supporting data (flight areas, times, frequencies)
- Maintain operator/pilot certification records
- Track aircraft maintenance and inspection status

---

### 6. Pilot/Operator Licensing

| Jurisdiction | BVLOS License | Notes |
|--------------|---------------|-------|
| USA | Part 107 + waiver (current), Part 108 (future) | [[2]](#ref-2) |
| EU | A2/STS certificate or Specific authorization | [[1]](#ref-1) |
| UK | GVC + Operational Authorization | |
| Japan | First Class (Level 4), Second Class (Level 3) | [[19]](#ref-19) |
| Canada | Advanced Pilot Certificate + SFOC | [[9]](#ref-9) |
| India | BVLOS pilot license (Apr 2025+) | [[25]](#ref-25) |

**Platform Requirements:**
- Pilot authentication before flight authorization
- License/certification verification and expiry tracking
- Per-pilot flight hour logging
- Role-based access (pilot vs. observer vs. fleet manager)

---

### 7. Flight Logging & Audit Trail

All jurisdictions require comprehensive flight records.

**Platform Requirements:**
- Complete telemetry recording (position, attitude, battery, sensors)
- Timestamped event logging (arm, takeoff, waypoints, land, disarm)
- Immutable storage (JetStream with retention policy)
- Export capability for regulatory submission
- Incident flagging and investigation support

---

### 8. Insurance & Liability

| Jurisdiction | Requirement |
|--------------|-------------|
| EU | Third-party liability insurance mandatory |
| USA | Insurance typically required by waiver conditions |
| All | Proof of insurance often required for authorization |

**Platform Requirements:**
- Insurance policy tracking per aircraft
- Coverage verification before flight authorization
- Incident documentation for claims

---

## External Integrations

BVLOS compliance requires real-time data exchange with external systems. These are the integrations our platform needs to support.

### Airspace & Traffic

| System | Purpose | Data Flow | Priority |
|--------|---------|-----------|----------|
| **ADS-B In** | Traffic awareness from manned aircraft | Receive → DAA system | High |
| **FLARM** | Traffic awareness from gliders/GA | Receive → DAA system | Medium |
| **UTM/USS** | UAS traffic management, flight authorization | Bidirectional | High |
| **LAANC** (USA) | Automated airspace authorization | Request/Response | High (USA) |
| **AIM/NOTAM** | Temporary flight restrictions, airspace changes | Receive → Geofence | High |

**UTM Service Suppliers (USS):**
- USA: Airspace Link, Aloft, DroneUp, FlyFreely (see [FAA LAANC list](https://www.faa.gov/uas/getting_started/laanc))
- EU: Various per member state
- Australia: Wing, Airservices Australia

*Note: AirMap shut down LAANC service in June 2023.*

---

### Remote ID Networks

| System | Purpose | Protocol | Priority |
|--------|---------|----------|----------|
| **OpenDroneID** | Broadcast Remote ID standard | Bluetooth 4/5, WiFi Beacon | Required |
| **Network Remote ID** | Fleet-wide ID via internet | USS APIs | Required (USA) |
| **FAA Remote ID** | US Remote ID compliance | USS integration | Required (USA) |

---

### Regulatory Databases

| System | Purpose | Data Flow | Priority |
|--------|---------|-----------|----------|
| **Pilot License DB** | Verify pilot certifications | Query → Authorization | High |
| **Aircraft Registry** | Verify registration status | Query → Pre-flight | High |
| **Insurance Verification** | Confirm coverage validity | Query → Pre-flight | Medium |
| **Maintenance Records** | Track airworthiness | Bidirectional | Medium |

---

### Weather & Environment

| System | Purpose | Data Flow | Priority |
|--------|---------|-----------|----------|
| **Aviation Weather (METAR/TAF)** | Wind, visibility, ceiling | Receive → Flight planning | High |
| **Radar/Precipitation** | Real-time weather hazards | Receive → In-flight alerts | Medium |
| **Solar/Geomagnetic** | GPS reliability forecasts | Receive → Flight planning | Low |

---

### Emergency Services

| System | Purpose | Data Flow | Priority |
|--------|---------|-----------|----------|
| **ATC (if required)** | Controlled airspace coordination | Voice/Data link | Situational |
| **Emergency Services (911/112)** | Incident reporting | Outbound alerts | Required |
| **NOTAM Filing** | Publish flight operations | Outbound → AIM | Situational |

---

### Integration Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         FLEET PLATFORM                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌─────────────┐    ┌─────────────┐    ┌─────────────┐             │
│   │   Gateway   │    │  JetStream  │    │   Fleet     │             │
│   │  (Vehicle)  │    │  (Logging)  │    │  Manager    │             │
│   └──────┬──────┘    └──────┬──────┘    └──────┬──────┘             │
│          │                  │                  │                     │
└──────────┼──────────────────┼──────────────────┼─────────────────────┘
           │                  │                  │
           ▼                  ▼                  ▼
┌──────────────────────────────────────────────────────────────────────┐
│                      INTEGRATION LAYER                               │
├──────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │   UTM    │  │  Remote  │  │ Weather  │  │ License  │             │
│  │   USS    │  │    ID    │  │   API    │  │   DB     │             │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘             │
│       │             │             │             │                    │
└───────┼─────────────┼─────────────┼─────────────┼────────────────────┘
        │             │             │             │
        ▼             ▼             ▼             ▼
   ┌─────────┐   ┌─────────┐   ┌─────────┐   ┌─────────┐
   │Airspace │   │   FAA   │   │  NOAA   │   │  CAA    │
   │  Link   │   │ DroneZone│  │ Aviation│   │ Registry│
   │ Aloft   │   │ NASA DIP│   │ Weather │   │         │
   └─────────┘   └─────────┘   └─────────┘   └─────────┘
```

---

### Integration Roadmap

| Integration | Status | Notes |
|-------------|--------|-------|
| OpenDroneID broadcast | ✅ Supported | Hardware module required |
| ADS-B In receiver | 🔧 Roadmap | uAvionix pingRX or similar |
| UTM/USS (Airspace Link) | 🔧 Roadmap | AirHub API, tiered pricing |
| UTM/USS (Aloft) | 🔧 Roadmap | REST API integration |
| LAANC | 🔧 Roadmap | Via USS partner |
| Aviation weather | 🔧 Roadmap | aviationweather.gov (open, 100 req/min) |
| NOTAM feed | 🔧 Roadmap | NASA DIP or FAA SWIM |
| Pilot license verification | 🔧 Roadmap | Per-jurisdiction APIs |
| Insurance verification | 🔧 Roadmap | Provider APIs |

---

## Compliance Matrix

| Requirement | Platform Component | Status |
|-------------|-------------------|--------|
| Remote ID broadcast | OpenDroneID module + Gateway | ✅ Supported |
| Remote ID logging | JetStream streams | ✅ Supported |
| Detect-and-Avoid | Jetson sensor fusion | 🔧 Roadmap |
| ADS-B In | Traffic receiver integration | 🔧 Roadmap |
| Geofencing (static) | PX4 `GF_*` parameters | ✅ Supported |
| Geofencing (dynamic) | Gateway policy updates | ✅ Supported |
| C2 link primary | LTE/5G via NATS | ✅ Supported |
| C2 link backup | Satellite integration | 🔧 Roadmap |
| Failsafe on link loss | PX4 failsafe config | ✅ Supported |
| Flight logging | JetStream telemetry streams | ✅ Supported |
| Audit export | Log export tooling | ✅ Supported |
| Pilot authentication | Authorization system | ✅ Supported |
| License tracking | Fleet management | 🔧 Roadmap |
| Insurance tracking | Fleet management | 🔧 Roadmap |

See [Safety Model]({{< relref "/fleet/safety" >}}) for how these integrate with the overall safety architecture.

---

## References

<span id="ref-1"></span>**[1]** EASA. "Drones (UAS) - Frequently Asked Questions." European Union Aviation Safety Agency. https://www.easa.europa.eu/en/the-agency/faqs/drones-uas

<span id="ref-2"></span>**[2]** Drone Pilot Ground School. "Part 108: What We Know So Far about the FAA's BVLOS Rule." https://www.dronepilotgroundschool.com/part-108/

<span id="ref-3"></span>**[3]** DroneLife. "FAA Progress on BVLOS Rules and Advanced Air Mobility Integration." Aug 2024. https://dronelife.com/2024/08/01/faa-progress-on-bvlos-rules-and-advanced-air-mobility-integration/

<span id="ref-4"></span>**[4]** FAA. "BEYOND Program." Federal Aviation Administration. https://www.faa.gov/uas/programs_partnerships/beyond

<span id="ref-5"></span>**[5]** DroneUp. "DroneUp Awarded Landmark FAA Approval for Beyond Visual Line of Sight (BVLOS) in the U.S." Jan 2024. https://www.droneup.com/news/droneup-awarded-bvlos

<span id="ref-6"></span>**[6]** DroneDJ. "Amazon gets FAA BVLOS approval for Prime Air drone delivery." May 2024. https://dronedj.com/2024/05/30/amazon-gets-faa-bvlos-approval-for-prime-air-drone-delivery/

<span id="ref-7"></span>**[7]** DroneDJ. "Asylon bags new FAA approval for BVLOS drone operations." Apr 2024. https://dronedj.com/2024/04/26/asylon-drone-faa-bvlos-dronesentry/

<span id="ref-8"></span>**[8]** DroneLife. "FAA Grants Historic BVLOS Approval for Multiple Operators in DFW Area." Jul 2024. https://dronelife.com/2024/07/31/faa-grants-historic-bvlos-approval-for-multiple-operators-in-dfw-area/

<span id="ref-9"></span>**[9]** MLT Aikins. "International drone regulations: Key updates on BVLOS, UTM and use cases." https://www.mltaikins.com/insights/international-drone-regulations-key-updates-on-bvlos-utm-and-use-cases/

<span id="ref-10"></span>**[10]** FlytBase. "Understanding EU Drone Regulations for Docked BVLOS Operations." https://www.flytbase.com/blog/eu-drone-regulations-for-bvlos-operations

<span id="ref-11"></span>**[11]** Elsight. "EASA Remote ID Requirements: Compliance for Drone Operators." https://www.elsight.com/blog/primer-on-easa-remote-id-regulations/

<span id="ref-12"></span>**[12]** Sparkco. "Autonomous Drone Delivery Regulatory Approval." https://sparkco.ai/blog/autonomous-drone-delivery-regulatory-approval

<span id="ref-13"></span>**[13]** sees.ai. "Why Us - First UK routine BVLOS permission." https://www.sees.ai/why-us/

<span id="ref-14"></span>**[14]** CASA. "BVLOS drone operations in regional Australia - Consultation." Civil Aviation Safety Authority. https://consultation.casa.gov.au/stakeholder-engagement-group/consultation.2023-10-05.5578154857/

<span id="ref-15"></span>**[15]** Measure Australia. "New CASA Drone Regulations Incoming!" https://www.measureaustralia.com.au/news/new-casa-drone-regulations-incoming

<span id="ref-16"></span>**[16]** Droneii. "Commercial Drone Regulation Report 2025." Drone Industry Insights. https://droneii.com/product/commercial-drone-regulation-report

<span id="ref-17"></span>**[17]** DroneLife. "Sphere Drones Secures BVLOS Area Approval to Advance Commercial Drone Use in Australia." May 2024. https://dronelife.com/2024/05/08/sphere-drones-secures-bvlos-area-approval-to-advance-commercial-drone-use-in-australia/

<span id="ref-18"></span>**[18]** Unmanned Airspace. "Japan's revised drone laws to permit BVLOS flights over people come into effect." https://www.unmannedairspace.info/emerging-regulations/japanese-revised-drone-laws-to-permit-bvlos-flights-over-people-come-into-effect/

<span id="ref-19"></span>**[19]** Fly Eye. "Japan Drone Laws (2025)." https://www.flyeye.io/japan-drone-laws/

<span id="ref-20"></span>**[20]** Drone Laws. "Singapore Drone Laws 2025." https://drone-laws.com/drone-laws-in-singapore/

<span id="ref-21"></span>**[21]** Heron AirBridge. "B-RID in Singapore: Full Guide to the New Drone Regulations." https://heron-airbridge.com/singapore-brid-regulation/

<span id="ref-22"></span>**[22]** Fly Eye. "Advanced Drone Operations in China." https://www.flyeye.io/chinese-drone-regulations-advanced-operations/

<span id="ref-23"></span>**[23]** Fly Eye. "Chinese Drone Laws (2025)." https://www.flyeye.io/chinese-drone-laws/

<span id="ref-24"></span>**[24]** Commercial UAV News. "India Makes a Giant Leap Toward BVLOS Operations." https://www.commercialuavnews.com/drone-delivery/india-makes-a-giant-leap-toward-bvlos-operations

<span id="ref-25"></span>**[25]** Fly and Tech. "Drone Rules in India 2025 – Latest DGCA Guidelines Update." https://flyandtech.com/drone-rules-india-2025/

<span id="ref-26"></span>**[26]** Sigma Chambers. "Drone Law Brief | Vol 1: September & October 2025." https://www.sigmachambers.in/post/drone-law-watch-vol-1-september-october-2025

<span id="ref-27"></span>**[27]** ITLN. "How will BVLOS drone delivery ops impact the Indian logistics industry?" https://www.itln.in/cargo-drones/how-will-bvlos-drone-delivery-ops-impact-the-indian-logistics-industry-1353269

---

## Back to Overview

[← Autonomous Vehicle Fleet Architecture]({{< relref "/fleet" >}})
