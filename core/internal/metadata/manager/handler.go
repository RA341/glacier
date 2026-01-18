package manager

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/metadata/v1"
	"github.com/ra341/glacier/generated/metadata/v1/v1connect"
	"github.com/ra341/glacier/internal/metadata/types"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}

	return v1connect.NewMetadataServiceHandler(h)
}

func (h *Handler) GetMetadataProviders(ctx context.Context, c *connect.Request[v1.MetadataProviderRequest]) (*connect.Response[v1.MetadataProviderResponse], error) {
	var provs []*v1.Provider

	h.srv.providers.Range(func(key types.ProviderType, value types.Provider) bool {
		provs = append(provs, &v1.Provider{
			Name: key.String(),
		})
		return true
	})

	return connect.NewResponse(&v1.MetadataProviderResponse{
		Providers: provs,
	}), nil
}
