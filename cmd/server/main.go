package main

import (
	"errors"
	"fmt"
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

func (e *MemStorage) Get(key string) (any, error) {
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
	typeMet := req.PathValue("type")
	nameMet := req.PathValue("name")
	valueMet := req.PathValue("value")
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
			local, err := strconv.Atoi(localCounter.(string))
			if err != nil {
				return
			}
			println(inRouter)
			println(local)
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
	localCounter, _ := memLocalStorage.Get(nameMet)
	fmt.Println("Converted string:", localCounter)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", JSONHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
