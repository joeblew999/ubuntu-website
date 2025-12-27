package airspace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ManifestGlobal is the top-level manifest structure.
type ManifestGlobal struct {
	Version int                       `json:"version"`
	Updated string                    `json:"updated"`
	Regions map[string]ManifestRegion `json:"regions"`
	Notes   map[string]string         `json:"notes,omitempty"`
}

// ManifestRegion describes a geographic region.
type ManifestRegion struct {
	Name          string    `json:"name"`
	BBox          []float64 `json:"bbox"`
	TilesPath     string    `json:"tiles_path"`
	ManifestFile  string    `json:"manifest_file"`
	DefaultLayers []string  `json:"default_layers"`
}

// ManifestUSA is the USA regional manifest structure.
type ManifestUSA struct {
	Region  string                   `json:"region"`
	Name    string                   `json:"name"`
	Version int                      `json:"version"`
	Updated string                   `json:"updated"`
	BBox    []float64                `json:"bbox"`
	Layers  map[string]ManifestLayer `json:"layers"`
	Source  ManifestSource           `json:"source"`
}

// ManifestLayer describes a data layer.
type ManifestLayer struct {
	Name           string        `json:"name"`
	File           string        `json:"file"`
	PMTilesLayer   string        `json:"pmtiles_layer"`
	GeomType       string        `json:"geom_type"` // polygon, point, line
	SizeMB         float64       `json:"size_mb"`
	Features       int           `json:"features"`
	ZoomRange      []int         `json:"zoom_range"`
	DefaultVisible bool          `json:"default_visible"`
	RenderRules    []RenderRule  `json:"render_rules"`
	Legend         []LegendEntry `json:"legend,omitempty"`
}

// RenderRule defines how to style features.
type RenderRule struct {
	FilterProp  string  `json:"filter_prop,omitempty"`
	FilterValue string  `json:"filter_value,omitempty"`
	Fill        string  `json:"fill"`
	Stroke      string  `json:"stroke,omitempty"`
	Opacity     float64 `json:"opacity,omitempty"`
	Width       float64 `json:"width,omitempty"`
	Radius      float64 `json:"radius,omitempty"` // For points
}

// LegendEntry for UI layer toggles.
type LegendEntry struct {
	Label string `json:"label"`
	Color string `json:"color"`
}

// ManifestSource describes the data source.
type ManifestSource struct {
	Authority   string            `json:"authority"`
	URLs        map[string]string `json:"urls"`
	UpdateCycle string            `json:"update_cycle"`
}

// LayerMetrics holds computed metrics for a layer.
type LayerMetrics struct {
	SizeMB   float64
	Features int
}

// GenerateManifests creates both global and USA manifests.
func GenerateManifests(dataDir, tilesDir, geoJSONDir string) error {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	// Collect metrics for each layer
	layerMetrics := make(map[string]LayerMetrics)
	for _, key := range DatasetOrder {
		ds := Datasets[key]
		pmTilesPath := filepath.Join(tilesDir, ds.PMTiles)
		geoJSONPath := filepath.Join(geoJSONDir, ds.GeoJSON)

		var sizeMB float64
		var features int

		if info, err := os.Stat(pmTilesPath); err == nil {
			sizeMB = float64(info.Size()) / (1024 * 1024)
		}
		features = CountGeoJSONFeatures(geoJSONPath)

		layerMetrics[key] = LayerMetrics{sizeMB, features}
	}

	// Create global manifest
	globalManifest := ManifestGlobal{
		Version: 1,
		Updated: timestamp,
		Regions: map[string]ManifestRegion{
			"usa": {
				Name:          "United States",
				BBox:          []float64{-125, 24, -66, 50},
				TilesPath:     "tiles",
				ManifestFile:  FileUSAManifest,
				DefaultLayers: []string{"boundary", "sua"},
			},
		},
		Notes: map[string]string{
			"bbox_format":    "[west, south, east, north]",
			"tiles_path":     "Relative to /airspace/ in R2",
			"future_regions": "europe, canada, australia, japan",
		},
	}

	// Create USA manifest with all layers
	usaManifest := createUSAManifest(timestamp, layerMetrics)

	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// Write manifests
	globalPath := filepath.Join(dataDir, FileManifest)
	usaPath := filepath.Join(dataDir, FileUSAManifest)

	if err := writeJSON(globalPath, globalManifest); err != nil {
		return err
	}

	if err := writeJSON(usaPath, usaManifest); err != nil {
		return err
	}

	// Copy to static directory for local dev
	staticGlobal := filepath.Join(geoJSONDir, FileManifest)
	staticUSA := filepath.Join(geoJSONDir, FileUSAManifest)

	copyFile(globalPath, staticGlobal)
	copyFile(usaPath, staticUSA)

	return nil
}

