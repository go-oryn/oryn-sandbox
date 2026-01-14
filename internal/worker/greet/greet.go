package greet

import (
	"context"
	"time"

	"github.com/go-oryn/oryn-sandbox/internal/domain/greet"
	"github.com/go-oryn/oryn-sandbox/pkg/worker"
)

var _ worker.Worker = (*GreetWorker)(nil)

type GreetWorker struct {
	service *greet.Service
}

func NewGreetWorker(service *greet.Service) *GreetWorker {
	return &GreetWorker{
		service: service,
	}
}

func (w *GreetWorker) Name() string {
	return "greet"
}

func (w *GreetWorker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			w.service.Greet(ctx)

			time.Sleep(10 * time.Second)
		}
	}
}
