package internal

import (
	"context"
	"testing"

	"github.com/go-oryn/oryn-sandbox/configs"
	internalapi "github.com/go-oryn/oryn-sandbox/internal/api"
	internaldomain "github.com/go-oryn/oryn-sandbox/internal/domain"
	internalworker "github.com/go-oryn/oryn-sandbox/internal/worker"
	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/db"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/go-oryn/oryn-sandbox/pkg/worker"

	"go.uber.org/fx"
)

var Bootstrapper = core.NewBootstrapper(
	// shared modules
	httpserver.Module,
	worker.Module,
	db.Module,
	// app modules
	internalapi.Module,
	internaldomain.Module,
	internalworker.Module,
	// app config
	config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
)

func Run(ctx context.Context, options ...fx.Option) {
	Bootstrapper.WithContext(ctx).RunApp(
		fx.Options(options...),
	)
}

func RunTest(tb testing.TB, options ...fx.Option) func() {
	tb.Helper()

	return Bootstrapper.WithContext(tb.Context()).RunTestApp(
		tb,
		fx.Options(options...),
	)
}
