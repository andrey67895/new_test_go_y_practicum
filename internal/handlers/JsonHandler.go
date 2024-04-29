package handlers

import (
	"encoding/json"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"net/http"
	"strings"
)

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
	tModel := model.JSONMetrics{}
	err := json.NewDecoder(req.Body).Decode(&tModel)
	if err != nil {
		println(err.Error())
		http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
		return
	}
	typeMet := tModel.MType
	nameMet := tModel.ID

	if typeMet == "gauge" {
		valueMet := tModel.GetValue()
		err = storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
		if err != nil {
			println(err.Error())
			return
		}
	} else if typeMet == "counter" {
		valueMet := tModel.GetDelta()
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, valueMet)
			if err != nil {
				println(err.Error())
				return
			}
		} else {
			tModel.SetDelta(localCounter + valueMet)
			err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, tModel.GetDelta())
			if err != nil {
				println(err.Error())
				return
			}
		}
	} else {
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	tJSON, _ := json.Marshal(tModel)
	_, err = w.Write(tJSON)
	if err != nil {
		println(err.Error())
		http.Error(w, "Ошибка при записи ответа", http.StatusBadRequest)
		return
	}
}

func SaveLocalData(tModel model.JSONMetrics) {
	typeMet := tModel.MType
	nameMet := tModel.ID

	if typeMet == "gauge" {
		valueMet := tModel.GetValue()
		err := storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
		if err != nil {
			println(err.Error())
			return
		}
	} else if typeMet == "counter" {
		valueMet := tModel.GetDelta()
		localCounter, err := storage.LocalNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := storage.LocalNewMemStorageCounter.SetCounter(nameMet, valueMet)
			if err != nil {
				println(err.Error())
				return
			}
		} else {
			tModel.SetDelta(localCounter + valueMet)
			err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, tModel.GetDelta())
			if err != nil {
				println(err.Error())
				return
			}
		}
	}
}

func JSONGetMetHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tModel := model.JSONMetrics{}
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
		tModel.SetValue(localGauge)
		println(localGauge)
		marshal, err := json.Marshal(tModel)
		println(string(marshal))
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

	} else {
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
