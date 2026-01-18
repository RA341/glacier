package manager

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/downloader/v1"
	"github.com/ra341/glacier/generated/downloader/v1/v1connect"
	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/pkg/listutils"
	"golang.org/x/exp/slices"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) (string, http.Handler) {
	h := &Handler{s: s}
	return v1connect.NewDownloaderServiceHandler(h)
}

func (h *Handler) GetActiveClients(ctx context.Context, c *connect.Request[v1.GetActiveClientsRequest]) (*connect.Response[v1.GetActiveClientsResponse], error) {
	clients := h.s.GetActive()
	slices.Sort(clients)

	result := listutils.ToMap(clients, func(t types.ClientType) *v1.Client {
		return &v1.Client{
			Name: t.String(),
		}
	})

	return connect.NewResponse(&v1.GetActiveClientsResponse{
		Clients: result,
	}), nil
}
