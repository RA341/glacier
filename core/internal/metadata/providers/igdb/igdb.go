package igdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ra341/glacier/internal/metadata/assets"
	"github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/listutils"
	"github.com/ra341/glacier/pkg/mapsct"
)

const (
	GamesBase  = "https://api.igdb.com/v4/games"
	TwitchAuth = "https://id.twitch.tv/oauth2/token"
)

type Config struct {
	ClientId     string
	ClientSecret string
	Debug        bool
}

type Client struct {
	config Config
	http   *http.Client
}

func New(input types.ProviderConfig) (types.Provider, error) {
	var conf Config
	if err := mapsct.ParseMap(&conf, input); err != nil {
		return nil, err
	}

	transport := &igdbTransport{
		clientId:     conf.ClientId,
		clientSecret: conf.ClientSecret,
		base:         http.DefaultTransport,
	}

	return &Client{
		config: conf,
		http:   &http.Client{Transport: transport},
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var gameFields = strings.Join([]string{
	"name",
	"genres.name",
	"url",
	"summary", "storyline",
	"cover.url",
	"videos.video_id",
	"aggregated_rating",
	"rating_count",
	"first_release_date",
	"status",
	"category",
	"platforms.name",
	"themes.name",
}, ", ")

func (ig *Client) GetFullMetadata(ctx context.Context, providerGameId string) (*types.Meta, error) {
	id, err := strconv.Atoi(providerGameId)
	if err != nil {
		return nil, fmt.Errorf("invalid IGDB game id: %w", err)
	}

	games, err := ig.doQuery(
		ctx,
		GamesBase,
		fmt.Sprintf("fields %s; where id = %d;", gameFields, id),
	)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, fmt.Errorf("no game found for id %s", providerGameId)
	}

	return new(gameToMeta(&games[0])), nil
}

func (ig *Client) GetMatches(ctx context.Context, query string) ([]types.Meta, error) {
	if query == "" {
		return nil, nil
	}

	games, err := ig.doQuery(
		ctx,
		GamesBase,
		fmt.Sprintf(`search "%s"; fields %s; limit 5;`, query, gameFields),
	)
	if err != nil {
		return nil, err
	}

	return listutils.ToMap(games, func(t Game) types.Meta {
		return gameToMeta(&t)
	}), nil
}

func (ig *Client) doQuery(ctx context.Context, endpoint, query string) ([]Game, error) {
	// todo add context
	resp, err := ig.http.Post(endpoint, "text/plain", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("igdb request failed: %w", err)
	}
	defer fileutil.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("igdb returned %d: %s", resp.StatusCode, body)
	}

	var games []Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("failed to decode igdb response: %w", err)
	}
	return games, nil
}

