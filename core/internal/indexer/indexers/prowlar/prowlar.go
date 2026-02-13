package prowlar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ra341/glacier/internal/indexer/types"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/mapsct"
)

type Config struct {
	Url    string
	ApiKey string
}

type Prowlar struct {
	config Config
	client *http.Client
}

func New(in map[string]any) (types.Indexer, error) {
	var cf Config

	err := mapsct.ParseMap(&cf, in)
	if err != nil {
		return nil, err
	}

	if cf.Url == "" {
		return nil, fmt.Errorf("prowlarr url cannot be empty")
	}
	if cf.ApiKey == "" {
		return nil, fmt.Errorf("prowlarr api key cannot be empty")
	}
	cf.Url = strings.TrimRight(cf.Url, "/")

	return &Prowlar{
		config: cf,
		client: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// IndexerReleaseResource is based on Prowlarr's API spec
type IndexerReleaseResource struct {
	Title       string    `json:"title"`
	DownloadUrl string    `json:"downloadUrl"`
	Size        int64     `json:"size"`
	PublishDate time.Time `json:"publishDate"`
}

func (p *Prowlar) Search(query string) ([]types.Source, error) {
	fullUrl, err := url.Parse(p.config.Url + "/api/v1/search")
	if err != nil {
		return nil, fmt.Errorf("invalid prowlarr url: %w", err)
	}

	q := fullUrl.Query()
	q.Set("query", query)
	fullUrl.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", fullUrl.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create prowlarr search request: %w", err)
	}

	req.Header.Set("X-Api-Key", p.config.ApiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prowlarr search request failed: %w", err)
	}
	defer fileutil.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("prowlarr search failed with status %s: %s", resp.Status, string(body))
	}

	var results []IndexerReleaseResource
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode prowlarr search response: %w", err)
	}

	var sources []types.Source
	for _, release := range results {
		sources = append(sources, types.Source{
			IndexerType: types.IndexerProwlarr,
			Title:       release.Title,
			DownloadUrl: release.DownloadUrl,
			FileSize:    fmt.Sprintf("%d", release.Size),
			CreatedISO:  release.PublishDate.Format(time.RFC3339),
		})
	}

	return sources, nil
}

func (p *Prowlar) Close() {
	// No background processes to stop
}
