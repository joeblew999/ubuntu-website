//go:build js && wasm

// Package main is a Cloudflare Worker for tile generation.
//
// This is experimental - the gotiler package compiles to ~118KB WASM,
// but processing 422MB of GeoJSON would exceed Worker CPU limits.
// Currently tiles are generated in GitHub Actions and served from R2.
//
// Build: GOOS=js GOARCH=wasm go build -o worker.wasm ./cmd/worker-tiler
package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/joeblew999/ubuntu-website/internal/airspace"
	"github.com/joeblew999/ubuntu-website/internal/airspace/gotiler"
	"github.com/syumai/workers-go"
)

func main() {
	http.HandleFunc("/tile", handleTile)
	http.HandleFunc("/health", handleHealth)
	workers.Serve(nil)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","engine":"gotiler"}`))
}

func handleTile(w http.ResponseWriter, r *http.Request) {
	// This is a proof of concept - in practice, processing 422MB
	// of GeoJSON would exceed Worker limits. Tiles should be
	// pre-generated in GitHub Actions and served from R2.

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	// Read GeoJSON from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For small inputs, we could generate tiles on the fly
	if len(body) > 1024*1024 { // 1MB limit for Worker
		http.Error(w, "input too large for Worker processing", http.StatusRequestEntityTooLarge)
		return
	}

	// Write input to temp buffer (Workers have no filesystem)
	// In practice, we'd need to refactor gotiler to work with io.Reader/Writer
	_ = body

	// For now, just return a placeholder
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"error":"tile generation not yet implemented in Worker"}`)
	_ = gotiler.New()
	_ = airspace.TileConfig{}
}
