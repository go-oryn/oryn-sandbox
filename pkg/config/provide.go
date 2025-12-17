package config

import (
	"go.uber.org/fx"
)

// AsConfigOptions provides the embed.FS for config files lookup.
func AsConfigOptions(options ...Option) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.ResultTags(`group:"config-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
