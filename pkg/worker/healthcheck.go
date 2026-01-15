package worker

import (
	"context"
	"errors"
	"log/slog"
)

type WorkersProbe struct {
	logger *slog.Logger
	pool   *Pool
}

func NewWorkersProbe(logger *slog.Logger, pool *Pool) *WorkersProbe {
	return &WorkersProbe{
		logger: logger,
		pool:   pool,
	}
}

func (p *WorkersProbe) Name() string {
	return "workers"
}

func (p *WorkersProbe) Probe(ctx context.Context) error {
	if !p.pool.Running() {
		msg := "worker pool is not running"

		p.logger.ErrorContext(ctx, msg)

		return errors.New(msg)
	}

	return nil
}
