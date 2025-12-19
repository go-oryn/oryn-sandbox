package metric

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
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
		fx.Annotate(
			ProvideMeterProvider,
			fx.As(fx.Self()),
			fx.As(new(metric.MeterProvider)),
		),
		ProvideMeter,
	),
)

type ProvideMeterProviderParams struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Resource  *resource.Resource
	Options   []sdkmetric.Option `group:"metric-provider-options"`
}

func ProvideMeterProvider(params ProvideMeterProviderParams) (*sdkmetric.MeterProvider, error) {
	mpOpts, err := MeterProviderOptions(context.Background(), params.Config, params.Resource, params.Options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create meter provider options: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(mpOpts...)

	params.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := mp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return mp.Shutdown(ctx)
		},
	})

	otel.SetMeterProvider(mp)

	return mp, nil
}

type ProvideMeterParams struct {
	fx.In
	Provider metric.MeterProvider
	Options  []metric.MeterOption `group:"metric-meter-options"`
}

func ProvideMeter(params ProvideMeterParams) metric.Meter {
	return params.Provider.Meter("github.com/go-oryn/oryn", params.Options...)
}
