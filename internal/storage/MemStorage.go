package storage

import (
	"errors"
	"log"
	"sync"
)

var LocalNewMemStorageGauge = NewMemStorageGauge()
var LocalNewMemStorageCounter = NewMemStorageCounter()

type MemStorageGauge struct {
	data map[string]float64
	mut  sync.RWMutex
}

func (e *MemStorageCounter) GetData() map[string]int64 {
	e.mut.RLock()
	defer e.mut.RUnlock()

	mapCopy := make(map[string]int64, len(e.data))
	for key, value := range e.data {
		mapCopy[key] = value
	}
	return mapCopy
}

func (e *MemStorageGauge) GetData() map[string]float64 {
	e.mut.RLock()
	defer e.mut.RUnlock()

	mapCopy := make(map[string]float64, len(e.data))
	for key, value := range e.data {
		mapCopy[key] = value
	}
	return mapCopy
}

type MemStorageCounter struct {
	data map[string]int64
	mut  sync.RWMutex
}

func (e *MemStorageCounter) SetCounter(key string, value int64) error {
	e.mut.Lock()
	e.data[key] = value
	e.mut.Unlock()
	return nil
}

func (e *MemStorageGauge) SetGauge(key string, value float64) error {
	e.mut.Lock()
	e.data[key] = value
	e.mut.Unlock()
	return nil
}

func (e *MemStorageGauge) GetGauge(key string) (float64, error) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	value, ok := e.data[key]
	if !ok {
		log.Println("Ошибка")
		return 0, errors.New("key not found")
	}
	return value, nil
}

func (e *MemStorageCounter) GetCounter(key string) (int64, error) {
	e.mut.RLock()
	defer e.mut.RUnlock()
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
