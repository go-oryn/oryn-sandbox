package httpclient

import (
	"net/http"
	"time"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const (
	ModuleName     = "httpclient"
	DefaultTimeout = 30 * time.Second
)

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(ProvideTransport, fx.As(fx.Self()), fx.As(new(http.RoundTripper))),
		ProvideClient,
	),
)

type ProvideTransportParams struct {
	fx.In
	Propagator     propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	Options        []otelhttp.Option `group:"httpclient-transport-options"`
}

func ProvideTransport(params ProvideTransportParams) *otelhttp.Transport {
	hcOpts := []otelhttp.Option{
		otelhttp.WithPropagators(params.Propagator),
		otelhttp.WithTracerProvider(params.TracerProvider),
		otelhttp.WithMeterProvider(params.MeterProvider),
	}

	hcOpts = append(hcOpts, params.Options...)

	return NewTransport(hcOpts...)
}

type ProvideClientParams struct {
	fx.In
	Config       *config.Config
	RoundTripper http.RoundTripper
}

func ProvideClient(params ProvideClientParams) *http.Client {
	return &http.Client{
		Transport: params.RoundTripper,
		Timeout:   params.Config.GetDurationOrDefault("httpclient.timeout", DefaultTimeout),
	}
}
