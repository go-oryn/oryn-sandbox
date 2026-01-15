package db

import (
	"context"
	"database/sql"
	"log/slog"
)

type DBProbe struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewDBProbe(logger *slog.Logger, db *sql.DB) *DBProbe {
	return &DBProbe{
		logger: logger,
		db:     db,
	}
}

func (p *DBProbe) Name() string {
	return "db"
}

func (p *DBProbe) Probe(ctx context.Context) error {
	err := p.db.PingContext(ctx)
	if err != nil {
		msg := "database ping error"

		p.logger.ErrorContext(ctx, msg, "error", err)

		return err
	}

	return nil
}
