package config

import (
	"go.uber.org/fx"
)

func AsConfigOptions(options ...Option) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(Option)),
				fx.ResultTags(`group:"config-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
