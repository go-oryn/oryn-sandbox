package trace

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TracerProviderOptions(
	ctx context.Context,
	cfg *config.Config,
	res *resource.Resource,
	opts ...trace.TracerProviderOption,
) ([]trace.TracerProviderOption, error) {
	// Resource option
	tpOpts := []trace.TracerProviderOption{
		trace.WithResource(res),
	}

	// Processors & exporters options
	for exporter := range cfg.GetStringMap("trace.exporters") {
		switch exporter {
		case "stdout":
			if cfg.GetBool("trace.exporters.stdout.enabled") {
				var expOpts []stdouttrace.Option

				if cfg.GetBool("trace.exporters.stdout.options.pretty_print") {
					expOpts = append(expOpts, stdouttrace.WithPrettyPrint())
				}

				if cfg.GetBool("trace.exporters.stdout.options.without_timestamp") {
					expOpts = append(expOpts, stdouttrace.WithoutTimestamps())
				}

				exp, err := stdouttrace.New(expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
				}

				tpOpts = append(tpOpts, trace.WithSyncer(exp))
			}
		case "otlp_grpc":
			if cfg.GetBool("trace.exporters.otlp_grpc.enabled") {
				expOpts := []otlptracegrpc.Option{
					otlptracegrpc.WithEndpoint(cfg.GetString("trace.exporters.otlp_grpc.options.endpoint")),
				}

				if cfg.GetBool("trace.exporters.otlp_grpc.options.insecure") {
					expOpts = append(expOpts, otlptracegrpc.WithInsecure())
				}

				exp, err := otlptracegrpc.New(ctx, expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create otlp_grpc trace exporter: %w", err)
				}

				tpOpts = append(tpOpts, trace.WithBatcher(exp))
			}

		}
	}

	return append(tpOpts, opts...), nil
}
