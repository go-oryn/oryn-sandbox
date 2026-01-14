package otel

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry interface {
	Logger() *slog.Logger
	Meter() metric.Meter
	Tracer() trace.Tracer
}

type TelemetryWrapper struct {
	logger *slog.Logger
	meter  metric.Meter
	tracer trace.Tracer
}

func NewTelemetryWrapper(logger *slog.Logger, meter metric.Meter, tracer trace.Tracer) *TelemetryWrapper {
	return &TelemetryWrapper{
		logger: logger,
		meter:  meter,
		tracer: tracer,
	}
}

func (w *TelemetryWrapper) Logger() *slog.Logger {
	return w.logger
}

func (w *TelemetryWrapper) Meter() metric.Meter {
	return w.meter
}

func (w *TelemetryWrapper) Tracer() trace.Tracer {
	return w.tracer
}
