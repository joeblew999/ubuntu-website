// Package airspace is the FAA airspace data behind the BVLOS demo (/fleet/airspace-demo/): FAA ArcGIS → GeoJSON →
// PMTiles → R2, which the demo's map reads; cmd/airspace runs it (mise run airspace:pipeline, then airspace:upload),
// weekly in .github/workflows/sync-airspace.yml.
//
// # The sync rules
//
// What decides whether a dataset is downloaded, tiled and uploaded. Sync, TileAll, GenerateManifests and UploadToR2
// each apply their part; this is the whole of it.
//
//  1. A dataset's version is its FAA FeatureServer layer's editingInfo.dataLastEditDate (Dataset.LayerURL, ?f=json):
//     when the FAA last changed its data. The FAA publishes no delta (the layers answer queries only: no sync, no
//     extractChanges, no editor tracking) and republishes each layer whole every 28-day AIRAC cycle, so the version
//     is all there is to know: a dataset changed or it didn't, and a changed one is downloaded whole.
//  2. What's published is the manifest in R2 (PublishedURL, airspace/manifest_usa.json): each layer's SourceVersion is
//     the version its tiles there were built from. It lives beside the tiles it describes, so it's the same for CI,
//     a laptop or a fresh clone; nothing is cached between runs.
//  3. Sync downloads a dataset when its version isn't the published one (or there's none: the first run), or with
//     --force. A dataset at the published version is skipped even when this machine has no file of it: R2 already
//     serves it. (For local files of an unchanged dataset, sync --force.)
//  4. A dataset whose version or download fails fails the pipeline: tiling around it would hide why.
//  5. The run's datasets are tiled; the manifest takes their size, features and version from this run's files, and
//     every other layer's from the published manifest, so it always describes all of R2.
//  6. Upload sends only the tiles this run built, then the manifests, last: the published versions move only once
//     the tiles they describe are up. A run that fails before that leaves the published manifest as it was, so the
//     next run downloads the same datasets again.
//
// Weekly, three runs in four find nothing changed and stop after six small requests.
package airspace
