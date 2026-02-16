package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "github.com/ra341/glacier/generated/config/v1"
	"github.com/ra341/glacier/generated/config/v1/v1connect"
	"github.com/ra341/glacier/shared/config"
	"github.com/rs/zerolog/log"

	"connectrpc.com/connect"
)

type Handler struct {
	srv *config.Service[Config]
}

func NewHandler(srv *config.Service[Config]) (string, http.Handler) {
	h := &Handler{
		srv: srv,
	}
	return v1connect.NewConfigServiceHandler(h)
}

func (h *Handler) Get(ctx context.Context, req *connect.Request[v1.GetRequest]) (*connect.Response[v1.GetResponse], error) {
	schema, err := h.srv.GetSchema()
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetResponse{
		ConfigSchema: string(schema),
	}), nil
}

func (h *Handler) GetField(ctx context.Context, c *connect.Request[v1.GetFieldRequest]) (*connect.Response[v1.GetFieldResponse], error) {
	url := h.srv.Get().Server.GlacierUrl

	var checkErr string
	if url != "" {
		err2 := checkGlacierUrl(url)
		if err2 != nil {
			checkErr = err2.Error()
		}
	}
	return connect.NewResponse(&v1.GetFieldResponse{
		Value: url,
		Error: checkErr,
	}), nil
}

func (h *Handler) SetField(ctx context.Context, c *connect.Request[v1.SetFieldRequest]) (*connect.Response[v1.SetFieldResponse], error) {
	value := c.Msg.Value
	if value == "" {
		return nil, fmt.Errorf("empty value")
	}

	validUrl, err := testGlacierUrls(value)
	if err != nil {
		return nil, err
	}
	log.Debug().Str("url", validUrl).Msg("found valid glacier url")

	get := h.srv.Get()
	get.Server.GlacierUrl = validUrl
	err = h.srv.Set(get)
	if err != nil {
		return nil, fmt.Errorf("could not set glacier url: %w", err)
	}

	return connect.NewResponse(&v1.SetFieldResponse{}), nil
}

func (h *Handler) Set(ctx context.Context, req *connect.Request[v1.SetRequest]) (*connect.Response[v1.SetResponse], error) {
	var conf Config

	err := json.Unmarshal([]byte(req.Msg.ConfigSchema), &conf)
	if err != nil {
		return nil, err
	}

	err = h.srv.Set(&conf)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.SetResponse{}), nil

}

func (h *Handler) ListFiles(ctx context.Context, c *connect.Request[v1.ListFilesRequest]) (*connect.Response[v1.ListFilesResponse], error) {
	files, err := h.srv.ListFiles(c.Msg.Base)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.ListFilesResponse{
		Files: files,
	}), nil
}
