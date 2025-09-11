package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPaymentsByPartyId(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]Payment, error) {
	query := `
	SELECT
		p.id,
		p.description AS description,
		p.amount AS amount,
		m.name AS payer_name,
		ARRAY_AGG(mm.name) AS payees
	FROM payment AS p
	LEFT JOIN member AS m
		ON p.payer_id = m.id
	LEFT JOIN member_payment AS mp
		ON mp.payment_id = p.id
	LEFT JOIN member AS mm
		ON mp.member_id = mm.id
	WHERE p.party_id = $1
	GROUP BY p.id, p.description, p.amount, m.name`
	args := []any{id}

	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return []Payment{}, err
	}

	payments, err := pgx.CollectRows(rows, pgx.RowToStructByName[Payment])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return []Payment{}, err
	}

	return payments, nil
}
