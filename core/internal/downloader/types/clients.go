package types

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=ClientType -output=enum_client_type.go
type ClientType int

const (
	ClientUnknown ClientType = iota
	ClientTransmission
)

type ClientConfig = map[string]any
