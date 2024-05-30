package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/router"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
)

func main() {
	config.InitServerConfig()
	if config.FileStoragePathServer != "" {
		if config.RestoreServer {
			CreateData(config.FileStoragePathServer)
		}
		go Save(config.FileStoragePathServer, config.StoreIntervalServer)
	}
	log.Fatal(http.ListenAndServe(":"+config.PortServer, router.GetRoutersForServer()))
}

func CreateData(fname string) {
	data, err := os.ReadFile(fname)
	if err != nil {
		println(err.Error())
		return
	}
	var tModel []model.JSONMetrics
	if err := json.Unmarshal(data, &tModel); err != nil {
		println(err.Error())
	}
	for i := 0; i < len(tModel); i++ {
		metric := tModel[i]
		SaveData(metric)
	}
}

func SaveData(tModel model.JSONMetrics) {
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

func Save(fname string, storeInterval int) {
	for {
		time.Sleep(time.Duration(storeInterval) * time.Second)
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
			println(err.Error())
			return
		}

		err = os.MkdirAll(filepath.Dir(fname), 0666)
		if err != nil {
			println(err.Error())
			return
		}
		_, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			println(err.Error())
			return
		}
		err = os.WriteFile(fname, data, 0666)
		if err != nil {
			println(err.Error())
			return
		}
	}
}

func getValueInEnv(env string) int {
	envInt, err := strconv.Atoi(env)
	if err != nil {
		log.Fatal(err)
	}
	return envInt
}
