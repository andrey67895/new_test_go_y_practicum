package storage

import "fmt"

type InMemStorage struct{}

func (mem InMemStorage) RetrySaveGauge(id string, delta float64) error {
	return LocalNewMemStorageGauge.SetGauge(id, delta)
}

func (mem InMemStorage) RetrySaveCounter(id string, value int64) error {
	return LocalNewMemStorageCounter.SetCounter(id, value)
}

func (mem InMemStorage) Ping() error {
	return nil
}

func (mem InMemStorage) GetCounter(id string) (int64, error) {
	return LocalNewMemStorageCounter.GetCounter(id)
}

func (mem InMemStorage) GetGauge(id string) (float64, error) {
	return LocalNewMemStorageGauge.GetGauge(id)
}

func (mem InMemStorage) GetData() (string, error) {
	var data = ""
	for k, v := range LocalNewMemStorageGauge.GetData() {
		data = fmt.Sprintf("%s Name: %s. Delta: %f \n", data, k, v)
	}
	for k, v := range LocalNewMemStorageCounter.GetData() {
		data = fmt.Sprintf("%s Name: %s. Value: %d", data, k, v)
	}
	return data, nil
}
