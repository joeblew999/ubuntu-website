// Package gotiler provides pure Go tile generation for CI environments.
//
// This replaces tippecanoe for GitHub Actions and other environments where
// installing C++ binaries is impractical. The generated PMTiles files are
// then uploaded to R2 for serving via Cloudflare CDN.
//
// Uses paulmach/orb for geometry and protomaps/go-pmtiles for output.
//
// # Architecture Decision: Why Not Cloudflare Workers?
//
// This package runs in GitHub Actions, not Cloudflare Workers, because:
//
//  1. WASM compatibility: go-pmtiles imports zombiezen.com/go/sqlite which
//     depends on modernc.org/libc - this doesn't compile to WASM.
//
//  2. Resource limits: Workers have 10-50ms CPU limits and memory constraints.
//     Processing 422MB of GeoJSON features would exceed these limits.
//
//  3. Batch vs streaming: Tile generation requires processing all features
//     to determine tile intersections - this is inherently a batch operation.
//
//  4. PMTiles design: PMTiles is optimized for static hosting with HTTP Range
//     requests. The browser's pmtiles.js library handles this efficiently.
//
// The current architecture (GitHub Actions → PMTiles → R2 → CDN → Browser)
// is optimal because FAA data updates weekly (AIRAC cycle) and PMTiles
// enables efficient CDN caching with Range request support.
package gotiler

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/simplify"
	"github.com/protomaps/go-pmtiles/pmtiles"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
)

// GoTiler implements airspace.Tiler using pure Go libraries.
type GoTiler struct{}

// New creates a new GoTiler.
func New() *GoTiler {
	return &GoTiler{}
}

// Name returns the engine name.
func (g *GoTiler) Name() string {
	return "go"
}

// Available always returns true (pure Go, no external deps).
func (g *GoTiler) Available() bool {
	return true
}

// Tile converts GeoJSON to PMTiles using pure Go.
func (g *GoTiler) Tile(inputPath, outputPath string, config airspace.TileConfig) error {
	// Read GeoJSON
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading geojson: %w", err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return fmt.Errorf("parsing geojson: %w", err)
	}

	// Determine zoom range
	minZoom := config.MinZoom
	maxZoom := config.MaxZoom
	if minZoom < 0 {
		minZoom = 0
	}
	if maxZoom < 0 || maxZoom > 14 {
		maxZoom = 14
	}

	// Generate tiles for each zoom level
	tiles := make(map[maptile.Tile][]byte)

	for z := minZoom; z <= maxZoom; z++ {
		zoomTiles := g.generateZoomLevel(fc, uint32(z), config.Layer)
		for tile, data := range zoomTiles {
			tiles[tile] = data
		}
	}

	// Write PMTiles
	return writePMTiles(outputPath, tiles, config)
}

// generateZoomLevel creates MVT tiles for a specific zoom level.
func (g *GoTiler) generateZoomLevel(fc *geojson.FeatureCollection, zoom uint32, layerName string) map[maptile.Tile][]byte {
	result := make(map[maptile.Tile][]byte)

	// Group features by tile
	tileFeatures := make(map[maptile.Tile][]*geojson.Feature)

	for _, f := range fc.Features {
		// Get all tiles that intersect this feature's bounds
		bounds := f.Geometry.Bound()
		tiles := tilesInBounds(bounds, zoom)

		for _, tile := range tiles {
			tileFeatures[tile] = append(tileFeatures[tile], f)
		}
	}

	// Generate MVT for each tile
	for tile, features := range tileFeatures {
		mvtData := g.createMVT(tile, features, layerName)
		if len(mvtData) > 0 {
			result[tile] = mvtData
		}
	}

	return result
}

// createMVT creates an MVT tile from features.
func (g *GoTiler) createMVT(tile maptile.Tile, features []*geojson.Feature, layerName string) []byte {
	// Create a FeatureCollection for this tile
	fc := geojson.NewFeatureCollection()
	tileBound := tile.Bound()

	for _, f := range features {
		// Skip if geometry doesn't intersect tile
		if !f.Geometry.Bound().Intersects(tileBound) {
			continue
		}

		// Clone the feature to avoid modifying the original
		clone := geojson.NewFeature(f.Geometry)
		for k, v := range f.Properties {
			clone.Properties[k] = v
		}
		fc.Append(clone)
	}

	// Skip empty tiles
	if len(fc.Features) == 0 {
		return nil
	}

	// Create layer from FeatureCollection
	layer := mvt.NewLayer(layerName, fc)

	// Simplify based on zoom level - less detail at lower zooms
	epsilon := simplifyEpsilon(tile.Z)
	if epsilon > 0 {
		layer.Simplify(simplify.DouglasPeucker(epsilon))
	}

	// Clip to tile bounds
	layer.Clip(tileBound)

	// Project to tile coordinates (0-4096 extent)
	layer.ProjectToTile(tile)

	// Remove empty features after clipping
	layer.RemoveEmpty(1.0, 1.0)

	// Skip if all features were removed
	if len(layer.Features) == 0 {
		return nil
	}

	// Encode to protobuf
	layers := mvt.Layers{layer}
	data, err := mvt.MarshalGzipped(layers)
	if err != nil {
		return nil
	}

	return data
}

