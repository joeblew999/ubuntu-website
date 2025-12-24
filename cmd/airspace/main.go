// Command airspace manages FAA airspace data for the BVLOS demo.
//
// Usage:
//
//	airspace download             # Download all datasets
//	airspace download -dataset uas
//	airspace status               # Show data file status
//	airspace serve                # Start demo server
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Dataset configuration
type Dataset struct {
	Name     string
	Filename string
	BaseURL  string
	// For FeatureServer APIs that require pagination
	IsPaginated bool
	PageSize    int
}

var datasets = map[string]Dataset{
	"uas": {
		Name:        "UAS Facility Map",
		Filename:    "faa_uas_facility_map.geojson",
		BaseURL:     "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/FAA_UAS_FacilityMap_Data/FeatureServer/0/query",
		IsPaginated: true,
		PageSize:    2000,
	},
	"boundary": {
		Name:     "Airspace Boundary",
		Filename: "faa_airspace_boundary.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/67885972e4e940b2aa6d74024901c561/geojson?layers=0",
	},
	"sua": {
		Name:     "Special Use Airspace",
		Filename: "faa_special_use_airspace.geojson",
		BaseURL:  "https://adds-faa.opendata.arcgis.com/api/download/v1/items/dd0d1b726e504137ab3c41b21835d05b/geojson?layers=0",
	},
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...) // Remove subcommand from args

	switch cmd {
	case "download":
		runDownload()
	case "status":
		runStatus()
	case "serve":
		runServe()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: airspace <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  download    Download FAA airspace data from ArcGIS APIs")
	fmt.Println("  status      Show data file status and age")
	fmt.Println("  serve       Start demo HTTP server")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  airspace download                    # Download all datasets")
	fmt.Println("  airspace download -dataset uas       # Download only UAS Facility Map")
	fmt.Println("  airspace status                      # Show data status")
	fmt.Println("  airspace serve -port 9091            # Start server on port 9091")
}

// ============================================================================
// Download Command
// ============================================================================

func runDownload() {
	fs := flag.NewFlagSet("download", flag.ExitOnError)
	outputDir := fs.String("output", "static/airspace", "Output directory for GeoJSON files")
	datasetFlag := fs.String("dataset", "", "Specific dataset (uas, boundary, sua). Empty = all")
	timeout := fs.Duration("timeout", 5*time.Minute, "HTTP request timeout")
	fs.Parse(os.Args[1:])

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
		os.Exit(1)
	}

	// Determine which datasets to download
	toDownload := make([]Dataset, 0)
	if *datasetFlag != "" {
		ds, ok := datasets[*datasetFlag]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown dataset: %s (valid: uas, boundary, sua)\n", *datasetFlag)
			os.Exit(1)
		}
		toDownload = append(toDownload, ds)
	} else {
		// Download all in consistent order
		for _, key := range []string{"uas", "boundary", "sua"} {
			toDownload = append(toDownload, datasets[key])
		}
	}

	client := &http.Client{Timeout: *timeout}

	for _, ds := range toDownload {
		fmt.Printf("Downloading %s...\n", ds.Name)
		outPath := filepath.Join(*outputDir, ds.Filename)

		var err error
		if ds.IsPaginated {
			err = downloadPaginated(client, ds, outPath)
		} else {
			err = downloadDirect(client, ds.BaseURL, outPath)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", ds.Name, err)
			os.Exit(1)
		}

		// Report file size
		if info, err := os.Stat(outPath); err == nil {
			fmt.Printf("  ✓ %s (%.1f MB)\n", ds.Filename, float64(info.Size())/(1024*1024))
		}
	}

	fmt.Println()
	fmt.Println("Done. Run 'task r2:airspace:upload' to sync to R2.")
}

// downloadDirect downloads a file directly (no pagination)
func downloadDirect(client *http.Client, url, outPath string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

// downloadPaginated handles ArcGIS FeatureServer pagination
func downloadPaginated(client *http.Client, ds Dataset, outPath string) error {
	type FeatureCollection struct {
		Type     string        `json:"type"`
		Features []interface{} `json:"features"`
	}

	collection := FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]interface{}, 0),
	}

	offset := 0
	for {
		params := url.Values{}
		params.Set("where", "1=1")
		params.Set("outFields", "*")
		params.Set("f", "geojson")
		params.Set("resultRecordCount", fmt.Sprintf("%d", ds.PageSize))
		params.Set("resultOffset", fmt.Sprintf("%d", offset))

		queryURL := ds.BaseURL + "?" + params.Encode()

		resp, err := client.Get(queryURL)
		if err != nil {
			return fmt.Errorf("fetch page at offset %d: %w", offset, err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return fmt.Errorf("HTTP %d at offset %d", resp.StatusCode, offset)
		}

		var page FeatureCollection
		if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
			resp.Body.Close()
			return fmt.Errorf("decode page at offset %d: %w", offset, err)
		}
		resp.Body.Close()

		collection.Features = append(collection.Features, page.Features...)
		fmt.Printf("  ... fetched %d features (total: %d)\n", len(page.Features), len(collection.Features))

		if len(page.Features) < ds.PageSize {
			break
		}
		offset += ds.PageSize
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(collection)
}

// ============================================================================
// Status Command
// ============================================================================

func runStatus() {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	outputDir := fs.String("output", "static/airspace", "Data directory")
	fs.Parse(os.Args[1:])

	fmt.Println("Airspace Data Status")
	fmt.Println("====================")
	fmt.Println()

	for _, key := range []string{"uas", "boundary", "sua"} {
		ds := datasets[key]
		path := filepath.Join(*outputDir, ds.Filename)
		info, err := os.Stat(path)
		if err != nil {
			fmt.Printf("  [%s] %s: NOT FOUND\n", key, ds.Name)
			continue
		}

		age := time.Since(info.ModTime())
		fmt.Printf("  [%s] %s: %.1f MB (%s old)\n",
			key, ds.Name, float64(info.Size())/(1024*1024), formatAge(age))
	}

	fmt.Println()
	fmt.Println("Run 'airspace download' to refresh data.")
}

func formatAge(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%d min", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(d.Hours()))
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}

// ============================================================================
// Serve Command
// ============================================================================

func runServe() {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 9091, "Port to listen on")
	staticDir := fs.String("static", "./static", "Static files directory")
	fs.Parse(os.Args[1:])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Correct Content-Type for GeoJSON
		if strings.HasSuffix(r.URL.Path, ".geojson") {
			w.Header().Set("Content-Type", "application/geo+json")
		}

		// Serve from static folder
		path := filepath.Join(*staticDir, r.URL.Path)
		if r.URL.Path == "/" {
			path = filepath.Join(*staticDir, "airspace-demo/index.html")
		}

		http.ServeFile(w, r, path)
	})

	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Airspace demo server starting on http://localhost%s\n", addr)
	fmt.Printf("  Demo:  http://localhost%s/airspace-demo/\n", addr)
	fmt.Printf("  Data:  http://localhost%s/airspace/\n", addr)
	fmt.Println()
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(addr, nil))
}
