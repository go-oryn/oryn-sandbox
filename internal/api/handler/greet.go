package handler

import (
	"net/http"

	"github.com/go-oryn/oryn-sandbox/internal/domain/greet"
	"github.com/labstack/echo/v4"
)

type GreetHandler struct {
	service *greet.GreetService
}

func NewGreetHandler(service *greet.GreetService) *GreetHandler {
	return &GreetHandler{
		service: service,
	}
}

func (h *GreetHandler) Handle() (echo.HandlerFunc, error) {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, h.service.Greet(c.Request().Context()))
	}, nil
}
