package library

import (
	"context"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/library/v1"
	"github.com/ra341/glacier/generated/library/v1/v1connect"
	"github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/pkg/listutils"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	svc := &Handler{srv: srv}
	return v1connect.NewLibraryServiceHandler(svc)
}

func (h *Handler) Edit(ctx context.Context, c *connect.Request[v1.EditRequest]) (*connect.Response[v1.EditResponse], error) {
	var game Game
	game.FromProto(c.Msg.Game)

	err := h.srv.store.Edit(ctx, &game)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.EditResponse{}), nil
}

func (h *Handler) ListWithState(ctx context.Context, c *connect.Request[v1.ListWithStateRequest]) (*connect.Response[v1.ListWithStateResponse], error) {
	list, err := h.srv.ListDownloading(ctx, c.Msg.State)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(list, func(t Game) *v1.Game {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.ListWithStateResponse{
		Game: res,
	}), nil
}

func (h *Handler) List(ctx context.Context, c *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	list, err := h.srv.List(ctx, c.Msg.Query, uint(c.Msg.Offset), uint(c.Msg.Limit))
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(list, func(t Game) *v1.Game {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.ListResponse{
		GameList: res,
	}), nil
}

func (h *Handler) Delete(ctx context.Context, c *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	u, err := user.GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}
	if u.Role != user.Omnissiah {
		return nil, user.ErrNotAuthorized
	}

	err = h.srv.Delete(ctx, uint(c.Msg.GameId))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteResponse{}), nil
}

func (h *Handler) TriggerTracker(ctx context.Context, req *connect.Request[v1.TriggerTrackerRequest]) (*connect.Response[v1.TriggerTrackerResponse], error) {
	h.srv.downloader.TriggerTracker()
	return connect.NewResponse(&v1.TriggerTrackerResponse{}), nil
}

func (h *Handler) Exists(ctx context.Context, req *connect.Request[v1.ExistsRequest]) (*connect.Response[v1.ExistsResponse], error) {
	typeString, err := types.ProviderTypeString(req.Msg.MetadataType)
	if err != nil {
		return nil, err
	}

	gameId, err := h.srv.store.Exists(typeString, req.Msg.MetadataGameId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.ExistsResponse{
		GameId: uint64(gameId),
	}), nil
}

func (h *Handler) GetGame(ctx context.Context, c *connect.Request[v1.GetGameRequest]) (*connect.Response[v1.GetGameResponse], error) {
	get, err := h.srv.Get(ctx, uint(c.Msg.GameId))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetGameResponse{
		Game: get.ToProto(),
	}), nil
}

func (h *Handler) Add(ctx context.Context, req *connect.Request[v1.AddRequest]) (*connect.Response[v1.AddResponse], error) {
	var game = &Game{}
	game.FromProto(req.Msg.Game)

	err := h.srv.Add(ctx, game)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.AddResponse{}), nil
}

func (h *Handler) ListFiles(ctx context.Context, c *connect.Request[v1.ListFilesRequest]) (*connect.Response[v1.ListFilesResponse], error) {
	u, err := user.GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	if u.Role != user.Omnissiah {
		return nil, user.ErrNotAuthorized
	}

	files, err := h.srv.ListFiles(
		ctx,
		uint(c.Msg.GameId),
		c.Msg.Downloaded,
		c.Msg.BasePath,
	)
	if err != nil {
		return nil, err
	}

	rpcFile, err := listutils.ToMapErr(files, func(t os.DirEntry) (*v1.File, error) {
		info, err := t.Info()
		if err != nil {
			return nil, err
		}
		return &v1.File{
			RelPath: t.Name(),
			IsDir:   t.IsDir(),
			Date:    info.ModTime().Format(time.RFC3339),
			Size:    uint64(info.Size()),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.ListFilesResponse{Files: rpcFile}), nil
}

func (h *Handler) DeleteFile(ctx context.Context, c *connect.Request[v1.DeleteFileRequest]) (*connect.Response[v1.DeleteFileResponse], error) {
	u, err := user.GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	if u.Role != user.Omnissiah {
		return nil, user.ErrNotAuthorized
	}

	err = h.srv.DeleteFile(
		ctx,
		uint(c.Msg.GameId),
		c.Msg.Path,
		c.Msg.Downloaded,
	)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteFileResponse{}), nil
}
