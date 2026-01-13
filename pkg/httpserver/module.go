package httpserver

import (
	"context"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "httpserver"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideRegistry,
		ProvideServer,
	),
)

type ProvideRegistryParams struct {
	fx.In
	Logger              *slog.Logger
	Handlers            []Handler           `group:"httpserver-handlers"`
	HandlersDefinitions []HandlerDefinition `group:"httpserver-handlers-definitions"`
}

func ProvideRegistry(params ProvideRegistryParams) *Registry {
	return NewRegistry(params.Logger, params.Handlers, params.HandlersDefinitions)
}

type ProvideServerParams struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Shutdown       fx.Shutdowner
	Config         *config.Config
	Logger         *slog.Logger
	Propagator     propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	Registry       *Registry
}

func ProvideServer(params ProvideServerParams) (*echo.Echo, error) {
	server := echo.New()
	server.HideBanner = true

	err := params.Registry.Register(server)
	if err != nil {
		return nil, err
	}

	server.Use(otelecho.Middleware(
		params.Config.GetString("app.name"),
		otelecho.WithTracerProvider(params.TracerProvider),
		otelecho.WithMeterProvider(params.MeterProvider),
		otelecho.WithPropagators(params.Propagator),
	))

	return server, nil
}

func RunServer() fx.Option {
	return fx.Invoke(
		func(
			lifecycle fx.Lifecycle,
			shutdown fx.Shutdowner,
			config *config.Config,
			logger *slog.Logger,
			server *echo.Echo,
		) {
			address := config.GetStringOrDefault("httpserver.address", ":8080")

			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := server.Start(address)
						if err != nil {

							logger.ErrorContext(ctx, "failed to start HTTP server", "error", err, "address", address)

							shutdown.Shutdown()
						}
					}()

					logger.DebugContext(ctx, "started HTTP server", "address", address)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					err := server.Shutdown(ctx)
					if err != nil {
						logger.ErrorContext(ctx, "failed to stop HTTP server", "error", err)

						return err
					}

					return nil
				},
			})
		},
	)
}
