// Package router... отвечает за работу с роутером
package router

import (
	"net/http/pprof"

	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// GetRoutersForServer инициализация всех роутеров
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
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	return r
}
