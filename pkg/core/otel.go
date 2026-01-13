package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	otellog "github.com/go-oryn/oryn-sandbox/pkg/otel/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/fx"
)

func ConfigureOTel() fx.Option {
	return fx.Options(
		// resource options decoration
		fx.Decorate(fx.Annotate(ConfigureOTelResourceOptions, fx.ResultTags(`group:"otel-resource-options"`))),
		// providers options decoration
		fx.Decorate(fx.Annotate(ConfigureOTelLoggerProviderOptions, fx.ResultTags(`group:"otel-log-provider-options"`))),
		fx.Decorate(fx.Annotate(ConfigureOTelMeterProviderOptions, fx.ResultTags(`group:"otel-metric-provider-options"`))),
		fx.Decorate(fx.Annotate(ConfigureOTelTracerProviderOptions, fx.ResultTags(`group:"otel-trace-provider-options"`))),
		// logger handler options decoration
		fx.Decorate(fx.Annotate(ConfigureOTelLoggerHandlerOptions, fx.ResultTags(`group:"otel-log-handler-options"`))),
		// logger handler decoration
		fx.Decorate(fx.Annotate(ConfigureOTelLoggerHandler, fx.As(new(slog.Handler)))),
	)
}

type ConfigureOTelResourceOptionsParams struct {
	fx.In
	Config  *config.Config
	Options []resource.Option `group:"otel-resource-options"`
}

func ConfigureOTelResourceOptions(params ConfigureOTelResourceOptionsParams) []resource.Option {
	return append(
		params.Options,
		resource.WithAttributes(
			semconv.ServiceName(params.Config.GetString("app.name")),
			semconv.ServiceVersion(params.Config.GetString("app.version")),
		),
	)
}

type ConfigureOTelLoggerProviderOptionsParams struct {
	fx.In
	Context context.Context
	Config  *config.Config
	Options []log.LoggerProviderOption `group:"otel-log-provider-options"`
}

func ConfigureOTelLoggerProviderOptions(params ConfigureOTelLoggerProviderOptionsParams) ([]log.LoggerProviderOption, error) {
	lpOpts := params.Options

	for exporter := range params.Config.GetStringMap("log.exporters") {
		switch exporter {
		case "stdout":
			if params.Config.GetBool("log.exporters.stdout.enabled") {
				var expOpts []stdoutlog.Option

				if params.Config.GetBool("log.exporters.stdout.options.pretty_print") {
					expOpts = append(expOpts, stdoutlog.WithPrettyPrint())
				}

				if params.Config.GetBool("log.exporters.stdout.options.without_timestamp") {
					expOpts = append(expOpts, stdoutlog.WithoutTimestamps())
				}

				exp, err := stdoutlog.New(expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create stdout log exporter: %w", err)
				}

				lpOpts = append(lpOpts, log.WithProcessor(log.NewSimpleProcessor(exp)))
			}
		case "otlp_grpc":
			if params.Config.GetBool("log.exporters.otlp_grpc.enabled") {
				expOpts := []otlploggrpc.Option{
					otlploggrpc.WithEndpoint(params.Config.GetString("log.exporters.otlp_grpc.options.endpoint")),
				}

				if params.Config.GetBool("log.exporters.otlp_grpc.options.insecure") {
					expOpts = append(expOpts, otlploggrpc.WithInsecure())
				}

				exp, err := otlploggrpc.New(params.Context, expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create otlp_grpc exporter: %w", err)
				}

				lpOpts = append(lpOpts, log.WithProcessor(log.NewBatchProcessor(exp)))
			}
		}
	}

	return lpOpts, nil
}

type ConfigureOTelMeterProviderOptionsParams struct {
	fx.In
	Context context.Context
	Config  *config.Config
	Options []metric.Option `group:"otel-metric-provider-options"`
}

