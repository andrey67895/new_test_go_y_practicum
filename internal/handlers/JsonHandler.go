package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
)

var log = logger.Log()

func JSONMetHandler(w http.ResponseWriter, req *http.Request) {
	contentEncoding := req.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")
	if sendsGzip {
		cr, err := newCompressReader(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		req.Body = cr.zr
		defer cr.zr.Close()
	}
	w.Header().Set("Content-Type", "application/json")
	var tModel model.JSONMetrics
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	typeMet := tModel.MType
	nameMet := tModel.ID
	switch typeMet {
	case "gauge":
		valueMet := tModel.GetValue()
		err = storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
		if err != nil {
			log.Error(err.Error())
			return
		}
	case "counter":
		valueMet := tModel.GetDelta()
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, valueMet)
			if err != nil {
				log.Error(err.Error())
				return
			}
		} else {
			tModel.SetDelta(localCounter + valueMet)
			err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, tModel.GetDelta())
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	default:
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	tJSON, _ := json.Marshal(tModel)
	_, err = w.Write(tJSON)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, "Ошибка при записи ответа", http.StatusBadRequest)
		return
	}
}

func JSONGetMetHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tModel model.JSONMetrics
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	typeMet := tModel.MType
	nameMet := tModel.ID
	switch typeMet {
	case "gauge":
		localGauge, err := storage.LocalNewMemStorageGauge.GetGauge(nameMet)
		if err != nil {
			http.Error(w, "Название метрики не найдено", http.StatusNotFound)
			return
		}
		tModel.SetValue(localGauge)
		marshal, err := json.Marshal(tModel)
		if err != nil {
			http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
			return
		}
		_, errWrite := w.Write(marshal)
		if errWrite != nil {
			return
		}
	case "counter":
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			http.Error(w, "Название метрики не найдено", http.StatusNotFound)
			return
		}
		tModel.SetDelta(localCounter)
		marshal, err := json.Marshal(tModel)
		if err != nil {
			http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
			return
		}
		_, errWrite := w.Write(marshal)
		if errWrite != nil {
			return
		}
	default:
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
