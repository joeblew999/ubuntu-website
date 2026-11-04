package airspace

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
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

// pageAttempts is how many times a page is fetched before the download fails, retryDelay the wait before the next
// (doubling): ArcGIS times out or errs now and then, and one bad page in hundreds shouldn't lose the whole download.
var (
	pageAttempts = 3
	retryDelay   = 2 * time.Second
)

// featurePage is a page of a FeatureServer query: its features, or the error ArcGIS answered in their place.
type featurePage struct {
	Type     string       `json:"type"`
	Features []any        `json:"features"`
	Error    *arcgisError `json:"error,omitempty"`
}

// DownloadPaginated downloads a FeatureServer layer through its query, a page at a time, each retried (pageAttempts).
// It fails, writing nothing, on a page that still errs after its retries (HTTP, a timeout, or an error ArcGIS answers
// in place of a page, with HTTP 200), and when the pages hold fewer features than the service counts: an empty or
// partial file would be tiled into nothing, or into a map with holes (doc.go, rule 4).
func DownloadPaginated(client *http.Client, ds Dataset, outPath string) error {
	want, err := featureCount(client, ds)
	if err != nil {
		return err
	}

	collection := featurePage{
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

		page, err := fetchPage(client, ds.BaseURL+"?"+params.Encode())
		for attempt, delay := 1, retryDelay; err != nil && attempt < pageAttempts; attempt, delay = attempt+1, delay*2 {
			time.Sleep(delay)
			page, err = fetchPage(client, ds.BaseURL+"?"+params.Encode())
		}
		if err != nil {
			return fmt.Errorf("page at offset %d, after %d attempts: %w", offset, pageAttempts, err)
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

// fetchPage is one page of a FeatureServer query.
func fetchPage(client *http.Client, pageURL string) (featurePage, error) {
	var page featurePage
	resp, err := client.Get(pageURL)
	if err != nil {
		return page, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return page, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return page, fmt.Errorf("decode: %w", err)
	}
	if page.Error != nil {
		return page, fmt.Errorf("ArcGIS error %d %s", page.Error.Code, page.Error.Message)
	}
	return page, nil
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
