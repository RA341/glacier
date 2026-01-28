package config_manager

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/service_config/v1"
	"github.com/ra341/glacier/generated/service_config/v1/v1connect"
	"github.com/ra341/glacier/pkg/listutils"
	"github.com/ra341/glacier/pkg/mapsct"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewServiceConfigServiceHandler(h)
}

func (h *Handler) GetSupportedValues(ctx context.Context, c *connect.Request[v1.GetSupportedValuesRequest]) (*connect.Response[v1.GetSupportedValuesResponse], error) {
	typeString, err := ServiceTypeString(c.Msg.ServiceType)
	if err != nil {
		return nil, err
	}

	values, err := h.srv.GetSupportedValues(typeString)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetSupportedValuesResponse{Values: values}), nil
}

func (h *Handler) GetSchema(ctx context.Context, c *connect.Request[v1.GetSchemaRequest]) (*connect.Response[v1.GetSchemaResponse], error) {
	typeString, err := ServiceTypeString(c.Msg.ServiceType)
	if err != nil {
		return nil, err
	}

	schema, err := h.srv.GetSchema(typeString, c.Msg.Flavour)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(schema, func(t mapsct.FieldSchema) *v1.FieldSchema {
		return &v1.FieldSchema{
			Name:      t.Name,
			Type:      t.Type,
			InsertKey: t.InsertKey,
			KeyType:   t.KeyType,
			ValueType: t.ValueType,
		}
	})

	return connect.NewResponse(&v1.GetSchemaResponse{
		Fields: res,
	}), nil

}

func (h *Handler) GetActiveService(ctx context.Context, req *connect.Request[v1.GetActiveServiceRequest]) (*connect.Response[v1.GetActiveServiceResponse], error) {
	serviceType := req.Msg.ServiceType
	typeString, err := ServiceTypeString(serviceType)
	if err != nil {
		return nil, err
	}

	enabled, err := h.srv.store.ListEnabled(typeString)
	if err != nil {
		return nil, err
	}

	res, err := listutils.ToMapErr(enabled, func(t ServiceConfig) (*v1.ServiceConfig, error) {
		// remove any sensitive config info this endpoint does not need it
		t.Config = nil
		return t.ToProto()
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetActiveServiceResponse{
		Names: res,
	}), nil
}

func (h *Handler) List(ctx context.Context, c *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	typeString, err := ServiceTypeString(c.Msg.ServiceType)
	if err != nil {
		return nil, err
	}

	all, err := h.srv.store.ListAll(typeString)
	if err != nil {
		return nil, err
	}

	res, err := listutils.ToMapErr(all, func(t ServiceConfig) (*v1.ServiceConfig, error) {
		return t.ToProto()
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.ListResponse{Conf: res}), nil
}

func (h *Handler) ListEnabled(ctx context.Context, c *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) New(ctx context.Context, c *connect.Request[v1.NewConfigRequest]) (*connect.Response[v1.NewConfigResponse], error) {
	var cf ServiceConfig
	err := cf.FromProto(c.Msg.Conf)
	if err != nil {
		return nil, err
	}

	err = h.srv.TestAndSave(&cf)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.NewConfigResponse{}), nil
}

func (h *Handler) Get(ctx context.Context, c *connect.Request[v1.GetRequest]) (*connect.Response[v1.GetResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (h *Handler) Edit(ctx context.Context, c *connect.Request[v1.EditRequest]) (*connect.Response[v1.EditResponse], error) {
	var cf ServiceConfig
	err := cf.FromProto(c.Msg.Conf)
	if err != nil {
		return nil, err
	}

	err = h.srv.store.Edit(&cf)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.EditResponse{}), nil
}

func (h *Handler) Delete(ctx context.Context, c *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	//TODO implement me
	panic("implement me")
}
