package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetGroupByID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) (Group, error) {
	query := "SELECT id, name FROM groups WHERE groups.id = @id"
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetGroupByID", "query", query, "args", args)
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return Group{}, err
	}

	group, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Group])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return Group{}, err
	}

	return group, nil
}

func CreateGroup(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, name string) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		query := "INSERT INTO groups (name) VALUES (@name) RETURNING id"
		args := pgx.StrictNamedArgs{
			"name": name,
		}

		var id int
		L.Info("CreateGroup", "query", query, "args", args)
		err := tx.QueryRow(ctx, query, args).Scan(&id)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		return id, nil
	})
}

func PatchGroup(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int, name string) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "UPDATE groups SET name = @name WHERE id = @id"
		args := pgx.StrictNamedArgs{
			"name": name,
			"id":   id,
		}

		L.Info("PatchGroup", "query", query, "args", args)
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
			L.Error(fmt.Sprintf("Patch failed: group %v does not exist", id))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})
	return err
}

func DeleteGroup(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		paymentQuery := "DELETE FROM payment WHERE group_id = @id"
		paymentArgs := pgx.StrictNamedArgs{
			"id": id,
		}

		L.Info("DeleteGroup.DeletePayments", "query", paymentQuery, "args", paymentArgs)
		_, err := tx.Exec(ctx, paymentQuery, paymentArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}

		query := "DELETE FROM groups WHERE id = @id"
		args := pgx.StrictNamedArgs{
			"id": id,
		}

		L.Info("DeleteGroup.DeleteGroup", "query", query, "args", args)
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
			L.Error(fmt.Sprintf("Delete failed: group %v does not exist", id))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})

	return err
}
