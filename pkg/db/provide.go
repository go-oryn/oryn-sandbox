package db

import (
	"go.uber.org/fx"
)

func AsMigratorOptions(options ...MigratorOption) fx.Option {
	fxOptions := []fx.Option{}

	for _, opt := range options {
		fxOptions = append(fxOptions, fx.Supply(
			fx.Annotate(
				opt,
				fx.ResultTags(`group:"db-migrator-options"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}

func AsSeeds(constructors ...any) fx.Option {
	fxOptions := []fx.Option{}

	for _, constructor := range constructors {
		fxOptions = append(fxOptions, fx.Provide(
			fx.Annotate(
				constructor,
				fx.As(new(Seed)),
				fx.ResultTags(`group:"db-seeder-seeds"`),
			),
		))
	}

	return fx.Options(fxOptions...)
}
