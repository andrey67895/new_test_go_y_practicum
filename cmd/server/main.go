package main

import (
	"flag"
	"fmt"
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func main() {
	port := flag.Int("a", 8080, "port for host")
	flag.Parse()
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetHandler)
	r.Get("/", handlers.GetAll)
	err := http.ListenAndServe(":"+string(fmt.Sprintf("%d", *port)), r)
	if err != nil {
		println(err.Error())
		return
	}
}