// tilesInBounds returns all tiles at a zoom level that intersect a bounding box.
func tilesInBounds(bounds orb.Bound, zoom uint32) []maptile.Tile {
	// Get corner tiles
	minTile := maptile.At(bounds.Min, maptile.Zoom(zoom))
	maxTile := maptile.At(bounds.Max, maptile.Zoom(zoom))

	// Ensure min/max are ordered correctly
	minX, maxX := minTile.X, maxTile.X
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	minY, maxY := minTile.Y, maxTile.Y
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	var tiles []maptile.Tile
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			tiles = append(tiles, maptile.New(x, y, maptile.Zoom(zoom)))
		}
	}
	return tiles
}

// simplifyEpsilon returns the simplification tolerance for a zoom level.
func simplifyEpsilon(zoom maptile.Zoom) float64 {
	// Higher zoom = less simplification
	switch {
	case zoom >= 14:
		return 0 // No simplification
	case zoom >= 10:
		return 0.00001
	case zoom >= 6:
		return 0.0001
	case zoom >= 4:
		return 0.001
	default:
		return 0.01
	}
}

// Ensure GoTiler implements Tiler.
var _ airspace.Tiler = (*GoTiler)(nil)

// writePMTiles writes tiles to a PMTiles file using the official go-pmtiles library.
// PMTiles v3 format: https://github.com/protomaps/PMTiles/blob/main/spec/v3/spec.md
func writePMTiles(path string, tiles map[maptile.Tile][]byte, config airspace.TileConfig) error {
	if len(tiles) == 0 {
		return fmt.Errorf("no tiles to write")
	}

	// Convert tiles to PMTiles entries, sorted by tile ID
	type tileEntry struct {
		id   uint64
		data []byte
	}
	var tileEntries []tileEntry

	for t, data := range tiles {
		// Use pmtiles ZxyToID for proper tile ID encoding
		id := pmtiles.ZxyToID(uint8(t.Z), uint32(t.X), uint32(t.Y))
		tileEntries = append(tileEntries, tileEntry{id: id, data: data})
	}

	// Sort by tile ID for clustered output
	sort.Slice(tileEntries, func(i, j int) bool {
		return tileEntries[i].id < tileEntries[j].id
	})

	// Build directory entries and collect tile data
	var entries []pmtiles.EntryV3
	var tileData bytes.Buffer
	currentOffset := uint64(0)

	for _, te := range tileEntries {
		entries = append(entries, pmtiles.EntryV3{
			TileID:    te.id,
			Offset:    currentOffset,
			Length:    uint32(len(te.data)),
			RunLength: 1,
		})
		tileData.Write(te.data)
		currentOffset += uint64(len(te.data))
	}

	// Build metadata JSON
	metadata := map[string]any{
		"name":        config.Layer,
		"format":      "pbf",
		"compression": "gzip",
		"minzoom":     config.MinZoom,
		"maxzoom":     config.MaxZoom,
	}
	metadataBytes, err := pmtiles.SerializeMetadata(metadata, pmtiles.Gzip)
	if err != nil {
		return fmt.Errorf("serializing metadata: %w", err)
	}

	// Serialize the root directory with gzip compression
	rootDirBytes := pmtiles.SerializeEntries(entries, pmtiles.Gzip)

	// Calculate offsets
	headerSize := uint64(pmtiles.HeaderV3LenBytes)
	rootDirOffset := headerSize
	rootDirLen := uint64(len(rootDirBytes))
	metadataOffset := rootDirOffset + rootDirLen
	metadataLen := uint64(len(metadataBytes))
	tileDataOffset := metadataOffset + metadataLen
	tileDataLen := uint64(tileData.Len())

	// Build header
	header := pmtiles.HeaderV3{
		SpecVersion:         3,
		RootOffset:          rootDirOffset,
		RootLength:          rootDirLen,
		MetadataOffset:      metadataOffset,
		MetadataLength:      metadataLen,
		LeafDirectoryOffset: 0, // No leaf directories for small files
		LeafDirectoryLength: 0,
		TileDataOffset:      tileDataOffset,
		TileDataLength:      tileDataLen,
		AddressedTilesCount: uint64(len(entries)),
		TileEntriesCount:    uint64(len(entries)),
		TileContentsCount:   uint64(len(entries)), // No deduplication
		Clustered:           true,
		InternalCompression: pmtiles.Gzip,
		TileCompression:     pmtiles.Gzip,
		TileType:            pmtiles.Mvt,
		MinZoom:             uint8(config.MinZoom),
		MaxZoom:             uint8(config.MaxZoom),
	}

	// Serialize header
	headerBytes := pmtiles.SerializeHeader(header)

	// Write file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header
	if _, err := f.Write(headerBytes); err != nil {
		return err
	}

	// Write root directory
	if _, err := f.Write(rootDirBytes); err != nil {
		return err
	}

	// Write metadata
	if _, err := f.Write(metadataBytes); err != nil {
		return err
	}

	// Write tile data
	if _, err := f.Write(tileData.Bytes()); err != nil {
		return err
	}

	return nil
}
