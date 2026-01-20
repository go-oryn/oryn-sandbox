package mcpserver

import (
	"go.uber.org/fx"
)

func AsCapability(constructor any) fx.Option {
	return fx.Provide(
		fx.Annotate(
			constructor,
			fx.As(new(Capability)),
			fx.ResultTags(`group:"mcpserver-capabilities"`),
		),
	)
}
