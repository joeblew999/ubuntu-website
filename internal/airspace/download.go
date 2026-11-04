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

// arcgisError is the error ArcGIS answers in place of features, with HTTP 200 (throttling, a bad query, the service
// down): {"error": {"code": 499, "message": "..."}}.
type arcgisError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DownloadPaginated handles ArcGIS FeatureServer pagination. It fails, writing nothing, on an error ArcGIS answers in
// place of a page, and when the pages hold fewer features than the service counts: an empty or partial file would be
// tiled into nothing, or into a map with holes.
func DownloadPaginated(client *http.Client, ds Dataset, outPath string) error {
	type FeatureCollection struct {
		Type     string       `json:"type"`
		Features []any        `json:"features"`
		Error    *arcgisError `json:"error,omitempty"`
	}

	want, err := featureCount(client, ds)
	if err != nil {
		return err
	}

	collection := FeatureCollection{
		Type:     "FeatureCollection",
		Features: make([]any, 0, want),
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
		if page.Error != nil {
			return fmt.Errorf("ArcGIS error at offset %d: %d %s", offset, page.Error.Code, page.Error.Message)
		}

		collection.Features = append(collection.Features, page.Features...)

		if len(page.Features) < ds.PageSize {
			break
		}
		offset += ds.PageSize
	}
	if len(collection.Features) < want {
		return fmt.Errorf("got %d of the %d features the service counts", len(collection.Features), want)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	return encoder.Encode(collection)
}

// featureCount is how many features the service holds (returnCountOnly), the number a full download must reach.
func featureCount(client *http.Client, ds Dataset) (int, error) {
	params := url.Values{"where": {"1=1"}, "returnCountOnly": {"true"}, "f": {"json"}}
	resp, err := client.Get(ds.BaseURL + "?" + params.Encode())
	if err != nil {
		return 0, fmt.Errorf("count features: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count features: HTTP %d", resp.StatusCode)
	}
	var body struct {
		Count *int         `json:"count"`
		Error *arcgisError `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("count features: %w", err)
	}
	if body.Error != nil {
		return 0, fmt.Errorf("count features: ArcGIS error %d %s", body.Error.Code, body.Error.Message)
	}
	if body.Count == nil || *body.Count == 0 {
		return 0, fmt.Errorf("count features: the service reports none")
	}
	return *body.Count, nil
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
