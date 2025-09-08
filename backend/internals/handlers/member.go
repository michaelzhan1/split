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

func GetMembers(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
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

		partyId, httpError := withPartyId(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyById(ctx, db, L, partyId)
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

		members, err := database.GetMembersByPartyId(ctx, db, L, partyId)
		if err != nil {
			httpError = &HttpError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
			return
		}

		res := toMemberList(members)
		data, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

func AddMember(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	type response struct {
		Id int `json:"id"`
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

		partyId, httpError := withPartyId(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyById(ctx, db, L, partyId)
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

		id, err := database.AddMemberToPartyById(ctx, db, L, partyId, body.Name)
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

func PatchMember(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
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

		partyId, httpError := withPartyId(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyById(ctx, db, L, partyId)
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

		memberId, httpError := withMemberId(r)
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

		err = database.PatchMember(ctx, db, L, partyId, memberId, *body.Name)
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

func DeleteMember(db *pgxpool.Pool, L *slog.Logger) http.HandlerFunc {
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

		partyId, httpError := withPartyId(r)
		if httpError != nil {
			return
		}

		_, err := database.GetPartyById(ctx, db, L, partyId)
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

		memberId, httpError := withMemberId(r)
		if httpError != nil {
			return
		}

		err = database.DeleteMember(ctx, db, L, partyId, memberId)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				httpError = &HttpError{
					Code:    http.StatusConflict,
					Message: "Cannot delete member with associated payments",
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