func gameToMeta(t *Game) types.Meta {
	t.Cover.Url = strings.TrimPrefix(t.Cover.Url, "//")
	t.Cover.Url = "https://" + t.Cover.Url

	allAssets := []assets.Asset{
		assets.NewRemoteAsset(
			strings.Replace(t.Cover.Url, "t_thumb", "t_cover_big", 1),
			assets.AssetThumbnail,
		),
		assets.NewRemoteAsset(
			strings.Replace(t.Cover.Url, "t_thumb", "t_1080p", 1),
			assets.AssetBanner,
		),
		assets.NewRemoteAsset(
			strings.Replace(t.Cover.Url, "t_thumb", "t_1080p", 1),
			assets.AssetBanner,
		),
	}
	allAssets = append(allAssets,
		listutils.ToMap(t.Videos, func(v Video) assets.Asset {
			return assets.NewRemoteAsset(
				fmt.Sprintf("https://www.youtube.com/watch?v=%s", v.VideoId),
				assets.AssetTrailer,
			)
		})...,
	)

	allAssets = append(allAssets,
		listutils.ToMap(t.Screenshots, func(s Screenshot) assets.Asset {
			return assets.NewRemoteAsset(
				fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_screenshot_huge/%s.jpg", s.ImageId),
				assets.AssetGameplayImage,
			)
		})...,
	)
	allAssets = append(allAssets,
		listutils.ToMap(t.Artworks, func(a Artwork) assets.Asset {
			return assets.NewRemoteAsset(
				fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", a.ImageId),
				assets.AssetArtwork,
			)
		})...,
	)

	return types.Meta{
		ProviderType: types.ProviderIGDB,
		GameDBID:     strconv.Itoa(t.Id),
		Name:         t.Name,
		ShortDesc:    t.Storyline,
		FullDesc:     t.Summary,
		URL:          t.Url,
		Genres: listutils.ToMap(t.Genres, func(g Genre) string {
			return g.Name
		}),
		Assets: allAssets,
		Platforms: listutils.ToMap(t.Platforms, func(p Platforms) string {
			return p.Name
		}),
		RatingCount: uint(t.RatingCount),
		ReleaseDate: time.Unix(int64(t.FirstReleaseDate), 0),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// models

type TwitchToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresInSeconds int    `json:"expires_in"`
}

type Game struct {
	Id               int          `json:"id"`
	Url              string       `json:"url"`
	AggregatedRating float64      `json:"aggregated_rating,omitempty"`
	Cover            Cover        `json:"cover"`
	FirstReleaseDate int          `json:"first_release_date,omitempty"`
	Genres           []Genre      `json:"genres"`
	Name             string       `json:"name"`
	Platforms        []Platforms  `json:"platforms"`
	RatingCount      int          `json:"rating_count,omitempty"`
	Summary          string       `json:"summary,omitempty"`
	Themes           []Theme      `json:"themes"`
	Videos           []Video      `json:"videos"`
	Storyline        string       `json:"storyline,omitempty"`
	Status           int          `json:"status,omitempty"`
	Screenshots      []Screenshot `json:"screenshots"`
	Artworks         []Artwork    `json:"artworks"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Video struct {
	Id      int    `json:"id"`
	VideoId string `json:"video_id"`
}

type Cover struct {
	Id  int    `json:"id"`
	Url string `json:"url"`
}

type Screenshot struct {
	Id      int    `json:"id"`
	ImageId string `json:"image_id"`
}

type Artwork struct {
	Id      int    `json:"id"`
	ImageId string `json:"image_id"`
}

type Platforms struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Theme struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// http util

// igdbTransport is a custom RoundTripper that injects a fresh auth token on
// every request, refreshing it transparently when it has expired.
type igdbTransport struct {
	clientId     string
	clientSecret string

	mu          sync.RWMutex
	accessToken string
	expiry      time.Time

	base http.RoundTripper
}

func (t *igdbTransport) token() (string, error) {
	t.mu.RLock()
	if time.Now().Before(t.expiry) {
		tok := t.accessToken
		t.mu.RUnlock()
		return tok, nil
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()

	if time.Now().Before(t.expiry) {
		return t.accessToken, nil
	}

	tok, err := t.fetchToken()
	if err != nil {
		return "", err
	}

	t.accessToken = tok.AccessToken
	t.expiry = time.Now().Add(time.Duration(tok.ExpiresInSeconds) * time.Second)
	return t.accessToken, nil
}

func (t *igdbTransport) fetchToken() (*TwitchToken, error) {
	params := url.Values{
		"client_id":     {t.clientId},
		"client_secret": {t.clientSecret},
		"grant_type":    {"client_credentials"},
	}

	resp, err := http.PostForm(TwitchAuth, params)
	if err != nil {
		return nil, fmt.Errorf("token fetch failed: %w", err)
	}
	defer fileutil.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token fetch returned %d: %s", resp.StatusCode, body)
	}

	var tok TwitchToken
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	return &tok, nil
}

func (t *igdbTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	tok, err := t.token()
	if err != nil {
		return nil, err
	}

	// Clone the request before mutating headers (RoundTripper contract)
	r := req.Clone(req.Context())
	r.Header.Set("Client-ID", t.clientId)
	r.Header.Set("Authorization", "Bearer "+tok)
	r.Header.Set("Content-Type", "text/plain")

	return t.base.RoundTrip(r)
}
