package main

import (
	_ "fmt"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"math/rand"
	"net/http"
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

func sendMetrics(pollInterval time.Duration) {
	for {
		time.Sleep(pollInterval * time.Second)
		for k, v := range metrics {
			body, err := http.Post("http://localhost:8080/update/gauge/"+string(k)+"/"+strconv.FormatFloat(v.GetMetrics(), 'f', -1, 64), "text/plain", nil)
			if err != nil {
				println(err)
			}
			errClose := body.Body.Close()
			if errClose != nil {
				println(errClose)
			}
		}
		body, err := http.Post("http://localhost:8080/update/counter/"+string(count.GetName())+"/"+strconv.Itoa(int(count.GetMetrics())), "text/plain", nil)
		if err != nil {
			println(err)
		}
		errClose := body.Body.Close()
		if errClose != nil {
			println(errClose)
		}
		count.ClearCount()
	}
}

func main() {
	go updateMetrics(2)
	go sendMetrics(10)
	server := http.Server{}
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}

func getMemByStats(name string) float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	if name == "Alloc" {
		return float64(memStats.Alloc)
	} else if name == "BuckHashSys" {
		return float64(memStats.BuckHashSys)
	} else if name == "Frees" {
		return float64(memStats.Frees)
	} else if name == "GCCPUFraction" {
		return memStats.GCCPUFraction
	} else if name == "GCSys" {
		return float64(memStats.GCSys)
	} else if name == "HeapAlloc" {
		return float64(memStats.HeapAlloc)
	} else if name == "HeapIdle" {
		return float64(memStats.HeapIdle)
	} else if name == "HeapInuse" {
		return float64(memStats.HeapInuse)
	} else if name == "HeapObjects" {
		return float64(memStats.HeapObjects)
	} else if name == "HeapReleased" {
		return float64(memStats.HeapReleased)
	} else if name == "HeapSys" {
		return float64(memStats.HeapSys)
	} else if name == "LastGC" {
		return float64(memStats.LastGC)
	} else if name == "Lookups" {
		return float64(memStats.Lookups)
	} else if name == "MCacheInuse" {
		return float64(memStats.MCacheInuse)
	} else if name == "MCacheSys" {
		return float64(memStats.MCacheSys)
	} else if name == "MSpanInuse" {
		return float64(memStats.MSpanInuse)
	} else if name == "MSpanSys" {
		return float64(memStats.MSpanSys)
	} else if name == "Mallocs" {
		return float64(memStats.Mallocs)
	} else if name == "NumForcedGC" {
		return float64(memStats.NumForcedGC)
	} else if name == "NumGC" {
		return float64(memStats.NumGC)
	} else if name == "OtherSys" {
		return float64(memStats.OtherSys)
	} else if name == "PauseTotalNs" {
		return float64(memStats.PauseTotalNs)
	} else if name == "StackInuse" {
		return float64(memStats.StackInuse)
	} else if name == "Sys" {
		return float64(memStats.Sys)
	} else if name == "TotalAlloc" {
		return float64(memStats.TotalAlloc)
	} else if name == "RandomValue" {
		return rand.Float64()
	} else {
		return 0
	}

}
