package airspace

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// UploadToR2 uploads PMTiles and manifests to Cloudflare R2.
// Reads from the manifest to know what files to upload.
func UploadToR2() error {
	manifestPath := filepath.Join(DirData, FileUSAManifest)

	// Check manifest exists
	if _, err := os.Stat(manifestPath); err != nil {
		return fmt.Errorf("manifest not found: %s (run 'airspace manifest' first)", manifestPath)
	}

	// Check wrangler is available
	if _, err := exec.LookPath("wrangler"); err != nil {
		return fmt.Errorf("wrangler not found in PATH (npm install -g wrangler)")
	}

	fmt.Printf("Uploading airspace data to r2://%s/airspace/\n\n", R2Bucket)

	// Load manifest
	manifest, err := loadUSAManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	// Upload PMTiles from manifest
	fmt.Println("=== PMTiles (from manifest) ===")
	for _, layer := range manifest.Layers {
		filePath := filepath.Join(DirPMTiles, layer.File)
		if _, err := os.Stat(filePath); err != nil {
			fmt.Printf("  SKIP: %s (file not found)\n", layer.File)
			continue
		}

		info, _ := os.Stat(filePath)
		sizeMB := float64(info.Size()) / (1024 * 1024)
		fmt.Printf("  Uploading: %s (%.1f MB)\n", layer.File, sizeMB)

		r2Key := fmt.Sprintf("%s/airspace/tiles/%s", R2Bucket, layer.File)
		if err := wranglerPut(r2Key, filePath); err != nil {
			return fmt.Errorf("uploading %s: %w", layer.File, err)
		}
	}

	// Upload manifests
	fmt.Println("\n=== Manifests ===")
	manifestFiles := []string{FileManifest, FileUSAManifest}
	for _, name := range manifestFiles {
		filePath := filepath.Join(DirData, name)
		if _, err := os.Stat(filePath); err != nil {
			continue
		}
		fmt.Printf("  Uploading: %s\n", name)
		r2Key := fmt.Sprintf("%s/airspace/%s", R2Bucket, name)
		if err := wranglerPut(r2Key, filePath); err != nil {
			return fmt.Errorf("uploading %s: %w", name, err)
		}
	}

	fmt.Printf("\nDone. Files available at: %s/airspace/\n", R2PublicURL)

	// Print endpoints
	fmt.Println("\nEndpoints:")
	for _, layer := range manifest.Layers {
		fmt.Printf("  %s/airspace/tiles/%s\n", R2PublicURL, layer.File)
	}

	return nil
}

// TestR2Endpoints tests that all R2 endpoints are accessible.
func TestR2Endpoints() error {
	manifestPath := filepath.Join(DirData, FileUSAManifest)

	manifest, err := loadUSAManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	fmt.Println("Testing R2 endpoints...")
	fmt.Println()

	client := &http.Client{Timeout: 10 * time.Second}
	allOK := true

	for _, layer := range manifest.Layers {
		url := fmt.Sprintf("%s/airspace/tiles/%s", R2PublicURL, layer.File)
		status, size := testURL(client, url)

		if status == 200 {
			fmt.Printf("  ✓ %s (%d, %s)\n", layer.Name, status, size)
		} else {
			fmt.Printf("  ✗ %s (%d)\n", layer.Name, status)
			allOK = false
		}
	}

	fmt.Println()
	if allOK {
		fmt.Println("All endpoints OK.")
		return nil
	}
	return fmt.Errorf("some endpoints failed")
}

func wranglerPut(r2Key, filePath string) error {
	cmd := exec.Command("wrangler", "r2", "object", "put", r2Key, "--file", filePath, "--remote")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func testURL(client *http.Client, url string) (int, string) {
	req, _ := http.NewRequest("HEAD", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return 0, "error"
	}
	defer resp.Body.Close()

	size := resp.Header.Get("Content-Length")
	if size == "" {
		size = "?"
	} else {
		// Convert to MB if large
		var bytes int64
		fmt.Sscanf(size, "%d", &bytes)
		if bytes > 1024*1024 {
			size = fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
		} else {
			size = fmt.Sprintf("%d bytes", bytes)
		}
	}

	return resp.StatusCode, size
}

func loadUSAManifest(path string) (*ManifestUSA, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest ManifestUSA
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}
