package router

import (
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRoutersForServer(iStorage storage.IStorageData) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP, handlers.WithLogging, middleware.Recoverer, handlers.GzipHandleResponse)
	r.Post("/update/{type}/{name}/{value}", handlers.SaveDataForPathParams(iStorage))
	r.Post("/update/", handlers.SaveMetDataForJson(iStorage))
	r.Post("/updates/", handlers.SaveArraysMetDataForJson(iStorage))
	r.Post("/value/", handlers.GetDataForJson(iStorage))
	r.Get("/value/{type}/{name}", handlers.GetDataByPathParams(iStorage))
	r.Get("/ping", handlers.GetPing(iStorage))
	r.Get("/", handlers.GetAllData(iStorage))
	return r
}
