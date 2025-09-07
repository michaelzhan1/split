package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPartyById(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) (Party, error) {
	query := "SELECT name FROM party WHERE party.id = $1"
	args := []any{id}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return Party{}, err
	}

	party, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Party])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return Party{}, err
	}

	return party, nil
}

func CreateParty(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, name string) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		query := "INSERT INTO party (name) VALUES ($1) RETURNING id"
		args := []any{name}

		var id int
		err := tx.QueryRow(ctx, query, args...).Scan(&id)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		return id, nil
	})
}

func PatchParty(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int, name string) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "UPDATE party SET name = $1 WHERE id = $2"
		args := []any{name, id}

		cmdTag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			L.Error(fmt.Sprintf("Patch failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() > 1 {
			L.Error("Patch failed: more than one row affected")
			return struct{}{}, errors.New("more than one row affected")
		}
		if cmdTag.RowsAffected() == 0 {
			L.Error(fmt.Sprintf("Patch failed: party %v does not exist", id))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})
	return err
}

func DeleteParty(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "DELETE FROM party WHERE id = $1"
		args := []any{id}

		cmdTag, err := tx.Exec(ctx, query, args...)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() > 1 {
			L.Error("Delete failed: more than one row affected")
			return struct{}{}, errors.New("more than one row affected")
		}
		if cmdTag.RowsAffected() == 0 {
			L.Error(fmt.Sprintf("Delete failed: party %v does not exist", id))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})

	return err
}
