package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateParty(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, name string) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	query := "INSERT INTO party (name) VALUES ($1) RETURNING id"
	args := []any{name}

	var id int
	err = tx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		L.Error(fmt.Sprintf("Insert failed: %v", err))
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		L.Error(fmt.Sprintf("Commit failed: %v\n", err))
		return 0, err
	}

	return id, nil
}

func DeleteParty(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	query := "DELETE FROM party WHERE id = $1"
	args := []any{id}

	cmdTag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		L.Error(fmt.Sprintf("Delete failed: %v", err))
		return err
	}
	if cmdTag.RowsAffected() > 1 {
		L.Error("Delete failed: more than one row affected")
		return errors.New("more than one row affected")
	}
	if cmdTag.RowsAffected() == 0 {
		L.Error(fmt.Sprintf("Delete failed: party %v does not exist", id))
		return pgx.ErrNoRows
	}

	err = tx.Commit(ctx)
	if err != nil {
		L.Error(fmt.Sprintf("Commit failed: %v\n", err))
		return err
	}
	return nil
}
