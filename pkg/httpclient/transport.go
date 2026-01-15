package httpclient

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewTransport(options ...otelhttp.Option) *otelhttp.Transport {
	// clone http.DefaultTransport to be data race free
	transport := http.DefaultTransport.(*http.Transport).Clone()

	// boost defaults settings
	transport.MaxIdleConns = 100
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	return otelhttp.NewTransport(transport, options...)
}
