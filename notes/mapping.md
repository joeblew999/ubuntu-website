# MAPPING

I was asked by a Drone shop about mapping.

We really need a good mapping solution for Drone.

## Use Case 

GEDTM30 + gedtm30api gives you a global 30 m terrain service you can hit over HTTP. Concrete ways it helps:

Terrain awareness/avoidance: fetch elevation along planned routes to enforce AGL ceilings, no-fly buffers, and geofences.
Line-of-sight and RF planning: compute sight lines and radio links for ground stations/relays using mean elevation and mask.
Landing/mission site selection: screen for flat/low-slope areas and avoid voids via the mask; use std-dev as a quality/confidence layer.
Terrain-relative nav (GPS-degraded): support terrain-matching or baro fusion by querying expected elevation profiles along the path.
Offline/edge tiling: prefetch tiles for AOIs so drones/field kits can operate without full 310 GB local data.
Analytics/visualization: overlay on your ops dashboard for quick terrain context for operators and partners.



https://github.com/bdon has soem good stuff for maps

https://github.com/akhenakh/gedtm30api lookks amazing for us — ready-made HTTP API to query GEDTM30 elevation/metadata tiles, so we can avoid building our own GDAL/tiler layer and just point clients at its endpoints (or reuse its code/config) for mean/std-dev/mask lookups. It also serves a basic web UI on the same HTTP port (8080 by default) to click a map for elevations/profiles.

- 310GB for the whole world. SO where to put it ?

- https://zenodo.org/records/15689805 for more info on how and why ?

Due to Zenodo's storage limitations, the original GEDTM30 dataset and its standard deviation map are provided via external links:

- https://s3.opengeohub.org/global/edtm/gedtm_rf_m_30m_s_20060101_20151231_go_epsg.4326.3855_v20250611.tif

- https://s3.opengeohub.org/global/edtm/gedtm_rf_std_30m_s_20060101_20151231_go_epsg.4326.3855_v20250611.tif

- https://s3.opengeohub.org/global/edtm/gedtm_mask_c_120m_s_20060101_20151231_go_epsg.4326.3855_v20250611.tif

These S3 endpoints already serve the rasters over HTTPS with `Accept-Ranges: bytes`, so we can stream/tiler directly (treat as COG or re-COG if needed) without downloading the full ~310GB locally.

Endpoint check (200 OK + `Accept-Ranges: bytes`):
- Mean: 333,361,957,920 bytes, last-mod 2025-06-12
- Std dev: 384,493,960,083 bytes, last-mod 2025-06-12
- Mask: 1,523,948,460 bytes, last-mod 2025-06-12

I wonder if this data can also be usd to help fly with no GPS too ? 

Next:
- GDAL is how we stream/inspect/convert these giant GeoTIFFs; keep it containerized (`osgeo/gdal:ubuntu-small-3.6.3`) so we don’t install locally.
- Validate COG compliance on the URLs (gdalinfo/cogeo-validator); if not COG, re-COG and re-host for faster tiling. Example: `docker run --rm -e GDAL_DISABLE_READDIR_ON_OPEN=YES -e CPL_VSIL_CURL_ALLOWED_EXTENSIONS=.tif osgeo/gdal:ubuntu-small-3.6.3 gdalinfo --mdd COG --checksum https://s3.opengeohub.org/global/edtm/gedtm_rf_m_30m_s_20060101_20151231_go_epsg.4326.3855_v20250611.tif`
- Reproject/COG if needed: convert to the projection your clients use (e.g., Web Mercator/WGS84) and ensure it’s tiled with overviews for fast HTTP range reads.
- Stand up a tile API (titiler/terracotta) pointing at these URLs; front with CDN; expose mean/std-dev/mask as layers.
- Confirm licensing/usage terms from the Zenodo record for drone ops.
- Plan storage/serving (object storage with range reads) and integrate tile distribution into Bento + NATS JetStream if needed.

## RUnning worklaods

https://github.com/akhenakh/geo-bento since we run bento !!!

https://github.com/akhenakh/bento-cbor since we run bento !!!

https://github.com/akhenakh/nlock since we use NATS Jetstream

https://github.com/akhenakh/narun since we use NATS Jetstream
