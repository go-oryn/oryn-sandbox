package otel

import (
	"context"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/otel/log"
	"github.com/go-oryn/oryn-sandbox/pkg/otel/metric"
	"github.com/go-oryn/oryn-sandbox/pkg/otel/trace"
	"go.opentelemetry.io/otel"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "otel"

var Module = fx.Module(
	ModuleName,
	// sub modules
	log.Module,
	metric.Module,
	trace.Module,
	// common dependencies
	fx.Provide(
		ProvideResource,
		ProvidePropagator,
		fx.Annotate(ProvideTelemetryWrapper, fx.As(fx.Self()), fx.As(new(Telemetry))),
	),
)

type ProvideResourceParams struct {
	fx.In
	Options []resource.Option `group:"otel-resource-options"`
}

func ProvideResource(params ProvideResourceParams) (*resource.Resource, error) {
	resOpts := []resource.Option{
		resource.WithSchemaURL(semconv.SchemaURL),
	}

	resOpts = append(resOpts, params.Options...)

	res, err := resource.New(context.Background(), resOpts...)
	if err != nil {
		return nil, err
	}

	return resource.Merge(resource.Default(), res)
}

func ProvidePropagator() propagation.TextMapPropagator {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	return propagator
}

type ProvideTelemetryWrapperParams struct {
	fx.In
	Logger *slog.Logger
	Meter  otelmetric.Meter
	Tracer oteltrace.Tracer
}

func ProvideTelemetryWrapper(params ProvideTelemetryWrapperParams) *TelemetryWrapper {
	return NewTelemetryWrapper(params.Logger, params.Meter, params.Tracer)
}
