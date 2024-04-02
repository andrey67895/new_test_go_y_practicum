package storage

import "errors"

var LocalNewMemStorageGauge = NewMemStorageGauge()
var LocalNewMemStorageCounter = NewMemStorageCounter()

type MemStorageGauge struct {
	data map[string]float64
}

func (e *MemStorageCounter) GetData() map[string]int64 {
	return e.data
}

func (e *MemStorageGauge) GetData() map[string]float64 {
	return e.data
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
