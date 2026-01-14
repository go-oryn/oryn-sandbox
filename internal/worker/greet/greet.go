package greet

import (
	"context"
	"time"

	"github.com/go-oryn/oryn-sandbox/internal/domain/greet"
)

type GreetWorker struct {
	service *greet.GreetService
}

func NewGreetWorker(service *greet.GreetService) *GreetWorker {
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

			time.Sleep(1 * time.Second)
		}
	}
}
