package config

import (
	"context"
	"encoding/json"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/config/v1"
	"github.com/ra341/glacier/generated/config/v1/v1connect"
	"github.com/ra341/glacier/shared/config"
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
