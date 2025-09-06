package database

import (
	"context"
	"fmt"
	"os"

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
