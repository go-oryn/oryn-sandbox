package mcpserver

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Capability interface {
	Register(mcpServer *mcp.Server) error
}

type Registry struct {
	capabilities []Capability
}

func NewRegistry(capabilities ...Capability) *Registry {
	return &Registry{
		capabilities: capabilities,
	}
}

func (r *Registry) Register(server *mcp.Server) error {
	for _, capability := range r.capabilities {
		err := capability.Register(server)
		if err != nil {
			return fmt.Errorf("cannot register MCP capability: %w", err)
		}
	}

	return nil
}
