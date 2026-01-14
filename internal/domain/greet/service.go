package greet

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/otel"
	"go.opentelemetry.io/otel/metric"
)

type GreetService struct {
	config    *config.Config
	telemetry otel.Telemetry
	counter   metric.Int64Counter
}

func NewGreetService(config *config.Config, telemetry otel.Telemetry) (*GreetService, error) {
	counter, err := telemetry.Meter().Int64Counter(
		"greet.counter",
		metric.WithDescription("Number of greets."),
		metric.WithUnit("{greet}"),
	)
	if err != nil {
		return nil, err
	}

	return &GreetService{
		config:    config,
		telemetry: telemetry,
		counter:   counter,
	}, nil
}

func (s *GreetService) Greet(ctx context.Context) string {
	ctx, span := s.telemetry.Tracer().Start(ctx, "Greet()")
	defer span.End()

	s.telemetry.Logger().DebugContext(ctx, "Greet() called")

	s.counter.Add(ctx, 1)

	return fmt.Sprintf("Greetings from %s", s.config.GetString("app.name"))
}
