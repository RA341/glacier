package http_client

import (
	"net/http"

	"github.com/ra341/glacier/frost/secrets"
	"github.com/ra341/glacier/internal/auth"
)

type FrostTransport struct {
	base   http.RoundTripper
	secret *secrets.Service
}

type HttpCliFactory func(customTransport *http.Transport) *http.Client

func NewFrostHttpClientFactory(sec *secrets.Service) HttpCliFactory {
	return func(customTransport *http.Transport) *http.Client {
		return &http.Client{
			Transport: NewFrostTransport(
				customTransport,
				sec,
			),
		}
	}
}

func NewFrostTransport(customTransport http.RoundTripper, sec *secrets.Service) *FrostTransport {
	if customTransport == nil {
		customTransport = http.DefaultTransport
	}
	return &FrostTransport{
		base:   customTransport,
		secret: sec,
	}
}

func (t *FrostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to stay thread-safe
	outReq := req.Clone(req.Context())

	outReq.Header.Set(auth.FrostReqHeader, "true")

	session, refresh := t.secret.GetSession()

	outReq.Header.Set(auth.HeaderFrostRefreshToken, refresh)
	outReq.Header.Set(auth.HeaderFrostSessionToken, session)

	return t.base.RoundTrip(outReq)
}
