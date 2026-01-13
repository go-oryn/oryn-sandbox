package core

import (
	"context"
	"testing"

	config2 "github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type Bootstrapper struct {
	context context.Context
	options []fx.Option
}

func NewBootstrapper(options ...fx.Option) *Bootstrapper {

	return &Bootstrapper{
		context: context.Background(),
		options: append([]fx.Option{Module}, options...),
	}
}

func (b *Bootstrapper) WithContext(ctx context.Context) *Bootstrapper {
	b.context = ctx

	return b
}

func (b *Bootstrapper) WithOptions(options ...fx.Option) *Bootstrapper {
	b.options = append(b.options, options...)

	return b
}

func (b *Bootstrapper) BootstrapApp(options ...fx.Option) *fx.App {
	return fx.New(
		fx.Supply(fx.Annotate(b.context, fx.As(new(context.Context)))),
		fx.Options(b.options...),
		fx.Options(options...),
	)
}

func (b *Bootstrapper) BootstrapTestApp(tb testing.TB, options ...fx.Option) *fxtest.App {
	tb.Helper()

	return fxtest.New(
		tb,
		config2.AsConfigOptions(config2.WithEnvironment("test")),
		fx.Supply(fx.Annotate(tb.Context(), fx.As(new(context.Context)))),
		fx.Options(b.options...),
		fx.Options(options...),
	)
}

func (b *Bootstrapper) RunApp(options ...fx.Option) {
	app := b.BootstrapApp(options...)

	app.Run()
}

func (b *Bootstrapper) RunTestApp(tb testing.TB, options ...fx.Option) func() {
	tb.Helper()

	app := b.BootstrapTestApp(tb, options...)

	app.RequireStart()

	return app.RequireStop
}
