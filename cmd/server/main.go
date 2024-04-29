package main

import (
	"encoding/json"
	"flag"
	"github.com/andrey67895/new_test_go_y_practicum/internal/handlers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var host string
var storeInterval int
var fileStoragePath string
var restore bool

func main() {
	flag.StringVar(&host, "a", "localhost:8080", "host for server")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		host = envRunAddr
	}
	port := strings.Split(host, ":")[1]

	flag.IntVar(&storeInterval, "i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	flag.StringVar(&fileStoragePath, "f", "tmp/metrics-db.json", "полное имя файла, куда сохраняются текущие значения ")
	flag.BoolVar(&restore, "r", true, "загружать или нет ранее сохранённые значения из указанного файла при старте сервера")
	flag.Parse()
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		storeInterval = getValueInEnv(envStoreInterval)
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		fileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		restore = getBool(envRestore)
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP, handlers.WithLogging, middleware.Recoverer, handlers.GzipHandleResponse)
	r.Post("/update/{type}/{name}/{value}", handlers.MetHandler)
	r.Post("/update/", handlers.JSONMetHandler)
	r.Post("/value/", handlers.JSONGetMetHandler)
	r.Get("/value/{type}/{name}", handlers.GetMetHandler)
	r.Get("/", handlers.GetAll)
	if fileStoragePath != "" {
		if restore {
			CreateData(fileStoragePath)
		}
		go Save(fileStoragePath, storeInterval)
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
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

func getBool(env string) bool {

	boolValue, err := strconv.ParseBool(env)
	if err != nil {
		log.Fatal(err)
	}
	return boolValue
}
