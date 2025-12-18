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
	"github.com/go-oryn/oryn-sandbox/pkg/core/trace"
	//"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otelsdklog "go.opentelemetry.io/otel/sdk/log"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	//traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	traceExporter, err := otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("oryn-lgtm:4317"),
	)
	if err != nil {
		panic(err)
	}

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
		trace.AsTracerProviderOptions(otelsdktrace.WithBatcher(traceExporter)),
		log.AsLoggerProviderOptions(otelsdklog.WithProcessor(otelsdklog.NewBatchProcessor(logExporter))),
		fx.Invoke(
			func(cfg *config.Config, tracer oteltrace.Tracer, logger *slog.Logger, shutdown fx.Shutdowner) error {
				ticker := time.NewTicker(3 * time.Second)
				defer ticker.Stop()

				for range ticker.C {
					ctx, span := tracer.Start(context.Background(), "main")
					span.End()

					logger.InfoContext(ctx, "some log info level")

					fmt.Printf("App name: %s, app env: %s\n", cfg.GetString("app.name"), cfg.Env())
				}

				return shutdown.Shutdown()
			},
		),
	).Run()
}
