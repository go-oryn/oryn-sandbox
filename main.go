package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Module,
		httpserver.Module,
		config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
		fx.Invoke(
			func(
				cfg *config.Config,
				logger *slog.Logger,
				tracer oteltrace.Tracer,
				meter metric.Meter,
				server *echo.Echo,
			) error {
				counter, err := meter.Int64Counter(
					"tick.counter",
					metric.WithDescription("Number of ticks."),
					metric.WithUnit("{tick}"),
				)
				if err != nil {
					return err
				}

				server.GET("/", func(c echo.Context) error {
					// trace example
					ctx, span := tracer.Start(c.Request().Context(), fmt.Sprintf("span-%d", time.Now().UnixNano()))
					defer span.End()

					// log example
					logger.DebugContext(ctx, "some debug")
					logger.InfoContext(ctx, "some info")

					// metric example
					counter.Add(ctx, 1)

					// response example
					return c.String(200, fmt.Sprintf("Hello, World from %s", cfg.GetString("app.name")))
				})

				return nil
			},
		),
	).Run()
}
