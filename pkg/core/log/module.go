package log

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
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
	Config    *config.Config
	Resource  *resource.Resource
	Options   []sdklog.LoggerProviderOption `group:"log-provider-options"`
}

func ProvideLoggerProvider(params ProvideLoggerProviderParams) (*sdklog.LoggerProvider, error) {
	lpOpts, err := LoggerProviderOptions(context.Background(), params.Config, params.Resource, params.Options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger provider options: %w", err)
	}

	lp := sdklog.NewLoggerProvider(lpOpts...)

	params.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := lp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return lp.Shutdown(ctx)
		},
	})

	global.SetLoggerProvider(lp)

	return lp, nil
}

type ProvideLoggerParams struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Provider  log.LoggerProvider
	Options   []otelslog.Option `group:"log-logger-options"`
}

func ProvideLogger(params ProvideLoggerParams) *slog.Logger {
	opts := []otelslog.Option{
		otelslog.WithLoggerProvider(params.Provider),
		otelslog.WithSource(params.Config.GetBool("log.source")),
	}

	handler := otelslog.NewHandler("github.com/go-oryn/oryn", append(opts, params.Options...)...)

	return slog.New(NewLeveledHandler(
		ParseLogLevel(params.Config.GetString("log.level")),
		handler,
	))
}
