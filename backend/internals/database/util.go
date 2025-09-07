package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WithTx[T any](ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) (T, error)) (res T, err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		var zero T
		return zero, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	res, err = fn(tx)
	if err != nil {
		return res, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}