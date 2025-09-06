package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main () {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	r := chi.NewRouter()

	r.Route("/parties", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			var name string
			err := pool.QueryRow(ctx, "select name from member where id=$1", 1).Scan(&name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Write([]byte(name))
		})
		r.Delete("/{party_id}", func(w http.ResponseWriter, r *http.Request) {})

		r.Get("/{party_id}/members", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/{party_id}/members", func(w http.ResponseWriter, r *http.Request) {})
		r.Patch("/{party_id}/members", func(w http.ResponseWriter, r *http.Request) {})
		r.Delete("/{party_id}/members", func(w http.ResponseWriter, r *http.Request) {})

		r.Post("/{party_id}/calculate", func(w http.ResponseWriter, r *http.Request) {})
	})

	r.Route("/payments", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/{payment_id}", func(w http.ResponseWriter, r *http.Request) {})
		r.Patch("/{payment_id}", func(w http.ResponseWriter, r *http.Request) {})
		r.Delete("/{payment_id}", func(w http.ResponseWriter, r *http.Request) {})
	})
	
	port := "3000"
	fmt.Printf("Serving on port %s\n", port)
	http.ListenAndServe(":" + port, r)
}