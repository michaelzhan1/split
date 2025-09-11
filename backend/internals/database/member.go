package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetMembersByPartyID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]Member, error) {
	query := "SELECT id, name, balance FROM member WHERE member.party_id = @id"
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	rows, err := db.Query(ctx, query, args)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return []Member{}, err
	}

	members, err := pgx.CollectRows(rows, pgx.RowToStructByName[Member])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return []Member{}, err
	}

	return members, nil
}

func AddMemberToPartyByID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int, name string) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		query := "INSERT INTO member (party_id, name) VALUES (@id, @name) RETURNING id"
		args := pgx.StrictNamedArgs{
			"id":   partyID,
			"name": name,
		}

		var id int
		err := tx.QueryRow(ctx, query, args).Scan(&id)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		return id, nil
	})
}

func PatchMember(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int, memberID int, name string) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "UPDATE member SET name = @name WHERE id = @id AND party_id = @partyID"
		args := pgx.StrictNamedArgs{
			"name":    name,
			"id":      memberID,
			"partyID": partyID,
		}

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
			L.Error(fmt.Sprintf("Patch failed: member %v does not exist", memberID))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})
	return err
}

func DeleteMember(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int, memberID int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		query := "DELETE FROM member WHERE id = @id AND party_id = @partyID"
		args := pgx.StrictNamedArgs{
			"id":      memberID,
			"partyID": partyID,
		}

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
			L.Error(fmt.Sprintf("Delete failed: member %v does not exist", memberID))
			return struct{}{}, pgx.ErrNoRows
		}

		return struct{}{}, nil
	})

	return err
}
