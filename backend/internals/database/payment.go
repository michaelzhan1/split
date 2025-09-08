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
		p.description AS description,
		p.amount AS amount,
		m.name AS payer_name
	FROM payment AS p
	LEFT JOIN member AS m
		ON p.payer_id = m.id
	WHERE p.id = $1`
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

	payeeQuery := `
	SELECT
		m.name AS payee_name
	FROM payment AS p
	LEFT JOIN member_payment AS mp
		ON mp.payment_id = p.id
	LEFT JOIN member AS m
		ON mp.member_id = m.id
	WHERE p.id = $1`
	payeeArgs := []any{id}

	rows, err = db.Query(ctx, payeeQuery, payeeArgs...)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return []Payment{}, err
	}

	payees, err := pgx.CollectRows(rows, pgx.RowToStructByName[[]string])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return []Payment{}, err
	}

	payments.Payees = payees
	return payments, nil
}
