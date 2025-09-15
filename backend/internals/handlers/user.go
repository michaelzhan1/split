package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelzhan1/split/internals/database"
)

func GetUsers(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
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

		users, err := database.GetUsersByPartyID(ctx, db, L, partyID)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
			return
		}

		res := toUserList(users)
		data, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func AddUser(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

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

		var body request
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Invalid JSON",
			}
			return
		}
		if body.Name == "" {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Empty name field",
			}
			return
		}

		id, err := database.AddUserToPartyByID(ctx, db, L, partyID, body.Name)
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

func PatchUser(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	type request struct {
		Name *string `json:"name"`
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

		userID, httpError := withUserID(r)
		if httpError != nil {
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
		if body.Name == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if body.Name != nil && *body.Name == "" {
			httpError = &HttpError{
				Code:    http.StatusBadRequest,
				Message: "Empty name field",
			}
			return
		}

		err = database.PatchUser(ctx, db, L, partyID, userID, *body.Name)
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

func DeleteUser(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
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

		userID, httpError := withUserID(r)
		if httpError != nil {
			return
		}

		err = database.DeleteUser(ctx, db, L, partyID, userID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				httpError = &HttpError{
					Code:    http.StatusConflict,
					Message: "Cannot delete user with associated payments",
				}
			} else if err == pgx.ErrNoRows {
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
