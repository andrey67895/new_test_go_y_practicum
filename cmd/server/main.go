package main

import (
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
