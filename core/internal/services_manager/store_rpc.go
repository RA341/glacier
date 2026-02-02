package services_manager

import (
	"encoding/json"

	v1 "github.com/ra341/glacier/generated/service_config/v1"
)

func (sc *ServiceConfig) ToProto() (*v1.ServiceConfig, error) {
	marshal, err := json.Marshal(sc.Config)
	if err != nil {
		return nil, err
	}

	return &v1.ServiceConfig{
		ID:          uint64(sc.ID),
		ServiceType: sc.ServiceType.String(),
		Name:        sc.Name,
		Enabled:     sc.Enabled,
		Flavour:     sc.Flavour,
		Config:      marshal,
	}, nil
}

func (sc *ServiceConfig) FromProto(pb *v1.ServiceConfig) error {
	var conf map[string]interface{}

	err := json.Unmarshal(pb.Config, &conf)
	if err != nil {
		return err
	}

	typeString, err := ServiceTypeString(pb.ServiceType)
	if err != nil {
		return err
	}

	sc.ID = uint(pb.ID)
	sc.ServiceType = typeString
	sc.Name = pb.Name
	sc.Enabled = pb.Enabled
	sc.Flavour = pb.Flavour
	sc.Config = conf

	return nil
}
