package download

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ra341/glacier/frost/local_library/download"
	v1 "github.com/ra341/glacier/generated/frost_library/v1"
	"github.com/ra341/glacier/generated/frost_library/v1/v1connect"
	"github.com/ra341/glacier/pkg/listutils"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewFrostLibraryServiceHandler(h)
}

func (h *Handler) Download(ctx context.Context, c *connect.Request[v1.DownloadRequest]) (*connect.Response[v1.DownloadResponse], error) {
	err := h.srv.Download(int(c.Msg.GameId), c.Msg.DownloadFolder)
	if err != nil {
		return nil, err
	}

	return &connect.Response[v1.DownloadResponse]{}, nil
}

func (h *Handler) ListDownloading(ctx context.Context, c *connect.Request[v1.ListDownloadingRequest]) (*connect.Response[v1.ListDownloadingResponse], error) {
	games, err := h.srv.ListDownloading(ctx)
	if err != nil {
		return nil, err
	}

	var res = map[uint64]*v1.FolderProgress{}

	for k, v := range games {
		res[uint64(k)] = &v1.FolderProgress{
			Files: listutils.ToMap(v, func(t download.FileProgress) *v1.FileProgress {
				return &v1.FileProgress{
					Name:     t.Name,
					Complete: uint64(t.Complete),
					Left:     uint64(t.Left),
				}
			}),
		}
	}

	return connect.NewResponse(&v1.ListDownloadingResponse{
		Downloads: res,
	}), nil
}
