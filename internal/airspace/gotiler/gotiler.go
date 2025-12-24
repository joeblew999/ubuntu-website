// Package gotiler provides pure Go tile generation for Cloudflare Workers.
//
// This replaces tippecanoe for environments where C++ binaries can't run.
// Uses paulmach/orb for geometry and protomaps/go-pmtiles for output.
package gotiler

import (
	"fmt"

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
	// TODO: Implement using paulmach/orb and protomaps/go-pmtiles
	return fmt.Errorf("gotiler not yet implemented - use tippecanoe for now")
}

// Ensure GoTiler implements Tiler.
var _ airspace.Tiler = (*GoTiler)(nil)
