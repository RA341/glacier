package download

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/ra341/glacier/internal/library/manifest"
	"github.com/ra341/glacier/pkg/fileutil"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Download struct {
	conf       *Config
	editStatus EditStatus
	onDone     OnDone
	cacheStore CacheStore

	ctx      context.Context
	cancel   context.CancelFunc
	done     chan struct{}
	isPaused *atomic.Bool

	gameId         int
	downloadFolder string

	metadataUrlBase string
	downloadUrlBase string
}

const MetadataFolder = ".frost.cache"

type EditStatus func(ctx context.Context, id int, down *LocalDownload) error
type OnDone func(id int)

func NewDownload(
	config *Config,
	editStatus EditStatus,
	onDone OnDone,
	gameId int,
	downloadFolder string,
	force bool,
) (*Download, error) {
	metaPath := filepath.Join(downloadFolder, MetadataFolder)

	if _, err := os.Stat(metaPath); err == nil {
		log.Info().Str("path", metaPath).Msg("Found previous metadata folder")
		if force {
			log.Info().Msg("force downloading, removing cache folder...")

			err := os.RemoveAll(downloadFolder)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := NewCacheStoreBadger(metaPath)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	a := atomic.Bool{}
	a.Store(false)

	d := &Download{
		conf:       config,
		onDone:     onDone,
		editStatus: editStatus,
		cacheStore: db,

		ctx:      ctx,
		cancel:   cancel,
		done:     make(chan struct{}),
		isPaused: &a,

		gameId:          gameId,
		downloadFolder:  downloadFolder,
		downloadUrlBase: fmt.Sprintf("%s/load/%d", config.UrlBase(), gameId),
		metadataUrlBase: fmt.Sprintf("%s/meta/%d", config.UrlBase(), gameId),
	}

	return d, nil
}

func (d *Download) Start() {
	defer func() {
		close(d.done)
		if !d.isPaused.Load() {
			// remove only if it is not paused from download tracker
			d.Close()
		}
	}()

	warnIfErr(d.editStatus(d.ctx, d.gameId, &LocalDownload{
		Status:        StatusMetadata,
		StatusMessage: "downloading manifest",
	}))

	var meta manifest.FolderManifest
	err := d.downloadMetadata(&meta)
	if err != nil {
		warnIfErr(d.editStatus(d.ctx, d.gameId, &LocalDownload{
			Status:        StatusError,
			StatusMessage: "could not download manifest",
		}))
		return
	}

	warnIfErr(d.editStatus(d.ctx, d.gameId, &LocalDownload{
		Status:        StatusDownloading,
		StatusMessage: "downloading files",
	}))

	eg := errgroup.Group{}
	eg.SetLimit(d.conf.MaxConcurrentFileChunks)

	for _, fi := range meta.FileInfo {
		eg.Go(func() error {
			select {
			case <-d.ctx.Done():
				return d.ctx.Err()
			default:
				err := d.setupFile(&fi)
				if err != nil {
					return fmt.Errorf("could not setup file metadata: %w", err)
				}

				err = d.downloadFile(&fi)
				if err != nil {
					return fmt.Errorf("could not download file: %w", err)
				}

				warnIfErr(d.editStatus(d.ctx, d.gameId, &LocalDownload{
					StatusMessage: "downloaded " + fi.RelPath,
				}))

				return nil
			}
		})
	}

	err = eg.Wait()
	if err != nil {
		if d.ctx.Err() != nil {
			log.Warn().Msg("download canceled by user")
			return
		}

		log.Error().Err(err).Msg("error downloading")
		warnIfErr(d.editStatus(d.ctx, d.gameId, &LocalDownload{
			Status:        StatusError,
			StatusMessage: "error downloading: " + err.Error(),
		}))
		return
	}

	log.Info().
		Int("game", d.gameId).
		Msg("download finished")

	warnIfErr(d.editStatus(d.ctx, d.gameId, &LocalDownload{
		Status:        StatusComplete,
		StatusMessage: "Download Complete",
		Done:          time.Now(),
	}))
}

func (d *Download) Close() {
	d.onDone(d.gameId)
	fileutil.Close(d.cacheStore)
}

func (d *Download) Cancel() {
	d.cancel()
	<-d.done // wait for download to stop
}

func (d *Download) Pause() {
	if d.isPaused.CompareAndSwap(false, true) {
		d.Cancel()

		warnIfErr(d.editStatus(context.Background(), d.gameId, &LocalDownload{
			Status: StatusPaused,
		}))

		log.Info().Msg("Download paused")
	}
}

func (d *Download) Resume() {
	if d.isPaused.CompareAndSwap(true, false) {
		ctx, cancel := context.WithCancel(context.Background())
		doneCh := make(chan struct{})

		d.ctx = ctx
		d.cancel = cancel
		d.done = doneCh

		warnIfErr(d.editStatus(context.Background(), d.gameId, &LocalDownload{
			Status: StatusDownloading,
		}))
		log.Info().Msg("Download resumed")
	}
}

func (d *Download) Progress() (complete []FileProgress, total error) {
	return d.cacheStore.Progress()
}

func warnIfErr(err error) {
	if err != nil {
		log.Warn().Err(err).Msg("error occurred while updating db")
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// metadata step

func (d *Download) downloadMetadata(meta *manifest.FolderManifest) error {
	req, err := http.NewRequestWithContext(d.ctx, "GET", d.metadataUrlBase, nil)
	if err != nil {
		return err
	}

	resp, err := d.conf.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer fileutil.Close(resp.Body)

	err = checkHttpErr(resp)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(resp.Body)
	return decoder.Decode(meta)
}

func (d *Download) setupFile(fm *manifest.FileManifest) error {
	if fm.RelPath == "" {
		log.Warn().Msg("relative path is empty THIS SHOULD NEVER HAPPEN")
		return nil
	}

	started := time.Now()

	fullPath := filepath.Join(d.downloadFolder, fm.RelPath)

	_, found, err := d.cacheStore.Get(fullPath)
	if err != nil {
		return err
	}
	if found {
		stat, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		isModified := fm.ModTime.After(stat.ModTime())
		if !isModified {
			// file is unmodified from the server
			// continue downloading
			return nil
		}
	}

	err = os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		return err
	}

	// file is either does not exist or is modified and needs to redownload the chunks
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fileutil.Close(file)

	// allocate file
	err = file.Truncate(fm.Size)
	if err != nil {
		return err
	}

	var chunkList []Chunk
	totalSize := fm.Size
	for start := int64(0); start < totalSize; start += d.conf.GetChunkSize() {
		end := start + d.conf.GetChunkSize() - 1
		// if it's the last chunk,
		// make sure not to overshoot the file size
		if end >= totalSize {
			end = totalSize - 1
		}

		chunk := Chunk{
			Start: start,
			End:   end,
			State: ChunkQueued,
		}

		chunkList = append(chunkList, chunk)
	}

	err = d.cacheStore.Add(fullPath, chunkList)
	if err != nil {
		return err
	}

	elapsed := time.Now().Sub(started)
	log.Info().Str("elapsed", elapsed.String()).
		Any("meta", chunkList).
		Str("file", filepath.Base(fullPath)).
		Msg("completed download setup")

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// download step

func (d *Download) downloadFile(fm *manifest.FileManifest) error {
	//log.Info().
	//	Str("file", met.RelPath).
	//	Str("size", humanize.Bytes(uint64(met.Size))).
	//	Msg("starting download")

	started := time.Now()

	fullPath := filepath.Join(d.downloadFolder, fm.RelPath)

	chunks, found, err := d.cacheStore.Get(fullPath)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("file %s not found in cache, THIS SHOULD NEVER HAPPEN", fm.RelPath)
	}

	file, err := os.OpenFile(fullPath, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer fileutil.Close(file)

	escaped := url.PathEscape(fm.RelPath)
	fileUrl := fmt.Sprintf("%s/%s", d.downloadUrlBase, escaped)

	eg := errgroup.Group{}

	eg.SetLimit(d.conf.MaxConcurrentFileChunks)

	for i, chunk := range chunks {
		eg.Go(func() error {
			select {
			case <-d.ctx.Done():
				return d.ctx.Err()
			default:
				if chunk.State == ChunkComplete {
					log.Debug().Str("file", filepath.Base(fullPath)).
						Any("chunk", chunks).
						Msg("chunk complete")
					return nil
				}

				errInner := d.downloadWithRange(fileUrl, &chunk, file, fm.ModTime)
				if errInner != nil {
					if errors.Is(errInner, context.Canceled) {
						//log.Debug().Msg("download cancelled")
						return nil
					}
					log.Error().Err(errInner).
						Int64("start", chunk.Start).Int64("end", chunk.End).
						Msg("could not download chunk")
					chunk.State = ChunkError
				} else {
					chunk.State = ChunkComplete
				}

				err := d.cacheStore.Update(fullPath, i, &chunk)
				if err != nil {
					log.Warn().Err(err).Msg("could not update chunk to cache")
				}
				return nil
			}
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	hash, err := manifest.GetHash(d.ctx, fullPath)
	if err != nil {
		return err
	}

	stat, err := file.Stat()
	if err != nil {
		log.Warn().Err(err).Msg("could not stat file, size info unavailable")
	}

	if fm.Checksum != hash {
		return fmt.Errorf(
			"checksum mismatch, expected: %s != got: %s\nExpected Size: %s, got size: %s",
			fm.Checksum,
			hash,
			humanize.Bytes(uint64(fm.Size)),
			humanize.Bytes(uint64(stat.Size())),
		)
	}

	elapsed := time.Now().Sub(started)
	log.Info().Str("elapsed", elapsed.String()).
		Str("size", humanize.Bytes(uint64(stat.Size()))).
		Str("file", filepath.Base(fullPath)).
		Msg("download complete")

	return nil
}

func (d *Download) downloadWithRange(url string, chunk *Chunk, writer io.WriterAt, modTime time.Time) error {
	req, err := http.NewRequestWithContext(d.ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", chunk.Start, chunk.End))
	req.Header.Set("If-Range", modTime.UTC().Format(http.TimeFormat))

	resp, err := d.conf.httpCli.Do(req)
	if err != nil {
		return err
	}
	// ensure body is closed to return connection to the pool
	defer fileutil.Close(resp.Body)

	err = checkHttpErr(resp)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server did not support range or file changed: status %s", resp.Status)
	}

	_, err = io.Copy(NewOffsetWriter(writer, chunk.Start), resp.Body)
	return err
}

func checkHttpErr(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not load body to get error message: %w", err)
	}
	return fmt.Errorf("error downloading: %d: %s", resp.StatusCode, string(all))
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// utils

// wrapper to make WriteAt behave like a standard Writer for io.Copy
type offsetWriter struct {
	w      io.WriterAt
	offset int64
}

func (ow *offsetWriter) Write(p []byte) (n int, err error) {
	n, err = ow.w.WriteAt(p, ow.offset)
	ow.offset += int64(n)
	return
}

func NewOffsetWriter(w io.WriterAt, offset int64) io.Writer {
	return &offsetWriter{w, offset}
}
