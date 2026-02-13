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

func (t *FrostTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	outReq := req.Clone(req.Context())
	outReq.Header.Set(auth.FrostReqHeader, "true")

	session, refresh := t.secret.GetSession()
	outReq.Header.Set(auth.HeaderFrostRefreshToken, refresh)
	outReq.Header.Set(auth.HeaderFrostSessionToken, session)

	return t.base.RoundTrip(outReq)
}
