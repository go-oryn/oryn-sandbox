package seeds

import (
	"context"
	"database/sql"
	"sort"

	"github.com/go-oryn/oryn-sandbox/pkg/config"
)

type UsersSeed struct {
	config *config.Config
}

func NewUsersSeed(config *config.Config) *UsersSeed {
	return &UsersSeed{
		config: config,
	}
}

func (s *UsersSeed) Name() string {
	return "users"
}

func (s *UsersSeed) Run(ctx context.Context, db *sql.DB) error {
	var txErr error

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	seedData := s.config.GetStringMapString("db.seeds.users")

	names := make([]string, 0, len(seedData))
	for name := range seedData {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		_, txErr = tx.ExecContext(ctx, "INSERT INTO users (name, job) VALUES (?, ?)", name, seedData[name])
	}

	if txErr != nil {
		return tx.Rollback()
	}

	return tx.Commit()
}
