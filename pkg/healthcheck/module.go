package healthcheck

import (
	"context"
	"log/slog"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"go.uber.org/fx"
)

const ModuleName = "healthcheck"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideChecker,
		ProvideServer,
	),
)

type ProvideCheckerParams struct {
	fx.In
	Logger *slog.Logger
	Probes []Probe `group:"healthcheck-probes"`
}

func ProvideChecker(params ProvideCheckerParams) *Checker {
	return NewChecker(params.Probes...)
}

type ProvideServerParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Shutdown  fx.Shutdowner
	Config    *config.Config
	Checker   *Checker
}

func ProvideServer(params ProvideServerParams) (*Server, error) {
	return NewServer(params.Config, params.Checker), nil
}

func RunServer() fx.Option {
	return fx.Invoke(
		func(
			lifecycle fx.Lifecycle,
			shutdown fx.Shutdowner,
			config *config.Config,
			logger *slog.Logger,
			server *Server,
		) {
			address := config.GetStringOrDefault("healthcheck.httpserver.address", ":8080")

			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := server.HTTPServer().Start(address)
						if err != nil {

							logger.ErrorContext(ctx, "failed to start healthcheck HTTP server", "error", err, "address", address)

							shutdown.Shutdown()
						}
					}()

					logger.DebugContext(ctx, "started healthcheck HTTP server", "address", address)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					err := server.HTTPServer().Shutdown(ctx)
					if err != nil {
						logger.ErrorContext(ctx, "failed to stop healthcheck HTTP server", "error", err)

						return err
					}

					return nil
				},
			})
		},
	)
}
