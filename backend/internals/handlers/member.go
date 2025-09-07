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
