package log

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/fx"
)

const ModuleName = "trace"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(
			ProvideLoggerProvider,
			fx.As(fx.Self()),
			fx.As(new(log.LoggerProvider)),
		),
		ProvideLogger,
	),
)

type ProvideLoggerProviderParams struct {
	fx.In
	LifeCycle fx.Lifecycle
	Resource  *resource.Resource
	Options   []sdklog.LoggerProviderOption `group:"log-provider-options"`
}

func ProvideLoggerProvider(params ProvideLoggerProviderParams) *sdklog.LoggerProvider {
	opts := []sdklog.LoggerProviderOption{
		sdklog.WithResource(params.Resource),
	}

	lp := sdklog.NewLoggerProvider(append(opts, params.Options...)...)

	params.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := lp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return lp.Shutdown(ctx)
		},
	})

	return lp
}

type ProvideLoggerParams struct {
	fx.In
	LifeCycle fx.Lifecycle
	Provider  log.LoggerProvider
	Options   []otelslog.Option `group:"log-logger-options"`
}

func ProvideLogger(params ProvideLoggerParams) *slog.Logger {
	opts := []otelslog.Option{
		otelslog.WithLoggerProvider(params.Provider),
	}

	return otelslog.NewLogger("github.com/go-oryn/oryn", append(opts, params.Options...)...)
}
