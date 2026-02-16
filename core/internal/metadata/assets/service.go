package assets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"

	"github.com/google/uuid"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/shared/config"
	"github.com/rs/zerolog/log"
)

type Service struct {
	store              Store
	getGameInstallPath GetInstallPathFunc
	config             config.Provider[Config]
}

type GetInstallPathFunc func(ctx context.Context, gameId uint) (string, error)

const MetadataDir = ".glacier"
const AssetDir = "assets"

func New(
	config config.Provider[Config],
	store Store,
	getPathFn GetInstallPathFunc,
) *Service {
	s := &Service{
		config:             config,
		store:              store,
		getGameInstallPath: getPathFn,
	}
	go s.DownloadUndownloaded()

	go s.orphanAssetCleaner()

	return s
}

func (s *Service) GetAsset(ctx context.Context, gameId int, assetPath string) (string, error) {
	return s.getAssetPath(ctx, uint(gameId), assetPath)
}

func (s *Service) SaveLocalAsset(ctx context.Context, asset *Asset, file io.Reader) error {
	filename := uuid.New().String()

	assetPath, err := s.getAssetPath(ctx, asset.GameID, filename)
	if err != nil {
		return err
	}

	openFile, err := os.OpenFile(assetPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fileutil.Close(openFile)

	if _, err = io.Copy(openFile, file); err != nil {
		return err
	}
	asset.LocalPath = filename // incase asset is new

	ifExists := func(oldAsset *Asset) error {
		if oldAsset.LocalPath == "" {
			return nil
		}

		join := filepath.Join(filepath.Dir(assetPath), oldAsset.LocalPath)
		oldAsset.LocalPath = filename // replace with new asset
		return os.RemoveAll(join)
	}

	return s.store.Save(ctx, asset, ifExists)
}

func (s *Service) GetAssetByType(ctx context.Context, gameId uint, assetType AssetType) (string, error) {
	id, err := s.store.GetById(ctx, gameId, assetType)
	if err != nil {
		return "", err
	}

	return s.getAssetPath(ctx, gameId, id.LocalPath)
}

func (s *Service) getAssetPath(ctx context.Context, gameId uint, filename string) (string, error) {
	dstDir, err := s.getGameInstallPath(ctx, gameId)
	if err != nil {
		return "", err
	}

	baseDir := filepath.Join(dstDir, MetadataDir, AssetDir)
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return "", err
	}

	assetPath := filepath.Join(baseDir, filename)
	return assetPath, nil
}

func (s *Service) delete(ctx context.Context, gameId uint) error {
	get, err := s.store.Get(ctx, gameId)
	if err != nil {
		return err
	}

	path, err := s.getAssetPath(ctx, get.GameID, get.LocalPath)
	if err != nil {
		return err
	}

	err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	return s.store.Delete(ctx, get)
}

func (s *Service) DownloadAssets(ctx context.Context, id uint) error {
	assets, err := s.store.ListUndownloadedByGame(ctx, id)
	if err != nil {
		return err
	}

	go s.downloadAssets(id, assets)

	return nil
}

func (s *Service) downloadAssets(gameId uint, assets []Asset) {
	cli := &http.Client{}

	// do sequentially to avoid any rate limiting shenanigans
	for _, asset := range assets {
		s.downloadSingleAsset(cli, &asset)
	}

	s.cleanupOrphanedAssets(gameId)
}

func (s *Service) orphanAssetCleaner() {
	ctx := context.Background()
	ids, err := s.store.GetGameIds(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get game ids")
		return
	}

	for _, id := range ids {
		s.cleanupOrphanedAssets(id)
	}
}

func (s *Service) cleanupOrphanedAssets(gameId uint) {
	ctx := context.Background()
	dstDir, err := s.getGameInstallPath(ctx, gameId)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to get game install path, to clean up files")
		return
	}

	gameAssets, err := s.store.GetByGame(ctx, gameId)
	if err != nil {
		log.Warn().Msg("Failed to get game assets")
		return
	}

	base := filepath.Join(dstDir, MetadataDir, AssetDir)
	dir, err := os.ReadDir(base)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to read metadata directory")
		return
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		if !slices.ContainsFunc(gameAssets, func(asset Asset) bool {
			return asset.LocalPath == file.Name()
		}) {
			err := os.RemoveAll(filepath.Join(base, file.Name()))
			if err != nil {
				log.Warn().Err(err).Msgf("Failed to remove orphaned asset %s", file.Name())
				return
			}
		}
	}
}

func (s *Service) downloadSingleAsset(cli *http.Client, asset *Asset) {
	if s.config().UseYtDlp {
		s.downloadWithYtDlp(cli, asset)
		return
	}

	s.downloadWithHttp(cli, asset)
	return
}

func (s *Service) DownloadUndownloaded() {
	ctx := context.Background()
	downloaded, err := s.store.ListUnDownloaded(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to list undownloaded assets")
		return
	}

	if len(downloaded) < 1 {
		log.Debug().Msg("No undownloaded assets found")
		return
	}

	for gameId, games := range downloaded {
		s.downloadAssets(gameId, games)
	}
}

func (s *Service) downloadWithHttp(cli *http.Client, asset *Asset) {
	resp, err := cli.Get(asset.RemoteURL)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to download asset")
		return
	}
	defer fileutil.Close(resp.Body)

	if checkHttpErr(resp, asset) {
		return
	}

	ctx := context.Background()
	err = s.SaveLocalAsset(ctx, asset, resp.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to save asset")
		return
	}
}

func checkHttpErr(resp *http.Response, asset *Asset) bool {
	if resp.StatusCode >= http.StatusBadRequest {
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Warn().Err(err).Msg("could not read response body")
		}

		log.Warn().Err(err).
			Str("url", asset.RemoteURL).
			Int("code", resp.StatusCode).
			Str("body", string(all)).
			Msg("asset download returned a non ok response code")
		return true
	}
	return false
}

func (s *Service) downloadWithYtDlp(cli *http.Client, asset *Asset) {
	ytRelay := s.config().YTRelayUrl
	encoded := url.QueryEscape(asset.RemoteURL)
	baseUrl := fmt.Sprintf("%s/download?url=%s", ytRelay, encoded)

	resp, err := cli.Get(baseUrl)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to download asset")
		return
	}
	defer fileutil.Close(resp.Body)

	if checkHttpErr(resp, asset) {
		return
	}

	err = s.SaveLocalAsset(
		context.Background(),
		asset,
		resp.Body,
	)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to save asset")
		return
	}
}

func (s *Service) Edit(ctx context.Context, a *Asset) error {
	return s.store.Update(ctx, a)
}
