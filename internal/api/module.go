package api

import (
	"net/http"

	"github.com/go-oryn/oryn-sandbox/internal/api/handler"
	"github.com/go-oryn/oryn-sandbox/pkg/httpserver"
	"go.uber.org/fx"
)

const ModuleName = "api"

var Module = fx.Module(
	ModuleName,
	// routes
	httpserver.AsHandler(http.MethodGet, "/greet", handler.NewGreetHandler),
)
