package httpserver

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Handle() (echo.HandlerFunc, error)
}

type Registry struct {
	logger              *slog.Logger
	handlers            []Handler
	handlersDefinitions []HandlerDefinition
}

func NewRegistry(logger *slog.Logger, handlers []Handler, handlersDefinitions []HandlerDefinition) *Registry {
	return &Registry{
		logger:              logger,
		handlers:            handlers,
		handlersDefinitions: handlersDefinitions,
	}
}

func (r *Registry) Register(srv *echo.Echo) error {
	for _, handlerDefinition := range r.handlersDefinitions {
		handler, err := r.lookupHandlerFromDefinition(handlerDefinition)
		if err != nil {
			return err
		}

		handlerFunc, err := handler.Handle()
		if err != nil {
			return fmt.Errorf("cannot register handler type %s func %w", handlerDefinition.Type, err)
		}

		r.logger.Debug("registered handler with type", "type", handlerDefinition.Type)
		srv.Add(handlerDefinition.Method, handlerDefinition.Path, handlerFunc)
	}

	return nil
}

func (r *Registry) lookupHandlerFromDefinition(definition HandlerDefinition) (Handler, error) {
	for _, handler := range r.handlers {
		if reflect.TypeOf(handler) == definition.Type {
			return handler, nil
		}
	}

	return nil, fmt.Errorf("cannot find handler type %s", definition.Type)
}
