package metric

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/core/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func MeterProviderOptions(
	ctx context.Context,
	cfg *config.Config,
	res *resource.Resource,
	opts ...metric.Option,
) ([]metric.Option, error) {
	// Resource option
	mpOpts := []metric.Option{
		metric.WithResource(res),
	}

	// Processors & exporters options
	for exporter := range cfg.GetStringMap("metric.exporters") {
		switch exporter {
		case "stdout":
			if cfg.GetBool("metric.exporters.stdout.enabled") {
				var expOpts []stdoutmetric.Option

				if cfg.GetBool("metric.exporters.stdout.options.pretty_print") {
					expOpts = append(expOpts, stdoutmetric.WithPrettyPrint())
				}

				if cfg.GetBool("metric.exporters.stdout.options.without_timestamp") {
					expOpts = append(expOpts, stdoutmetric.WithoutTimestamps())
				}

				exp, err := stdoutmetric.New(expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create stdout metric exporter: %w", err)
				}

				mpOpts = append(mpOpts, metric.WithReader(
					metric.NewPeriodicReader(exp, metric.WithInterval(cfg.GetDuration("metric.interval"))),
				))
			}
		case "otlp_grpc":
			if cfg.GetBool("metric.exporters.otlp_grpc.enabled") {
				expOpts := []otlpmetricgrpc.Option{
					otlpmetricgrpc.WithEndpoint(cfg.GetString("metric.exporters.otlp_grpc.options.endpoint")),
				}

				if cfg.GetBool("metric.exporters.otlp_grpc.options.insecure") {
					expOpts = append(expOpts, otlpmetricgrpc.WithInsecure())
				}

				exp, err := otlpmetricgrpc.New(ctx, expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create otlp_grpc metric exporter: %w", err)
				}

				mpOpts = append(mpOpts, metric.WithReader(
					metric.NewPeriodicReader(exp, metric.WithInterval(cfg.GetDuration("metric.interval"))),
				))
			}

		}
	}

	return append(mpOpts, opts...), nil
}
