# BVLOS Integration TODO

Go libraries and APIs for BVLOS compliance. See [BVLOS Regulations](content/english/fleet/bvlos-regulations.md) for full regulatory requirements.

---

## §1 Remote ID

### What It Is

Remote ID is the "license plate" for drones. The aircraft continuously broadcasts its identity, location, altitude, velocity, and operator location so that authorities and other airspace users can identify it.

### Why It Matters

Without Remote ID, regulators can't tell who's flying what. It's the foundation for integrating drones into shared airspace - law enforcement, airports, and other pilots need to know who that drone belongs to. Required for any commercial BVLOS operation.

### How It Works

Two methods:
1. **Broadcast Remote ID** - drone transmits via Bluetooth/WiFi beacon (local, ~1km range)
2. **Network Remote ID** - drone reports to internet service (global, requires connectivity)

OpenDroneID is the open standard. Messages include BASIC_ID (serial number), LOCATION (position/velocity), SYSTEM (operator location), OPERATOR_ID.

### Library: gomavlib

| | |
|-|-|
| **Repo** | https://github.com/bluenviron/gomavlib |
| **Purpose** | MAVLink library with OpenDroneID support |
| **Messages** | BASIC_ID, LOCATION, SYSTEM, OPERATOR_ID, etc. |
| **Status** | Evaluate |

**Integration:** Gateway generates OpenDroneID messages, sends to broadcast module hardware (Bluetooth/WiFi transmitter).

---

## §2 Detect-and-Avoid (DAA)

### What It Is

DAA replaces the "see and avoid" responsibility pilots have in manned aircraft. Since BVLOS operators can't see traffic, the system must detect other aircraft and either alert the operator or autonomously maneuver to stay clear.

### Why It Matters

This is the hardest BVLOS requirement. A drone can't see a Cessna coming - it needs sensors. Without DAA, you're flying blind into airspace shared with manned aircraft. Regulators won't approve BVLOS without some DAA solution.

### How It Works

Two approaches:
1. **Cooperative** - other aircraft broadcast their position (ADS-B, FLARM, transponder). You receive and track.
2. **Non-cooperative** - detect aircraft that aren't transmitting (radar, cameras, acoustic). Much harder.

Most practical: **ADS-B In receiver**. All aircraft in controlled airspace must have ADS-B Out (since 2020 in USA). You receive their broadcasts and know where they are.

**Limitation:** Small GA aircraft in uncontrolled airspace may not have ADS-B. Ground-based radar or camera systems fill this gap for some operations.

### Library: skypies/adsb

| | |
|-|-|
| **Repo** | https://github.com/skypies/adsb |
| **Purpose** | ADS-B message type definitions and parsing |
| **Status** | Evaluate |

### Library: skypies/pi

| | |
|-|-|
| **Repo** | https://github.com/skypies/pi |
| **Purpose** | Raspberry Pi ADS-B receiver with dump1090 integration |
| **Hardware** | RTL-SDR dongle (~$30) + Raspberry Pi |
| **Status** | Evaluate |

**Integration:** Pi + RTL-SDR receives ADS-B. Parse with skypies/adsb. Feed traffic positions to Jetson for fusion with other sensors. Alert operator or trigger avoidance maneuver.

---

## §3 Geofencing

### What It Is

Invisible boundaries that prevent the drone from entering restricted areas. Implemented in the flight controller - if the drone approaches a boundary, it refuses to cross (returns, lands, or holds position).

### Why It Matters

Prevents flyaways into airports, prisons, military bases, TFRs, or beyond your operational area. Regulators require it because it's the last line of defense when comms fail or operator makes a mistake. Hardcoded into PX4/ArduPilot.

### How It Works

Two types:
1. **Static geofence** - operational boundary defined before flight (your approved area)
2. **Dynamic geofence** - temporary restrictions that change (TFRs, NOTAMs, stadium events)

Data sources:
- FAA UDDS for permanent airspace
- NOTAM feeds for temporary restrictions
- Mission-specific boundaries from flight planning

### Data Sources by Region

#### USA 🇺🇸

