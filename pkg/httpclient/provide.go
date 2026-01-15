package httpclient

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"
)

func AsTransportOptions(options ...otelhttp.Option) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.As(new(otelhttp.Option)),
				fx.ResultTags(`group:"httpclient-transport-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
