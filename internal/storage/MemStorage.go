package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
)

// LocalNewMemStorageGauge глобавльная переменная для хранилища данных для Storage в формате Gauge
var LocalNewMemStorageGauge = NewMemStorageGauge()

// LocalNewMemStorageCounter глобавльная переменная для хранилища данных для Storage в формате Counter
var LocalNewMemStorageCounter = NewMemStorageCounter()

// MemStorageGauge хранилище данных для Storage в формате Gauge
type MemStorageGauge struct {
	data map[string]float64
	mut  sync.RWMutex
}

// GetData получение всех данных из MemStorageCounter
func (e *MemStorageCounter) GetData() map[string]int64 {
	e.mut.RLock()
	defer e.mut.RUnlock()

	mapCopy := make(map[string]int64, len(e.data))
	for key, value := range e.data {
		mapCopy[key] = value
	}
	return mapCopy
}

// GetData получение всех данных из MemStorageGauge
func (e *MemStorageGauge) GetData() map[string]float64 {
	e.mut.RLock()
	defer e.mut.RUnlock()

	mapCopy := make(map[string]float64, len(e.data))
	for key, value := range e.data {
		mapCopy[key] = value
	}
	return mapCopy
}

// MemStorageCounter хранилище данных для Storage в формате Counter
type MemStorageCounter struct {
	data map[string]int64
	mut  sync.RWMutex
}

// SetCounter сохранение данных в MemStorageCounter
func (e *MemStorageCounter) SetCounter(key string, value int64) error {
	e.mut.Lock()
	e.data[key] = value
	e.mut.Unlock()
	return nil
}

// SetGauge сохранение данных в MemStorageGauge
func (e *MemStorageGauge) SetGauge(key string, value float64) error {
	e.mut.Lock()
	e.data[key] = value
	e.mut.Unlock()
	return nil
}

// GetGauge получение данных из MemStorageGauge
func (e *MemStorageGauge) GetGauge(key string) (float64, error) {
	e.mut.RLock()
	defer e.mut.RUnlock()
	value, ok := e.data[key]
	if !ok {
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

func SaveDataInFile(fname string) error {
	var tModel []model.JSONMetrics
	for k, v := range LocalNewMemStorageGauge.GetData() {
		tJSON := model.JSONMetrics{}
		tJSON.ID = k
		tJSON.SetValue(v)
		tJSON.MType = "gauge"
		tModel = append(tModel, tJSON)
	}
	for k, v := range LocalNewMemStorageCounter.GetData() {
		tJSON := model.JSONMetrics{}
		tJSON.ID = k
		tJSON.SetDelta(v)
		tJSON.MType = "counter"
		tModel = append(tModel, tJSON)
	}
	data, err := json.MarshalIndent(tModel, "", "   ")
	if err != nil {
		log.Error(err.Error())
		return err
	}

	if config.StoreIntervalServer == 0 {
		err = os.MkdirAll(filepath.Dir(fname), 0666)
		if err != nil {
			log.Error(err.Error())
		}
		_, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			log.Error(err.Error())
		}
		err = os.WriteFile(fname, data, 0666)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return nil
}
