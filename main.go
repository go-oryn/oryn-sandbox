package main

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"github.com/go-oryn/oryn-sandbox/pkg/core/trace"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	otelsdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func main() {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	}

	fx.New(
		core.Module,
		config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
		trace.AsTracerProviderOptions(otelsdktrace.WithBatcher(traceExporter)),
		fx.Invoke(func(cfg *config.Config, tracer oteltrace.Tracer, shutdown fx.Shutdowner) error {
			_, span := tracer.Start(context.Background(), "main")
			defer span.End()

			fmt.Printf("App name: %s, app env: %s\n", cfg.GetString("app.name"), cfg.Env())

			return shutdown.Shutdown()
		}),
	).Run()
}
