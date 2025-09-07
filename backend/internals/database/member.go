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
