package db

import (
	"context"
	"database/sql"
	"embed"
	"log/slog"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	logger  *slog.Logger
	embedFS *embed.FS
	db      *sql.DB
}

func NewMigrator(logger *slog.Logger, db *sql.DB, options ...MigratorOption) *Migrator {
	mOpts := &MigratorOptions{}

	for _, opt := range options {
		opt(mOpts)
	}

	return &Migrator{
		logger:  logger,
		embedFS: mOpts.embedFS,
		db:      db,
	}
}

func (m *Migrator) Run(ctx context.Context, dialect string, command string, args ...string) error {
	// set base BF
	if m.embedFS != nil {
		goose.SetBaseFS(m.embedFS)
	}
	// set dialect
	err := goose.SetDialect(dialect)
	if err != nil {
		m.logger.ErrorContext(ctx, "db migration dialect error", "dialect", dialect, "error", err)

		return err
	}

	// run migration
	err = goose.RunContext(ctx, command, m.db, ".", args...)
	if err != nil {
		m.logger.ErrorContext(ctx, "db migration error", "dialect", dialect, "command", command, "error", err)

		return err
	}

	m.logger.InfoContext(ctx, "db migration success", "dialect", dialect, "command", command)

	return nil
}
