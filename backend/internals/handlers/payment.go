package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelzhan1/split/internals/database"
)

func GetPayments(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var httpError *HttpError
		defer func() {
			if httpError != nil {
				data, _ := json.Marshal(httpError)
				L.Info(httpError.Message, "code", httpError.Code)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpError.Code)
				w.Write(data)
			}
		}()

		partyID, httpError := withPartyID(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyByID(ctx, db, L, partyID)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		payments, err := database.GetPaymentsByPartyID(ctx, db, L, partyID)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
			return
		}

		res := toPaymentList(payments)
		data, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func AddPayment(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	type request = database.InsertPayment

	type response struct {
		ID int `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var httpError *HttpError
		defer func() {
			if httpError != nil {
				data, _ := json.Marshal(httpError)
				L.Info(httpError.Message, "code", httpError.Code)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpError.Code)
				w.Write(data)
			}
		}()

		partyId, httpError := withPartyID(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyByID(ctx, db, L, partyId)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		var body request
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Invalid JSON",
			}
			return
		}
		if body.Amount <= 0 {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Non-positive balance",
			}
			return
		}
		if len(body.PayeeIDs) == 0 {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "No payees in payment",
			}
			return
		}

		id, err := database.AddPaymentByPartyId(ctx, db, L, partyId, body)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
			return
		}

		res := response{id}
		data, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}

func PatchPayment(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	type request struct {
		Amount      *float32    `json:"amount"`
		Description *string `json:"description"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var httpError *HttpError
		defer func() {
			if httpError != nil {
				data, _ := json.Marshal(httpError)
				L.Info(httpError.Message, "code", httpError.Code)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpError.Code)
				w.Write(data)
			}
		}()

		partyId, httpError := withPartyID(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyByID(ctx, db, L, partyId)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		paymentID, httpError := withPaymentID(r)
		if httpError != nil {
			return
		}

		payment, err := database.GetPaymentByID(ctx, db, L, paymentID)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		var body request
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Invalid JSON",
			}
			return
		}
		if body.Amount == nil && body.Description == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if body.Amount != nil && *body.Amount <= 0 {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Invalid amount field",
			}
			return
		}

		err = database.PatchPayment(ctx, db, L, payment, body.Amount, body.Description)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}

func DeletePayment(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var httpError *HttpError
		defer func() {
			if httpError != nil {
				data, _ := json.Marshal(httpError)
				L.Info(httpError.Message, "code", httpError.Code)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpError.Code)
				w.Write(data)
			}
		}()

		partyId, httpError := withPartyID(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyByID(ctx, db, L, partyId)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		paymentID, httpError := withPaymentID(r)
		if httpError != nil {
			return
		}

		payment, err := database.GetPaymentByID(ctx, db, L, paymentID)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		err = database.DeletePayment(ctx, db, L, payment)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}

func DeleteAllPayments(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var httpError *HttpError
		defer func() {
			if httpError != nil {
				data, _ := json.Marshal(httpError)
				L.Info(httpError.Message, "code", httpError.Code)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpError.Code)
				w.Write(data)
			}
		}()

		partyId, httpError := withPartyID(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyByID(ctx, db, L, partyId)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		err = database.DeleteAllPayments(ctx, db, L, partyId)
		if err != nil {
			if err == pgx.ErrNoRows {
				httpError = &HttpError{
					Code:    http.StatusNotFound,
					Message: "Not found",
				}
			} else {
				httpError = &HttpError{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				}
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
