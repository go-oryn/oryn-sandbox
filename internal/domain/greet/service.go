package greet

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/otel"
	"go.opentelemetry.io/otel/metric"
)

type Service struct {
	client    *http.Client
	repo      *Repository
	config    *config.Config
	telemetry otel.Telemetry
	counter   metric.Int64Counter
}

func NewService(client *http.Client, repo *Repository, config *config.Config, telemetry otel.Telemetry) (*Service, error) {
	counter, err := telemetry.Meter().Int64Counter(
		"greet.counter",
		metric.WithDescription("Number of greets."),
		metric.WithUnit("{greet}"),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:    client,
		repo:      repo,
		config:    config,
		telemetry: telemetry,
		counter:   counter,
	}, nil
}

func (s *Service) Greet(ctx context.Context) string {
	ctx, span := s.telemetry.Tracer().Start(ctx, "Greet()")
	defer span.End()

	s.telemetry.Logger().DebugContext(ctx, "Greet() called")

	s.counter.Add(ctx, 1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com", nil)
	if err != nil {
		s.telemetry.Logger().ErrorContext(ctx, "cannot prepare http request", err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		s.telemetry.Logger().ErrorContext(ctx, "cannot send http request", err)
	}
	defer res.Body.Close()

	dbTime, err := s.repo.Time(ctx)
	if err != nil {
		s.telemetry.Logger().ErrorContext(ctx, "cannot retrieve db time", "error", err)

		return fmt.Sprintf(
			"Greetings from %s.",
			s.config.GetString("app.name"),
		)
	}

	return fmt.Sprintf(
		"Greetings from %s, it is now %s on the db.",
		s.config.GetString("app.name"),
		dbTime.Format(time.RFC3339),
	)
}
