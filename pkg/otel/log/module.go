package log

import (
	"context"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.uber.org/fx"
)

const ModuleName = "log"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		fx.Annotate(ProvideLoggerProvider, fx.As(fx.Self()), fx.As(new(log.LoggerProvider))),
		fx.Annotate(ProvideLoggerHandler, fx.As(fx.Self()), fx.As(new(slog.Handler))),
		ProvideLogger,
	),
)

type ProvideLoggerProviderParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Config    *config.Config
	Resource  *resource.Resource
	Options   []sdklog.LoggerProviderOption `group:"otel-log-provider-options"`
}

func ProvideLoggerProvider(params ProvideLoggerProviderParams) (*sdklog.LoggerProvider, error) {
	lpOpts := []sdklog.LoggerProviderOption{
		sdklog.WithResource(params.Resource),
	}

	lpOpts = append(lpOpts, params.Options...)

	lp := sdklog.NewLoggerProvider(lpOpts...)

	global.SetLoggerProvider(lp)

	params.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			err := lp.ForceFlush(ctx)
			if err != nil {
				return err
			}

			return lp.Shutdown(ctx)
		},
	})

	return lp, nil
}

type ProvideLoggerHandlerParams struct {
	fx.In
	Provider log.LoggerProvider
	Options  []otelslog.Option `group:"otel-log-handler-options"`
}

func ProvideLoggerHandler(params ProvideLoggerHandlerParams) *otelslog.Handler {
	lhOpts := []otelslog.Option{
		otelslog.WithLoggerProvider(params.Provider),
	}

	lhOpts = append(lhOpts, params.Options...)

	return otelslog.NewHandler("github.com/go-oryn/oryn/otel/log", lhOpts...)
}

type ProvideLoggerParams struct {
	fx.In
	Handler slog.Handler
}

func ProvideLogger(params ProvideLoggerParams) *slog.Logger {
	//	otelslog.WithSource(params.Config.GetBool("log.source")),

	return slog.New(params.Handler)
}
