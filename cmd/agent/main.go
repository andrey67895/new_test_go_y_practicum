package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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

func updateMetrics(pollInterval time.Duration) {
	for {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			metrics.SetDataMetricsForMap(metricsName)
		}()
		go func() {
			defer wg.Done()
			getMemByGopsutil()
		}()
		wg.Wait()
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
	metrics.SetDataMetrics("TotalMemory", model.NewGauge("TotalMemory", float64(v.Total)))
	metrics.SetDataMetrics("FreeMemory", model.NewGauge("FreeMemory", float64(v.Free)))
	c, _ := cpu.Percent(0, true)
	for i, percent := range c {
		metrics.SetDataMetrics(fmt.Sprintf("CPUutilization%d", i+1), model.NewGauge(fmt.Sprintf("CPUutilization%d", i+1), percent))
	}
}
