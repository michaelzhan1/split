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

func GetPaymentsByPartyID(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int) ([]Payment, error) {
	query := `
	SELECT
		p.id,
		p.description         AS description,
		p.amount              AS amount,
		m.id                  AS payer_id,
		m.name                AS payer_name,
		m.balance             AS payer_balance,
		ARRAY_AGG(mm.id)      AS payee_ids,
		ARRAY_AGG(mm.name)    AS payee_names,
		ARRAY_AGG(mm.balance) AS payee_balances
	FROM payment AS p
	LEFT JOIN member AS m
		ON p.payer_id = m.id
	LEFT JOIN member_payment AS mp
		ON mp.payment_id = p.id
	LEFT JOIN member AS mm
		ON mp.member_id = mm.id
	WHERE p.party_id = @id
	GROUP BY p.id, p.description, p.amount, m.name, m.id, m.balance`
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetPaymentsByPartyID", "query", query, "args", args)
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
		m.id                  AS payer_id,
		m.name                AS payer_name,
		m.balance             AS payer_balance,
		ARRAY_AGG(mm.id)      AS payee_ids,
		ARRAY_AGG(mm.name)    AS payee_names,
		ARRAY_AGG(mm.balance) AS payee_balances
	FROM payment AS p
	LEFT JOIN member AS m
		ON p.payer_id = m.id
	LEFT JOIN member_payment AS mp
		ON mp.payment_id = p.id
	LEFT JOIN member AS mm
		ON mp.member_id = mm.id
	WHERE p.id = @id
	GROUP BY p.id, p.description, p.amount, m.name, m.id, m.balance`
	args := pgx.StrictNamedArgs{
		"id": id,
	}

	L.Info("GetPaymentsByPartyID", "query", query, "args", args)
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

func AddPaymentByPartyId(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, id int, body InsertPayment) (int, error) {
	return WithTx(ctx, db, func(tx pgx.Tx) (int, error) {
		// insert payment
		paymentQuery := `INSERT INTO payment (party_id, description, amount, payer_id)
VALUES (@id, @description, @amount, @payer_id)
RETURNING id`
		paymentArgs := pgx.StrictNamedArgs{
			"id":          id,
			"description": body.Description,
			"amount":      body.Amount,
			"payer_id":    body.PayerID,
		}

		var paymentID int
		L.Info("AddPaymentByPartyId.payment", "query", paymentQuery, "args", paymentArgs)
		err := tx.QueryRow(ctx, paymentQuery, paymentArgs).Scan(&paymentID)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}

		// insert junction
		mpArgs := []any{}
		mpValues := []string{}
		for _, payeeID := range body.PayeeIDs {
			mpArgs = append(mpArgs, payeeID)
			mpArgs = append(mpArgs, paymentID)
			mpValues = append(mpValues, "($"+strconv.Itoa(len(mpArgs)-1)+", $"+strconv.Itoa(len(mpArgs))+")")
		}
		mpQuery := "INSERT INTO member_payment (member_id, payment_id) VALUES " + strings.Join(mpValues, ", ")
		L.Info("AddPaymentByPartyId.member_payment", "query", mpQuery, "args", mpArgs)
		cmdTag, err := tx.Exec(ctx, mpQuery, mpArgs...)
		if err != nil {
			L.Error(fmt.Sprintf("Insert failed: %v", err))
			return 0, err
		}
		if cmdTag.RowsAffected() != int64(len(body.PayeeIDs)) {
			L.Error("Unexpected number of rows affected in member_payment table")
			return 0, errors.New("unexpected number of rows affected")
		}

		// update balance
		payeeBalance := body.Amount / float32(len(body.PayeeIDs))
		payeeQuery := "UPDATE member SET balance = balance + @payeeBalance WHERE id = ANY(@payeeIDs)"
		payeeArgs := pgx.StrictNamedArgs{
			"payeeBalance": payeeBalance,
			"payeeIDs":     body.PayeeIDs,
		}
		L.Info("AddPaymentByPartyId.payees", "query", payeeQuery, "args", payeeArgs)
		cmdTag, err = tx.Exec(ctx, payeeQuery, payeeArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return 0, err
		}
		if cmdTag.RowsAffected() != int64(len(body.PayeeIDs)) {
			L.Error("Unexpected number of rows affected in member table")
			return 0, errors.New("unexpected number of rows affected")
		}

		payerQuery := "UPDATE member SET balance = balance - @amount WHERE id = @id"
		payerArgs := pgx.StrictNamedArgs{
			"amount": body.Amount,
			"id":     body.PayerID,
		}
		L.Info("AddPaymentByPartyId.payer", "query", payerQuery, "args", payerArgs)
		cmdTag, err = tx.Exec(ctx, payerQuery, payerArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return 0, err
		}
		if cmdTag.RowsAffected() != 1 {
			L.Error("Unexpected number of rows affected in member table")
			return 0, errors.New("unexpected number of rows affected")
		}

		return paymentID, nil
	})
}

func DeletePayment(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, payment Payment) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		// update payees
		payeeBalance := payment.Amount / float32(len(payment.PayeeIDs))
		payeeQuery := "UPDATE member SET balance = balance - @payeeBalance WHERE id = ANY(@payeeIDs)"
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
			L.Error("Unexpected number of rows affected in member table")
			return struct{}{}, errors.New("unexpected number of rows affected")
		}

		// update payer
		payerQuery := "UPDATE member SET balance = balance - @amount WHERE id = @id"
		payerArgs := pgx.StrictNamedArgs{
			"amount": payment.Amount,
			"id":     payment.PayerID,
		}
		L.Info("AddPaymentByPartyId.payer", "query", payerQuery, "args", payerArgs)
		cmdTag, err = tx.Exec(ctx, payerQuery, payerArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Update failed: %v", err))
			return struct{}{}, err
		}
		if cmdTag.RowsAffected() != 1 {
			L.Error("Unexpected number of rows affected in member table")
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
			L.Error("Unexpected number of rows affected in member table")
			return struct{}{}, errors.New("unexpected number of rows affected")
		}

		return struct{}{}, nil
	})
	return err
}

func DeleteAllPayments(ctx context.Context, db *pgxpool.Pool, L *slog.Logger, partyID int) error {
	_, err := WithTx(ctx, db, func(tx pgx.Tx) (struct{}, error) {
		// remove all payments
		deleteQuery := "DELETE FROM payment WHERE party_id = @id"
		deleteArgs := pgx.StrictNamedArgs{
			"id": partyID,
		}
		L.Info("DeleteAllPayments.delete", "query", deleteQuery, "args", deleteArgs)
		_, err := tx.Exec(ctx, deleteQuery, deleteArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}

		// clear all balances
		memberQuery := "UPDATE member SET balance = 0 WHERE party_id = @id"
		memberArgs := pgx.StrictNamedArgs{
			"id": partyID,
		}
		L.Info("DeleteAllPayments.delete", "query", memberQuery, "args", memberArgs)
		_, err = tx.Exec(ctx, memberQuery, memberArgs)
		if err != nil {
			L.Error(fmt.Sprintf("Delete failed: %v", err))
			return struct{}{}, err
		}

		return struct{}{}, nil
	})

	return err
}
