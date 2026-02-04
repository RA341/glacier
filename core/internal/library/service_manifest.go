package library

import (
	"context"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type ManifestService struct {
	gameStore       Store
	folderMetaStore StoreGameManifest
}

func NewManifestService(gameStore Store, folderMetaStore StoreGameManifest) *ManifestService {
	m := &ManifestService{
		folderMetaStore: folderMetaStore,
		gameStore:       gameStore,
	}

	go func() {
		err := m.CheckManifest(context.Background())
		if err != nil {
			log.Warn().Err(err).Msg("Failed to update manifest")
		}
	}()

	return m
}

func (s *ManifestService) CheckManifest(ctx context.Context) error {
	gameIds, err := s.folderMetaStore.ListGamesWithoutManifest(ctx)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}

	for _, gid := range gameIds {
		eg.Go(func() error {
			_, err := s.GenerateManifest(ctx, gid)
			return err
		})
	}

	return eg.Wait()
}

func (s *ManifestService) FileDownload(ctx context.Context, id int, file string) (*os.File, error) {
	game, err := s.gameStore.GetById(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(game.Download.DownloadPath, file)
	return os.Open(filePath)
}

func (s *ManifestService) GetDownloadManifest(ctx context.Context, gameId int, writer io.Writer) error {
	meta, err := s.GenerateManifest(ctx, gameId)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(writer)
	return encoder.Encode(meta)
}

func (s *ManifestService) GenerateManifest(ctx context.Context, gameId int) (FolderManifest, error) {
	game, err := s.gameStore.GetById(ctx, uint(gameId))
	if err != nil {
		return FolderManifest{}, err
	}

	if game.Download.State != types.Complete {
		return FolderManifest{}, fmt.Errorf("game is not complete")
	}

	prevMeta, err := s.folderMetaStore.Get(ctx, gameId)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get previous manifest")
	}

	eg := errgroup.Group{}
	eg.SetLimit(-1)
	metadataChan := make(chan MetaResult, 5)

	err = filepath.WalkDir(game.Download.DownloadPath, func(path string, d fs.DirEntry, err error) error {
		if game.Download.DownloadPath == path || d.IsDir() {
			// process files inside dir directly with their paths
			return nil
		}

		eg.Go(func() error {
			return s.gatherMeta(metadataChan, path, &game, &prevMeta)
		})

		return nil
	})
	if err != nil {
		return FolderManifest{}, err
	}

	go func() {
		_ = eg.Wait()
		close(metadataChan)
	}()

	var finalMeta FolderManifest

	finalMeta.ID = prevMeta.ID
	finalMeta.CreatedAt = prevMeta.CreatedAt
	finalMeta.UpdatedAt = prevMeta.UpdatedAt
	finalMeta.GameID = prevMeta.GameID
	finalMeta.TotalSize = 0
	finalMeta.FileInfo = make([]FileManifest, len(prevMeta.FileInfo))

	for me := range metadataChan {
		log.Info().Any("data", me.meta).Msg("got metadata")

		finalMeta.TotalSize += me.meta.Size
		if me.Update {
			// something changed
			finalMeta.FileInfo[me.InsertIndex] = me.meta
			continue
		}

		// new metadata
		finalMeta.FileInfo = append(finalMeta.FileInfo, me.meta)
	}

	err = eg.Wait()
	if err != nil {
		return FolderManifest{}, err
	}

	err = s.folderMetaStore.Add(ctx, gameId, &finalMeta)
	if err != nil {
		return FolderManifest{}, err
	}

	if finalMeta.ID == 0 {
		return FolderManifest{}, fmt.Errorf("metadata DB id was 0, THIS SHOULD NEVER HAPPEN: %v", finalMeta)
	}

	log.Debug().Int("game", gameId).Msg("completed metadata extraction")

	return finalMeta, nil
}

type MetaResult struct {
	InsertIndex uint
	meta        FileManifest
	Update      bool
}

func (s *ManifestService) gatherMeta(
	metadataChan chan MetaResult,
	path string,
	game *Game,
	prevMeta *FolderManifest,
) error {
	relPath, err := filepath.Rel(game.Download.DownloadPath, path)
	if err != nil {
		return err
	}

	curStat, err := os.Stat(path)
	if err != nil {
		return err
	}

	var res MetaResult

	var missReason = "file not found"
	for i, prevState := range prevMeta.FileInfo {
		if prevState.RelPath == relPath {

			// update at index
			res.InsertIndex = uint(i)
			res.Update = true

			if prevState.ModTime.Equal(curStat.ModTime()) {
				// file is not modified, meta does not need to be updated
				log.Info().Str("file", relPath).Msg("using cached metadata")
				res.meta = prevState

				metadataChan <- res
				return nil
			}

			missReason = fmt.Sprintf(
				"file modified prev %s, cur: %s",
				prevState.ModTime.String(),
				curStat.ModTime().String(),
			)

		}
	}

	log.Info().
		Str("reason", missReason).
		Str("file", relPath).
		Msg("metadata cache miss")

	hash, err := GetHash(path)
	if err != nil {
		return err
	}

	res.meta = FileManifest{
		RelPath:  relPath,
		Size:     curStat.Size(),
		ModTime:  curStat.ModTime(),
		Checksum: hash,
	}

	log.Debug().Str("file", relPath).Msg("metadata done")

	metadataChan <- res

	return nil
}

func GetHash(path string) (string, error) {
	open, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fileutil.Close(open)

	_, err = open.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	h := md5.New()
	if _, err := io.Copy(h, open); err != nil {
		return "", err
	}

	checksum := hex.EncodeToString(h.Sum(nil))

	return checksum, nil
}
