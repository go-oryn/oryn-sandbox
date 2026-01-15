package healthcheck

import (
	"net/http"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
	"github.com/labstack/echo/v4"
)

type Server struct {
	httpServer *echo.Echo
}

func NewServer(config *config.Config, checker *Checker) *Server {
	server := echo.New()
	server.HideBanner = true

	server.GET(config.GetStringOrDefault("healthcheck.httpserver.path", "/health"), func(c echo.Context) error {
		verbose := c.QueryParam("verbose")

		res := checker.Check(c.Request().Context())
		if !res.Healthy() {
			if verbose == "" {
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusInternalServerError, res)
		}

		if verbose == "" {
			return c.NoContent(http.StatusOK)
		}

		return c.JSON(http.StatusOK, res)
	})

	return &Server{
		httpServer: server,
	}
}

func (s *Server) HTTPServer() *echo.Echo {
	return s.httpServer
}
