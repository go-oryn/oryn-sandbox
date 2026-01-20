package mcpserver

import (
	"context"
	"log/slog"
	"net"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

const ModuleName = "mcpserver"

var Module = fx.Module(
	ModuleName,
	fx.Provide(
		ProvideRegistry,
		ProvideServer,
		ProvideStreamableHTTPServer,
	),
)

type ProvideRegistryParams struct {
	fx.In
	Capabilities []Capability `group:"mcpserver-capabilities"`
}

func ProvideRegistry(params ProvideRegistryParams) *Registry {
	return NewRegistry(params.Capabilities...)
}

type ProvideServerParams struct {
	fx.In
	Logger   *slog.Logger
	Config   *config.Config
	Registry *Registry
}

func ProvideServer(params ProvideServerParams) (*mcp.Server, error) {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    params.Config.GetString("app.name"),
			Version: params.Config.GetString("app.version"),
		},
		&mcp.ServerOptions{
			Logger: params.Logger,
			Capabilities: &mcp.ServerCapabilities{
				Tools: &mcp.ToolCapabilities{},
			},
		},
	)

	err := params.Registry.Register(server)
	if err != nil {
		return nil, err
	}

	return server, nil
}

type ProvideStreamableHTTPHandlerParams struct {
	fx.In
	Config         *config.Config
	Server         *mcp.Server
	Propagator     propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
}

func ProvideStreamableHTTPServer(params ProvideStreamableHTTPHandlerParams) *StreamableHTTPServer {
	return NewStreamableHTTPServer(
		params.Config,
		params.Server,
		params.Propagator,
		params.TracerProvider,
		params.MeterProvider,
	)
}

func RunStreamableHTTPServer() fx.Option {
	return fx.Invoke(
		func(
			lifecycle fx.Lifecycle,
			shutdown fx.Shutdowner,
			logger *slog.Logger,
			server *StreamableHTTPServer,
		) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					addr := server.HTTPServer().Addr

					lis, err := net.Listen("tcp", addr)
					if err != nil {
						logger.ErrorContext(ctx, "failed to start MCP streamable HTTP server listener",
							"error", err,
							"address", addr,
						)

						return err
					}

					go func() {
						err := server.HTTPServer().Serve(lis)
						if err != nil {
							logger.ErrorContext(ctx, "failed to start MCP streamable HTTP server",
								"error", err,
								"address", addr,
							)

							shutdown.Shutdown()
						}
					}()

					logger.DebugContext(ctx, "started MCP streamable HTTP server", "address", addr)

					return nil
				},
				OnStop: func(ctx context.Context) error {
					err := server.HTTPServer().Shutdown(ctx)
					if err != nil {
						logger.ErrorContext(ctx, "failed to stop MCP streamable HTTP server", "error", err)

						return err
					}

					return nil
				},
			})
		},
	)
}
