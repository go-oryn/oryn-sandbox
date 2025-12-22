package httpserver

import "github.com/labstack/echo/v4"

type Handler interface {
	Handle() echo.HandlerFunc
}

type Registry struct{}
