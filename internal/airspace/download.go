package airspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// DownloadDirect downloads a file directly (no pagination).
func DownloadDirect(client *http.Client, downloadURL, outPath string) error {
	resp, err := client.Get(downloadURL)
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

// DownloadPaginated handles ArcGIS FeatureServer pagination.
func DownloadPaginated(client *http.Client, ds Dataset, outPath string) error {
	type FeatureCollection struct {
		Type     string `json:"type"`
		Features []any  `json:"features"`
	}

	collection := FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]any, 0),
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

// Download downloads all specified datasets.
func Download(client *http.Client, outputDir string, datasetKeys []string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	for _, key := range datasetKeys {
		ds, ok := Datasets[key]
		if !ok {
			return fmt.Errorf("unknown dataset: %s", key)
		}

		outPath := outputDir + "/" + ds.GeoJSON

		var err error
		if ds.IsPaginated {
			err = DownloadPaginated(client, ds, outPath)
		} else {
			err = DownloadDirect(client, ds.BaseURL, outPath)
		}

		if err != nil {
			return fmt.Errorf("downloading %s: %w", ds.Name, err)
		}
	}

	return nil
}
