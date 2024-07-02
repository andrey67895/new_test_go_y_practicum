package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/router"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
)

var log = logger.Log()

func main() {
	config.InitServerConfig()
	var st storage.IStorageData
	if config.DatabaseDsn != "" {
		ctx := context.Background()
		st = storage.InitDB(ctx)
	} else {
		st = storage.InMemStorage{}
		if config.FileStoragePathServer != "" {
			if config.RestoreServer {
				RestoringDataFromFile(config.FileStoragePathServer)
			}
			go SaveDataForInterval(config.FileStoragePathServer, config.StoreIntervalServer)
		}
	}
	log.Fatal(http.ListenAndServe(":"+config.PortServer, router.GetRoutersForServer(st)))
}

func RestoringDataFromFile(fname string) {
	data, err := os.ReadFile(fname)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var tModel []model.JSONMetrics
	if err := json.Unmarshal(data, &tModel); err != nil {
		log.Error(err.Error())
	}
	for i := 0; i < len(tModel); i++ {
		metric := tModel[i]
		SaveData(metric)
	}
}

func SaveData(tModel model.JSONMetrics) {
	typeMet := tModel.MType
	nameMet := tModel.ID

	switch typeMet {
	case "gauge":
		valueMet := tModel.GetValue()
		err := storage.LocalNewMemStorageGauge.SetGauge(nameMet, valueMet)
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
			} else {
				tModel.SetDelta(localCounter + valueMet)
				err = storage.LocalNewMemStorageCounter.SetCounter(nameMet, tModel.GetDelta())
				if err != nil {
					log.Error(err.Error())
					return
				}
			}
		}
	}

}

func SaveDataForInterval(fname string, storeInterval int) {
	if storeInterval > 0 {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		for range ticker.C {
			SaveDataInFile(fname)
			log.Infoln("Save Data file at: ", time.Now())
		}

	} else {
		for {
			SaveDataInFile(fname)
			log.Infoln("Save Data file at: ", time.Now())
		}
	}

}

func SaveDataInFile(fname string) {
	var tModel []model.JSONMetrics
	for k, v := range storage.LocalNewMemStorageGauge.GetData() {
		tJSON := model.JSONMetrics{}
		tJSON.ID = k
		tJSON.SetValue(v)
		tJSON.MType = "gauge"
		tModel = append(tModel, tJSON)
	}
	for k, v := range storage.LocalNewMemStorageCounter.GetData() {
		tJSON := model.JSONMetrics{}
		tJSON.ID = k
		tJSON.SetDelta(v)
		tJSON.MType = "counter"
		tModel = append(tModel, tJSON)
	}
	data, err := json.MarshalIndent(tModel, "", "   ")
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = os.MkdirAll(filepath.Dir(fname), 0666)
	if err != nil {
		log.Error(err.Error())
		return
	}
	_, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)

	if err != nil {
		log.Error(err.Error())
		return
	}
	err = os.WriteFile(fname, data, 0666)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
