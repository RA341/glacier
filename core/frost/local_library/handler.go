package download

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ncruces/zenity"
	"github.com/ra341/glacier/frost/local_library/download"
	v1 "github.com/ra341/glacier/generated/frost_library/v1"
	"github.com/ra341/glacier/generated/frost_library/v1/v1connect"
	"github.com/ra341/glacier/pkg/listutils"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewFrostLibraryServiceHandler(h)
}

func (h *Handler) ThrottleSpeed(ctx context.Context, c *connect.Request[v1.ThrottleSpeedRequest]) (*connect.Response[v1.ThrottleSpeedResponse], error) {
	newLm := rate.Limit(c.Msg.Limit)
	h.srv.downloader.SpeedLimiter.Set(newLm)
	return &connect.Response[v1.ThrottleSpeedResponse]{}, nil
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

func (h *Handler) LaunchFilePicker(ctx context.Context, c *connect.Request[v1.LaunchFilePickerRequest]) (*connect.Response[v1.LaunchFilePickerResponse], error) {
	file, err := zenity.SelectFile(
		zenity.Filename(c.Msg.BaseDir),
		zenity.Title("Pick the exe that launches the game"),
		zenity.DisallowEmpty(),
	)
	if err != nil && !errors.Is(err, zenity.ErrCanceled) {
		return nil, err
	}

	return connect.NewResponse(&v1.LaunchFilePickerResponse{Path: file}), nil
}

func (h *Handler) Edit(ctx context.Context, c *connect.Request[v1.EditRequest]) (*connect.Response[v1.EditResponse], error) {
	var ll LocalGame
	ll.FromProto(c.Msg.LocalGame)

	err := h.srv.store.Edit(ctx, &ll)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.EditResponse{}), nil
}

func (h *Handler) Launch(ctx context.Context, c *connect.Request[v1.LaunchRequest]) (*connect.Response[v1.LaunchResponse], error) {
	err := h.srv.Launch(ctx, int(c.Msg.Id))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.LaunchResponse{}), nil
}

func (h *Handler) GetByGameId(ctx context.Context, c *connect.Request[v1.GetByGameIdRequest]) (*connect.Response[v1.GetByGameIdResponse], error) {
	id, err := h.srv.store.GetByGameId(ctx, c.Msg.Id, c.Msg.LocalDownload)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetByGameIdResponse{
		Download: id.ToProto(),
	}), nil
}

func (h *Handler) Cancel(_ context.Context, c *connect.Request[v1.CancelRequest]) (*connect.Response[v1.CancelResponse], error) {
	err := h.srv.Cancel(int(c.Msg.Id))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.CancelResponse{}), nil
}

func (h *Handler) Pause(ctx context.Context, c *connect.Request[v1.PauseRequest]) (*connect.Response[v1.PauseResponse], error) {
	//TODO implement me
	return nil, fmt.Errorf("pause implement me")
}

func (h *Handler) Delete(ctx context.Context, c *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	err := h.srv.Delete(ctx, c.Msg.Id)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteResponse{}), nil
}

func (h *Handler) ListFiles(ctx context.Context, c *connect.Request[v1.ListFilesRequest]) (*connect.Response[v1.ListFilesResponse], error) {
	//TODO implement me
	return nil, fmt.Errorf("ListFiles implement me")
}

func (h *Handler) Download(ctx context.Context, c *connect.Request[v1.DownloadRequest]) (*connect.Response[v1.DownloadResponse], error) {
	err := h.srv.Download(
		ctx,
		int(c.Msg.GameId),
		c.Msg.DownloadFolder,
		c.Msg.Recheck,
		c.Msg.Force,
	)
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

	res := listutils.ToMap(games, func(t LocalGame) *v1.DownloadProgress {
		var progress []download.FileProgress
		value, ok := h.srv.downloader.ActiveDownloads.Load(t.GameId)
		if ok {
			progress, err = value.Progress()
			if err != nil {
				log.Warn().Msg("could not get download progress")
			}
		}

		var totalLeft int64 = 0
		var totalComplete int64 = 0

		toMap := listutils.ToMap(progress, func(t download.FileProgress) *v1.FileProgress {
			totalLeft += t.Left
			totalComplete += t.Complete

			return &v1.FileProgress{
				Name:     t.Name,
				Complete: uint64(t.Complete),
				Left:     uint64(t.Left),
			}
		})

		return &v1.DownloadProgress{
			ID:        uint64(t.ID),
			Thumbnail: t.Game.Meta.ThumbnailURL,
			Title:     t.Game.Meta.Name,
			Download:  t.Download.ToProto(),
			Progress: &v1.FolderProgress{
				Complete: totalComplete,
				Left:     totalLeft,
				Files:    toMap,
			},
		}
	})

	return connect.NewResponse(&v1.ListDownloadingResponse{
		Downloads: res,
	}), nil
}
