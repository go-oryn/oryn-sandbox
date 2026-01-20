package internal

import (
	"context"
	"testing"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/db/migrations"
	internalapi "github.com/go-oryn/oryn-sandbox/internal/api"
	internaldomain "github.com/go-oryn/oryn-sandbox/internal/domain"
	internalinfra "github.com/go-oryn/oryn-sandbox/internal/infra"
	internalmcp "github.com/go-oryn/oryn-sandbox/internal/mcp"
	internalworker "github.com/go-oryn/oryn-sandbox/internal/worker"
	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/db"
	"github.com/go-oryn/oryn-sandbox/pkg/healthcheck"
	"github.com/go-oryn/oryn-sandbox/pkg/httpclient"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/go-oryn/oryn-sandbox/pkg/mcpserver"
	"github.com/go-oryn/oryn-sandbox/pkg/worker"

	"go.uber.org/fx"
)

var Bootstrapper = core.NewBootstrapper(
	// shared modules
	db.Module,
	healthcheck.Module,
	httpclient.Module,
	httpserver.Module,
	mcpserver.Module,
	worker.Module,
	// app modules
	internalapi.Module,
	internaldomain.Module,
	internalmcp.Module,
	internalinfra.Module,
	internalworker.Module,
	// app config
	config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
	// app migrations
	db.AsMigratorOptions(db.WithMigrationsEmbedFS(migrations.MigrationsFS)),
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
		// run db seeds
		db.RunSeeds(),
		// apply test options
		fx.Options(options...),
	)
}
