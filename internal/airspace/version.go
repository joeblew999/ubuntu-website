package airspace

// A dataset's version and what's published: rules 1 and 2 of the sync rules (doc.go).

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PublishedURL is the manifest in R2 that records what's published (rule 2).
var PublishedURL = R2PublicURL + "/airspace/" + FileUSAManifest

// SourceVersion is the version of a dataset's data at the FAA (rule 1): its FeatureServer layer's
// editingInfo.dataLastEditDate, as an RFC 3339 time.
func SourceVersion(client *http.Client, ds Dataset) (string, error) {
	resp, err := client.Get(ds.LayerURL + "?f=json")
	if err != nil {
		return "", fmt.Errorf("version: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("version: HTTP %d from %s", resp.StatusCode, ds.LayerURL)
	}
	var layer struct {
		EditingInfo struct {
			DataLastEditDate int64 `json:"dataLastEditDate"`
		} `json:"editingInfo"`
		Error *arcgisError `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&layer); err != nil {
		return "", fmt.Errorf("version: %w", err)
	}
	if layer.Error != nil {
		return "", fmt.Errorf("version: ArcGIS error %d %s", layer.Error.Code, layer.Error.Message)
	}
	if layer.EditingInfo.DataLastEditDate == 0 {
		return "", fmt.Errorf("version: %s reports no dataLastEditDate", ds.LayerURL)
	}
	return time.UnixMilli(layer.EditingInfo.DataLastEditDate).UTC().Format(time.RFC3339Nano), nil
}

// FetchPublished is the published manifest (rule 2), or nil when R2 has none yet (every dataset then counts as
// changed).
func FetchPublished(client *http.Client, url string) (*ManifestUSA, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("published manifest: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("published manifest: HTTP %d from %s", resp.StatusCode, url)
	}
	var m ManifestUSA
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("published manifest %s: %w", url, err)
	}
	return &m, nil
}

// publishedVersion is the version a dataset's tiles in R2 were built from, or "" if none is recorded.
func (m *ManifestUSA) publishedVersion(key string) string {
	l, _ := m.layer(key)
	return l.SourceVersion
}
