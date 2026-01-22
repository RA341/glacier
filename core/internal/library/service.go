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

type Downloader interface {
	Add(ctx context.Context, gameId *Game) (err error)
	TriggerTracker()
}

type Service struct {
	config     ConfigLoader
	downloader Downloader

	store           Store
	folderMetaStore StoreFolderMetadata
}

type ConfigLoader func() *Config

func New(
	store Store,
	folderMetaStore StoreFolderMetadata,
	downloader Downloader,
	config ConfigLoader,
) *Service {
	return &Service{
		downloader: downloader,
		config:     config,

		store:           store,
		folderMetaStore: folderMetaStore,
	}
}

func (s *Service) List(ctx context.Context, query string, offset, limit uint) ([]Game, error) {
	return s.store.List(ctx, limit, offset)
}

func (s *Service) ListDownloading(ctx context.Context, state string) ([]Game, error) {
	dState, err := types.DownloadStateString(state)
	if err != nil {
		return nil, err
	}

	return s.store.ListDownloadState(ctx, dState)
}

func (s *Service) Add(ctx context.Context, game *Game) error {
	game.Download.State = types.Queued
	game.Download.DownloadPath = filepath.Join(
		s.config().GameDir,
		filepath.Clean(game.Meta.Name),
	)

	err := s.store.Add(ctx, game)
	if err != nil {
		return err
	}

	return s.downloader.Add(ctx, game)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.store.DeleteGame(ctx, id)
}

func (s *Service) Get(ctx context.Context, id uint) (Game, error) {
	return s.store.GetById(ctx, id)
}

func (s *Service) GetDownloadMetadata(ctx context.Context, gameId int, writer io.Writer) error {
	game, err := s.store.GetById(ctx, uint(gameId))
	if err != nil {
		return err
	}

	if game.Download.State != types.Complete {
		return fmt.Errorf("game is not complete")
	}

	var meta FolderMetadata

	if folderMeta, err := s.folderMetaStore.Get(ctx, gameId); err == nil {
		// meta previously exists
		meta = folderMeta
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
			return s.gatherMeta(metadataChan, path, &game, &meta)
		})

		return nil
	})
	if err != nil {
		return err
	}

	go func() {
		_ = eg.Wait()
		close(metadataChan)
	}()

	for me := range metadataChan {
		log.Info().Any("data", me.meta).Msg("got metadata")

		meta.TotalSize += me.meta.Size
		if me.Update {
			// something changed
			meta.FileInfo[me.InsertIndex] = me.meta
			continue
		}

		// new metadata
		meta.FileInfo = append(meta.FileInfo, me.meta)
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	err = s.folderMetaStore.Add(ctx, gameId, &meta)
	if err != nil {
		return err
	}

	if meta.ID == 0 {
		return fmt.Errorf("metadata DB id was 0, THIS SHOULD NEVER HAPPEN: %v", meta)
	}

	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(meta)
	if err != nil {
		return err
	}

	return nil
}

type MetaResult struct {
	InsertIndex uint
	meta        FileMetadata
	Update      bool
}

func (s *Service) gatherMeta(
	metadataChan chan MetaResult, path string,
	game *Game, prevMeta *FolderMetadata,
) error {
	relPath, err := filepath.Rel(game.Download.DownloadPath, path)
	if err != nil {
		return err
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	var res MetaResult

	for i, m := range prevMeta.FileInfo {
		if m.RelPath == relPath {
			if m.ModTime == stat.ModTime() {
				// file is not modified, meta does not need to be updated
				return nil
			}
			// update at index
			res.InsertIndex = uint(i)
			res.Update = true
		}
	}

	hash, err := GetHash(path)
	if err != nil {
		return err
	}

	res.meta = FileMetadata{
		RelPath:  relPath,
		Size:     stat.Size(),
		ModTime:  stat.ModTime(),
		Checksum: hash,
	}

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

func (s *Service) FileDownload(ctx context.Context, id int, file string) (*os.File, error) {
	game, err := s.store.GetById(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(game.Download.DownloadPath, file)
	return os.Open(filePath)
}
