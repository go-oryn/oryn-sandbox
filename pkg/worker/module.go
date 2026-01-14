package worker

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

const ModuleName = "worker"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		ProvidePool,
	),
)

type ProvidePoolParams struct {
	fx.In
	Logger  *slog.Logger
	Workers []Worker `group:"worker-workers"`
}

func ProvidePool(params ProvidePoolParams) *Pool {
	return NewWorkerPool(params.Logger, params.Workers...)
}

func RunWorkers() fx.Option {
	return fx.Invoke(
		func(lifecycle fx.Lifecycle, pool *Pool) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return pool.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return pool.Stop(ctx)
				},
			})
		},
	)
}
