package db

import (
	"context"
	"database/sql"
	"log/slog"
)

type Seed interface {
	Name() string
	Run(ctx context.Context, db *sql.DB) error
}

type Seeder struct {
	logger *slog.Logger
	db     *sql.DB
	seeds  []Seed
}

func NewSeeder(logger *slog.Logger, db *sql.DB, seeds ...Seed) *Seeder {
	return &Seeder{
		logger: logger,
		db:     db,
		seeds:  seeds,
	}
}

func (m *Seeder) Run(ctx context.Context, names ...string) error {
	var seedsToExecute []Seed

	if len(names) == 0 {
		seedsToExecute = m.seeds
	} else {
		for _, name := range names {
			for _, seed := range m.seeds {
				if name == seed.Name() {
					seedsToExecute = append(seedsToExecute, seed)
				}
			}
		}
	}

	for _, seedToExecute := range seedsToExecute {
		err := seedToExecute.Run(ctx, m.db)
		if err != nil {
			m.logger.ErrorContext(ctx, "seed failure", "seed", seedToExecute.Name(), "error", err)

			return err
		}

		m.logger.DebugContext(ctx, "seed success", "seed", seedToExecute.Name())
	}

	return nil
}
