package log

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

func LoggerProviderOptions(
	ctx context.Context,
	cfg *config.Config,
	res *resource.Resource,
	opts ...log.LoggerProviderOption,
) ([]log.LoggerProviderOption, error) {
	// Resource option
	lpOpts := []log.LoggerProviderOption{
		log.WithResource(res),
	}

	// Processors & exporters options
	for exporter := range cfg.GetStringMap("log.exporters") {
		switch exporter {
		case "stdout":
			if cfg.GetBool("log.exporters.stdout.enabled") {
				var expOpts []stdoutlog.Option

				if cfg.GetBool("log.exporters.stdout.options.pretty_print") {
					expOpts = append(expOpts, stdoutlog.WithPrettyPrint())
				}

				if cfg.GetBool("log.exporters.stdout.options.without_timestamp") {
					expOpts = append(expOpts, stdoutlog.WithoutTimestamps())
				}

				exp, err := stdoutlog.New(expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
				}

				lpOpts = append(lpOpts, log.WithProcessor(log.NewSimpleProcessor(exp)))
			}
		case "otlp_grpc":
			if cfg.GetBool("log.exporters.otlp_grpc.enabled") {
				expOpts := []otlploggrpc.Option{
					otlploggrpc.WithEndpoint(cfg.GetString("log.exporters.otlp_grpc.options.endpoint")),
				}

				if cfg.GetBool("log.exporters.otlp_grpc.options.insecure") {
					expOpts = append(expOpts, otlploggrpc.WithInsecure())
				}

				exp, err := otlploggrpc.New(ctx, expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create otlp_grpc exporter: %w", err)
				}

				lpOpts = append(lpOpts, log.WithProcessor(log.NewBatchProcessor(exp)))
			}

		}
	}

	return append(lpOpts, opts...), nil
}
