package infra

import (
	"github.com/go-oryn/oryn-sandbox/internal/mcp/tool"
	"github.com/go-oryn/oryn-sandbox/pkg/mcpserver"
	"go.uber.org/fx"
)

const ModuleName = "mcp"

var Module = fx.Module(
	ModuleName,
	// MCP tools
	mcpserver.AsCapability(tool.NewGreetTool),
)
