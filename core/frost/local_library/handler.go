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

func (h *Handler) Get(ctx context.Context, c *connect.Request[v1.GetRequest]) (*connect.Response[v1.GetResponse], error) {
	get, err := h.srv.store.Get(ctx, int(c.Msg.Id))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetResponse{
		Lg: get.ToProto(),
	}), nil
}

func (h *Handler) Delete(ctx context.Context, c *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) ListFiles(ctx context.Context, c *connect.Request[v1.ListFilesRequest]) (*connect.Response[v1.ListFilesResponse], error) {
	//TODO implement me
	panic("implement me")
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewFrostLibraryServiceHandler(h)
}

func (h *Handler) Download(ctx context.Context, c *connect.Request[v1.DownloadRequest]) (*connect.Response[v1.DownloadResponse], error) {
	err := h.srv.Download(ctx, int(c.Msg.GameId), c.Msg.DownloadFolder)
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
		var totalLeft int64 = 0
		var totalComplete int64 = 0

		toMap := listutils.ToMap(v, func(t download.FileProgress) *v1.FileProgress {
			totalLeft += t.Left
			totalComplete += t.Complete

			return &v1.FileProgress{
				Name:     t.Name,
				Complete: uint64(t.Complete),
				Left:     uint64(t.Left),
			}
		})

		res[uint64(k)] = &v1.FolderProgress{
			Complete: totalComplete,
			Left:     totalLeft,
			Files:    toMap,
		}
	}

	return connect.NewResponse(&v1.ListDownloadingResponse{
		Downloads: res,
	}), nil
}
