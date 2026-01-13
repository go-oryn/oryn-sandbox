package log

import (
	"go.opentelemetry.io/contrib/bridges/otelslog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/fx"
)

func AsLoggerProviderOptions(options ...sdklog.LoggerProviderOption) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(sdklog.LoggerProviderOption)),
				fx.ResultTags(`group:"otel-log-provider-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}

func AsLoggerHandlerOptions(options ...otelslog.Option) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(otelslog.Option)),
				fx.ResultTags(`group:"otel-log-handler-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
