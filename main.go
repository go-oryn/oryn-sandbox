package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"github.com/go-oryn/oryn-sandbox/pkg/core/log"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otelsdklog "go.opentelemetry.io/otel/sdk/log"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	logExporter, err := otlploggrpc.New(
		context.Background(),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithEndpoint("oryn-lgtm:4317"),
	)
	if err != nil {
		panic(err)
	}

	fx.New(
		core.Module,
		config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
		log.AsLoggerProviderOptions(otelsdklog.WithProcessor(otelsdklog.NewBatchProcessor(logExporter))),
		fx.Invoke(
			func(cfg *config.Config, tracer oteltrace.Tracer, logger *slog.Logger, shutdown fx.Shutdowner) error {

				ticker := time.NewTicker(3 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					ctx, span := tracer.Start(context.Background(), fmt.Sprintf("span-%d", time.Now().UnixNano()))
					span.End()

					logger.DebugContext(ctx, "some log debug level")
					logger.InfoContext(ctx, "some log info level")

					fmt.Printf("App name: %s, app env: %s\n", cfg.GetString("app.name"), cfg.Env())
				}

				return shutdown.Shutdown()
			},
		),
	).Run()
}
