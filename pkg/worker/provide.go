package worker

import (
	"go.uber.org/fx"
)

func AsWorker(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Worker)),
			fx.ResultTags(`group:"worker-workers"`),
		),
	)
}
