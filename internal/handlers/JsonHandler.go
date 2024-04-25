package handlers

import (
	"encoding/json"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"net/http"
)

func JsonMetHandler(w http.ResponseWriter, req *http.Request) {
	tModel := model.JsonMetrics{}
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	typeMet := tModel.MType
	nameMet := tModel.ID

	if typeMet == "gauge" {
		valueMet := tModel.Value
		err = storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
		if err != nil {
			return
		}
	} else if typeMet == "counter" {
		valueMet := tModel.Delta
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, valueMet)
			if err != nil {
				return
			}
		} else {
			tModel.Delta = localCounter + valueMet
			err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, tModel.Delta)
			if err != nil {
				return
			}
		}
	} else {
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	tJson, _ := json.Marshal(tModel)
	_, err = w.Write(tJson)
	if err != nil {
		http.Error(w, "Ошибка при записи ответа", http.StatusBadRequest)
		return
	}
}

func JsonGetMetHandler(w http.ResponseWriter, req *http.Request) {
	tModel := model.JsonMetrics{}
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	typeMet := tModel.MType
	nameMet := tModel.ID
	if typeMet == "gauge" {
		localGauge, err := storage.LocalNewMemStorageGauge.GetGauge(nameMet)
		if err != nil {
			http.Error(w, "Название метрики не найдено", http.StatusNotFound)
			return
		}
		tModel.Value = localGauge
		marshal, err := json.Marshal(tModel)
		if err != nil {
			http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
			return
		}
		_, errWrite := w.Write(marshal)
		if errWrite != nil {
			return
		}
	} else if typeMet == "counter" {

		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			http.Error(w, "Название метрики не найдено", http.StatusNotFound)
			return
		}
		tModel.Delta = localCounter
		marshal, err := json.Marshal(tModel)
		if err != nil {
			http.Error(w, "Ошибка записи ответа", http.StatusNotFound)
			return
		}
		_, errWrite := w.Write(marshal)
		if errWrite != nil {
			return
		}

	} else {
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
