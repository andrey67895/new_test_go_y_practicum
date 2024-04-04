package main

import (
	"flag"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"
)

var metricsName = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
	"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
	"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "RandomValue",
}

var count = model.NewCount("PollCount", 0)
var metrics = make(map[string]model.Gauge)

func updateMetrics(pollInterval time.Duration) {
	for {
		for _, statName := range metricsName {
			metrics[statName] = model.NewGauge(statName, getMemByStats(statName))
		}
		count.UpdateCountPlusOne()
		time.Sleep(pollInterval * time.Second)
	}
}

func sendMetrics(pollInterval time.Duration, host string) {
	for {
		time.Sleep(pollInterval * time.Second)
		for k, v := range metrics {
			sendRequest(host, "gauge", k, strconv.FormatFloat(v.GetMetrics(), 'f', -1, 64))
		}
		sendRequest(host, "counter", count.GetName(), strconv.Itoa(int(count.GetMetrics())))
		count.ClearCount()

	}
}

func sendRequest(host string, typeMetr string, nameMetr string, metrics string) {
	url := "http://" + host + "/update/" + typeMetr + "/" + nameMetr + "/" + metrics
	body, err := http.Post(url, "text/plain", nil)
	if err != nil {
		println(err.Error())
	} else {
		errClose := body.Body.Close()
		if errClose != nil {
			println(errClose.Error())
		}
	}
}

var host string
var reportInterval int
var pollInterval int

func main() {
	flag.StringVar(&host, "a", "localhost:8080", "host for server")
	flag.IntVar(&reportInterval, "r", 10, "reportInterval for send metrics to server")
	flag.IntVar(&pollInterval, "p", 2, "pollInterval for update metrics")
	flag.Parse()
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		host = envRunAddr
	}
	getEnvByKey("REPORT_INTERVAL", reportInterval)
	getEnvByKey("POLL_INTERVAL", pollInterval)
	go updateMetrics(time.Duration(pollInterval))
	go sendMetrics(time.Duration(reportInterval), host)
	server := http.Server{}
	log.Fatal(server.ListenAndServe())
}

func getEnvByKey(key string, value int) {
	if envInt := os.Getenv(key); envInt != "" {
		value, _ = strconv.Atoi(envInt)
	}
}

func getMemByStats(name string) float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	switch name {
	case "Alloc":
		return float64(memStats.Alloc)
	case "BuckHashSys":
		return float64(memStats.BuckHashSys)
	case "Frees":
		return float64(memStats.Frees)
	case "GCCPUFraction":
		return memStats.GCCPUFraction
	case "GCSys":
		return float64(memStats.GCSys)
	case "HeapAlloc":
		return float64(memStats.HeapAlloc)
	case "HeapIdle":
		return float64(memStats.HeapIdle)
	case "HeapInuse":
		return float64(memStats.HeapInuse)
	case "HeapObjects":
		return float64(memStats.HeapObjects)
	case "HeapReleased":
		return float64(memStats.HeapReleased)
	case "HeapSys":
		return float64(memStats.HeapSys)
	case "LastGC":
		return float64(memStats.LastGC)
	case "Lookups":
		return float64(memStats.Lookups)
	case "MCacheInuse":
		return float64(memStats.MCacheInuse)
	case "MCacheSys":
		return float64(memStats.MCacheSys)
	case "MSpanInuse":
		return float64(memStats.MSpanInuse)
	case "MSpanSys":
		return float64(memStats.MSpanSys)
	case "Mallocs":
		return float64(memStats.Mallocs)
	case "NumForcedGC":
		return float64(memStats.NumForcedGC)
	case "NumGC":
		return float64(memStats.NumGC)
	case "OtherSys":
		return float64(memStats.OtherSys)
	case "PauseTotalNs":
		return float64(memStats.PauseTotalNs)
	case "StackInuse":
		return float64(memStats.StackInuse)
	case "Sys":
		return float64(memStats.Sys)
	case "TotalAlloc":
		return float64(memStats.TotalAlloc)
	case "RandomValue":
		return rand.Float64()
	default:
		return 0
	}
}
