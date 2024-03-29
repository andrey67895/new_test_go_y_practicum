package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var localNewMemStorageGauge = NewMemStorageGauge()
var localNewMemStorageCounter = NewMemStorageCounter()

type Gauge float64
type Counter int64

type MemStorageGauge struct {
	data map[string]Gauge
}

type MemStorageCounter struct {
	data map[string]Counter
}

func (e *MemStorageCounter) SetCounter(key string, value Counter) error {
	e.data[key] = value
	return nil
}

func (e *MemStorageGauge) SetGauge(key string, value Gauge) error {
	e.data[key] = value
	return nil
}

func (e *MemStorageGauge) GetGauge(key string) (Gauge, error) {
	value, ok := e.data[key]
	if !ok {
		return 0, errors.New("key not found")
	}
	return value, nil
}

func (e *MemStorageCounter) GetCounter(key string) (Counter, error) {
	value, ok := e.data[key]
	if !ok {
		return 0, errors.New("key not found")
	}
	return value, nil
}

func NewMemStorageGauge() *MemStorageGauge {
	return &MemStorageGauge{
		data: make(map[string]Gauge),
	}
}

func NewMemStorageCounter() *MemStorageCounter {
	return &MemStorageCounter{
		data: make(map[string]Counter),
	}
}

func MetHandler(w http.ResponseWriter, req *http.Request) {
	println(req.URL.Path)
	typeMet := chi.URLParam(req, "type")
	nameMet := chi.URLParam(req, "name")
	valueMet, err := strconv.Atoi(chi.URLParam(req, "value"))
	if err != nil {
		http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
		return
	}
	if typeMet == "gauge" {
		err := localNewMemStorageGauge.SetGauge(nameMet, Gauge(valueMet))
		if err != nil {
			return
		}
	} else if typeMet == "counter" {
		localCounter, err := localNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := localNewMemStorageCounter.SetCounter(nameMet, Counter(valueMet))
			if err != nil {
				return
			}
		} else {
			err = localNewMemStorageCounter.SetCounter(nameMet, Counter(int(localCounter)+valueMet))
			if err != nil {
				return
			}
		}
	} else {
		http.Error(w, "Неверный тип метрики! Допустимые значения: gauge, counter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", MetHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
