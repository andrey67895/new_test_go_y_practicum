package storage

import (
	"context"
	"fmt"
)

// InMemStorage инициализация InMemory Storage
type InMemStorage struct{}

// RetrySaveGauge сохранение Gauge InMemory
func (mem InMemStorage) RetrySaveGauge(_ context.Context, id string, delta float64) error {
	return LocalNewMemStorageGauge.SetGauge(id, delta)
}

// RetrySaveCounter сохранение Counter InMemory
func (mem InMemStorage) RetrySaveCounter(_ context.Context, id string, value int64) error {
	return LocalNewMemStorageCounter.SetCounter(id, value)
}

// Ping заглушка для InMemory
func (mem InMemStorage) Ping() error {
	return nil
}

// GetCounter получение данных Counter по id для InMemory
func (mem InMemStorage) GetCounter(_ context.Context, id string) (int64, error) {
	return LocalNewMemStorageCounter.GetCounter(id)
}

// GetGauge получение данных Gauge по id для InMemory
func (mem InMemStorage) GetGauge(_ context.Context, id string) (float64, error) {
	return LocalNewMemStorageGauge.GetGauge(id)
}

// GetData получение всех данных хранимых в InMemory
func (mem InMemStorage) GetData(_ context.Context) (string, error) {
	var data = ""
	for k, v := range LocalNewMemStorageGauge.GetData() {
		data = fmt.Sprintf("%s Name: %s. Delta: %f \n", data, k, v)
	}
	for k, v := range LocalNewMemStorageCounter.GetData() {
		data = fmt.Sprintf("%s Name: %s. Value: %d", data, k, v)
	}
	return data, nil
}
