package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-oryn/oryn-sandbox/internal"
	config2 "github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
)

func TestGreetHandler(t *testing.T) {
	var server *echo.Echo

	stop := internal.RunTest(
		t,
		config2.AsConfigOptions(config2.WithValues(map[string]any{
			"app.name": "test app",
		})),
		fx.Populate(&server),
	)
	defer stop()

	req := httptest.NewRequest(http.MethodGet, "/greet", nil)
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Greetings from test app")
}
