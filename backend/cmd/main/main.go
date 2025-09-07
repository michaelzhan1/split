package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michaelzhan1/split/internals/handlers"
	"github.com/michaelzhan1/split/internals/logs"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	L := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	r := chi.NewRouter()
	r.Use(logs.RequestLogger(L))

	r.Route("/parties", func(r chi.Router) {
		r.Post("/", handlers.CreateParty(db, L))
		r.Get("/{party_id}", handlers.GetParty(db, L))
		r.Patch("/{party_id}", handlers.PatchParty(db, L))
		r.Delete("/{party_id}", handlers.DeleteParty(db, L))

		r.Get("/{party_id}/members", handlers.GetMembers(db, L))
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
	L.Info(fmt.Sprintf("Serving on port %s", port))
	http.ListenAndServe(":"+port, r)
}
