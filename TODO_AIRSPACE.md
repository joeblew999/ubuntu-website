# Airspace Data Architecture

This document outlines the data pipeline architecture for FAA airspace data used in the BVLOS demo at `/fleet/airspace-demo/`.

---

## Current State (Problems)

| Issue | Impact |
|-------|--------|
| **Laptop as middleman** | `task r2:airspace:upload` downloads 600MB+ to local, then uploads to R2. Slow, wasteful. |
| **No incremental sync** | Re-uploads everything on each run. No diffing or checksums. |
| **Web GUI loads all data** | Browser fetches 50MB+ of GeoJSON on page load. Will break with more datasets. |
| **Manual refresh** | No automation. Data goes stale unless someone runs the task. |
| **Large files excluded** | Obstacles (575MB) skipped because R2 upload times out. |

---

## Target Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          DATA SOURCES (FAA ArcGIS)                          │
│  https://adds-faa.opendata.arcgis.com                                       │
│  https://udds-faa.opendata.arcgis.com                                       │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      PIPELINE (GitHub Actions or CF Worker)                  │
│                                                                              │
│  1. Fetch metadata (checksums, last-modified)                               │
│  2. Compare with R2 object metadata                                          │
│  3. Download only changed files                                              │
│  4. Process/transform (tile, simplify, compress)                            │
│  5. Upload to R2 with metadata tags                                          │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              R2 STORAGE                                      │
│                                                                              │
│  /airspace/                                                                  │
│    manifest.json              <- Index of all datasets with checksums        │
│    faa_airports.geojson       <- Full dataset                                │
│    faa_airports.pmtiles       <- Vector tiles for large datasets             │
│    faa_uas_facility_map/      <- Tiled directory for huge files              │
│      0/0/0.pbf                                                               │
│      ...                                                                     │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              WEB GUI                                         │
│                                                                              │
│  1. Fetch manifest.json (list of available layers)                          │
│  2. Load only visible viewport via vector tiles                              │
│  3. Progressive loading as user pans/zooms                                   │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Diff-Based Sync

**Goal:** Only download/upload changed data. No laptop middleman.

### 1.1 Manifest File

Create `manifest.json` in R2 that tracks all datasets:

```json
{
  "version": 1,
  "updated": "2024-12-24T04:30:00Z",
  "datasets": {
    "airports": {
      "filename": "faa_airports.geojson",
      "source_url": "https://adds-faa.opendata.arcgis.com/...",
      "size_bytes": 13347449,
      "etag": "abc123...",
      "last_modified": "2024-12-20T00:00:00Z",
      "feature_count": 19847
    }
  }
}
```

### 1.2 Sync Logic

```go
// cmd/airspace/sync.go

func runSync() {
    // 1. Fetch current manifest from R2
    manifest := fetchManifest(r2URL + "/airspace/manifest.json")

    // 2. For each dataset, check if source has changed
    for _, ds := range datasets {
        sourceEtag := headRequest(ds.BaseURL).Header.Get("ETag")

        if manifest.Datasets[ds.Key].Etag == sourceEtag {
            fmt.Printf("  [%s] Up to date\n", ds.Key)
            continue
        }

        // 3. Download changed dataset
        data := download(ds.BaseURL)

        // 4. Upload directly to R2 (no local storage)
        uploadToR2(data, ds.Filename)

        // 5. Update manifest
        manifest.Datasets[ds.Key] = DatasetMeta{...}
    }

    // 6. Upload updated manifest
    uploadToR2(manifest, "manifest.json")
}
```

### 1.3 GitHub Actions Cron

```yaml
# .github/workflows/monitor-airspace.yml
name: Airspace Data Sync

on:
  schedule:
    - cron: '0 6 * * 0'  # Weekly on Sunday 6am UTC
  workflow_dispatch:      # Manual trigger

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Sync airspace data
        env:
          CLOUDFLARE_API_TOKEN: ${{ secrets.CLOUDFLARE_API_TOKEN }}
          CLOUDFLARE_ACCOUNT_ID: ${{ secrets.CF_ACCOUNT_ID }}
          R2_BUCKET: ubuntu-website-assets
        run: |
          go run ./cmd/airspace sync
```

