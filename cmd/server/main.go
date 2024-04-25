package main

import (
	"flag"
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"strings"
)

var host string

func main() {
	flag.StringVar(&host, "a", "localhost:8080", "host for server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		host = envRunAddr
	}
	port := strings.Split(host, ":")[1]
	r := chi.NewRouter()
	r.Use(middleware.RealIP, handlers.WithLogging, middleware.Recoverer)
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	r.Post("/update/", handlers.JsonMetHandler)
	r.Post("/value/", handlers.JsonGetMetHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetHandler)
	r.Get("/", handlers.GetAll)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
