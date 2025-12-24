# Airspace Data Pipeline

FAA airspace data for the BVLOS demo at `/fleet/airspace-demo/`.

## Current Status

| Component | Status | Notes |
|-----------|--------|-------|
| FAA Data Sync | ✅ Done | Chunked (2000 features/request) |
| PMTiles Generation | ⚠️ Blocked | tippecanoe (C++) can't run in Cloudflare |
| LAANC Layer | ✅ Done | 378K features, 129MB PMTiles |
| Regional Manifest | ✅ Done | `manifest_usa.json` drives Hugo template |
| Local Dev | ✅ Done | Port 1313 detection for LAN access |
| R2 Hosting | ✅ Done | Production tiles served from R2 |

---

## Priority: Chunked Cloudflare Pipeline

**Constraint:** All processing must run in Cloudflare Workers (128MB memory limit).

### Current State

Download is already chunked:
```go
// cmd/airspace/main.go - downloads 2000 features at a time
PageSize: 2000
```

Tile generation is NOT chunked - tippecanoe loads all 378K features into memory.

### Solution: Chunked Tile Generation

#### Architecture: Spatial Chunking

```
┌─────────────────────────────────────────────────────────────┐
│  Cloudflare Queue (triggers on FAA data change)             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Coordinator Worker                                          │
│  - Divides USA into grid cells (e.g., 10° × 10° boxes)      │
│  - Queues one job per cell                                  │
│  - Tracks completion in KV                                  │
└─────────────────────────────────────────────────────────────┘
                              │
           ┌──────────────────┼──────────────────┐
           ▼                  ▼                  ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  Tile Worker 1  │ │  Tile Worker 2  │ │  Tile Worker N  │
│  Cell: [-130,-120,40,50] │ │ Cell: [-120,-110,40,50] │ │  ...  │
│  - Reads features in bbox  │ │                    │ │       │
│  - Generates tiles z0-z14  │ │                    │ │       │
│  - Writes to R2/tiles/     │ │                    │ │       │
└─────────────────┘ └─────────────────┘ └─────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│  Merge Worker (triggered when all cells complete)           │
│  - Reads all cell PMTiles from R2                           │
│  - Merges into single PMTiles (go-pmtiles)                  │
│  - Writes final faa_*.pmtiles to R2                         │
└─────────────────────────────────────────────────────────────┘
```

#### Step 1: Spatial Index in R2

Store GeoJSON pre-indexed by spatial grid:

```
R2:
/airspace/raw/
├── faa_uas_facility_map.geojson          (full file)
└── indexed/
    ├── cell_-130_-120_40_50.geojson      (features in this bbox)
    ├── cell_-120_-110_40_50.geojson
    └── ...
```

Or use streaming: read full GeoJSON, filter by bbox on the fly.

#### Step 2: Per-Cell Tile Worker

```go
// Runs in Cloudflare Worker (Go → WASM via syumai/workers)
func processTileCell(bbox [4]float64) {
    // 1. Read features from R2 (filtered by bbox)
    features := readFeaturesInBbox(r2, "faa_uas.geojson", bbox)

    // 2. Generate tiles for this cell only
    tiles := generateTiles(features, minZoom, maxZoom)  // orb/encoding/mvt

    // 3. Write cell PMTiles to R2
    writePMTiles(r2, fmt.Sprintf("tiles/cell_%v.pmtiles", bbox), tiles)

    // 4. Mark cell complete in KV
    kv.Put(fmt.Sprintf("cell:%v", bbox), "done")
}
```

#### Step 3: Merge Worker

```go
// Triggered when all cells complete
func mergeCells() {
    // go-pmtiles can merge multiple PMTiles files
    cells := listR2("tiles/cell_*.pmtiles")
    merged := pmtiles.Merge(cells...)
    writeR2("tiles/faa_uas_facility_map.pmtiles", merged)

    // Cleanup temp cells
    deleteR2("tiles/cell_*.pmtiles")
}
```

### Memory Budget

