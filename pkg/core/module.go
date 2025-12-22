package core

import (
	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"github.com/go-oryn/oryn-sandbox/pkg/core/log"
	"github.com/go-oryn/oryn-sandbox/pkg/core/metric"
	"github.com/go-oryn/oryn-sandbox/pkg/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/fx"
)

const ModuleName = "core"

var Module = fx.Module(
	ModuleName,
	// Core sub modules
	config.Module,
	log.Module,
	metric.Module,
	trace.Module,
	// Core common dependencies
	fx.Provide(
		ProvideOTelResource,
		ProvideOTelPropagator,
	),
)

type ProvideOTELResourceParams struct {
	fx.In
	Config *config.Config
}

func ProvideOTelResource(params ProvideOTELResourceParams) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(params.Config.GetString("app.name")),
			semconv.ServiceVersion(params.Config.GetString("app.version")),
		),
	)
}

func ProvideOTelPropagator() propagation.TextMapPropagator {
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)

	otel.SetTextMapPropagator(propagator)

	return propagator
}
