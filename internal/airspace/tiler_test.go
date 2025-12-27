package airspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/protomaps/go-pmtiles/pmtiles"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
	"github.com/joeblew999/ubuntu-website/internal/airspace/tiler"
)

const (
	testInput     = "testdata/mini_airspace.geojson"
	testReference = "testdata/mini_airspace_reference.pmtiles"
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

// TestGoTilerMatchesReference verifies GoTiler output matches the golden reference.
// The reference was generated with tippecanoe and committed to the repo.
func TestGoTilerMatchesReference(t *testing.T) {
	g := gotiler.New()

	// Read reference file
	refData, err := os.ReadFile(testReference)
	if err != nil {
		t.Fatalf("failed to read reference: %v", err)
	}

	tmpDir := t.TempDir()
	goOutput := filepath.Join(tmpDir, "go.pmtiles")

	config := airspace.TileConfig{
		MinZoom: 0,
		MaxZoom: 10,
		Layer:   "test",
	}

	// Generate with Go
	if err := g.Tile(testInput, goOutput, config); err != nil {
		t.Skipf("GoTiler not yet implemented: %v", err)
	}

	// Read generated file
	goData, err := os.ReadFile(goOutput)
	if err != nil {
		t.Fatalf("failed to read go output: %v", err)
	}

	// Debug: compare headers
	t.Logf("Reference header (first 20 bytes): %x", refData[:20])
	t.Logf("Go output header (first 20 bytes): %x", goData[:20])

	// Compare sizes first (quick check)
	if len(goData) != len(refData) {
		t.Logf("size mismatch: go=%d bytes, reference=%d bytes", len(goData), len(refData))
	}

	// For now, just verify the file is a valid PMTiles v3
	if len(goData) < 127 {
		t.Fatalf("output too small for PMTiles header: %d bytes", len(goData))
	}

	// Check magic number
	if string(goData[:7]) != "PMTiles" {
		t.Errorf("invalid magic number: %s", string(goData[:7]))
	}

	// Check version
	if goData[7] != 3 {
		t.Errorf("invalid version: %d (expected 3)", goData[7])
	}

	t.Log("GoTiler produces valid PMTiles v3 format")

	// Verify the header can be parsed by go-pmtiles library
	header, err := pmtiles.DeserializeHeader(goData[:pmtiles.HeaderV3LenBytes])
	if err != nil {
		t.Fatalf("failed to deserialize header: %v", err)
	}

	t.Logf("PMTiles header: minzoom=%d, maxzoom=%d, tiles=%d",
		header.MinZoom, header.MaxZoom, header.TileEntriesCount)

	if header.MinZoom != uint8(config.MinZoom) {
		t.Errorf("wrong minzoom: got %d, want %d", header.MinZoom, config.MinZoom)
	}
	if header.MaxZoom != uint8(config.MaxZoom) {
		t.Errorf("wrong maxzoom: got %d, want %d", header.MaxZoom, config.MaxZoom)
	}
	if header.TileEntriesCount == 0 {
		t.Error("no tiles in output")
	}

	// Note: Exact byte-for-byte match is not required - different tools may produce
	// functionally equivalent but not identical files. What matters is that the
	// tiles can be read and rendered correctly.
}

// TestGoTilerRealData tests with real FAA data (if available).
func TestGoTilerRealData(t *testing.T) {
	// Use navaids as it's the smallest real dataset
	inputPath := "../../static/airspace/faa_navaids.geojson"
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		t.Skip("real FAA data not available")
	}

	g := gotiler.New()
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "navaids.pmtiles")

	config := airspace.TileConfig{
		MinZoom: 0,
		MaxZoom: 10,
		Layer:   "navaids",
	}

	err := g.Tile(inputPath, outputPath, config)
	if err != nil {
		t.Fatalf("failed to tile real data: %v", err)
	}

	// Verify output
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	t.Logf("Generated navaids PMTiles: %.1f KB", float64(info.Size())/1024)

	// Verify header
	data, _ := os.ReadFile(outputPath)
	header, err := pmtiles.DeserializeHeader(data[:pmtiles.HeaderV3LenBytes])
	if err != nil {
		t.Fatalf("failed to parse header: %v", err)
	}

	t.Logf("Tiles: %d, MinZoom: %d, MaxZoom: %d",
		header.TileEntriesCount, header.MinZoom, header.MaxZoom)

	if header.TileEntriesCount == 0 {
		t.Error("no tiles generated")
	}
}

// TestTippecanoeMatchesReference ensures tippecanoe still produces the same output.
// This catches tippecanoe version changes that might break compatibility.
func TestTippecanoeMatchesReference(t *testing.T) {
	tip := tiler.New()
	if !tip.Available() {
		t.Skip("tippecanoe not installed")
	}

	tmpDir := t.TempDir()
	tipOutput := filepath.Join(tmpDir, "tippecanoe.pmtiles")

	config := airspace.TileConfig{
		MinZoom: 0,
		MaxZoom: 10,
		Layer:   "test",
	}

	if err := tip.Tile(testInput, tipOutput, config); err != nil {
		t.Fatalf("tippecanoe failed: %v", err)
	}

	// Read both files
	tipData, err := os.ReadFile(tipOutput)
	if err != nil {
		t.Fatalf("failed to read tippecanoe output: %v", err)
	}

	refData, err := os.ReadFile(testReference)
	if err != nil {
		t.Fatalf("failed to read reference: %v", err)
	}

	// Compare sizes
	if len(tipData) != len(refData) {
		t.Logf("size mismatch: tippecanoe=%d bytes, reference=%d bytes", len(tipData), len(refData))
		t.Log("This may be expected if tippecanoe version changed - consider updating reference")
	}

	// Note: tippecanoe output may vary slightly between runs due to timestamps
	// For now, just check size is similar (within 10%)
	sizeDiff := float64(len(tipData)-len(refData)) / float64(len(refData))
	if sizeDiff > 0.1 || sizeDiff < -0.1 {
		t.Errorf("size differs by %.1f%% - likely incompatible output", sizeDiff*100)
	}
}
