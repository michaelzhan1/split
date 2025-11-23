package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/michaelzhan1/split/internals/handlers"
	"github.com/michaelzhan1/split/internals/logs"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	L := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	r := chi.NewRouter()
	r.Use(logs.RequestLogger(L))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONTEND_URL")},
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	r.Route("/groups", func(r chi.Router) {
		r.Post("/", handlers.CreateGroup(db, L))
		r.Get("/{group_id}", handlers.GetGroup(db, L))
		r.Patch("/{group_id}", handlers.PatchGroup(db, L))
		r.Delete("/{group_id}", handlers.DeleteGroup(db, L))

		r.Get("/{group_id}/users", handlers.GetUsers(db, L))
		r.Post("/{group_id}/users", handlers.AddUser(db, L))
		r.Patch("/{group_id}/users/{user_id}", handlers.PatchUser(db, L))
		r.Delete("/{group_id}/users/{user_id}", handlers.DeleteUser(db, L))

		r.Get("/{group_id}/payments", handlers.GetPayments(db, L))
		r.Post("/{group_id}/payments", handlers.AddPayment(db, L))
		r.Patch("/{group_id}/payments/{payment_id}", handlers.PatchPayment(db, L))
		r.Delete("/{group_id}/payments/{payment_id}", handlers.DeletePayment(db, L))
		r.Delete("/{group_id}/payments", handlers.DeleteAllPayments(db, L)) // delete all

		r.Post("/{group_id}/calculate", handlers.Calculate(db, L))
	})

	port := "3000"
	L.Info(fmt.Sprintf("Serving on port %s", port))
	http.ListenAndServe(":"+port, r)
}
