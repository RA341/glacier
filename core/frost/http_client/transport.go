package http_client

import (
	"io"
	"net/http"
	"sync/atomic"

	"github.com/ra341/glacier/frost/secrets"
	"github.com/ra341/glacier/internal/auth"
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

type FrostTransport struct {
	base   http.RoundTripper
	secret *secrets.Service
}

func (t *FrostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	outReq := req.Clone(req.Context())
	outReq.Header.Set(auth.FrostReqHeader, "true")

	session, refresh := t.secret.GetSession()
	outReq.Header.Set(auth.HeaderFrostRefreshToken, refresh)
	outReq.Header.Set(auth.HeaderFrostSessionToken, session)

	return t.base.RoundTrip(outReq)
}

type HttpCliFactory func(customTransport *http.Transport) *http.Client

func NewFrostHttpClientFactory(sec *secrets.Service, globalCounter *uint64) HttpCliFactory {
	return func(customTransport *http.Transport) *http.Client {
		// Wrap with Network layer
		var transport http.RoundTripper = customTransport
		if transport == nil {
			transport = http.DefaultTransport
		}

		// Wrap with Speed Tracking
		if globalCounter != nil {
			transport = &SpeedTrackingTransport{
				Base:    transport,
				Counter: globalCounter,
			}
		}

		// Wrap with Frost Auth
		transport = &FrostTransport{
			base:   transport,
			secret: sec,
		}

		return &http.Client{
			Transport: transport,
		}
	}
}
