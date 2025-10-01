package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPaymentsByGroupID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]Payment, error) {
	query := `
	SELECT
		p.id,
		p.description         AS description,
		p.amount              AS amount,
		u.id                  AS payer_id,
		u.name                AS payer_name,
		u.balance             AS payer_balance,
		ARRAY_AGG(uu.id)      AS payee_ids,
		ARRAY_AGG(uu.name)    AS payee_names,
		ARRAY_AGG(uu.balance) AS payee_balances
	FROM payment AS p
	LEFT JOIN users AS u
		ON p.payer_id = u.id
	LEFT JOIN users_payment AS up
		ON up.payment_id = p.id
	LEFT JOIN users AS uu
		ON up.user_id = uu.id
	WHERE p.group_id = @id
	GROUP BY p.id, p.description, p.amount, u.name, u.id, u.balance`
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetPaymentsByGroupID", "query", query, "args", args)
	rows, err := db.Query(ctx, query, args)
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

func GetPaymentByID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) (Payment, error) {
	query := `
	SELECT
		p.id,
		p.description         AS description,
		p.amount              AS amount,
		u.id                  AS payer_id,
		u.name                AS payer_name,
		u.balance             AS payer_balance,
		ARRAY_AGG(uu.id)      AS payee_ids,
		ARRAY_AGG(uu.name)    AS payee_names,
		ARRAY_AGG(uu.balance) AS payee_balances
	FROM payment AS p
	LEFT JOIN users AS u
		ON p.payer_id = u.id
	LEFT JOIN users_payment AS up
		ON up.payment_id = p.id
	LEFT JOIN users AS uu
		ON up.user_id = uu.id
	WHERE p.id = @id
	GROUP BY p.id, p.description, p.amount, u.name, u.id, u.balance`
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetPaymentsByID", "query", query, "args", args)
	rows, err := db.Query(ctx, query, args)
	if err != nil {
		L.Error(fmt.Sprintf("Get failed: %v", err))
		return Payment{}, err
	}

	payment, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Payment])
	if err != nil {
		L.Error(fmt.Sprintf("Binding failed: %v", err))
		return Payment{}, err
	}

	return payment, nil
}

type InsertPayment struct {
	Description *string `json:"description"`
	Amount      float32 `json:"amount"`
	PayerID     int     `json:"payer_id"`
	PayeeIDs    []int   `json:"payee_ids"`
}

func AddPaymentByGroupId(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int, body InsertPayment) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		// insert payment
		paymentQuery := `INSERT INTO payment (group_id, description, amount, payer_id)
VALUES (@id, @description, @amount, @payer_id)
RETURNING id`
		paymentArgs := pgx.StrictNamedArgs{
			"id":          id,
			"description": body.Description,
			"amount":      body.Amount,
			"payer_id":    body.PayerID,
		}

		var paymentID int
		L.Info("AddPaymentByGroupId.payment", "query", paymentQuery, "args", paymentArgs)
		err := tx.QueryRow(ctx, paymentQuery, paymentArgs).Scan(&paymentID)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		// insert junction
		upArgs := []any{}
		upValues := []string{}
		for _, payeeID := range body.PayeeIDs {
			upArgs = append(upArgs, payeeID)
			upArgs = append(upArgs, paymentID)
			upValues = append(upValues, "($"+strconv.Itoa(len(upArgs)-1)+", $"+strconv.Itoa(len(upArgs))+")")
		}
		upQuery := "INSERT INTO users_payment (user_id, payment_id) VALUES " + strings.Join(upValues, ", ")
		L.Info("AddPaymentByGroupId.users_payment", "query", upQuery, "args", upArgs)
		cmdTag, err := tx.Exec(ctx, upQuery, upArgs...)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}
		if cmdTag.RowsAffected() != int64(len(body.PayeeIDs)) {
			L.Error("Unexpected number of rows affected in users_payment table")
			return 0, errors.New("unexpected number of rows affected")
		}

		// update balance
		payeeBalance := body.Amount / float32(len(body.PayeeIDs))
		payeeQuery := "UPDATE users SET balance = balance + @payeeBalance WHERE id = ANY(@payeeIDs)"
		payeeArgs := pgx.StrictNamedArgs{
			"payeeBalance": payeeBalance,
			"payeeIDs":     body.PayeeIDs,
		}
		L.Info("AddPaymentByGroupId.payees", "query", payeeQuery, "args", payeeArgs)
		cmdTag, err = tx.Exec(ctx, payeeQuery, payeeArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return 0, err
		}
		if cmdTag.RowsAffected() != int64(len(body.PayeeIDs)) {
			L.Error("Unexpected number of rows affected in users table")
			return 0, errors.New("unexpected number of rows affected")
		}

		payerQuery := "UPDATE users SET balance = balance - @amount WHERE id = @id"
		payerArgs := pgx.StrictNamedArgs{
			"amount": body.Amount,
			"id":     body.PayerID,
		}
		L.Info("AddPaymentByGroupId.payer", "query", payerQuery, "args", payerArgs)
		cmdTag, err = tx.Exec(ctx, payerQuery, payerArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return 0, err
		}
		if cmdTag.RowsAffected() != 1 {
			L.Error("Unexpected number of rows affected in users table")
			return 0, errors.New("unexpected number of rows affected")
		}

		return paymentID, nil
	})
}

