package manifest

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/ra341/glacier/internal/metadata/assets"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	//gameStore       Store
	folderMetaStore StoreGameManifest

	activeChecks syncmap.Map[int, *WorkResult]
}

func New(folderMetaStore StoreGameManifest) *Service {
	m := &Service{
		folderMetaStore: folderMetaStore,
	}

	go func() {
		err := m.CheckManifest(context.Background())
		if err != nil {
			log.Warn().Err(err).Msg("Failed to update manifest")
		}
	}()

	return m
}

func (s *Service) CheckManifest(ctx context.Context) error {
	infs, err := s.folderMetaStore.ListGamesWithoutManifest(ctx)
	if err != nil {
		return err
	}

	eg := errgroup.Group{}

	for _, in := range infs {
		eg.Go(func() error {
			_, err := s.GenerateManifest(ctx, in.ID, in.DownloadPath)
			if err != nil && !errors.Is(err, context.Canceled) {
				return err
			}
			return nil
		})
	}

	return eg.Wait()
}

// GetGameManifest generates then sends the manifest via io.writer
func (s *Service) GetGameManifest(ctx context.Context, gameId int, downloadPath string, writer io.Writer) error {
	meta, err := s.GenerateManifest(ctx, gameId, downloadPath)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
		return nil
	}

	encoder := gob.NewEncoder(writer)
	return encoder.Encode(meta)
}

type WorkResult struct {
	result *FolderManifest
	err    error
	done   chan struct{}
}

func (s *Service) GenerateManifest(ctx context.Context, gameId int, downloadPath string) (*FolderManifest, error) {
	val, ok := s.activeChecks.Load(gameId)
	if ok {
		log.Debug().Int("game", gameId).Msg("manifest already being generated, waiting...")
		// manifest is being generated wait for it complete
		<-val.done
		return val.result, val.err
	}

	log.Debug().Int("game", gameId).Msg("generating manifest")

	var workResult WorkResult
	workResult.done = make(chan struct{})
	s.activeChecks.Store(gameId, &workResult)

	workResult.result, workResult.err = s.realGenerateManifest(ctx, gameId, downloadPath)

	close(workResult.done)
	s.activeChecks.Delete(gameId)

	return workResult.result, workResult.err
}

func (s *Service) realGenerateManifest(ctx context.Context, gameId int, downloadPath string) (*FolderManifest, error) {
	prevMeta, err := s.folderMetaStore.Get(ctx, gameId)
	if err != nil {
		log.Debug().Err(err).Msg("previous manifest not found")
	}

	eg := errgroup.Group{}
	eg.SetLimit(-1)
	metadataChan := make(chan MetaResult, 5)

	err = filepath.WalkDir(downloadPath, func(path string, d fs.DirEntry, err error) error {
		if downloadPath == path || d.IsDir() {
			// process files inside dir directly with their paths
			return nil
		}

		rel, err := filepath.Rel(downloadPath, path)
		if err != nil {
			return err
		}

		if strings.HasPrefix(rel, assets.MetadataDir) {
			return filepath.SkipDir
		}

		eg.Go(func() error {
			return s.gatherMeta(ctx, metadataChan, path, downloadPath, &prevMeta)
		})

		return nil
	})
	if err != nil {
		return nil, err
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
		return nil, err
	}

	cleaned := make([]FileManifest, 0)
	for _, f := range finalMeta.FileInfo {
		// remove any empty slots in case files were added/deleted
		if f.RelPath != "" {
			cleaned = append(cleaned, f)
		}
	}
	finalMeta.FileInfo = cleaned

	err = s.folderMetaStore.Add(ctx, gameId, &finalMeta)
	if err != nil {
		return nil, err
	}

	if finalMeta.ID == 0 {
		return nil, fmt.Errorf("metadata DB id was 0, THIS SHOULD NEVER HAPPEN: %v", finalMeta)
	}

	log.Debug().Int("game", gameId).Msg("completed metadata extraction")

	return &finalMeta, nil
}

type MetaResult struct {
	InsertIndex uint
	meta        FileManifest
	Update      bool
}

func (s *Service) gatherMeta(ctx context.Context, metadataChan chan MetaResult, path string, downloadPath string, prevMeta *FolderManifest) error {
	relPath, err := filepath.Rel(downloadPath, path)
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

	hash, err := GetHash(ctx, path)
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

// GetHash now accepts a context.Context
func GetHash(ctx context.Context, path string) (string, error) {
	// 1. Check if context is already canceled before starting
	if err := ctx.Err(); err != nil {
		return "0", err
	}

	f, err := os.Open(path)
	if err != nil {
		return "0", err
	}
	defer fileutil.Close(f)

	h := xxhash.New()
	buf := make([]byte, 1024*1024)

	reader := &contextReader{ctx: ctx, r: f}

	if _, err := io.CopyBuffer(h, reader, buf); err != nil {
		return "0", err
	}

	sum64 := h.Sum64()
	return strconv.FormatUint(sum64, 10), nil
}

// contextReader intercepts Read calls to check for context cancellation
type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (cr *contextReader) Read(p []byte) (n int, err error) {
	// Check context before performing the read
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	default:
		return cr.r.Read(p)
	}
}
