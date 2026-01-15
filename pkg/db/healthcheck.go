package db

import (
	"context"
	"database/sql"
)

type DBProbe struct {
	db *sql.DB
}

func NewDBProbe(db *sql.DB) *DBProbe {
	return &DBProbe{
		db: db,
	}
}

func (p *DBProbe) Name() string {
	return "db"
}

func (p *DBProbe) Probe(ctx context.Context) (string, error) {
	err := p.db.PingContext(ctx)
	if err != nil {
		return "database ping error", err
	}

	return "database ping success", nil
}
