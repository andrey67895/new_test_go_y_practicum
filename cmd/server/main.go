package main

import (
	"flag"
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strings"
)

var host = flag.String("a", "localhost:8080", "host for server")

func main() {

	flag.Parse()
	port := strings.Split(*host, ":")[1]
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetHandler)
	r.Get("/", handlers.GetAll)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
