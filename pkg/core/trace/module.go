package trace

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
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
	Options   []sdktrace.TracerProviderOption `group:"trace-provider-options"`
}

func ProvideTracerProvider(params ProvideTracerProviderParams) (*sdktrace.TracerProvider, error) {
	tpOpts, err := TracerProviderOptions(context.Background(), params.Config, params.Resource, params.Options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer provider options: %w", err)
	}

	tp := sdktrace.NewTracerProvider(tpOpts...)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := tp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return tp.Shutdown(ctx)
		},
	})

	otel.SetTracerProvider(tp)

	return tp, nil
}

type ProvideTracerParams struct {
	fx.In
	Provider trace.TracerProvider
	Options  []trace.TracerOption `group:"trace-tracer-options"`
}

func ProvideTracer(params ProvideTracerParams) trace.Tracer {
	return params.Provider.Tracer("github.com/go-oryn/oryn", params.Options...)
}
