package storage

import (
	"context"
	"fmt"
)

type InMemStorage struct{}

func (mem InMemStorage) RetrySaveGauge(_ context.Context, id string, delta float64) error {
	return LocalNewMemStorageGauge.SetGauge(id, delta)
}

func (mem InMemStorage) RetrySaveCounter(_ context.Context, id string, value int64) error {
	return LocalNewMemStorageCounter.SetCounter(id, value)
}

func (mem InMemStorage) Ping() error {
	return nil
}

func (mem InMemStorage) GetCounter(_ context.Context, id string) (int64, error) {
	return LocalNewMemStorageCounter.GetCounter(id)
}

func (mem InMemStorage) GetGauge(_ context.Context, id string) (float64, error) {
	return LocalNewMemStorageGauge.GetGauge(id)
}

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
