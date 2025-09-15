package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUsersByPartyID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]User, error) {
	query := "SELECT id, name, balance FROM user WHERE user.party_id = @id"
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetUsersByPartyID", "query", query, "args", args)
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return []User{}, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return []User{}, err
	}

	return users, nil
}

func AddUserToPartyByID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int, name string) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		query := "INSERT INTO user (party_id, name) VALUES (@id, @name) RETURNING id"
		args := pgx.StrictNamedArgs{
			"id":   partyID,
			"name": name,
		}

		var id int
		L.Info("AddUserToPartyByID", "query", query, "args", args)
		err := tx.QueryRow(ctx, query, args).Scan(&id)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		return id, nil
	})
}

func PatchUser(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int, userID int, name string) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "UPDATE user SET name = @name WHERE id = @id AND party_id = @partyID"
		args := pgx.StrictNamedArgs{
			"name":    name,
			"id":      userID,
			"partyID": partyID,
		}

		L.Info("PatchUser", "query", query, "args", args)
		cmdTag, err := tx.Exec(ctx, query, args)
		if err != nil {
			L.Error(fmt.Sprintf("Patch failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() > 1 {
			L.Error("Patch failed: more than one row affected")
			return struct{}{}, errors.New("more than one row affected")
		}
		if cmdTag.RowsAffected() == 0 {
			L.Error(fmt.Sprintf("Patch failed: user %v does not exist", userID))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})
	return err
}

func DeleteUser(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int, userID int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "DELETE FROM user WHERE id = @id AND party_id = @partyID"
		args := pgx.StrictNamedArgs{
			"id":      userID,
			"partyID": partyID,
		}

		L.Info("DeleteUser", "query", query, "args", args)
		cmdTag, err := tx.Exec(ctx, query, args)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() > 1 {
			L.Error("Delete failed: more than one row affected")
			return struct{}{}, errors.New("more than one row affected")
		}
		if cmdTag.RowsAffected() == 0 {
			L.Error(fmt.Sprintf("Delete failed: user %v does not exist", userID))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})

	return err
}
