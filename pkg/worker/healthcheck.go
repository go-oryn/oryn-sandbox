package worker

import (
	"context"
	"errors"
)

type WorkersProbe struct {
	pool *Pool
}

func NewWorkersProbe(pool *Pool) *WorkersProbe {
	return &WorkersProbe{
		pool: pool,
	}
}

func (p *WorkersProbe) Name() string {
	return "workers"
}

func (p *WorkersProbe) Probe(ctx context.Context) (string, error) {
	if !p.pool.Running() {
		msg := "workers are not running"

		return msg, errors.New(msg)
	}

	return "workers are running", nil
}
