package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		core.Module,
		config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
		fx.Invoke(
			func(
				cfg *config.Config,
				logger *slog.Logger,
				tracer oteltrace.Tracer,
				meter metric.Meter,
				shutdown fx.Shutdowner,
			) error {
				ctx := context.Background()

				counter, _ := meter.Int64Counter(
					"tick.counter",
					metric.WithDescription("Number of ticks."),
					metric.WithUnit("{tick}"),
				)

				ticker := time.NewTicker(3 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					// trace example
					ctx, span := tracer.Start(ctx, fmt.Sprintf("span-%d", time.Now().UnixNano()))

					// log example
					logger.DebugContext(ctx, "some log debug level")
					logger.InfoContext(ctx, "some log info level")

					// metric example
					counter.Add(ctx, 1)

					span.End()
				}

				return shutdown.Shutdown()
			},
		),
	).Run()
}
