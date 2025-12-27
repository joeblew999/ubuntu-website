package airspace

// Dataset describes a FAA data source.
type Dataset struct {
	Name        string // Human-readable name
	Key         string // Dataset key (uas, boundary, etc.)
	GeoJSON     string // GeoJSON filename
	PMTiles     string // PMTiles filename
	Layer       string // PMTiles layer name
	BaseURL     string // Download URL
	IsPaginated bool   // For FeatureServer APIs that require pagination
	PageSize    int    // Page size for paginated APIs
	ETagURL     string // URL to check for ETag/Last-Modified (for paginated APIs)
}

// Datasets is the registry of all FAA data sources.
var Datasets = map[string]Dataset{
	"uas": {
		Name:        "UAS Facility Map",
		Key:         "uas",
		GeoJSON:     GeoJSONUAS,
		PMTiles:     PMTilesUAS,
		Layer:       LayerUAS,
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    2000,
		ETagURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0",
	},
	"boundary": {
		Name:    "Airspace Boundary",
		Key:     "boundary",
		GeoJSON: GeoJSONBoundary,
		PMTiles: PMTilesBoundary,
		Layer:   LayerBoundary,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/67885972e4e940b2aa6d74024901c561/geojson?layers=0",
	},
	"sua": {
		Name:    "Special Use Airspace",
		Key:     "sua",
		GeoJSON: GeoJSONSUA,
		PMTiles: PMTilesSUA,
		Layer:   LayerSUA,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/dd0d1b726e504137ab3c41b21835d05b/geojson?layers=0",
	},
	"airports": {
		Name:    "Airports",
		Key:     "airports",
		GeoJSON: GeoJSONAirports,
		PMTiles: PMTilesAirports,
		Layer:   LayerAirports,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/e747ab91a11045e8b3f8a3efd093d3b5/geojson?layers=0",
	},
	"navaids": {
		Name:    "Navigation Aids",
		Key:     "navaids",
		GeoJSON: GeoJSONNavaids,
		PMTiles: PMTilesNavaids,
		Layer:   LayerNavaids,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/990e238991b44dd08af27d7b43e70b92/geojson?layers=0",
	},
	"obstacles": {
		Name:    "Obstacles",
		Key:     "obstacles",
		GeoJSON: GeoJSONObstacles,
		PMTiles: PMTilesObstacles,
		Layer:   LayerObstacles,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/c6a62360338e408cb1512366ad61559e/geojson?layers=0",
	},
}

// TileConfigs holds tile generation settings per dataset.
var TileConfigs = map[string]TileConfig{
	"boundary":  {MinZoom: -1, MaxZoom: -1, DropDensest: false},                                        // -zg (auto zoom)
	"sua":       {MinZoom: -1, MaxZoom: -1, DropDensest: false},                                        // -zg
	"uas":       {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true}, // explicit zoom, no reduction
	"airports":  {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"navaids":   {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"obstacles": {MinZoom: -1, MaxZoom: -1, DropDensest: true},                                         // -zg --drop-densest-as-needed
}