| Stage | Memory Use | Within 128MB? |
|-------|-----------|---------------|
| Read cell GeoJSON (~30K features) | ~30MB | ✅ |
| Generate MVT tiles | ~50MB | ✅ |
| Write PMTiles | ~20MB | ✅ |
| **Total per cell** | ~100MB | ✅ |

### Grid Size

USA bbox: `[-125, 24, -66, 50]` = 59° × 26° = 1,534 sq degrees

| Grid Size | Cells | Features/Cell | Memory/Cell |
|-----------|-------|---------------|-------------|
| 10° × 10° | ~15 | ~25K | ~80MB ✅ |
| 5° × 5° | ~60 | ~6K | ~20MB ✅ |
| 2° × 2° | ~375 | ~1K | ~5MB ✅ |

Start with 10° × 10° (~15 cells), adjust if needed.

---

## Implementation Steps

1. **Proof of concept**: Pure Go tile generation locally (no WASM yet)
   - Use `paulmach/orb/encoding/mvt` to generate tiles
   - Use `protomaps/go-pmtiles` to package
   - Verify output matches tippecanoe

2. **Add spatial chunking**: Split by bbox, merge results

3. **Port to Cloudflare**: Compile to WASM with `syumai/workers`

4. **Add Queues + KV**: Coordinate parallel cell processing

---

## Go Libraries

| Library | Purpose | WASM? |
|---------|---------|-------|
| [paulmach/orb](https://github.com/paulmach/orb) | Geometry + GeoJSON | ✅ Pure Go |
| [paulmach/orb/encoding/mvt](https://github.com/paulmach/orb/tree/master/encoding/mvt) | MVT encoder | ✅ Pure Go |
| [protomaps/go-pmtiles](https://github.com/protomaps/go-pmtiles) | PMTiles read/write/merge | ✅ Pure Go |
| [syumai/workers](https://github.com/syumai/workers) | Go → Cloudflare Workers | ✅ |

---

## Tile Serving (Already Solved)

PMTiles Cloudflare Worker for serving:
- https://github.com/protomaps/PMTiles/tree/main/serverless/cloudflare
- https://docs.protomaps.com/deploy/cloudflare

---

## Research: go-pmtiles Issues

Key findings from GitHub issue research (2024-12-24):

### Issue #105: `pmtiles merge` - Now Available! (v1.29.1)

**Good news**: `pmtiles merge` was added in v1.29.1 for **Case 1: disjoint tilesets**.

```bash
pmtiles merge INPUT_1.pmtiles INPUT_2.pmtiles OUTPUT.pmtiles
```

**Constraints:**
- Tilesets must be **disjoint** (no overlapping tiles)
- All inputs must have same TileType and compression
- All inputs must be clustered

**Our use case**: Perfect! Spatial grid cells produce disjoint tiles by design.

**Not supported** (won't work for our case anyway):
- Case 2: Raster tile overlap (needs pixel merging - requires image libs)
- Case 3: Vector overlap with disjoint layers (PBF concatenation - planned)
- Case 4: Vector overlap with same layer names (complex - not planned)

### Issue #240: Multi-part PMTiles (Cloudflare Caching)

User wanted to split large files for Cloudflare's 512MB cache limit.

**bdon's response**: Use Cloudflare Workers integration (what we're doing). Multi-file format breaks ETag invalidation - not planned. Would need content-addressable storage (like casync) to do correctly.

**For us**: Not relevant - our final PMTiles (~129MB) is under Cloudflare limits.

### Alternative: pmtiles_mosaic

Mentioned in #240: https://github.com/ramSeraph/pmtiles_mosaic - third-party tool for serving multiple PMTiles as one logical tileset. Client-side solution, not what we need.

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2024-12-24 | Must run in Cloudflare | No external CI, edge-native |
| 2024-12-24 | Replace tippecanoe with Go | C++ can't run in Workers |
| 2024-12-24 | Spatial chunking | 128MB Worker limit, 378K features |
| 2024-12-24 | 10° grid cells | ~15 cells, ~25K features each |
| 2024-12-24 | Queue + KV coordination | Parallel cell processing |
| 2024-12-24 | Use `pmtiles merge` | v1.29.1 supports disjoint merge |
