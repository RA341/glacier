package download

import (
	"context"
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

type Download struct {
	ctx    context.Context
	cancel context.CancelFunc

	conf           Config
	gameId         int
	baseUrl        string
	downloadFolder string
	met            *library.FolderMetadata

	cacheStore CacheStore
}

const MetadataFolder = ".frost.metadata"

func NewDownload(config Config, baseUrl string, gameId int, downloadFolder string, met *library.FolderMetadata) (*Download, error) {
	metaPath := filepath.Join(downloadFolder, MetadataFolder)
	db, err := NewCacheStoreBadger(metaPath)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	d := &Download{
		ctx:     ctx,
		cancel:  cancel,
		baseUrl: fmt.Sprintf("%s/load/%d", baseUrl, gameId),
		gameId:  gameId,
		met:     met,
		conf:    config,

		cacheStore:     db,
		downloadFolder: downloadFolder,
	}
	d.Start()

	return d, nil
}

func (d *Download) Start() {
	defer fileutil.Close(d.cacheStore)

	start := time.Now()

	eg := errgroup.Group{}
	eg.SetLimit(d.conf.getMaxConcurrentFileChunks())

	for _, fi := range d.met.FileInfo {
		eg.Go(func() error {
			err := d.gatherMeta(&fi)
			if err != nil {
				log.Error().Err(err).Str("file", fi.RelPath).Msg("metadata cache err")
				// todo how to handle err
			}

			err = d.downloadFile(&fi)
			if err != nil {
				log.Error().Err(err).Str("file", fi.RelPath).Msg("download err")
				// todo how to handle err
			}

			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		log.Info().Err(err).Msg("download err")
	}

	elapsed := time.Since(start)
	log.Info().
		Str("elapsed", elapsed.String()).
		Int("game", d.gameId).
		Msg("download finished")
}

func (d *Download) Close() {
	fileutil.Close(d.cacheStore)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// metadata step

func (d *Download) gatherMeta(fm *library.FileMetadata) error {
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
			State: Queued,
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
// downloads

func (d *Download) downloadFile(fm *library.FileMetadata) error {
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
	fileUrl := fmt.Sprintf("%s/%s", d.baseUrl, escaped)

	eg := errgroup.Group{}
	eg.SetLimit(d.conf.getMaxConcurrentFileChunks())

	for i, chunk := range chunks {
		eg.Go(func() error {
			if chunk.State == Complete {
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
				chunk.State = Error
			} else {
				chunk.State = Complete
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
	if resp.StatusCode >= 400 {
		all, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return fmt.Errorf("could not load body to get error message: %w", err2)
		}
		return fmt.Errorf("error downloading: %d: %s", resp.StatusCode, string(all))
	}

	if resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server did not support range or file changed: status %s", resp.Status)
	}

	_, err = io.Copy(NewOffsetWriter(writer, chunk.Start), resp.Body)

	return err
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
