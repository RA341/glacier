package download

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func NewDownloadClient(maxWorkers, maxConn int) *http.Client {
	transport := &http.Transport{
		// MaxIdleConns is the total connections across all hosts
		MaxIdleConns: maxConn,

		// MaxIdleConnsPerHost must be >= your worker count.
		// The default is only 2 If you have 30 workers, 28 will
		// constantly create new TCP connections.
		MaxIdleConnsPerHost: maxWorkers,

		// Time to keep an idle connection in the pool
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0,
	}
}

type Downloader struct {
	BaseUrl    string
	httpClient *http.Client

	maxConcurrentFiles      int
	maxConcurrentFileChunks int
	chunkSize               int64
}

const MB = 1024 * 1024

func NewDownloader(maxConcurrentFiles, maxConcurrentFileChunks int) *Downloader {
	return &Downloader{
		httpClient: NewDownloadClient(100, 0),
		BaseUrl:    "http://localhost:6699/api/server/library/download",

		maxConcurrentFiles:      maxConcurrentFiles,
		maxConcurrentFileChunks: maxConcurrentFileChunks,
		chunkSize:               MB * 250,
	}
}

func (d *Downloader) SetupDownloadMetadata() {

}

func (d *Downloader) WriteDownloadMetadata() {

}

func (d *Downloader) Download(downloadFolder string, gameId int, met *library.FolderMetadata) error {

	return nil
}

func (d *Downloader) StartDownload(downloadFolder string, gameId int, met *library.FolderMetadata) error {
	start := time.Now()

	eg := errgroup.Group{}
	eg.SetLimit(d.maxConcurrentFiles)

	baseUrl := fmt.Sprintf("%s/load/%d", d.BaseUrl, gameId)

	for _, fi := range met.FileInfo {
		eg.Go(func() error {
			err := d.downloadFile(downloadFolder, baseUrl, &fi)
			if err != nil {
				log.Error().Err(err).Str("file", fi.RelPath).Msg("download err")
			}
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	log.Info().
		Str("elapsed", elapsed.String()).
		Int("game", gameId).
		Msg("download finished")

	return nil
}

func (d *Downloader) downloadFile(downloadFolder string, baseurl string, met *library.FileMetadata) error {
	//log.Info().
	//	Str("file", met.RelPath).
	//	Str("size", humanize.Bytes(uint64(met.Size))).
	//	Msg("starting download")

	started := time.Now()

	fullPath := filepath.Join(downloadFolder, met.RelPath)
	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fileutil.Close(file)

	// allocate file
	err = file.Truncate(met.Size)
	if err != nil {
		return err
	}

	escaped := url.PathEscape(met.RelPath)
	fileUrl := fmt.Sprintf("%s/%s", baseurl, escaped)

	totalSize := met.Size

	eg := errgroup.Group{}
	eg.SetLimit(d.maxConcurrentFileChunks)

	for start := int64(0); start < totalSize; start += d.chunkSize {
		end := start + d.chunkSize - 1
		// if it's the last chunk,
		// make sure not to overshoot the file size
		if end >= totalSize {
			end = totalSize - 1
		}

		eg.Go(func() error {
			errInner := d.downloadWithRange(fileUrl, start, end, file, met.ModTime)
			if errInner != nil {
				log.Error().Err(errInner).
					Int64("start", start).Int64("end", end).
					Msg("could not download chunk")
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

	if met.Checksum != hash {
		return fmt.Errorf(
			"checksum mismatch, expected: %s != got: %s\nExpected Size: %s, got size: %s",
			met.Checksum,
			hash,
			humanize.Bytes(uint64(met.Size)),
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

func (d *Downloader) downloadWithRange(url string, start, end int64, writer io.WriterAt, modTime time.Time) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	req.Header.Set("If-Range", modTime.UTC().Format(http.TimeFormat))

	resp, err := d.httpClient.Do(req)
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

	_, err = io.Copy(NewOffsetWriter(writer, start), resp.Body)

	return err
}

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