func createUSAManifest(timestamp string, metrics map[string]LayerMetrics) ManifestUSA {
	usaManifest := ManifestUSA{
		Region:  "usa",
		Name:    "United States",
		Version: 1,
		Updated: timestamp,
		BBox:    []float64{-125, 24, -66, 50},
		Layers:  make(map[string]ManifestLayer),
		Source: ManifestSource{
			Authority:   "FAA",
			URLs:        map[string]string{"adds": "https://adds-faa.opendata.arcgis.com", "udds": "https://udds-faa.opendata.arcgis.com"},
			UpdateCycle: "28-day AIRAC",
		},
	}

	// Add layer definitions with render rules
	usaManifest.Layers["boundary"] = ManifestLayer{
		Name:           "Airspace Boundary",
		File:           PMTilesBoundary,
		PMTilesLayer:   LayerBoundary,
		GeomType:       "polygon",
		SizeMB:         metrics["boundary"].SizeMB,
		Features:       metrics["boundary"].Features,
		ZoomRange:      []int{4, 14},
		DefaultVisible: true,
		RenderRules: []RenderRule{
			{FilterProp: "CLASS", FilterValue: "A", Fill: "#0066cc", Stroke: "#0066cc", Opacity: 0.15, Width: 1},
			{FilterProp: "CLASS", FilterValue: "C", Fill: "#cc00cc", Stroke: "#cc00cc", Opacity: 0.2, Width: 2},
			{FilterProp: "CLASS", FilterValue: "D", Fill: "#0099cc", Stroke: "#0099cc", Opacity: 0.15, Width: 2},
			{FilterProp: "CLASS", FilterValue: "E", Fill: "#00cc99", Stroke: "#00cc99", Opacity: 0.1, Width: 1},
			{FilterProp: "CLASS", FilterValue: "G", Fill: "#999999", Stroke: "#999999", Opacity: 0.05, Width: 1},
			{Fill: "#666666", Stroke: "#666666", Opacity: 0.1, Width: 1}, // fallback
		},
		Legend: []LegendEntry{
			{Label: "Class A", Color: "#0066cc"},
			{Label: "Class C", Color: "#cc00cc"},
			{Label: "Class D", Color: "#0099cc"},
			{Label: "Class E", Color: "#00cc99"},
		},
	}

	usaManifest.Layers["sua"] = ManifestLayer{
		Name:           "Special Use Airspace",
		File:           PMTilesSUA,
		PMTilesLayer:   LayerSUA,
		GeomType:       "polygon",
		SizeMB:         metrics["sua"].SizeMB,
		Features:       metrics["sua"].Features,
		ZoomRange:      []int{4, 14},
		DefaultVisible: true,
		RenderRules: []RenderRule{
			{FilterProp: "TYPE_CODE", FilterValue: "R", Fill: "#cc0000", Stroke: "#cc0000", Opacity: 0.3, Width: 2},
			{FilterProp: "TYPE_CODE", FilterValue: "P", Fill: "#ff0000", Stroke: "#ff0000", Opacity: 0.4, Width: 2},
			{FilterProp: "TYPE_CODE", FilterValue: "MOA", Fill: "#ff9900", Stroke: "#ff9900", Opacity: 0.2, Width: 1},
			{FilterProp: "TYPE_CODE", FilterValue: "A", Fill: "#ffcc00", Stroke: "#ffcc00", Opacity: 0.2, Width: 1},
			{FilterProp: "TYPE_CODE", FilterValue: "W", Fill: "#996600", Stroke: "#996600", Opacity: 0.15, Width: 1},
			{Fill: "#666666", Stroke: "#666666", Opacity: 0.1, Width: 1}, // fallback
		},
		Legend: []LegendEntry{
			{Label: "Restricted", Color: "#cc0000"},
			{Label: "Prohibited", Color: "#ff0000"},
			{Label: "MOA", Color: "#ff9900"},
			{Label: "Alert", Color: "#ffcc00"},
			{Label: "Warning", Color: "#996600"},
		},
	}

	usaManifest.Layers["laanc"] = ManifestLayer{
		Name:           "LAANC/UAS Facility Map",
		File:           PMTilesUAS,
		PMTilesLayer:   LayerUAS,
		GeomType:       "polygon",
		SizeMB:         metrics["uas"].SizeMB,
		Features:       metrics["uas"].Features,
		ZoomRange:      []int{6, 14},
		DefaultVisible: false,
		RenderRules: []RenderRule{
			{Fill: "#ffff00", Stroke: "#cc9900", Opacity: 0.5, Width: 1},
		},
		Legend: []LegendEntry{
			{Label: "LAANC Grid", Color: "#ffff00"},
		},
	}

	usaManifest.Layers["airports"] = ManifestLayer{
		Name:           "Airports",
		File:           PMTilesAirports,
		PMTilesLayer:   LayerAirports,
		GeomType:       "point",
		SizeMB:         metrics["airports"].SizeMB,
		Features:       metrics["airports"].Features,
		ZoomRange:      []int{0, 10},
		DefaultVisible: false,
		RenderRules: []RenderRule{
			{Fill: "#00ff00", Stroke: "#006600", Width: 1, Radius: 5},
		},
		Legend: []LegendEntry{
			{Label: "Airport", Color: "#00ff00"},
		},
	}

	usaManifest.Layers["navaids"] = ManifestLayer{
		Name:           "Navigation Aids",
		File:           PMTilesNavaids,
		PMTilesLayer:   LayerNavaids,
		GeomType:       "point",
		SizeMB:         metrics["navaids"].SizeMB,
		Features:       metrics["navaids"].Features,
		ZoomRange:      []int{0, 10},
		DefaultVisible: false,
		RenderRules: []RenderRule{
			{Fill: "#ff00ff", Stroke: "#660066", Width: 1, Radius: 4},
		},
		Legend: []LegendEntry{
			{Label: "VOR/NDB", Color: "#ff00ff"},
		},
	}

	return usaManifest
}

// CountGeoJSONFeatures counts features in a GeoJSON file.
func CountGeoJSONFeatures(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return strings.Count(string(data), `"type":"Feature"`)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
