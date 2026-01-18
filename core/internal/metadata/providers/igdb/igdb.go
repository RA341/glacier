package igdb

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/pkg/listutils"
	"github.com/ra341/glacier/pkg/mapsct"
	"resty.dev/v3"
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

	rw          sync.RWMutex
	accessToken string
	expiry      time.Time
}

func New(input types.ProviderConfig) (types.Provider, error) {
	var conf Config
	err := mapsct.ParseMap(&conf, input)
	if err != nil {
		return nil, err
	}

	return &Client{
		config: conf,
		expiry: time.Now(),
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// auth stuff

type TwitchToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresInSeconds int    `json:"expires_in"`
}

func (ig *Client) getAccessToken() (string, error) {
	ig.rw.RLock()
	now := time.Now()
	if !now.After(ig.expiry) {
		return ig.accessToken, nil
	}
	ig.rw.RUnlock()

	ig.rw.Lock()
	defer ig.rw.Unlock()

	token := &TwitchToken{}
	err := ig.fetchNewToken(token)
	if err != nil {
		return "", err
	}

	now = time.Now()
	ig.accessToken = token.AccessToken
	ig.expiry = now.Add(time.Second * time.Duration(token.ExpiresInSeconds))

	return ig.accessToken, nil
}

func (ig *Client) fetchNewToken(token *TwitchToken) error {
	post, err := resty.New().NewRequest().
		SetQueryParams(map[string]string{
			"client_id":     ig.config.ClientId,
			"client_secret": ig.config.ClientSecret,
			"grant_type":    "client_credentials",
		}).
		SetDebug(ig.config.Debug).
		SetResult(token).
		Post(TwitchAuth)
	if err != nil {
		return err
	}

	if post.IsError() {
		return fmt.Errorf("%v", post.String())
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// core interface definitions

type Game struct {
	Id               int         `json:"id"`
	Url              string      `json:"url"`
	AggregatedRating float64     `json:"aggregated_rating,omitempty"`
	Cover            Cover       `json:"cover"`
	FirstReleaseDate int         `json:"first_release_date,omitempty"`
	Genres           []Genre     `json:"genres"`
	Name             string      `json:"name"`
	Platforms        []Platforms `json:"platforms"`
	RatingCount      int         `json:"rating_count,omitempty"`
	Summary          string      `json:"summary,omitempty"`
	Themes           []Theme     `json:"themes"`
	Videos           []Video     `json:"videos"`
	Storyline        string      `json:"storyline,omitempty"`
	Status           int         `json:"status,omitempty"`
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

type Platforms struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Theme struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (ig *Client) GetFullMetadata(id string) (*types.Meta, error) {
	//TODO implement me
	panic("implement me")
}

func (ig *Client) GetMatches(query string) ([]types.Meta, error) {
	games, err := ig.searchGames(query)
	if err != nil {
		return nil, err
	}

	return listutils.ToMap(games, func(t Game) types.Meta {
		return types.Meta{
			ProviderType: types.ProviderIGDB,
			GameDBID:     strconv.Itoa(t.Id),
			Name:         t.Name,
			Summary:      t.Summary,
			Description:  t.Storyline,
			URL:          t.Url,
			Genres: listutils.ToMap(t.Genres, func(t Genre) string {
				return t.Name
			}),
			ThumbnailURL: t.Cover.Url,
			Videos: listutils.ToMap(t.Videos, func(t Video) string {
				return t.VideoId
			}),
			Platforms: listutils.ToMap(t.Platforms, func(t Platforms) string {
				return t.Name
			}),
			RatingCount: uint(t.RatingCount),
			ReleaseDate: time.Unix(int64(t.FirstReleaseDate), 0),
			// todo
			//Rating:        (t.AggregatedRating),
			//ReleaseStatus: t.Status,
			//Category:      t.,
		}
	}), nil
}

func (ig *Client) searchGames(query string) ([]Game, error) {
	if query == "" {
		return nil, nil
	}

	token, err := ig.getAccessToken()
	if err != nil {
		return nil, err
	}

	var gameResults []Game

	fields := strings.Join([]string{
		"name",
		"genres",
		"url",
		"genres.name",
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

	req, err := resty.New().
		SetHeaders(map[string]string{
			"Client-ID":     ig.config.ClientId,
			"Authorization": "Bearer " + token,
		}).R().
		SetBody([]byte(
			fmt.Sprintf(`search "%s"; fields %s; limit 5;`, query, fields),
		)).
		SetResult(&gameResults).
		SetDebug(ig.config.Debug).
		Post(GamesBase)
	if err != nil {
		return nil, err
	}

	if req.IsError() {
		return nil, fmt.Errorf("%v", req.String())
	}

	return gameResults, nil
}
