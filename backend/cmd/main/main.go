package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main () {
	r := chi.NewRouter()

	r.Route("/groups", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {})
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