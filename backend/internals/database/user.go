package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUsersByGroupID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]User, error) {
	query := "SELECT id, name, balance FROM users WHERE users.group_id = @id ORDER BY id"
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetUsersByGroupID", "query", query, "args", args)
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

func AddUserToGroupByID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, groupID int, name string) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		query := "INSERT INTO users (group_id, name) VALUES (@id, @name) RETURNING id"
		args := pgx.StrictNamedArgs{
			"id":   groupID,
			"name": name,
		}

		var id int
		L.Info("AddUserToGroupByID", "query", query, "args", args)
		err := tx.QueryRow(ctx, query, args).Scan(&id)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		return id, nil
	})
}

func PatchUser(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, groupID int, userID int, name string) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "UPDATE users SET name = @name WHERE id = @id AND group_id = @groupID"
		args := pgx.StrictNamedArgs{
			"name":    name,
			"id":      userID,
			"groupID": groupID,
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

func DeleteUser(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, groupID int, userID int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "DELETE FROM users WHERE id = @id AND group_id = @groupID"
		args := pgx.StrictNamedArgs{
			"id":      userID,
			"groupID": groupID,
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
