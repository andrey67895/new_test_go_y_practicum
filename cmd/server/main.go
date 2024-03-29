package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var localNewMemStorageGauge = NewMemStorageGauge()
var localNewMemStorageCounter = NewMemStorageCounter()

type MemStorageGauge struct {
	data map[string]float64
}

type MemStorageCounter struct {
	data map[string]int64
}

func (e *MemStorageCounter) SetCounter(key string, value int64) error {
	e.data[key] = value
	return nil
}

func (e *MemStorageGauge) SetGauge(key string, value float64) error {
	e.data[key] = value
	return nil
}

func (e *MemStorageGauge) GetGauge(key string) (float64, error) {
	value, ok := e.data[key]
	if !ok {
		return 0, errors.New("key not found")
	}
	return value, nil
}

func (e *MemStorageCounter) GetCounter(key string) (int64, error) {
	value, ok := e.data[key]
	if !ok {
		return 0, errors.New("key not found")
	}
	return value, nil
}

func NewMemStorageGauge() *MemStorageGauge {
	return &MemStorageGauge{
		data: make(map[string]float64),
	}
}

func NewMemStorageCounter() *MemStorageCounter {
	return &MemStorageCounter{
		data: make(map[string]int64),
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
		err := localNewMemStorageGauge.SetGauge(nameMet, float64(valueMet))
		if err != nil {
			return
		}
	} else if typeMet == "counter" {
		localCounter, err := localNewMemStorageCounter.GetCounter(nameMet)
		if err != nil {
			err := localNewMemStorageCounter.SetCounter(nameMet, int64(valueMet))
			if err != nil {
				return
			}
		} else {
			err = localNewMemStorageCounter.SetCounter(nameMet, int64(int(localCounter)+valueMet))
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
