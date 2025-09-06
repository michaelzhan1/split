package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

func main () {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	r := chi.NewRouter()

	r.Route("/groups", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var name string
			err = conn.QueryRow(ctx, "select name from member where id=$1", 1).Scan(&name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write([]byte(name))
		})
		r.Delete("/{group_id}", func(w http.ResponseWriter, r *http.Request) {})

		r.Get("/{group_id}/members", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/{group_id}/members", func(w http.ResponseWriter, r *http.Request) {})
		r.Patch("/{group_id}/members", func(w http.ResponseWriter, r *http.Request) {})
		r.Delete("/{group_id}/members", func(w http.ResponseWriter, r *http.Request) {})

		r.Post("/{group_id}/calculate", func(w http.ResponseWriter, r *http.Request) {})
	})

	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/{transaction_id}", func(w http.ResponseWriter, r *http.Request) {})
		r.Patch("/{transaction_id}", func(w http.ResponseWriter, r *http.Request) {})
		r.Delete("/{transaction_id}", func(w http.ResponseWriter, r *http.Request) {})
	})
	
	port := "3000"
	fmt.Printf("Serving on port %s\n", port)
	http.ListenAndServe(":" + port, r)
}