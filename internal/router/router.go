package router

import (
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRoutersForServer() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RealIP, handlers.WithLogging, WithCrypto, middleware.Recoverer, handlers.GzipHandleResponse)
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	r.Post("/update/", handlers.JSONMetHandler)
	r.Post("/updates/", handlers.JSONMetHandlerUpdates)
	r.Post("/value/", handlers.JSONGetMetHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetHandler)
	r.Get("/ping", helpers.GetPing)
	r.Get("/", handlers.GetAll)
	return r
}
