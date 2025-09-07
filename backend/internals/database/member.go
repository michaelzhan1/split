package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetMembersByPartyId(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]Member, error) {
	query := "SELECT name FROM member WHERE member.party_id = $1"
	args := []any{id}

	rows, err := db.Query(ctx, query, args...)
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

func AddMemberToPartyById(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyId int, name string) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		query := "INSERT INTO member (party_id, name) values ($1, $2) RETURNING id"
		args := []any{partyId, name}

		var id int
		err := tx.QueryRow(ctx, query, args...).Scan(&id)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		return id, nil
	})
}
