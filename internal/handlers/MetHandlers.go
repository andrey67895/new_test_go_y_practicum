package handlers

import (
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func MetHandler(w http.ResponseWriter, req *http.Request) {
	typeMet := chi.URLParam(req, "type")
	nameMet := chi.URLParam(req, "name")
	valueMet, err := strconv.Atoi(chi.URLParam(req, "value"))
	if err != nil {

		http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
		return
	}
	if typeMet == "gauge" {
		err := storage.LocalNewMemStorageGauge.SetGauge(nameMet, float64(valueMet))
		if err != nil {
			return
		}
	} else if typeMet == "counter" {
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, int64(valueMet))
			if err != nil {
				return
			}
		} else {
			err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, int64(int(localCounter)+valueMet))
			if err != nil {
				return
			}
		}
	} else {
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
