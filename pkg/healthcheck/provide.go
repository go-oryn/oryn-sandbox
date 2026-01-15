package healthcheck

import (
	"go.uber.org/fx"
)

func AsProbe(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Probe)),
			fx.ResultTags(`group:"healthcheck-probes"`),
		),
	)
}
