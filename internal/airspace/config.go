// Package airspace provides tile generation for FAA airspace data.
package airspace

// =============================================================================
// Directory Paths
// =============================================================================

// Directory paths - single source of truth for all file locations.
const (
	DirGeoJSON = "static/airspace"       // GeoJSON output directory
	DirPMTiles = "static/airspace/tiles" // PMTiles output directory
	DirData    = "data/airspace"         // Data/metadata directory
)

// =============================================================================
// Data Files (in DirData)
// =============================================================================

// Data files use underscores for Hugo data access compatibility.
const (
	FileSyncETags   = "sync_etags.json"   // ETag cache for change detection
	FileSyncResult  = "sync_result.json"  // Last sync result (for pipeline idempotency)
	FileSyncHistory = "sync_history.json" // Rolling sync history
	FileManifest    = "manifest.json"     // Global manifest
	FileUSAManifest = "manifest_usa.json" // USA regional manifest
)

// =============================================================================
// GeoJSON Filenames (in DirGeoJSON)
// =============================================================================

const (
	GeoJSONBoundary  = "faa_airspace_boundary.geojson"
	GeoJSONSUA       = "faa_special_use_airspace.geojson"
	GeoJSONUAS       = "faa_uas_facility_map.geojson"
	GeoJSONAirports  = "faa_airports.geojson"
	GeoJSONNavaids   = "faa_navaids.geojson"
	GeoJSONObstacles = "faa_obstacles.geojson"
)

// =============================================================================
// PMTiles Filenames (in DirPMTiles)
// =============================================================================

const (
	PMTilesBoundary  = "faa_airspace_boundary.pmtiles"
	PMTilesSUA       = "faa_special_use_airspace.pmtiles"
	PMTilesUAS       = "faa_uas_facility_map.pmtiles"
	PMTilesAirports  = "faa_airports.pmtiles"
	PMTilesNavaids   = "faa_navaids.pmtiles"
	PMTilesObstacles = "faa_obstacles.pmtiles"
	PMTilesCombined  = "faa_airspace_combined.pmtiles"
)

// =============================================================================
// Layer Names (for tippecanoe and map rendering)
// =============================================================================

const (
	LayerBoundary  = "boundary"
	LayerSUA       = "sua"
	LayerUAS       = "uas"
	LayerAirports  = "airports"
	LayerNavaids   = "navaids"
	LayerObstacles = "obstacles"
)

// =============================================================================
// R2 Storage Configuration
// =============================================================================

const (
	R2Bucket    = "ubuntu-website-assets"
	R2PublicURL = "https://pub-97cfaeb734ae474c80c79c3e3cc6dbee.r2.dev"
)

// =============================================================================
// Sync Configuration
// =============================================================================

const (
	MaxHistoryRuns = 20 // Maximum sync history entries to keep
)

// =============================================================================
// Dataset Processing Order
// =============================================================================

// DatasetOrder defines the default processing order.
// Note: obstacles is excluded by default due to large file size.
var DatasetOrder = []string{"uas", "boundary", "sua", "airports", "navaids"}

// AllDatasets includes obstacles for commands that need it.
var AllDatasets = []string{"uas", "boundary", "sua", "airports", "navaids", "obstacles"}
