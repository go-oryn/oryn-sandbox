package mcpserver

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel/propagation"
)

func NewReceivingMiddleware(propagator propagation.TextMapPropagator) mcp.Middleware {
	return func(mh mcp.MethodHandler) mcp.MethodHandler {
		return func(ctx context.Context, method string, req mcp.Request) (result mcp.Result, err error) {
			params := req.GetParams()
			meta := params.GetMeta()

			fmt.Printf("*** params %+#v\n", params)
			fmt.Printf("*** meta %+#v\n", meta)

			if meta == nil {
				meta = make(map[string]any)
				params.SetMeta(meta)
			}

			propagator.Inject(ctx, stringAnyMapCarrier{meta})

			return mh(ctx, method, req)
		}

	}
}

type stringAnyMapCarrier struct {
	m map[string]any
}

func (a stringAnyMapCarrier) Get(key string) string {
	if a.m == nil {
		return ""
	}
	if s, ok := a.m[key].(string); ok {
		return s
	}
	return ""
}

// Keys implements propagation.TextMapCarrier.
func (a stringAnyMapCarrier) Keys() []string {
	if len(a.m) == 0 {
		return nil
	}
	keys := make([]string, 0, len(a.m))
	for k := range a.m {
		keys = append(keys, k)
	}
	return keys
}

// Set implements propagation.TextMapCarrier.
func (a stringAnyMapCarrier) Set(key string, value string) {
	a.m[key] = value
}
