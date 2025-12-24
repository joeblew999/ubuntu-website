---
title: "US Airspace Viewer"
meta_title: "Interactive US Airspace Map | Ubuntu Software"
description: "Interactive map visualizing FAA airspace data including controlled airspace, special use airspace, and LAANC ceiling altitudes for drone operations."
image: "/images/robotics.svg"
draft: false
---

## Interactive Airspace Visualization

Understanding airspace structure is fundamental to BVLOS drone operations. This demo visualizes FAA airspace data to show where drones can fly and what restrictions apply.

<a href="/airspace-demo/" target="_blank" class="btn btn-primary">Open Full-Screen Demo →</a>

---

## What You'll See

The map displays three layers of FAA airspace data:

| Layer | Description | Data Source |
|-------|-------------|-------------|
| **Airspace Boundaries** | Class B/C/D/E controlled airspace around airports | FAA UDDS |
| **Special Use Airspace** | Restricted areas, Prohibited areas, MOAs | FAA UDDS |
| **LAANC Ceilings** | Pre-approved altitude limits for automated authorization | FAA UAS Facility Map |

### Airspace Classes

| Class | Color | Description |
|-------|-------|-------------|
| **B** | Blue | Busiest airports (ATL, LAX, JFK, etc.) |
| **C** | Purple | Towered airports with radar service |
| **D** | Cyan | Smaller towered airports |
| **E** | Green | Controlled airspace (surface or above) |

### Special Use Airspace

| Type | Color | Restrictions |
|------|-------|--------------|
| **Restricted (R)** | Red | Military operations, weapons testing |
| **Prohibited (P)** | Red | No flight (White House, nuclear plants) |
| **MOA** | Orange | Military Operations Areas (active during exercises) |

---

## LAANC Integration

The **Low Altitude Authorization and Notification Capability** (LAANC) grid shows pre-approved ceiling altitudes for Part 107 operations. When a grid cell shows a ceiling (e.g., 200 ft), drone operators can get near-instant authorization up to that altitude through a LAANC provider.

| Ceiling | Meaning |
|---------|---------|
| 0 ft | No LAANC authorization available (requires manual approval) |
| 50-400 ft | Automated authorization available up to that altitude |

---

## Data Source

All data comes from the **FAA** via their public ArcGIS REST APIs:

| Dataset | Source | Size |
|---------|--------|------|
| [UAS Facility Map](https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer) | FAA UAS Data Exchange | 2.2 MB |
| [Airspace Boundary](https://adds-faa.opendata.arcgis.com/datasets/67885972e4e940b2aa6d74024901c561) | FAA ADDS Open Data | 14 MB |
| [Special Use Airspace](https://adds-faa.opendata.arcgis.com/datasets/dd0d1b726e504137ab3c41b21835d05b) | FAA ADDS Open Data | 28 MB |

Data is refreshed periodically and stored in Cloudflare R2 for fast global delivery (files exceed Cloudflare Pages' 25MB limit).

---

## Platform Integration

Our fleet platform integrates airspace data for:

- **Pre-flight planning** — Validate routes against airspace restrictions
- **Geofence enforcement** — Hardware-level boundaries on flight controller
- **LAANC automation** — Request authorization through USS providers
- **Dynamic restrictions** — Update geofences for TFRs and NOTAMs

See [BVLOS Regulations]({{< relref "/fleet/bvlos-regulations" >}}) for how airspace compliance fits into the overall regulatory framework.

---

## Back to Overview

[← Autonomous Vehicle Fleet Architecture]({{< relref "/fleet" >}})
