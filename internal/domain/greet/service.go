package greet

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.opentelemetry.io/otel/metric"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type GreetService struct {
	config  *config.Config
	logger  *slog.Logger
	tracer  oteltrace.Tracer
	counter metric.Int64Counter
}

func NewGreetService(
	config *config.Config,
	logger *slog.Logger,
	tracer oteltrace.Tracer,
	meter metric.Meter,
) (*GreetService, error) {
	counter, err := meter.Int64Counter(
		"greet.counter",
		metric.WithDescription("Number of greets."),
		metric.WithUnit("{greet}"),
	)
	if err != nil {
		return nil, err
	}

	return &GreetService{
		config:  config,
		logger:  logger,
		tracer:  tracer,
		counter: counter,
	}, nil
}

func (s *GreetService) Greet(ctx context.Context) string {
	ctx, span := s.tracer.Start(ctx, "Greet()")
	defer span.End()

	s.logger.DebugContext(ctx, "Greet() called")

	s.counter.Add(ctx, 1)

	return fmt.Sprintf("Greetings from %s", s.config.GetString("app.name"))
}
