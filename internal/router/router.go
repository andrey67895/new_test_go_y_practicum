package router

import (
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRoutersForServer(iStorage storage.IStorageData) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RealIP, handlers.WithLogging, middleware.Recoverer, handlers.GzipHandleResponse, handlers.WithSendsGzip, handlers.CheckHeaderCrypto, handlers.ResponseAddHeaderCrypto)
	r.Post("/update/{type}/{name}/{value}", handlers.SaveDataForPathParams(iStorage))
	r.Post("/update/", handlers.SaveMetDataForJSON(iStorage))
	r.Post("/updates/", handlers.SaveArraysMetDataForJSON(iStorage))
	r.Post("/value/", handlers.GetDataForJSON(iStorage))
	r.Get("/value/{type}/{name}", handlers.GetDataByPathParams(iStorage))
	r.Get("/ping", handlers.GetPing(iStorage))
	r.Get("/", handlers.GetAllData(iStorage))
	return r
}
