package config

import "go.uber.org/fx"

const ModuleName = "config"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideConfig,
	),
)

type ProvideConfigParams struct {
	fx.In
	Options []Option `group:"config-options"`
}

func ProvideConfig(params ProvideConfigParams) (*Config, error) {
	return NewConfig(params.Options...)
}
