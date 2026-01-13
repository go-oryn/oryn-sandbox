package otel

import (
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/fx"
)

func AsResourceOptions(options ...resource.Option) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(resource.Option)),
				fx.ResultTags(`group:"otel-resource-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