func ConfigureOTelMeterProviderOptions(params ConfigureOTelMeterProviderOptionsParams) ([]metric.Option, error) {
	mpOpts := params.Options

	for exporter := range params.Config.GetStringMap("metric.exporters") {
		switch exporter {
		case "stdout":
			if params.Config.GetBool("metric.exporters.stdout.enabled") {
				var expOpts []stdoutmetric.Option

				if params.Config.GetBool("metric.exporters.stdout.options.pretty_print") {
					expOpts = append(expOpts, stdoutmetric.WithPrettyPrint())
				}

				if params.Config.GetBool("metric.exporters.stdout.options.without_timestamp") {
					expOpts = append(expOpts, stdoutmetric.WithoutTimestamps())
				}

				exp, err := stdoutmetric.New(expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create stdout metric exporter: %w", err)
				}

				mpOpts = append(mpOpts, metric.WithReader(
					metric.NewPeriodicReader(exp, metric.WithInterval(params.Config.GetDuration("metric.interval"))),
				))
			}
		case "otlp_grpc":
			if params.Config.GetBool("metric.exporters.otlp_grpc.enabled") {
				expOpts := []otlpmetricgrpc.Option{
					otlpmetricgrpc.WithEndpoint(params.Config.GetString("metric.exporters.otlp_grpc.options.endpoint")),
				}

				if params.Config.GetBool("metric.exporters.otlp_grpc.options.insecure") {
					expOpts = append(expOpts, otlpmetricgrpc.WithInsecure())
				}

				exp, err := otlpmetricgrpc.New(params.Context, expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create otlp_grpc metric exporter: %w", err)
				}

				mpOpts = append(mpOpts, metric.WithReader(
					metric.NewPeriodicReader(exp, metric.WithInterval(params.Config.GetDuration("metric.interval"))),
				))
			}
		}
	}

	return mpOpts, nil
}

type ConfigureOTelTracerProviderOptionsParams struct {
	fx.In
	Context context.Context
	Config  *config.Config
	Options []trace.TracerProviderOption `group:"otel-trace-provider-options"`
}

func ConfigureOTelTracerProviderOptions(params ConfigureOTelTracerProviderOptionsParams) ([]trace.TracerProviderOption, error) {
	tpOpts := params.Options

	for exporter := range params.Config.GetStringMap("trace.exporters") {
		switch exporter {
		case "stdout":
			if params.Config.GetBool("trace.exporters.stdout.enabled") {
				var expOpts []stdouttrace.Option

				if params.Config.GetBool("trace.exporters.stdout.options.pretty_print") {
					expOpts = append(expOpts, stdouttrace.WithPrettyPrint())
				}

				if params.Config.GetBool("trace.exporters.stdout.options.without_timestamp") {
					expOpts = append(expOpts, stdouttrace.WithoutTimestamps())
				}

				exp, err := stdouttrace.New(expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create stdout trace exporter: %w", err)
				}

				tpOpts = append(tpOpts, trace.WithSyncer(exp))
			}
		case "otlp_grpc":
			if params.Config.GetBool("trace.exporters.otlp_grpc.enabled") {
				expOpts := []otlptracegrpc.Option{
					otlptracegrpc.WithEndpoint(params.Config.GetString("trace.exporters.otlp_grpc.options.endpoint")),
				}

				if params.Config.GetBool("trace.exporters.otlp_grpc.options.insecure") {
					expOpts = append(expOpts, otlptracegrpc.WithInsecure())
				}

				exp, err := otlptracegrpc.New(params.Context, expOpts...)
				if err != nil {
					return nil, fmt.Errorf("failed to create otlp_grpc trace exporter: %w", err)
				}

				tpOpts = append(tpOpts, trace.WithBatcher(exp))
			}
		}
	}

	return tpOpts, nil
}

type ConfigureOTelLoggerHandlerOptionsParams struct {
	fx.In
	Config  *config.Config
	Options []otelslog.Option `group:"otel-log-handler-options"`
}

func ConfigureOTelLoggerHandlerOptions(params ConfigureOTelLoggerHandlerOptionsParams) []otelslog.Option {
	return append(params.Options, otelslog.WithSource(params.Config.GetBool("log.source")))
}

type ConfigureOTelLoggerHandlerParams struct {
	fx.In
	Config  *config.Config
	Handler slog.Handler
}

func ConfigureOTelLoggerHandler(params ConfigureOTelLoggerHandlerParams) *otellog.LeveledHandler {
	return otellog.NewLeveledHandler(
		otellog.ParseLogLevel(params.Config.GetString("log.level")),
		params.Handler,
	)
}