| Data Type | Source | Access | Notes |
|-----------|--------|--------|-------|
| **Permanent airspace** | [FAA UDDS](https://udds-faa.opendata.arcgis.com) | **Free, open** | Class B/C/D/E boundaries, UAS Facility Maps |
| **TFRs / NOTAMs** | [NASA DIP](https://dip.amesaero.nasa.gov) | Partner registration | Parsed SWIM feed with geospatial data |
| **Combined** | [Airspace Link API](https://developers.airspacelink.com/) | Developer signup | Static + dynamic, LAANC-ready |

#### Europe 🇪🇺

| Data Type | Source | Access | Notes |
|-----------|--------|--------|-------|
| **EU-wide common rules** | [EASA Easy Access Rules](https://www.easa.europa.eu/en/document-library/easy-access-rules) | **Free** | Regulations only, no geodata |
| **Germany** | [DFS AIS Portal](https://aip.dfs.de) | Registration | AIP, NOTAM, airspace charts |
| **France** | [SIA France](https://www.sia.aviation-civile.gouv.fr) | Registration | AIP, SUP AIP, NOTAMs |
| **UK** | [NATS AIS](https://nats-uk.ead-it.com) | Registration | AIP, NOTAMs, drone zones |
| **Spain** | [ENAIRE AIS](https://ais.enaire.es) | Registration | AIP, NOTAMs |
| **Italy** | [ENAV AIS](https://www.enav.it/servizi/ais) | Registration | AIP, NOTAMs |

**EU Reality:** No single source. Each CAA publishes their own AIP (Aeronautical Information Publication). Eurocontrol's Network Manager doesn't provide public geodata API.

#### Asia-Pacific 🌏

| Country | Source | Access | Notes |
|---------|--------|--------|-------|
| **Australia** | [CASA Airspace](https://www.casa.gov.au/drones/rules/drone-safety-map) | **Free** | OpenSky map for recreational |
| **Australia (commercial)** | [Airservices AIS](https://www.airservicesaustralia.com/aip/aip.asp) | Registration | Full AIP, NOTAMs |
| **Japan** | [MLIT AIS](https://www.mlit.go.jp/koku/koku_fr10_000003.html) | Registration | Restricted areas require permission |
| **Singapore** | [CAAS](https://www.caas.gov.sg/operations-safety/unmanned-aircraft) | Permit required | Most of Singapore is controlled |

#### Global Aggregators

| Source | Access | Coverage | Notes |
|--------|--------|----------|-------|
| **[OpenAIP](https://www.openaip.net)** | **Free** (CC BY-NC-SA) | 170+ countries | Community-maintained, not official |
| **[OpenFlightMaps](https://www.openflightmaps.org)** | **Free** | Europe, partial global | VFR charts, not UAS-specific |
| **[SkyVector](https://skyvector.com)** | **Free** | Worldwide | Charts, not machine-readable |

**OpenAIP is the closest thing to "global airspace data"** - but it's community-maintained, not authoritative. For commercial BVLOS, you need official CAA data for your jurisdiction.

#### Strategy

For multi-region operations:

```
┌─────────────────────────────────────────────────────────────────┐
│  Priority 1: Official CAA data for your operating region        │
│  - USA: FAA UDDS (free) + NASA DIP (NOTAMs)                    │
│  - EU: Per-country AIS portal                                   │
│  - AU: Airservices AIS                                         │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Priority 2: Commercial aggregator (if multi-country)          │
│  - Airspace Link, Aloft - combine static + dynamic             │
│  - Handle regional complexity for you                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Priority 3: OpenAIP (supplementary reference)                  │
│  - Good for development/testing                                 │
│  - Not sufficient alone for production                         │
└─────────────────────────────────────────────────────────────────┘
```

### Library: skypies/geo

| | |
|-|-|
| **Repo** | https://github.com/skypies/geo |
| **Purpose** | Lat/lon calculations, Class B airspace boundary data |
| **Status** | Evaluate |

**Integration:** Load airspace polygons from FAA UDDS. Use skypies/geo for point-in-polygon checks. Push boundaries to PX4 via MAVLink. Gateway monitors position and can add dynamic restrictions.

---

## External: Weather (Pre-flight)

### What It Is

Aviation weather reports from airports. METAR = current conditions (wind, visibility, clouds, temp). TAF = forecast for next 24-30 hours.

### Why It Matters

Drones can't fly in bad weather - wind limits, visibility minimums, icing, thunderstorms. Pre-flight weather check is part of every approval. Some operations require real-time weather monitoring during flight.

### How It Works

Weather stations at airports report hourly (or more). Data flows to Aviation Weather Center. You query by airport code (ICAO 4-letter, e.g., KJFK).

### API: Aviation Weather Center (METAR/TAF)

| | |
|-|-|
| **Endpoint** | `https://aviationweather.gov/api/data/metar` |
| **Example** | `GET /api/data/metar?ids=KMCI&format=json` |
| **Access** | **Open - no API key required** |
| **Rate Limit** | 100 requests/minute |
| **Max Results** | 400 per request |
| **Formats** | JSON, XML, CSV, GeoJSON, raw |
| **Data Refresh** | Every minute |

### Library: lus/awc.go

| | |
|-|-|
| **Repo** | https://github.com/lus/awc.go |
| **Purpose** | Go client for Aviation Weather Center |
| **Features** | METAR queries (TAF not yet implemented) |
| **Status** | Evaluate - ready to use |

**Integration:** Query nearest airport METAR before flight authorization. Check wind < limits, visibility > minimums, no thunderstorms. Display in dashboard. Optionally poll during flight.

---

## External: NOTAM/TFR

### What It Is

NOTAMs (Notices to Air Missions) are temporary changes to airspace - closed runways, military exercises, VIP movement, stadium TFRs, rocket launches, wildfires. TFRs are the no-fly zones within NOTAMs.

### Why It Matters

The Super Bowl creates a 30nm no-fly zone. The President's travel closes airspace. A wildfire spawns a TFR. Your static geofence doesn't know about these - you need real-time NOTAM data or you'll fly into restricted space.

### How It Works

FAA publishes NOTAMs through SWIM (System Wide Information Management). NASA's DIP repackages SWIM data with better structure. NOTAMs are text blobs with embedded coordinates and times - parsing them is non-trivial.

### API: NASA Digital Information Platform (DIP)

| | |
|-|-|
| **Endpoint** | `https://dip.amesaero.nasa.gov` |
| **Protocol** | REST over HTTPS |
| **Format** | JSON |
| **Access** | **Partner registration required** |
| **Data Source** | Redistributes FAA SWIM feed with value-added processing |
| **Features** | Extracts geospatial, temporal, regulatory info from NOTAMs |

**Registration Status:** ✅ NASA Launchpad account created, awaiting DIP API access approval

**Registration Steps:**
1. [x] Go to https://dip.amesaero.nasa.gov
2. [x] Redirected to NASA Launchpad (auth.launchpad.nasa.gov)
3. [x] Created NASA Guest account
4. [ ] Complete partner onboarding form (if required)
5. [ ] Wait for approval (NASA reviews access requests)
6. [ ] Get API credentials from DIP portal

**Note:** NASA Launchpad uses multi-level authentication. Guest accounts work for DIP access.

**Integration:** Poll NOTAMs for your operational area. Parse TFR polygons. Push to Gateway as dynamic geofence updates. Alert operator if new TFR affects planned route.

---

## External: Airspace Data

### What It Is

Static airspace structure - where controlled airspace exists, what altitudes are allowed, special use airspace (military), and UAS-specific data like facility maps showing max LAANC altitudes.

### Why It Matters

You need to know the rules before you ask permission. UAS Facility Maps tell you "LAANC can approve up to 200ft here" vs "LAANC can't help, you need a waiver." This data drives your geofencing and flight planning.

### How It Works

FAA publishes authoritative geospatial data. ArcGIS-based portal, downloadable in standard formats. Updates regularly as airspace changes.

### API: FAA UAS Data Delivery System (UDDS) ⭐ Direct FAA Access

| | |
|-|-|
| **Portal** | `https://udds-faa.opendata.arcgis.com` |
| **Access** | **Open - free, no registration** |
| **Formats** | CSV, JSON, GeoJSON, KML, Shapefile |

**Available Datasets & API Endpoints:**

| Dataset | FeatureServer URL | Description |
|---------|-------------------|-------------|
| **UAS Facility Maps** | `https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0` | LAANC ceiling altitudes in controlled airspace |
| **Special Use Airspace** | `https://adds-faa.opendata.arcgis.com/datasets/faa::special-use-airspace` | MOAs, Restricted Areas, Prohibited Areas |
| **Class Airspace** | `https://adds-faa.opendata.arcgis.com/datasets/class-airspace` | Class B/C/D/E boundaries |
| **FRIA** | Via UDDS portal | FAA-Recognized Identification Areas |
| **National Security** | Via UDDS portal | UAS Flight Restrictions |

**UAS Facility Map Fields:**
- `CEILING` - Max altitude (ft AGL) for LAANC authorization
- `APT1_FAAID`, `APT1_ICAO`, `APT1_NAME` - Associated airport(s)
- `APT1_LAANC` - LAANC availability status
- `REGION` - FAA region

**GeoJSON Query Format:**
```
{FeatureServer URL}/query?outFields=*&where=1=1&f=geojson
```

**Example - Download all UAS Facility Map data as GeoJSON:**
```bash
curl "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0/query?outFields=*&where=1=1&f=geojson" > uas_facility_map.geojson
```

**Integration:** Download airspace layers. Load into PostGIS or similar. Query vehicle position against polygons. Drive geofencing and flight planning UI.

---

## External: LAANC Authorization

### What It Is

LAANC (Low Altitude Authorization and Notification Capability) is near-real-time approval to fly in controlled airspace. Instead of waiting weeks for a waiver, you get automated approval in seconds - if your request fits within pre-approved parameters.

### Why It Matters

Most interesting places to fly (cities, near airports) are controlled airspace. Without LAANC, you'd need FAA waivers taking 90+ days. LAANC makes commercial operations practical. Required for any controlled airspace BVLOS.

### How It Works

You submit: where, when, how high. FAA checks against UAS Facility Maps (pre-approved altitudes) and current airspace status. If it fits, instant approval. If not, you need a waiver.

**Catch:** FAA doesn't give developers direct API access. You must go through approved USS (UAS Service Supplier) partners.

### Architecture Reality

```
┌─────────────────────────────────────────────────────────┐
│  FAA LAANC Automation Platform (api.faa.gov/laanc)      │
│  ⚠️ USS-ONLY - requires FAA approval, annual testing    │
└─────────────────────────────────────────────────────────┘
                          ▲
                          │ USS-FAA API
                          │
┌─────────────────────────────────────────────────────────┐
│  Approved USS Partners (middlemen)                       │
│  Airspace Link, Aloft, DroneUp, FlyFreely               │
└─────────────────────────────────────────────────────────┘
                          ▲
                          │ Their APIs
                          │
┌─────────────────────────────────────────────────────────┐
│  Your Application                                        │
└─────────────────────────────────────────────────────────┘
```

**Why USS middlemen exist:** FAA wants tested, insured intermediaries. Becoming a USS requires: annual application (May), technical interview, formal testing, MOA signing, insurance. Not practical for most developers.

### Option 1: Use USS API (practical)

**Airspace Link (AirHub):**

| | |
|-|-|
| **Docs** | https://developers.airspacelink.com/ |
| **Sandbox** | `https://airhub-api-sandbox.airspacelink.com` |
| **Production** | `https://airhub-api.airspacelink.com` |
| **Auth** | OAuth2 bearer token |
| **LAANC endpoint** | `POST /v1/operations/{id}/laanc` |

**How to Register:**
1. Go to https://airspacelink.com/developers
2. Request developer access
3. Get sandbox credentials

**Aloft:** Contact sales (no public API docs).

### Option 2: Become a USS (overkill for most)

| | |
|-|-|
| **Apply** | May each year |
| **Process** | [faa.gov/uas/programs_partnerships/data_exchange/how-to-apply-2024](https://www.faa.gov/uas/programs_partnerships/data_exchange/how-to-apply-2024) |
| **Requirements** | Technical demo, formal testing, MOA, insurance |

**Note:** AirMap shut down LAANC service (June 2023).

---

## Next Steps

### Phase 1: Data Foundation (no middleman) ⬅️ CURRENT

| # | Task | Status | Action |
|---|------|--------|--------|
| 1 | **FAA UDDS** | [x] | Download airspace data from https://udds-faa.opendata.arcgis.com |
| | | [x] | - Identify available datasets (see [§ Airspace Data](#external-airspace-data)) |
| | | [x] | - Download UAS Facility Map GeoJSON (2.2MB, 2000 grid cells) |
| | | [x] | - Download Class Airspace GeoJSON (14MB - `faa_airspace_boundary.geojson`) |
| | | [x] | - Download Special Use Airspace GeoJSON (28MB - MOAs, Restricted, Prohibited) |
| | | [x] | - Stored locally in `data/airspace/` |
| 2 | **lus/awc.go** | [ ] | Build METAR client |
| | | [ ] | - Clone https://github.com/lus/awc.go |
| | | [ ] | - Write simple test: fetch METAR for KJFK |
| | | [ ] | - Verify JSON parsing works |
| 3 | **OpenAIP** | [ ] | Download global airspace for dev/testing |
| | | [ ] | - Register at https://www.openaip.net |
| | | [ ] | - Download airspace data (GeoJSON) |
| | | [ ] | - Use for non-US testing |

### Phase 2: Libraries to Evaluate

| # | Task | Status | Action |
|---|------|--------|--------|
| 4 | **gomavlib** | [ ] | Wire OpenDroneID messages in Gateway (§1 Remote ID) |
| | | [ ] | - Clone https://github.com/bluenviron/gomavlib |
| | | [ ] | - Find OpenDroneID message definitions |
| | | [ ] | - Write example: generate BASIC_ID message |
| 5 | **skypies/geo** | [ ] | Integrate airspace boundaries for geofencing (§3) |
| | | [ ] | - Clone https://github.com/skypies/geo |
| | | [ ] | - Test point-in-polygon with Class B boundary |
| 6 | **skypies/adsb + pi** | [ ] | Evaluate for DAA traffic awareness (§2) |
| | | [ ] | - Review https://github.com/skypies/adsb |
| | | [ ] | - Review https://github.com/skypies/pi |
| | | [ ] | - Document hardware requirements (RTL-SDR) |

### Phase 3: Requires Registration

| # | Task | Status | Action |
|---|------|--------|--------|
| 7 | **NASA DIP** | [~] | NOTAM/TFR feed (see [§ NOTAM/TFR](#external-notamtfr)) |
| | | [x] | - Go to https://dip.amesaero.nasa.gov |
| | | [x] | - Redirected to NASA Launchpad |
| | | [x] | - Created NASA Guest account |
| | | [ ] | - Complete partner onboarding form (if required) |
| | | [ ] | - Wait for approval |
| | | [ ] | - Get API credentials from DIP portal |
| | | [ ] | - Test API: fetch NOTAMs for a location |
| 8 | **Airspace Link** | [ ] | Request developer access (for LAANC authorization) |
| | | [ ] | - Go to https://airspacelink.com/developers |
| | | [ ] | - Request developer access |
| | | [ ] | - Get sandbox credentials |
| | | [ ] | - Test API: query airspace for a location |

### Phase 4: Web GUI (Visualization & Testing)

Build a map-based dashboard to visualize all data layers. Serves two purposes:
1. **Development testing** - Verify data parsing and integration is correct
2. **Demo** - Show a real working system to stakeholders/regulators

#### MVP Features

| Layer | Data Source | Visualization |
|-------|-------------|---------------|
| **Airspace boundaries** | FAA UDDS / OpenAIP | Colored polygons (Class B=blue, C=magenta, etc.) |
| **UAS Facility Maps** | FAA UDDS | Altitude ceiling overlays |
| **TFRs** | NASA DIP (or mock) | Red restricted zones |
| **Weather** | aviationweather.gov | METAR pins at airports |
| **Live traffic** | ADS-B receiver | Moving aircraft icons |
| **Fleet vehicles** | JetStream telemetry | Drone/vehicle positions |
| **Geofence** | Mission config | Editable polygon overlay |

#### Tech Stack Options

| Approach | Pros | Cons |
|----------|------|------|
| **Leaflet.js + Go backend** | Simple, lightweight, open-source | Manual layer management |
| **Mapbox GL JS** | Beautiful, performant, UAS-friendly | Costs at scale |
| **deck.gl** | Great for real-time data viz | Learning curve |
| **CesiumJS** | 3D airspace visualization | Heavy, complex |

**Recommended:** Start with Leaflet.js for MVP. Swap to Mapbox if you need more polish.

#### Implementation Tasks

| # | Task | Status | Action |
|---|------|--------|--------|
| 9 | **Map backend** | [ ] | Go service serving GeoJSON from FAA UDDS |
| 10 | **Airspace layer** | [ ] | Parse and render Class B/C/D/E polygons |
| 11 | **Weather overlay** | [ ] | METAR pins at airports with wind/visibility |
| 12 | **TFR layer** | [ ] | Red zones from NOTAM feed |
| 13 | **Vehicle layer** | [ ] | Real-time positions from JetStream |
| 14 | **Geofence editor** | [ ] | Draw operational boundary on map |

---

## Data Access Requirements (Valley of Death Analysis)

For an **Australian company** trying to access these data providers:

| Provider | Difficulty | Demo Required? | What You Need | Contact |
|----------|-----------|----------------|---------------|---------|
| **FAA UDDS** | 🟢 Easy | No | Nothing - just download | Open API |
| **aviationweather.gov** | 🟢 Easy | No | Nothing - open API | Open API |
| **OpenAIP** | 🟢 Easy | No | Free account, API key from profile | https://openaip.net |
| **NASA DIP** | 🟡 Medium | Likely | Partner agreement, show use case | NASA Launchpad |
| **Airspace Link** | 🟡 Medium | Maybe | Contact sales, DPP for LAANC | hello@airspacelink.com |
| **Airservices AU** | 🟡 Medium | Unknown | Paid data subscription | data@airservicesaustralia.com |
| **CASA Platform** | 🔴 Hard | **Yes** | Working app showing accurate airspace | Wing-built, verification required |
| **AU FIMS (May 2026)** | 🔴 Hard | **Yes** | Become approved USS | USSonboarding@airservicesaustralia.com |

### Strategy: Build Demo First

The hard providers (CASA, FIMS) require a **working system** before they'll onboard you. Solution:

```
┌─────────────────────────────────────────────────────────────────┐
│  Phase 1-2: Build with FREE data                                │
│  - FAA UDDS (US airspace) - validates your map rendering        │
│  - OpenAIP (global) - shows AU airspace for dev/demo           │
│  - aviationweather.gov - proves weather integration            │
│  This becomes your DEMO for harder access                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 3: Apply for gated access with working demo              │
│  - NASA DIP: "Here's our platform, we need NOTAM data"         │
│  - Airservices: "Here's our USS platform, need AU data"        │
│  - CASA: "Here's our app, ready for verification"              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  Phase 4: USS Onboarding (if targeting AU production)           │
│  - FIMS launches May 2026                                       │
│  - First USS cohort: AvSoft, OneSky, Yarra Drones              │
│  - Contact: USSonboarding@airservicesaustralia.com             │
└─────────────────────────────────────────────────────────────────┘
```

### What Your Company Has (Ubuntu Software - AU)

| Asset | Status | Useful For |
|-------|--------|------------|
| ABN | ✅ 95 595 575 880 | Company verification |
| D-U-N-S | ✅ 891770992 | Enterprise credibility |
| Domain | ✅ ubuntusoftware.net | Professional presence |
| Website | ✅ Fleet architecture docs | Shows technical capability |
| Working demo | ✅ MVP built | `task dev` → http://localhost:1313/airspace-demo/ |

---

## Progress Summary

| Phase | Status | Progress |
|-------|--------|----------|
| Phase 1: Data Foundation | ✅ Complete | FAA UDDS downloaded: 44MB in `data/airspace/` |
| Phase 2: Libraries | ⚪ Not started | 0/3 complete |
| Phase 3: Registration | 🟡 In progress | NASA DIP account created, awaiting approval |
| Phase 4: Web GUI | ✅ MVP Complete | `static/airspace-demo/` + `cmd/airspace-server/` |

### Downloaded Data Files

| File | Size | Contents |
|------|------|----------|
| `static/airspace/faa_uas_facility_map.geojson` | 2.2MB | 2000 LAANC grid cells with ceiling altitudes |
| `static/airspace/faa_airspace_boundary.geojson` | 14MB | Class B/C/D/E airspace polygons |
| `static/airspace/faa_special_use_airspace.geojson` | 28MB | MOAs, Restricted Areas, Prohibited Areas |

**Next action:** Register for OpenAIP and NASA DIP API access using working demo as credibility

---

## Reference

| Regulation | Platform Requirement | API/Library | Access |
|------------|---------------------|-------------|--------|
| §1 Remote ID | Broadcast module | gomavlib | Open source |
| §2 DAA | Traffic awareness | skypies/adsb, skypies/pi | Open source |
| §3 Geofencing | Airspace boundaries | FAA UDDS + skypies/geo | **Free, open** |
| §7 Flight Logging | Telemetry storage | JetStream | Native |
| External | Weather | aviationweather.gov | **Free, open** |
| External | NOTAMs | NASA DIP | Registration |
| External | LAANC Auth | Airspace Link | USS API (middleman) |
