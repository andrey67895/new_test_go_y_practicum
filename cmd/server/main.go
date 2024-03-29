package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var memLocalStorage = NewMemStorage()

type MemStorage struct {
	data map[string]string
}

func (e *MemStorage) Set(key string, value string) error {
	e.data[key] = value
	return nil
}

func (e *MemStorage) Get(key string) (string, error) {
	value, ok := e.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]string),
	}
}

func JSONHandler(w http.ResponseWriter, req *http.Request) {
	println(req.URL.Path)
	typeMet := chi.URLParam(req, "type")
	nameMet := chi.URLParam(req, "name")
	valueMet := chi.URLParam(req, "value")
	if _, err := strconv.Atoi(valueMet); err != nil {
		http.Error(w, "Неверный значение метрики! Допустимые числовые значения!", http.StatusBadRequest)
	}
	if typeMet == "gauge" {
		err := memLocalStorage.Set(nameMet, valueMet)
		if err != nil {
			return
		}
	} else if typeMet == "counter" {
		localCounter, _ := memLocalStorage.Get(nameMet)
		if localCounter == "" {
			err := memLocalStorage.Set(nameMet, valueMet)
			if err != nil {
				return
			}
		} else {
			inRouter, err := strconv.Atoi(valueMet)
			if err != nil {
				return
			}
			local, err := strconv.Atoi(localCounter)
			if err != nil {
				return
			}
			value := strconv.Itoa(local + inRouter)

			err = memLocalStorage.Set(nameMet, value)
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
	r.Post("/update/{type}/{name}/{value}", JSONHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return
	}
}
