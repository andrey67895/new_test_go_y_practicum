package main

import (
	"flag"
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strings"
)

func main() {
	host := flag.String("a", "localhost:8080", "host for server")
	port := strings.Split(*host, ":")[1]
	flag.Parse()
	println("HOST :::: ", host)
	println("PORT :::: ", port)
	r := chi.NewRouter()
	r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetHandler)
	r.Get("/", handlers.GetAll)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		println(err.Error())
		return
	}
}
