package internal

import (
	"context"
	"testing"

	"github.com/go-oryn/oryn-sandbox/configs"
	"github.com/go-oryn/oryn-sandbox/internal/api"
	"github.com/go-oryn/oryn-sandbox/internal/domain"
	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/core"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"github.com/go-oryn/oryn-sandbox/pkg/otel/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/fx"
)

var Bootstrapper = core.NewBootstrapper(
	// blueprint modules
	httpserver.Module,
	// app modules
	api.Module,
	domain.Module,
	// app config
	config.AsConfigOptions(config.WithEmbedFS(configs.ConfigFS)),
	// to remove
	log.AsLoggerHandlerOptions(otelslog.WithAttributes(attribute.String("foo", "bar"))),
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
