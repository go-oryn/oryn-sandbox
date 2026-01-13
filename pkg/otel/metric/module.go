package metric

import (
	"context"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/fx"
)

const ModuleName = "metric"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(ProvideMeterProvider, fx.As(fx.Self()), fx.As(new(metric.MeterProvider))),
		ProvideMeter,
	),
)

type ProvideMeterProviderParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
	Resource  *resource.Resource
	Options   []sdkmetric.Option `group:"otel-metric-provider-options"`
}

func ProvideMeterProvider(params ProvideMeterProviderParams) (*sdkmetric.MeterProvider, error) {
	mpOpts := []sdkmetric.Option{
		sdkmetric.WithResource(params.Resource),
	}

	mpOpts = append(mpOpts, params.Options...)

	mp := sdkmetric.NewMeterProvider(mpOpts...)

	otel.SetMeterProvider(mp)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := mp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return mp.Shutdown(ctx)
		},
	})

	return mp, nil
}

type ProvideMeterParams struct {
	fx.In
	Provider metric.MeterProvider
	Options  []metric.MeterOption `group:"otel-metric-meter-options"`
}

func ProvideMeter(params ProvideMeterParams) metric.Meter {
	return params.Provider.Meter("github.com/go-oryn/oryn/otel/metric", params.Options...)
}
