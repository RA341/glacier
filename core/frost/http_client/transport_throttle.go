package http_client

import (
	"context"
	"io"
	"net/http"

	"golang.org/x/time/rate"
)

type LimiterService struct {
	lm *rate.Limiter
}

const DefaultBurst = 5 * 1024 * 1024

func NewLimiterService(initialLimit rate.Limit) *LimiterService {
	if initialLimit <= 0 {
		initialLimit = rate.Inf
	}

	ls := &LimiterService{}
	ls.lm = rate.NewLimiter(initialLimit, DefaultBurst)
	return ls
}

func (s *LimiterService) Set(limit rate.Limit) {
	if limit <= 0 {
		limit = rate.Inf
	}
	s.lm.SetLimit(limit)
}

func (s *LimiterService) Get() *rate.Limiter {
	return s.lm
}

type LimiterFn func() *rate.Limiter

type throttledReader struct {
	io.ReadCloser
	limiter LimiterFn
	ctx     context.Context
}

func (tr *throttledReader) Read(p []byte) (n int, err error) {
	n, err = tr.ReadCloser.Read(p)
	if n > 0 {
		// WaitN blocks until the limiter allows us to proceed with 'n' bytes.
		// This is shared across all goroutines using the same limiter.
		if err := tr.limiter().WaitN(tr.ctx, n); err != nil {
			return n, err
		}
	}
	return
}

type ThrottledTransport struct {
	Base    http.RoundTripper
	Limiter LimiterFn
}

func (t *ThrottledTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	// Wrap the body with the throttler
	resp.Body = &throttledReader{
		ReadCloser: resp.Body,
		limiter:    t.Limiter,
		ctx:        req.Context(),
	}
	return resp, nil
}
