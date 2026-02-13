package http_client

import (
	"net/http"

	"github.com/ra341/glacier/frost/secrets"
)

type HttpCliFactory func(customTransport *http.Transport) *http.Client

func NewFrostHttpClientFactory(sec *secrets.Service, globalCounter *uint64, limiter LimiterFn) HttpCliFactory {
	return func(customTransport *http.Transport) *http.Client {
		// transport middleware in order of priority LIFO will run first

		// Wrap with Network layer
		var transport http.RoundTripper = customTransport
		if transport == nil {
			transport = http.DefaultTransport
		}

		if limiter != nil {
			transport = &ThrottledTransport{
				Base:    transport,
				Limiter: limiter,
			}
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
