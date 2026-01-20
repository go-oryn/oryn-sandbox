package tool

import (
	"context"
	"fmt"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/go-oryn/oryn-sandbox/pkg/otel"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetInput struct {
	User string `json:"name" jsonschema:"the user to greet"`
}

type GreetOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to send to the user"`
}

type GreetTool struct {
	config    *config.Config
	telemetry otel.Telemetry
}

func NewGreetTool(config *config.Config, telemetry otel.Telemetry) *GreetTool {
	return &GreetTool{
		config:    config,
		telemetry: telemetry,
	}
}

func (t *GreetTool) Register(server *mcp.Server) error {
	fmt.Println("****************** in register")
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "greet",
			Description: "Greet a user",
		},
		t.Greet,
	)

	return nil
}

func (t *GreetTool) Greet(ctx context.Context, _ *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	ctx, span := t.telemetry.Tracer().Start(ctx, "tool.GreetTool::Greet()")
	defer span.End()

	t.telemetry.Logger().DebugContext(ctx, "MCP Greet() called!")

	return nil, GreetOutput{
		Greeting: fmt.Sprintf("Greeting, %s!", input.User),
	}, nil
}