func PatchPayment(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, payment Payment, amount *float32, description *string) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		if description != nil {
			query := "UPDATE payment SET description = @description WHERE id = @id"
			args := pgx.StrictNamedArgs{
				"description": *description,
				"id":          payment.ID,
			}
			L.Info("PatchPayment.description", "query", query, "args", args)
			cmdTag, err := tx.Exec(ctx, query, args)
			if err != nil {
				L.Error(fmt.Sprintf("Update failed: %v", err))
				return struct{}{}, err
			}
			if cmdTag.RowsAffected() != 1 {
				L.Error("Unexpected number of rows affected in users table")
				return struct{}{}, errors.New("unexpected number of rows affected")
			}
		}

		if amount != nil {
			// update payment
			paymentQuery := "UPDATE payment SET amount = @amount WHERE id = @id"
			paymentArgs := pgx.StrictNamedArgs{
				"amount": amount,
				"id":     payment.ID,
			}
			L.Info("PatchPayment.paymentAmount", "query", paymentQuery, "args", paymentArgs)
			cmdTag, err := tx.Exec(ctx, paymentQuery, paymentArgs)
			if err != nil {
				L.Error(fmt.Sprintf("Update failed: %v", err))
				return struct{}{}, err
			}
			if cmdTag.RowsAffected() != 1 {
				L.Error("Unexpected number of rows affected in users table")
				return struct{}{}, errors.New("unexpected number of rows affected")
			}

			// update payer balance
			amtDiff := *amount - payment.Amount
			payerQuery := "UPDATE users SET balance = balance - @diff WHERE id = @id"
			payerArgs := pgx.StrictNamedArgs{
				"diff": amtDiff,
				"id":   payment.PayerID,
			}
			L.Info("PatchPayment.payerBalance", "query", payerQuery, "args", payerArgs)
			cmdTag, err = tx.Exec(ctx, payerQuery, payerArgs)
			if err != nil {
				L.Error(fmt.Sprintf("Update failed: %v", err))
				return struct{}{}, err
			}
			if cmdTag.RowsAffected() != 1 {
				L.Error("Unexpected number of rows affected in users table")
				return struct{}{}, errors.New("unexpected number of rows affected")
			}

			// update payee balances
			amtDiffPer := amtDiff / float32(len(payment.PayeeIDs))
			payeeQuery := "UPDATE users SET balance = balance + @amtDiffPer WHERE id = ANY(@ids)"
			payeeArgs := pgx.StrictNamedArgs{
				"amtDiffPer": amtDiffPer,
				"ids":        payment.PayeeIDs,
			}
			L.Info("PatchPayment.payeeBalance", "query", payeeQuery, "args", payeeArgs)
			cmdTag, err = tx.Exec(ctx, payeeQuery, payeeArgs)
			if err != nil {
				L.Error(fmt.Sprintf("Update failed: %v", err))
				return struct{}{}, err
			}
			if cmdTag.RowsAffected() != int64(len(payment.PayeeIDs)) {
				L.Error("Unexpected number of rows affected in users table")
				return struct{}{}, errors.New("unexpected number of rows affected")
			}
		}

		return struct{}{}, nil
	})

	return err
}

func DeletePayment(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, payment Payment) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		// update payees
		payeeBalance := payment.Amount / float32(len(payment.PayeeIDs))
		payeeQuery := "UPDATE users SET balance = balance - @payeeBalance WHERE id = ANY(@payeeIDs)"
		payeeArgs := pgx.StrictNamedArgs{
			"payeeBalance": payeeBalance,
			"payeeIDs":     payment.PayeeIDs,
		}
		L.Info("DeletePayment.payees", "query", payeeQuery, "args", payeeArgs)
		cmdTag, err := tx.Exec(ctx, payeeQuery, payeeArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() != int64(len(payment.PayeeIDs)) {
			L.Error("Unexpected number of rows affected in users table")
			return struct{}{}, errors.New("unexpected number of rows affected")
		}

		// update payer
		payerQuery := "UPDATE users SET balance = balance - @amount WHERE id = @id"
		payerArgs := pgx.StrictNamedArgs{
			"amount": payment.Amount,
			"id":     payment.PayerID,
		}
		L.Info("AddPaymentByGroupId.payer", "query", payerQuery, "args", payerArgs)
		cmdTag, err = tx.Exec(ctx, payerQuery, payerArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() != 1 {
			L.Error("Unexpected number of rows affected in users table")
			return struct{}{}, errors.New("unexpected number of rows affected")
		}

		// remove payment
		deleteQuery := "DELETE FROM payment WHERE id = @id"
		deleteArgs := pgx.StrictNamedArgs{
			"id": payment.ID,
		}
		L.Info("DeletePayment.delete", "query", deleteQuery, "args", deleteArgs)
		cmdTag, err = tx.Exec(ctx, deleteQuery, deleteArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() != 1 {
			L.Error("Unexpected number of rows affected in users table")
			return struct{}{}, errors.New("unexpected number of rows affected")
		}

		return struct{}{}, nil
	})
	return err
}

func DeleteAllPayments(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, groupID int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		// remove all payments
		deleteQuery := "DELETE FROM payment WHERE group_id = @id"
		deleteArgs := pgx.StrictNamedArgs{
			"id": groupID,
		}
		L.Info("DeleteAllPayments.delete", "query", deleteQuery, "args", deleteArgs)
		_, err := tx.Exec(ctx, deleteQuery, deleteArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}

		// clear all balances
		userQuery := "UPDATE users SET balance = 0 WHERE group_id = @id"
		userArgs := pgx.StrictNamedArgs{
			"id": groupID,
		}
		L.Info("DeleteAllPayments.deleteAll", "query", userQuery, "args", userArgs)
		_, err = tx.Exec(ctx, userQuery, userArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}

		return struct{}{}, nil
	})

	return err
}
