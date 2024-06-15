package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/helpers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
)

func GetAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	for k, v := range storage.LocalNewMemStorageGauge.GetData() {
		_, err := w.Write([]byte("Name: " + k + ". Value: " + fmt.Sprint(v) + "\n"))
		if err != nil {
			return
		}
	}
	for k, v := range storage.LocalNewMemStorageCounter.GetData() {
		_, err := w.Write([]byte("Name: " + k + ". Value: " + fmt.Sprint(v)))
		if err != nil {
			return
		}
	}
}

func GetMetHandler(w http.ResponseWriter, req *http.Request) {
	typeMet := chi.URLParam(req, "type")
	nameMet := chi.URLParam(req, "name")
	if typeMet == "gauge" {
		localGauge, err := storage.LocalNewMemStorageGauge.GetGauge(nameMet)
		if err != nil {
			http.Error(w, "Название метрики не найдено", http.StatusNotFound)
			return
		}
		_, errWrite := w.Write([]byte(fmt.Sprint(localGauge)))
		if errWrite != nil {
			return
		}
	} else if typeMet == "counter" {

		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
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

func MetHandler(w http.ResponseWriter, req *http.Request) {
	contentEncoding := req.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")
	if sendsGzip {
		cr, err := newCompressReader(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Body = cr.r
		defer cr.zr.Close()
	}
	typeMet := chi.URLParam(req, "type")
	nameMet := chi.URLParam(req, "name")
	addHash(w, req)
	if typeMet == "gauge" {
		valueMet, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
			return
		}
		if config.DatabaseDsn != "" {
			helpers.RetrySaveGaugeInDB(nameMet, valueMet)
		}
		err = storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
		if err != nil {
			return
		}
	} else if typeMet == "counter" {
		valueMet, err := strconv.Atoi(chi.URLParam(req, "value"))
		if err != nil {

			http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
			return
		}
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			if config.DatabaseDsn != "" {
				helpers.RetrySaveCounterInDB(nameMet, int64(valueMet))
			}
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, int64(valueMet))
			if err != nil {
				return
			}
		} else {
			if config.DatabaseDsn != "" {
				helpers.RetrySaveCounterInDB(nameMet, int64(int(localCounter)+valueMet))
			}
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
