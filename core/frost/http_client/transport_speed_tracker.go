package http_client

import (
	"io"
	"net/http"
	"sync/atomic"
)

type trackingReader struct {
	io.ReadCloser
	counter *uint64
}

func (tr *trackingReader) Read(p []byte) (n int, err error) {
	n, err = tr.ReadCloser.Read(p)
	atomic.AddUint64(tr.counter, uint64(n))
	return
}

type SpeedTrackingTransport struct {
	Base    http.RoundTripper
	Counter *uint64
}

func (t *SpeedTrackingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	resp.Body = &trackingReader{ReadCloser: resp.Body, counter: t.Counter}
	return resp, nil
}
