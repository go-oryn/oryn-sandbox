package log

import (
	"context"
	"log/slog"
)

type LeveledHandler struct {
	level   slog.Level
	handler slog.Handler
}

func NewLeveledHandler(level slog.Level, handler slog.Handler) *LeveledHandler {
	return &LeveledHandler{
		level:   level,
		handler: handler,
	}
}

func (h *LeveledHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *LeveledHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handler.Handle(ctx, record)
}

func (h *LeveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewLeveledHandler(h.level, h.handler.WithAttrs(attrs))
}

func (h *LeveledHandler) WithGroup(name string) slog.Handler {
	return NewLeveledHandler(h.level, h.handler.WithGroup(name))
}
