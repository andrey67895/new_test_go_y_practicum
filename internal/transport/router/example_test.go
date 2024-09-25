package router

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/helpers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
)

var client = &http.Client{}

func Example_ping() {
	url := "http://" + config.HostServer + "/ping"
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	body, _ := client.Do(r)
	body.Body.Close()
}

func Example_allDataMain() {
	url := "http://" + config.HostServer + "/"
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	body, _ := client.Do(r)
	body.Body.Close()
}

func Example_valuePathParam() {
	tType := "gauge"
	name := "NAME_METRICS"
	url := "http://" + config.HostServer + "value/" + tType + "/" + name
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	body, _ := client.Do(r)
	body.Body.Close()
}

func Example_value() {
	tJSON := model.JSONMetrics{
		ID:    "NAME_METRICS",
		MType: "gauge",
	}
	tJSON.SetValue(10.0)
	tModel, _ := json.Marshal(tJSON)
	url := "http://" + config.HostServer + "/value/"
	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(helpers.Compress(tModel)))
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	body, _ := client.Do(r)
	body.Body.Close()
}

func Example_updates() {
	var tJSON []model.JSONMetrics
	tJSON = append(tJSON, model.JSONMetrics{
		ID:    "NAME_METRICS",
		MType: "gauge",
	})
	tJSON = append(tJSON, model.JSONMetrics{
		ID:    "NAME_METRICS2",
		MType: "gauge",
	})
	url := "http://" + config.HostServer + "/updates/"
	tModel, _ := json.Marshal(tJSON)
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(helpers.Compress(tModel)))
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	body, _ := client.Do(r)
	body.Body.Close()
}

func Example_update() {
	tJSON := model.JSONMetrics{}
	url := "http://" + config.HostServer + "/update/"
	tModel, _ := json.Marshal(tJSON)
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(helpers.Compress(tModel)))
	r.Header.Add("Content-Encoding", "gzip")
	r.Header.Add("Content-Type", "application/json")
	body, _ := client.Do(r)
	body.Body.Close()
}
