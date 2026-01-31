package download

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/pkg/fileutil"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type Config interface {
	getMaxConcurrentFiles() int
	getChunkSize() int64
	getMaxConcurrentFileChunks() int
	getHttpClient() *http.Client
}

type ProgressUpdater interface {
	EditStatus(ctx context.Context, id int, down *Info) error
}

type OnDone func(id int)

type Download struct {
	ctx    context.Context
	cancel context.CancelFunc

	conf   Config
	gameId int

	metadataUrlBase string
	downloadUrlBase string

	downloadFolder string
	OnDone         OnDone
	cacheStore     CacheStore
	progress       ProgressUpdater
}

const MetadataFolder = ".frost.cache"

func NewDownload(
	config Config,
	OnDone OnDone,
	progress ProgressUpdater,
	baseUrl, downloadFolder string,
	gameId int,
) (*Download, error) {
	metaPath := filepath.Join(downloadFolder, MetadataFolder)
	db, err := NewCacheStoreBadger(metaPath)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := &Download{
		ctx:    ctx,
		cancel: cancel,

		OnDone: OnDone,

		downloadUrlBase: fmt.Sprintf("%s/load/%d", baseUrl, gameId),
		metadataUrlBase: fmt.Sprintf("%s/meta/%d", baseUrl, gameId),
		gameId:          gameId,
		downloadFolder:  downloadFolder,

		conf:     config,
		progress: progress,

		cacheStore: db,
	}
	go d.Start()

	return d, nil
}

func (d *Download) Start() {
	defer fileutil.Close(d.cacheStore)

	warnIfErr(d.progress.EditStatus(d.ctx, d.gameId, &Info{
		Status:        StatusMetadata,
		StatusMessage: "Downloading Metadata",
	}))

	var meta library.FolderManifest
	err := d.downloadMetadata(&meta)
	if err != nil {
		warnIfErr(d.progress.EditStatus(d.ctx, d.gameId, &Info{
			Status:        StatusError,
			StatusMessage: "could not download metadata",
		}))
		return
	}

	warnIfErr(d.progress.EditStatus(d.ctx, d.gameId, &Info{
		Status:        StatusDownloading,
		StatusMessage: "starting file download",
	}))

	eg := errgroup.Group{}
	eg.SetLimit(d.conf.getMaxConcurrentFileChunks())

	for _, fi := range meta.FileInfo {
		eg.Go(func() error {
			err := d.setupFile(&fi)
			if err != nil {
				return fmt.Errorf("could not setup file metadata: %w", err)
			}

			err = d.downloadFile(&fi)
			if err != nil {
				return fmt.Errorf("could not download file: %w", err)
			}

			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("error downloading")
		warnIfErr(d.progress.EditStatus(d.ctx, d.gameId, &Info{
			Status:        StatusError,
			StatusMessage: "error downloading: " + err.Error(),
		}))
		return
	}

	log.Info().
		Int("game", d.gameId).
		Msg("download finished")

	// remove from download tracker
	d.OnDone(d.gameId)

	warnIfErr(d.progress.EditStatus(d.ctx, d.gameId, &Info{
		Status:        StatusComplete,
		StatusMessage: "Download Complete",
		Done:          time.Now(),
	}))
}

func (d *Download) Progress() (complete []FileProgress, total error) {
	return d.cacheStore.Progress()
}

func warnIfErr(err error) {
	if err != nil {
		log.Warn().Err(err).Msg("error occurred while updating db")
	}
}

func (d *Download) Close() {
	fileutil.Close(d.cacheStore)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// metadata step

func (d *Download) downloadMetadata(meta *library.FolderManifest) error {
	resp, err := d.conf.getHttpClient().Get(d.metadataUrlBase)
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

func (d *Download) setupFile(fm *library.FileManifest) error {
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
	for start := int64(0); start < totalSize; start += d.conf.getChunkSize() {
		end := start + d.conf.getChunkSize() - 1
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

func (d *Download) downloadFile(fm *library.FileManifest) error {
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
	eg.SetLimit(d.conf.getMaxConcurrentFileChunks())

	for i, chunk := range chunks {
		eg.Go(func() error {
			if chunk.State == ChunkComplete {
				log.Debug().Str("file", filepath.Base(fullPath)).
					Any("chunk", chunks).
					Msg("chunk complete")
				return nil
			}

			errInner := d.downloadWithRange(fileUrl, &chunk, file, fm.ModTime)
			if errInner != nil {
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
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	hash, err := library.GetHash(fullPath)
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", chunk.Start, chunk.End))
	req.Header.Set("If-Range", modTime.UTC().Format(http.TimeFormat))

	resp, err := d.conf.getHttpClient().Do(req)
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
