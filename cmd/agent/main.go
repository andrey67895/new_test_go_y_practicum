package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/helpers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var log = logger.Log()
var metricsName = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
	"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
	"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
	"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc", "RandomValue",
}

var count = model.NewCount("PollCount", 0)
var metrics = model.NewMetrics()

func worker(name string, results chan<- int) {
	err := metrics.SetDataMetrics(name, model.NewGauge(name, getMemByStats(name)))
	if err != nil {
		log.Error(err)
	}
	results <- 1
}

func updateMetrics(pollInterval time.Duration) {
	for {
		numJobs := len(metricsName)
		jobs := make(chan int, numJobs)
		results := make(chan int, numJobs)
		for _, statName := range metricsName {
			statName := statName
			go worker(statName, results)
		}
		for j := 1; j <= numJobs; j++ {
			jobs <- j
		}
		close(jobs)
		for a := 1; a <= numJobs; a++ {
			<-results
		}

		go getMemByGopsutil()
		count.UpdateCountPlusOne()
		time.Sleep(pollInterval * time.Second)
	}
}

func sendMetrics(pollInterval time.Duration, host string) {
	for {
		time.Sleep(pollInterval * time.Second)

		var tJSON []model.JSONMetrics
		for k, v := range metrics.GetDataMetrics() {
			gauge := v.GetMetrics()
			tJSON = append(tJSON, model.JSONMetrics{
				ID:    k,
				MType: "gauge",
				Value: &gauge,
			})
		}
		err := retrySendRequestJSONFloatAll(host, tJSON)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		err = retrySendRequestJSONInt(host, "counter", count.GetName(), count.GetMetrics())
		if err != nil {
			log.Error(err.Error())
			continue
		}
		count.ClearCount()

	}
}

func sendRequestJSONFloatAll(host string, tJSON []model.JSONMetrics) error {
	url := "http://" + host + "/updates/"
	tModel, _ := json.Marshal(tJSON)
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(helpers.Compress(tModel)))
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	sendHashKey(r, tModel)
	body, err := client.Do(r)
	if err != nil {
		log.Error(err.Error())
		return err
	} else {
		errClose := body.Body.Close()
		if errClose != nil {
			log.Error(errClose.Error())
			return errClose
		}
	}
	return err
}

func retrySendRequestJSONFloatAll(host string, tJSON []model.JSONMetrics) error {
	err := sendRequestJSONFloatAll(host, tJSON)
	if err != nil {
		for i := 1; i <= 5; i = i + 2 {
			timer := time.NewTimer(time.Duration(i) * time.Second)
			t := <-timer.C
			log.Info(t.Local())
			err = sendRequestJSONFloatAll(host, tJSON)
			if err == nil {
				break
			}
		}
	}
	return err
}

func retrySendRequestJSONInt(host string, typeMetr string, nameMetr string, metrics int64) error {
	err := sendRequestJSONInt(host, typeMetr, nameMetr, metrics)
	if err != nil {
		for i := 1; i <= 5; i = i + 2 {
			timer := time.NewTimer(time.Duration(i) * time.Second)
			t := <-timer.C
			log.Info(t.Local())
			err := sendRequestJSONInt(host, typeMetr, nameMetr, metrics)
			if err == nil {
				break
			}
		}
	}
	return err
}

func sendHashKey(r *http.Request, data []byte) {
	if config.HashKeyAgent != "" {
		hBody := bytes.Clone(data)
		h := hmac.New(sha256.New, []byte(config.HashKeyServer))
		h.Write(hBody)
		r.Header.Add("HashSHA256", fmt.Sprintf("%x", h))
	}
}

func sendRequestJSONInt(host string, typeMetr string, nameMetr string, metrics int64) error {
	url := "http://" + host + "/update/"
	tJSON := model.JSONMetrics{}
	tJSON.ID = nameMetr
	tJSON.MType = typeMetr
	tJSON.SetDelta(metrics)
	tModel, _ := json.Marshal(tJSON)
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(helpers.Compress(tModel)))
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	sendHashKey(r, tModel)
	body, err := client.Do(r)
	if err != nil {
		log.Error(err.Error())
		return err
	} else {
		errClose := body.Body.Close()
		if errClose != nil {
			log.Error(errClose.Error())
			return errClose
		}
	}
	return err
}

func main() {
	config.InitAgentConfig()
	go updateMetrics(time.Duration(config.PollIntervalAgent))
	go sendMetrics(time.Duration(config.ReportIntervalAgent), config.HostAgent)
	server := http.Server{}
	log.Fatal(server.ListenAndServe())
}

func getMemByGopsutil() {
	v, _ := mem.VirtualMemory()
	err := metrics.SetDataMetrics("TotalMemory", model.NewGauge("TotalMemory", float64(v.Total)))
	if err != nil {
		log.Error(err.Error())
	}
	err = metrics.SetDataMetrics("FreeMemory", model.NewGauge("FreeMemory", float64(v.Free)))
	if err != nil {
		log.Error(err.Error())
	}
	c, _ := cpu.Percent(0, true)
	for i, percent := range c {
		err = metrics.SetDataMetrics(fmt.Sprintf("CPUutilization%d", i+1), model.NewGauge(fmt.Sprintf("CPUutilization%d", i+1), percent))
	}
	if err != nil {
		log.Error(err.Error())
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
