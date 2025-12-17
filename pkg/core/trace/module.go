package trace

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "trace"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(
			ProvideTracerProvider,
			fx.As(fx.Self()),
			fx.As(new(trace.TracerProvider)),
		),
		ProvideTracer,
	),
)

type ProvideTracerProviderParams struct {
	fx.In
	LifeCycle fx.Lifecycle
	Resource  *resource.Resource
	Options   []sdktrace.TracerProviderOption `group:"trace-provider-options"`
}

func ProvideTracerProvider(params ProvideTracerProviderParams) *sdktrace.TracerProvider {
	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(params.Resource),
	}

	tp := sdktrace.NewTracerProvider(append(opts, params.Options...)...)

	params.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := tp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return tp.Shutdown(ctx)
		},
	})

	return tp
}

type ProvideTracerParams struct {
	fx.In
	TracerProvider trace.TracerProvider
	Options        []trace.TracerOption `group:"trace-tracer-options"`
}

func ProvideTracer(params ProvideTracerParams) trace.Tracer {
	return params.TracerProvider.Tracer("github.com/go-oryn/oryn", params.Options...)
}
