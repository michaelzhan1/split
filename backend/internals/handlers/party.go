package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelzhan1/split/internals/database"
)

func PostCreateParty(db *pgxpool.Pool) http.HandlerFunc {
	type request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		var httpErrCode int
		var httpErrMsg string
		defer func() {
			if httpErrCode != 0 {
				httpErr := HttpError{
					Code:    httpErrCode,
					Message: httpErrMsg,
				}
				data, _ := json.Marshal(httpErr)
	
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpErrCode)
				w.Write(data)
			}
		}()

		var body request
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			httpErrCode = http.StatusBadRequest
			httpErrMsg = "Invalid JSON"
			return
		}
		if body.Name == "" {
			httpErrCode = http.StatusBadRequest
			httpErrMsg = "Empty name field"
		}

		id, err := database.CreateParty(ctx, db, body.Name)
		if err != nil {
			httpErrCode = http.StatusInternalServerError
			httpErrMsg = "Internal server error"
			return
		}

		res := CreatePartyResponse{id}
		data, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(data)
	}
}
