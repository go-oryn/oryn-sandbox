package trace

import (
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func AsTracerProviderOptions(options ...sdktrace.TracerProviderOption) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(sdktrace.TracerProviderOption)),
				fx.ResultTags(`group:"trace-provider-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}

func AsTracerOptions(options ...trace.TracerOption) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(trace.TracerOption)),
				fx.ResultTags(`group:"trace-tracer-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
