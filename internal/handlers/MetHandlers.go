package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
)

func GetAllData(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		data, err := iStorage.GetData(req.Context())
		if err != nil {
			log.Error(err.Error())
			return
		}
		_, err = w.Write([]byte(data))
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}

func GetDataByPathParams(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		typeMet := chi.URLParam(req, "type")
		nameMet := chi.URLParam(req, "name")
		if typeMet == "gauge" {
			localGauge, err := iStorage.GetGauge(req.Context(), nameMet)
			if err != nil {
				http.Error(w, "Название метрики не найдено", http.StatusNotFound)
				return
			}
			_, errWrite := w.Write([]byte(fmt.Sprint(localGauge)))
			if errWrite != nil {
				return
			}
		} else if typeMet == "counter" {

			localCounter, err := iStorage.GetCounter(req.Context(), nameMet)
			if err != nil {
				http.Error(w, "Название метрики не найдено", http.StatusNotFound)
				return
			}
			_, errWrite := w.Write([]byte(fmt.Sprint(localCounter)))
			if errWrite != nil {
				return
			}

		} else {
			http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func SaveDataForPathParams(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		typeMet := chi.URLParam(req, "type")
		nameMet := chi.URLParam(req, "name")

		if typeMet == "gauge" {
			valueMet, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
			if err != nil {
				http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
				return
			}
			err = iStorage.RetrySaveGauge(req.Context(), nameMet, valueMet)
			if err != nil {
				return
			}
		} else if typeMet == "counter" {
			valueMet, err := strconv.Atoi(chi.URLParam(req, "value"))
			if err != nil {

				http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
				return
			}
			localCounter, err := iStorage.GetCounter(req.Context(), nameMet)
			if err != nil {
				err = iStorage.RetrySaveCounter(req.Context(), nameMet, int64(valueMet))
				if err != nil {
					return
				}
			} else {
				err = iStorage.RetrySaveCounter(req.Context(), nameMet, int64(int(localCounter)+valueMet))
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
}
