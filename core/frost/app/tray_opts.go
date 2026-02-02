package app

import "github.com/ra341/glacier/shared/api"

type TrayConfig struct {
	disableUI  bool
	serverBase api.ServerBase
}

type Opt func(config *TrayConfig)

func WithServerBase(opts ...api.ServerOpt) Opt {
	return func(config *TrayConfig) {
		var serverBase api.ServerBase
		api.ParseOpts(&serverBase, opts...)
		config.serverBase = serverBase
	}
}

func WithDisableTrayUI() Opt {
	return func(config *TrayConfig) {
		config.disableUI = true
	}
}
