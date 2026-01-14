package greet

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Time(ctx context.Context) (time.Time, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT CURRENT_TIMESTAMP`)
	if err != nil {
		return time.Time{}, err
	}
	defer rows.Close()

	var dbTime time.Time
	for rows.Next() {
		err = rows.Scan(&dbTime)
		if err != nil {
			return time.Time{}, err
		}
	}

	if err = rows.Err(); err != nil {
		return time.Time{}, err
	}

	return dbTime, nil
}
