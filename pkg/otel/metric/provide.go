package metric

import (
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

func AsMeterProviderOptions(options ...sdkmetric.Option) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(sdkmetric.Option)),
				fx.ResultTags(`group:"otel-metric-provider-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}

func AsMeterOptions(options ...metric.MeterOption) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(metric.MeterOption)),
				fx.ResultTags(`group:"otel-metric-meter-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
