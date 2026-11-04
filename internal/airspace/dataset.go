package airspace

// Dataset describes a FAA data source.
type Dataset struct {
	Name        string // Human-readable name
	Key         string // Dataset key (uas, boundary, etc.)
	GeoJSON     string // GeoJSON filename
	PMTiles     string // PMTiles filename
	Layer       string // PMTiles layer name
	BaseURL     string // Download URL: its FeatureServer layer's query (paginated), or a file
	IsPaginated bool   // A FeatureServer query, downloaded in pages
	PageSize    int    // The layer's maxRecordCount: the most a page holds
	LayerURL    string // Its FAA FeatureServer layer, whose editingInfo.dataLastEditDate is its version (doc.go, rule 1)
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
		LayerURL:    "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0",
	},
	"boundary": {
		Name:        "Airspace Boundary",
		Key:         "boundary",
		GeoJSON:     GeoJSONBoundary,
		PMTiles:     PMTilesBoundary,
		Layer:       LayerBoundary,
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Boundary_Airspace/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    2000,
		LayerURL:    "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Boundary_Airspace/FeatureServer/0",
	},
	"sua": {
		Name:        "Special Use Airspace",
		Key:         "sua",
		GeoJSON:     GeoJSONSUA,
		PMTiles:     PMTilesSUA,
		Layer:       LayerSUA,
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Special_Use_Airspace/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    2000,
		LayerURL:    "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Special_Use_Airspace/FeatureServer/0",
	},
	"airports": {
		Name:        "Airports",
		Key:         "airports",
		GeoJSON:     GeoJSONAirports,
		PMTiles:     PMTilesAirports,
		Layer:       LayerAirports,
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/US_Airport/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    1000,
		LayerURL:    "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/US_Airport/FeatureServer/0",
	},
	"navaids": {
		Name:        "Navigation Aids",
		Key:         "navaids",
		GeoJSON:     GeoJSONNavaids,
		PMTiles:     PMTilesNavaids,
		Layer:       LayerNavaids,
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/NAVAIDSystem/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    1000,
		LayerURL:    "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/NAVAIDSystem/FeatureServer/0",
	},
	"obstacles": {
		Name:    "Obstacles",
		Key:     "obstacles",
		GeoJSON: GeoJSONObstacles,
		PMTiles: PMTilesObstacles,
		Layer:   LayerObstacles,
		BaseURL: "https://adds-faa.opendata.arcgis.com/api/download/v1/items/c6a62360338e408cb1512366ad61559e/geojson?layers=0",
		// The FAA item behind BaseURL is Class_Airspace, not obstacles (TODO.md)
		LayerURL: "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Class_Airspace/FeatureServer/0",
	},
}

// TileConfigs holds tile generation settings per dataset.
var TileConfigs = map[string]TileConfig{
	"boundary":  {MinZoom: -1, MaxZoom: -1, DropDensest: false},                                        // -zg (auto zoom)
	"sua":       {MinZoom: -1, MaxZoom: -1, DropDensest: false},                                        // -zg
	"uas":       {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true}, // explicit zoom, no reduction
	"airports":  {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"navaids":   {MinZoom: 0, MaxZoom: 10, ReduceRate: 1, NoFeatureLimit: true, NoTileSizeLimit: true},
	"obstacles": {MinZoom: -1, MaxZoom: -1, DropDensest: true}, // -zg --drop-densest-as-needed
}
