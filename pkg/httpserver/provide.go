package httpserver

import (
	"reflect"

	"go.uber.org/fx"
)

type HandlerDefinition struct {
	Method string
	Path   string
	Type   reflect.Type
}

func AsHandler(method string, path string, constructor any) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor,
				fx.As(new(Handler)),
				fx.ResultTags(`group:"httpserver-handlers"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				HandlerDefinition{
					Method: method,
					Path:   path,
					Type:   reflect.TypeOf(constructor).Out(0),
				},
				fx.ResultTags(`group:"httpserver-handlers-definitions"`),
			),
		),
	)
}
