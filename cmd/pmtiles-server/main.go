// Package main provides a simple PMTiles server for local development.
//
// Hugo's dev server doesn't properly handle concurrent HTTP range requests,
// which PMTiles requires for efficient tile loading. This server runs alongside
// Hugo to serve PMTiles files with proper range request support.
//
// Usage:
//
//	pmtiles-server [options]
//	pmtiles-server -dir static/airspace/tiles -port 8765
//
// The server adds proper CORS headers so the frontend can access tiles
// from a different port than Hugo (1313).
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	dir := flag.String("dir", "static/airspace/tiles", "Directory containing PMTiles files")
	port := flag.Int("port", 8765, "Port to listen on")
	flag.Parse()

	// Verify directory exists
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", *dir)
	}

	// List available PMTiles files
	files, _ := filepath.Glob(filepath.Join(*dir, "*.pmtiles"))
	if len(files) == 0 {
		log.Printf("Warning: No .pmtiles files found in %s", *dir)
	} else {
		log.Printf("Serving %d PMTiles files from %s", len(files), *dir)
		for _, f := range files {
			log.Printf("  - %s", filepath.Base(f))
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle CORS preflight
		if r.Method == "OPTIONS" {
			setCORSHeaders(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Only serve .pmtiles files
		if !strings.HasSuffix(r.URL.Path, ".pmtiles") {
			http.NotFound(w, r)
			return
		}

		path := filepath.Join(*dir, filepath.Clean(r.URL.Path))

		// Security: ensure path is within dir
		absDir, _ := filepath.Abs(*dir)
		absPath, _ := filepath.Abs(path)
		if !strings.HasPrefix(absPath, absDir) {
			http.NotFound(w, r)
			return
		}

		f, err := os.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		size := stat.Size()

		// Set CORS headers
		setCORSHeaders(w)

		// Handle range request
		rangeHeader := r.Header.Get("Range")
		if rangeHeader != "" && strings.HasPrefix(rangeHeader, "bytes=") {
			start, end, err := parseRange(rangeHeader, size)
			if err != nil {
				http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
				return
			}

			length := end - start + 1

			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Length", strconv.FormatInt(length, 10))
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusPartialContent)

			if _, err := f.Seek(start, io.SeekStart); err != nil {
				return
			}
			io.CopyN(w, f, length)

			log.Printf("206 %s [%d-%d]", r.URL.Path, start, end)
		} else {
			// Full file request
			w.Header().Set("Accept-Ranges", "bytes")
			w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
			w.Header().Set("Content-Type", "application/octet-stream")

			http.ServeContent(w, r, path, stat.ModTime(), f)
			log.Printf("200 %s", r.URL.Path)
		}
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("PMTiles server listening on http://localhost%s", addr)
	log.Printf("Configure frontend to fetch tiles from http://localhost%s/<file>.pmtiles", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Range")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Range, Accept-Ranges, Content-Length")
}

func parseRange(rangeHeader string, size int64) (start, end int64, err error) {
	// Parse "bytes=start-end" or "bytes=start-"
	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")
	parts := strings.Split(rangeSpec, "-")

	start, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	end = size - 1
	if len(parts) > 1 && parts[1] != "" {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	if start < 0 || end >= size || start > end {
		return 0, 0, fmt.Errorf("invalid range")
	}

	return start, end, nil
}
