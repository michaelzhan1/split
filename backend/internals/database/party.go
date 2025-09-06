package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateParty(ctx context.Context, db *pgxpool.Pool, name string) (int, error) {
	query := "INSERT INTO party (name) VALUES ($1) RETURNING id"
	args := []any{name}
	
	var id int
	err := db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert failed: %v\n", err)
		return 0, err
	}
	return id, nil
}

func DeleteParty(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := "DELETE FROM party WHERE id = $1"
	args := []any{id}

	cmdTag, err := db.Exec(ctx, query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete failed: %v\n", err)
		return err
	}
	if cmdTag.RowsAffected() > 1 {
		fmt.Fprintf(os.Stderr, "Delete failed: more than one row affected\n")
		return errors.New("more than one row affected")
	}
	if cmdTag.RowsAffected() == 0 {
		fmt.Fprintf(os.Stderr, "Delete failed: party %v does not exist\n", id)
		return pgx.ErrNoRows
	}
	return nil
}
