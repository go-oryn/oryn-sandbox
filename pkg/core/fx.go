package core

import (
	"log/slog"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func ConfigureFx() fx.Option {
	return fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
		return &fxevent.SlogLogger{
			Logger: logger,
		}
	})
}
