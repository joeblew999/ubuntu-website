package airspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
	"github.com/joeblew999/ubuntu-website/internal/airspace/tiler"
)

func TestTippecanoeAvailable(t *testing.T) {
	tip := tiler.New()
	if !tip.Available() {
		t.Skip("tippecanoe not installed - skipping")
	}
	t.Log("tippecanoe is available")
}

func TestGoTilerAvailable(t *testing.T) {
	g := gotiler.New()
	if !g.Available() {
		t.Error("GoTiler should always be available")
	}
}

func TestTippecanoeGeneratesTiles(t *testing.T) {
	tip := tiler.New()
	if !tip.Available() {
		t.Skip("tippecanoe not installed")
	}

	// Use mini test data
	inputPath := filepath.Join("testdata", "mini_airspace.geojson")
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("test data not found: %s", inputPath)
	}

	// Create temp output
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.pmtiles")

	config := airspace.TileConfig{
		MinZoom: 0,
		MaxZoom: 10,
		Layer:   "test",
	}

	err := tip.Tile(inputPath, outputPath, config)
	if err != nil {
		t.Fatalf("tippecanoe failed: %v", err)
	}

	// Verify output exists and has content
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}

	t.Logf("Generated %s (%.1f KB)", outputPath, float64(info.Size())/1024)
}

func TestGoTilerGeneratesTiles(t *testing.T) {
	g := gotiler.New()

	inputPath := filepath.Join("testdata", "mini_airspace.geojson")
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("test data not found: %s", inputPath)
	}

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test.pmtiles")

	config := airspace.TileConfig{
		MinZoom: 0,
		MaxZoom: 10,
		Layer:   "test",
	}

	err := g.Tile(inputPath, outputPath, config)
	if err != nil {
		// Expected to fail until implemented
		t.Skipf("GoTiler not yet implemented: %v", err)
	}

	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

// TestTilersEquivalent verifies both engines produce equivalent output.
// This is the key test for the Go replacement.
func TestTilersEquivalent(t *testing.T) {
	tip := tiler.New()
	g := gotiler.New()

	if !tip.Available() {
		t.Skip("tippecanoe not installed - can't compare")
	}

	inputPath := filepath.Join("testdata", "mini_airspace.geojson")
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Fatalf("test data not found: %s", inputPath)
	}

	tmpDir := t.TempDir()
	tipOutput := filepath.Join(tmpDir, "tippecanoe.pmtiles")
	goOutput := filepath.Join(tmpDir, "go.pmtiles")

	config := airspace.TileConfig{
		MinZoom: 0,
		MaxZoom: 10,
		Layer:   "test",
	}

	// Generate with tippecanoe (reference)
	if err := tip.Tile(inputPath, tipOutput, config); err != nil {
		t.Fatalf("tippecanoe failed: %v", err)
	}

	// Generate with Go
	if err := g.Tile(inputPath, goOutput, config); err != nil {
		t.Skipf("GoTiler not yet implemented: %v", err)
	}

	// TODO: Compare tile contents
	// - Parse both PMTiles
	// - Extract tiles at same z/x/y coordinates
	// - Compare features (geometry, properties)
	// - Allow for minor floating-point differences

	t.Log("Equivalence test requires PMTiles comparison - TODO")
}
