package trace

import (
	"context"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "trace"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(ProvideTracerProvider, fx.As(fx.Self()), fx.As(new(trace.TracerProvider))),
		ProvideTracer,
	),
)

type ProvideTracerProviderParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
	Resource  *resource.Resource
	Options   []sdktrace.TracerProviderOption `group:"otel-trace-provider-options"`
}

func ProvideTracerProvider(params ProvideTracerProviderParams) (*sdktrace.TracerProvider, error) {
	tpOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(params.Resource),
	}

	tpOpts = append(tpOpts, params.Options...)

	tp := sdktrace.NewTracerProvider(tpOpts...)

	otel.SetTracerProvider(tp)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := tp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return tp.Shutdown(ctx)
		},
	})

	return tp, nil
}

type ProvideTracerParams struct {
	fx.In
	Provider trace.TracerProvider
	Options  []trace.TracerOption `group:"otel-trace-tracer-options"`
}

func ProvideTracer(params ProvideTracerParams) trace.Tracer {
	return params.Provider.Tracer("github.com/go-oryn/oryn/otel/trace", params.Options...)
}
