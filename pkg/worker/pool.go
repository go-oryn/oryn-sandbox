package worker

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type Worker interface {
	Name() string
	Run(ctx context.Context) error
}

type Pool struct {
	running      bool
	logger       *slog.Logger
	workers      []Worker
	errGrp       *errgroup.Group
	errGrpCancel context.CancelFunc
}

func NewWorkerPool(logger *slog.Logger, workers ...Worker) *Pool {
	return &Pool{
		logger:  logger,
		workers: workers,
	}
}

func (p *Pool) Start(ctx context.Context) error {
	errGrpCtx, errGrpCancel := context.WithCancel(context.Background())
	errGrp, errGrpCtx := errgroup.WithContext(errGrpCtx)

	p.errGrp = errGrp
	p.errGrpCancel = errGrpCancel

	p.logger.DebugContext(errGrpCtx, "starting workers pool")

	for _, w := range p.workers {
		p.errGrp.Go(func() error {
			p.logger.DebugContext(errGrpCtx, "starting worker", "worker", w.Name())

			err := w.Run(errGrpCtx)
			if err != nil {
				p.logger.ErrorContext(errGrpCtx, "worker stopped with error", "worker", w.Name(), "error", err)
			} else {
				p.logger.DebugContext(errGrpCtx, "worker stopped with success", "worker", w.Name())
			}

			return err
		})
	}

	p.running = true

	select {
	case <-ctx.Done():
		p.errGrpCancel()
		p.running = false

		return ctx.Err()
	default:
		return nil
	}
}

func (p *Pool) Stop(ctx context.Context) error {
	if p.errGrp == nil || p.running == false {
		p.logger.WarnContext(ctx, "workers pool is not started")

		return nil
	}

	p.logger.DebugContext(ctx, "stopping workers pool")

	if p.errGrpCancel != nil {
		p.errGrpCancel()
	}

	done := make(chan error, 1)
	go func() {
		done <- p.errGrp.Wait()
	}()

	select {
	case err := <-done:
		p.running = false

		return err
	case <-ctx.Done():
		p.running = false

		return ctx.Err()
	}
}

func (p *Pool) Running() bool {
	return p.running
}
