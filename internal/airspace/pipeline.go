package airspace

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// PipelineOptions configures the pipeline.
type PipelineOptions struct {
	Force     bool
	TilerName string // "auto", "tippecanoe", "gotiler"
	Verbose   bool
}

// PipelineResult contains the outcome of a pipeline run.
type PipelineResult struct {
	SyncResult *SyncResult
	TileCount  int
	Skipped    bool // True if no changes and not forced
}

// Pipeline runs the full sync → tile → manifest pipeline.
func Pipeline(opts PipelineOptions, tiler Tiler) (*PipelineResult, error) {
	result := &PipelineResult{}

	// Step 1: Sync
	syncOpts := DefaultSyncOptions()
	syncOpts.Force = opts.Force

	syncResult, err := Sync(syncOpts)
	if err != nil {
		return nil, fmt.Errorf("sync: %w", err)
	}
	result.SyncResult = syncResult

	// Check if we should continue
	if !syncResult.HasChanges && !opts.Force {
		result.Skipped = true
		return result, nil
	}

	// Step 2: Tile
	tileCount, err := TileAll(tiler, opts.Force)
	if err != nil {
		return nil, fmt.Errorf("tile: %w", err)
	}
	result.TileCount = tileCount

	// Step 3: Manifest
	if err := GenerateManifests(DirData, DirPMTiles, DirGeoJSON); err != nil {
		return nil, fmt.Errorf("manifest: %w", err)
	}

	return result, nil
}

// TileOptions configures tile generation.
type TileOptions struct {
	Force       bool
	DatasetKeys []string // nil = all datasets
}

// TileAll generates PMTiles for all datasets.
func TileAll(tiler Tiler, force bool) (int, error) {
	if err := os.MkdirAll(DirPMTiles, 0755); err != nil {
		return 0, fmt.Errorf("creating tiles dir: %w", err)
	}

	tiled := 0
	for _, key := range DatasetOrder {
		ds := Datasets[key]
		geoJSONPath := filepath.Join(DirGeoJSON, ds.GeoJSON)
		pmTilesPath := filepath.Join(DirPMTiles, ds.PMTiles)

		// Check if GeoJSON exists
		geoJSONInfo, err := os.Stat(geoJSONPath)
		if err != nil {
			continue // Skip missing
		}

		// Check if PMTiles is newer than GeoJSON (skip if up-to-date)
		if !force {
			if pmTilesInfo, err := os.Stat(pmTilesPath); err == nil {
				if pmTilesInfo.ModTime().After(geoJSONInfo.ModTime()) {
					continue // Up to date
				}
			}
		}

		// Build config
		cfg := TileConfigs[key]
		cfg.Layer = ds.Layer

		// For gotiler, use sensible defaults when auto-zoom is requested
		if tiler.Name() == "go" && cfg.MinZoom < 0 {
			cfg.MinZoom = 0
		}
		if tiler.Name() == "go" && cfg.MaxZoom < 0 {
			cfg.MaxZoom = 10
		}

		if err := tiler.Tile(geoJSONPath, pmTilesPath, cfg); err != nil {
			return tiled, fmt.Errorf("tiling %s: %w", key, err)
		}
		tiled++
	}

	return tiled, nil
}

// TileOne generates PMTiles for a single dataset.
func TileOne(tiler Tiler, key string, force bool) error {
	ds, ok := Datasets[key]
	if !ok {
		return fmt.Errorf("unknown dataset: %s", key)
	}

	geoJSONPath := filepath.Join(DirGeoJSON, ds.GeoJSON)
	pmTilesPath := filepath.Join(DirPMTiles, ds.PMTiles)

	// Check if GeoJSON exists
	geoJSONInfo, err := os.Stat(geoJSONPath)
	if err != nil {
		return fmt.Errorf("GeoJSON not found: %s", ds.GeoJSON)
	}

	// Check if PMTiles is newer than GeoJSON
	if !force {
		if pmTilesInfo, err := os.Stat(pmTilesPath); err == nil {
			if pmTilesInfo.ModTime().After(geoJSONInfo.ModTime()) {
				return nil // Up to date
			}
		}
	}

	// Build config
	cfg := TileConfigs[key]
	cfg.Layer = ds.Layer

	// For gotiler, use sensible defaults when auto-zoom is requested
	if tiler.Name() == "go" && cfg.MinZoom < 0 {
		cfg.MinZoom = 0
	}
	if tiler.Name() == "go" && cfg.MaxZoom < 0 {
		cfg.MaxZoom = 10
	}

	if err := os.MkdirAll(DirPMTiles, 0755); err != nil {
		return fmt.Errorf("creating tiles dir: %w", err)
	}

	return tiler.Tile(geoJSONPath, pmTilesPath, cfg)
}

// SelectTiler returns the appropriate tiler based on name.
// "auto" tries tippecanoe first, falls back to gotiler.
func SelectTiler(name string, tippecanoe, gotiler Tiler) (Tiler, error) {
	switch name {
	case "tippecanoe":
		if !tippecanoe.Available() {
			return nil, fmt.Errorf("tippecanoe not found in PATH")
		}
		return tippecanoe, nil
	case "gotiler", "go":
		return gotiler, nil
	case "auto", "":
		if tippecanoe.Available() {
			return tippecanoe, nil
		}
		return gotiler, nil
	default:
		return nil, fmt.Errorf("unknown tiler: %s (valid: auto, tippecanoe, gotiler)", name)
	}
}

// IsTippecanoeAvailable checks if tippecanoe is installed.
func IsTippecanoeAvailable() bool {
	_, err := exec.LookPath("tippecanoe")
	return err == nil
}
