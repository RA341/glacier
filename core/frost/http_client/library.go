package http_client

import (
	"net/http"

	glacier "github.com/ra341/glacier/generated/library/v1/v1connect"
)

type LibraryRpcClient struct {
	cachedCli glacier.LibraryServiceClient
	cachedUrl string
	baseurl   func() string
	fac       HttpCliFactory
}

func NewLibCli(baseurl func() string, cli HttpCliFactory) *LibraryRpcClient {
	return &LibraryRpcClient{
		baseurl: baseurl,
		fac:     cli,
	}
}

func (c *LibraryRpcClient) Get() glacier.LibraryServiceClient {
	if c.cachedCli != nil && c.baseurl() == c.cachedUrl {
		return c.cachedCli
	}

	c.cachedUrl = c.baseurl()
	return glacier.NewLibraryServiceClient(
		c.fac(&http.Transport{}),
		c.cachedUrl,
	)
}
