package mcpserver

import (
	"net/http"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type StreamableHTTPServer struct {
	httpServer *http.Server
}

func NewStreamableHTTPServer(
	config *config.Config,
	mcpServer *mcp.Server,
	propagator propagation.TextMapPropagator,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
) *StreamableHTTPServer {
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return mcpServer
		},
		&mcp.StreamableHTTPOptions{
			JSONResponse: true,
			Stateless:    true,
		},
	)

	mux := http.NewServeMux()
	mux.Handle(
		config.GetStringOrDefault("mcpserver.transport.options.path", "/mcp"),
		otelhttp.NewHandler(
			handler,
			"MCP Request",
			otelhttp.WithServerName(config.GetString("app.name")),
			otelhttp.WithPropagators(propagator),
			otelhttp.WithTracerProvider(tracerProvider),
			otelhttp.WithMeterProvider(meterProvider),
		),
	)

	httpServer := &http.Server{
		Addr:    config.GetStringOrDefault("mcpserver.transport.options.address", ":8080"),
		Handler: mux,
	}

	return &StreamableHTTPServer{
		httpServer: httpServer,
	}
}

func (s *StreamableHTTPServer) HTTPServer() *http.Server {
	return s.httpServer
}
