package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
)

var log = logger.Log()

func GetPing(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := iStorage.Ping()
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка ping DB: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func SaveMetDataForJSON(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
			tErr := iStorage.RetrySaveGauge(req.Context(), nameMet, valueMet)
			if tErr != nil {
				log.Error(err.Error())
				return
			}
		case "counter":
			valueMet := tModel.GetDelta()
			localCounter, tErr := iStorage.GetCounter(req.Context(), nameMet)
			if tErr != nil {
				ttErr := iStorage.RetrySaveCounter(req.Context(), nameMet, valueMet)
				if ttErr != nil {
					return
				}
			} else {
				tModel.SetDelta(localCounter + valueMet)
				ttErr := iStorage.RetrySaveCounter(req.Context(), nameMet, tModel.GetDelta())
				if ttErr != nil {
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
}

func SaveArraysMetDataForJSON(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var tModels []model.JSONMetrics
		err := json.NewDecoder(req.Body).Decode(&tModels)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, "Ошибка десериализации!", http.StatusBadRequest)
			return
		}
		for _, tModel := range tModels {
			typeMet := tModel.MType
			nameMet := tModel.ID
			switch typeMet {
			case "gauge":
				valueMet := tModel.GetValue()
				err := iStorage.RetrySaveGauge(req.Context(), nameMet, valueMet)
				if err != nil {
					log.Error(err.Error())
					return
				}
			case "counter":
				valueMet := tModel.GetDelta()
				localCounter, err := iStorage.GetCounter(req.Context(), nameMet)
				if err != nil {
					err := iStorage.RetrySaveCounter(req.Context(), nameMet, valueMet)
					if err != nil {
						log.Error(err.Error())
						return
					}
				} else {
					tModel.SetDelta(localCounter + valueMet)
					err := iStorage.RetrySaveCounter(req.Context(), nameMet, tModel.GetDelta())
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
	}
}

func GetDataForJSON(iStorage storage.IStorageData) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
			localGauge, err := iStorage.GetGauge(req.Context(), nameMet)
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
			localCounter, err := iStorage.GetCounter(req.Context(), nameMet)
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
}