### 1.4 Taskfile Updates

```yaml
# taskfiles/Taskfile.airspace.yml

sync:
  desc: Sync changed datasets to R2 (diff-based, no local storage)
  cmds:
    - go run ./cmd/airspace sync

sync:check:
  desc: Check which datasets have changed (dry run)
  cmds:
    - go run ./cmd/airspace sync -dry-run
```

---

## Phase 2: Vector Tiles for Large Datasets

**Goal:** Stream data progressively instead of loading entire GeoJSON.

### 2.1 The Problem

| Dataset | Size | Features | Load Time |
|---------|------|----------|-----------|
| Airports | 13 MB | 19,847 | ~2s |
| Airspace Boundary | 14 MB | ~3,000 | ~2s |
| Special Use Airspace | 28 MB | ~5,000 | ~4s |
| UAS Facility Map | 2 MB | 2,000 | <1s |
| Obstacles | 575 MB | ~600,000 | ❌ Unusable |

Loading 50MB+ on page load is already slow. Adding more datasets will break the demo.

### 2.2 Solution: PMTiles

[PMTiles](https://github.com/protomaps/PMTiles) is a single-file format for vector tiles that works directly from static hosting (R2).

```
Original: faa_obstacles.geojson (575 MB)
    ↓ tippecanoe
Tiled:   faa_obstacles.pmtiles (50-100 MB)
    ↓ served from R2
Browser: Only loads visible tiles (~100KB per viewport)
```

### 2.3 Build Pipeline

```bash
# Install tippecanoe (one-time)
brew install tippecanoe

# Convert GeoJSON to PMTiles
tippecanoe -o faa_obstacles.pmtiles \
  --maximum-zoom=14 \
  --minimum-zoom=4 \
  --drop-densest-as-needed \
  --extend-zooms-if-still-dropping \
  faa_obstacles.geojson
```

### 2.4 Leaflet Integration

```javascript
// Use protomaps-leaflet for PMTiles
import { PMTiles, leafletLayer } from 'pmtiles';

const obstacles = new PMTiles(`${DATA_BASE}/faa_obstacles.pmtiles`);
const obstaclesLayer = leafletLayer(obstacles, {
    style: { color: '#ff6600', weight: 1 }
});
```

### 2.5 Thresholds

| Dataset Size | Strategy |
|--------------|----------|
| < 5 MB | GeoJSON (load all) |
| 5-50 MB | GeoJSON with lazy loading |
| > 50 MB | PMTiles (vector tiles) |

---

## Phase 3: Cloudflare Worker Pipeline

**Goal:** Run sync entirely on Cloudflare infrastructure (no GitHub Actions).

### 3.1 Architecture

```
┌──────────────────┐     ┌──────────────────┐     ┌──────────────────┐
│  Cron Trigger    │────▶│  CF Worker       │────▶│  R2 Bucket       │
│  (every Sunday)  │     │  (fetch & sync)  │     │  (storage)       │
└──────────────────┘     └──────────────────┘     └──────────────────┘
```

### 3.2 Worker Code (Future)

```typescript
// workers/airspace-sync/src/index.ts
export default {
  async scheduled(controller, env, ctx) {
    const datasets = [
      { key: 'airports', url: '...' },
      // ...
    ];

    for (const ds of datasets) {
      const current = await env.R2.head(`airspace/${ds.key}.geojson`);
      const source = await fetch(ds.url, { method: 'HEAD' });

      if (current?.etag !== source.headers.get('etag')) {
        const data = await fetch(ds.url);
        await env.R2.put(`airspace/${ds.key}.geojson`, data.body);
      }
    }
  }
};
```

---

## Phase 4: Smart Web GUI

**Goal:** Load data progressively based on viewport.

### 4.1 Lazy Loading Strategy

```javascript
// Only load layers when toggled ON
document.getElementById('toggle-airports').addEventListener('change', async (e) => {
    if (e.target.checked && !airportsLoaded) {
        showLoading('Loading airports...');
        const data = await fetch(`${DATA_BASE}/faa_airports.geojson`);
        // ... add to layer
        airportsLoaded = true;
        hideLoading();
    }
    e.target.checked ? map.addLayer(airportsLayer) : map.removeLayer(airportsLayer);
});
```

### 4.2 Viewport-Based Loading

```javascript
// For large datasets, only load visible area
map.on('moveend', async () => {
    const bounds = map.getBounds();
    const bbox = `${bounds.getWest()},${bounds.getSouth()},${bounds.getEast()},${bounds.getNorth()}`;

    // Query R2/Worker with bbox filter
    const data = await fetch(`${DATA_BASE}/faa_obstacles.geojson?bbox=${bbox}`);
    // This requires a Worker to filter, or pre-tiled data
});
```

### 4.3 Progressive Enhancement

| Zoom Level | Data Loaded |
|------------|-------------|
| 1-4 | Airspace classes only (outlines) |
| 5-8 | + Major airports, MOAs |
| 9-12 | + All airports, navaids |
| 13+ | + Obstacles, detailed boundaries |

---

## Implementation Roadmap

### Now (Phase 1A) - Immediate Fixes

- [ ] Add `sync` subcommand to `cmd/airspace` with dry-run
- [ ] Create `manifest.json` schema
- [ ] Update taskfile with `sync` and `sync:check` tasks
- [ ] Add ETag/Last-Modified tracking

### Soon (Phase 1B) - GitHub Actions

- [ ] Create `.github/workflows/monitor-airspace.yml`
- [ ] Direct R2 upload from CI (no laptop)
- [ ] Slack/email notification on sync

### Later (Phase 2) - Vector Tiles

- [ ] Install tippecanoe in CI
- [ ] Convert large datasets to PMTiles
- [ ] Update web GUI to use PMTiles for obstacles
- [ ] Add zoom-level filtering

### Future (Phase 3) - Full Automation

- [ ] Cloudflare Worker for sync
- [ ] Worker for bbox queries
- [ ] Real-time TFR overlay (dynamic data)

---

## Datasets Inventory

| Key | Name | Size | Strategy | Status |
|-----|------|------|----------|--------|
| `uas` | UAS Facility Map | 2 MB | GeoJSON | ✅ In R2 |
| `boundary` | Airspace Boundary | 14 MB | GeoJSON | ✅ In R2 |
| `sua` | Special Use Airspace | 28 MB | GeoJSON | ✅ In R2 |
| `airports` | Airports | 13 MB | GeoJSON | ✅ In R2 |
| `navaids` | Navigation Aids | 0.8 MB | GeoJSON | ✅ In R2 |
| `obstacles` | Obstacles | 575 MB | PMTiles (Phase 2) | ❌ Too large |

---

## Data Sources Reference

| Dataset | Source URL | Update Frequency |
|---------|------------|------------------|
| All FAA aeronautical | https://adds-faa.opendata.arcgis.com | 28-day AIRAC cycle |
| UAS-specific | https://udds-faa.opendata.arcgis.com | Varies |

**AIRAC Cycle:** Aeronautical data updates every 28 days on specific dates. See [AIRAC calendar](https://www.nm.eurocontrol.int/RAD/common/airac_dates.html).

---

---

## Protomaps Ecosystem (bdon's work)

Brandon Liu ([bdon](https://github.com/bdon)) has built the entire stack we need. **Don't reinvent - integrate.**

### Key Projects

| Project | Purpose | Our Use |
|---------|---------|---------|
| [protomaps/PMTiles](https://github.com/protomaps/PMTiles) | Single-file vector tile format | Store tiled airspace data in R2 |
| [protomaps/go-pmtiles](https://github.com/protomaps/go-pmtiles) | CLI for PMTiles operations | `pmtiles serve`, `pmtiles upload`, `pmtiles extract` |
| [protomaps/protomaps-leaflet](https://github.com/protomaps/protomaps-leaflet) | Vector rendering for Leaflet | Replace GeoJSON loading with progressive tiles |
| [felt/tippecanoe](https://github.com/felt/tippecanoe) | GeoJSON → vector tiles | Convert FAA data to PMTiles |
| [bdon/cng-storage-guide](https://github.com/bdon/cng-storage-guide) | Cloud storage comparison | Confirms R2 is ideal (zero egress) |

### Why R2 is Perfect (from bdon's research)

```yaml
# From bdon/cng-storage-guide
- name: Cloudflare R2
  cost_per_gb_stored: 0.015
  cost_per_gb_egress: 0        # <-- THIS IS KEY
  cost_per_1k_gets: 0.00036
```

PMTiles on R2 = **zero bandwidth costs**. A 575MB obstacles file served to thousands of users costs only storage + requests, not bandwidth.

### PMTiles CLI Commands We Need

```bash
# Install
go install github.com/protomaps/go-pmtiles@latest

# Convert GeoJSON to PMTiles (via tippecanoe)
tippecanoe -zg -o faa_obstacles.pmtiles faa_obstacles.geojson

# Upload directly to R2 (S3-compatible)
pmtiles upload faa_obstacles.pmtiles s3://ubuntu-website-assets/airspace/

# Serve locally for testing
pmtiles serve ./static/airspace/ --port 8081

# Extract subset by bounding box
pmtiles extract source.pmtiles subset.pmtiles --bbox=-122.5,37.5,-122.0,38.0
```

### Updated Architecture with Protomaps

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          FAA DATA SOURCES                                    │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      GITHUB ACTIONS PIPELINE                                 │
│                                                                              │
│  1. Download GeoJSON from FAA ArcGIS                                        │
│  2. Convert to PMTiles via tippecanoe                                        │
│  3. Upload to R2 via `pmtiles upload`                                        │
│  4. Update manifest.json                                                     │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                              R2 STORAGE                                      │
│                                                                              │
│  /airspace/                                                                  │
│    manifest.json                                                             │
│    faa_airspace.pmtiles        <- All airspace in one file!                 │
│    faa_airports.pmtiles                                                      │
│    faa_obstacles.pmtiles       <- 575MB → ~50MB tiled                       │
└───────────────────────────────────┬─────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      WEB GUI (protomaps-leaflet)                             │
│                                                                              │
│  import { PMTiles, leafletLayer } from 'pmtiles';                           │
│  const tiles = new PMTiles('https://r2.../airspace/faa_airspace.pmtiles');  │
│  map.addLayer(leafletLayer(tiles));                                          │
│                                                                              │
│  → Only loads visible tiles via HTTP Range requests                          │
│  → Progressive loading as user pans/zooms                                    │
│  → No 50MB+ initial download!                                               │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Revised Implementation Roadmap

| Phase | Task | Tools |
|-------|------|-------|
| **1A** | Add `sync` command with ETag tracking | Go, R2 API |
| **1B** | GitHub Actions for FAA→R2 sync | GH Actions |
| **2A** | Install tippecanoe in CI | `brew install tippecanoe` |
| **2B** | Convert large datasets to PMTiles | tippecanoe, go-pmtiles |
| **2C** | Update web GUI to protomaps-leaflet | protomaps-leaflet |
| **3** | Combine all layers into single PMTiles | tippecanoe --layer flag |

### Resources

- Protomaps Docs: https://docs.protomaps.com
- PMTiles Spec: https://github.com/protomaps/PMTiles/blob/main/spec/v3/spec.md
- Cloud Storage Guide: https://bdon.github.io/cng-storage-guide/
- Tippecanoe: https://github.com/felt/tippecanoe

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2024-12-24 | Skip obstacles for now | 575MB too large for GeoJSON; needs vector tiles |
| 2024-12-24 | Diff sync before more sources | Adding data without smart sync = exponential pain |
| 2024-12-24 | GitHub Actions first | Simpler than CF Worker; already have CI |
| 2024-12-24 | Use Protomaps stack | bdon already solved this; don't reinvent |
| 2024-12-24 | R2 for storage | Zero egress costs, HTTP Range support for PMTiles |
