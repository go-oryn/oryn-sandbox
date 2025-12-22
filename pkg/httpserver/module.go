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
		ProvideHTTPServer,
	),
)

type ProvideHTTPServerParams struct {
	fx.In
	Lifecycle      fx.Lifecycle
	Shutdown       fx.Shutdowner
	Config         *config.Config
	Logger         *slog.Logger
	Propagator     propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
}

func ProvideHTTPServer(params ProvideHTTPServerParams) (*echo.Echo, error) {
	server := echo.New()
	server.HideBanner = true

	server.Use(otelecho.Middleware(
		params.Config.GetString("app.name"),
		otelecho.WithTracerProvider(params.TracerProvider),
		otelecho.WithMeterProvider(params.MeterProvider),
		otelecho.WithPropagators(params.Propagator),
	))

	logger := params.Logger.With("module", ModuleName)
	serverAddress := params.Config.GetStringOrDefault("httpserver.address", ":8080")

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := server.Start(serverAddress)
				if err != nil {

					logger.ErrorContext(ctx, "failed to start HTTP server", "error", err, "address", serverAddress)

					params.Shutdown.Shutdown()
				}
			}()

			logger.DebugContext(ctx, "started HTTP server", "address", serverAddress)

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

	return server, nil
}
