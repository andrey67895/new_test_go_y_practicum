package model

import (
	"sync"

	"github.com/andrey67895/new_test_go_y_practicum/internal/helpers"
)

// Metrics создание структура
type Metrics struct {
	data map[string]Gauge
	mut  sync.RWMutex
}

// NewMetrics инициализация структуры для объекта Metrics
func NewMetrics() *Metrics {
	return &Metrics{
		data: make(map[string]Gauge),
	}
}

// SetDataMetricsForMap сохранение данных в Metrics
func (e *Metrics) SetDataMetricsForMap(metricsName []string) {
	e.mut.Lock()
	for _, statName := range metricsName {
		e.data[statName] = NewGauge(statName, helpers.GetMemByStats(statName))
	}
	e.mut.Unlock()
}

// SetDataMetrics сохранение данных по ключу в Metrics
func (e *Metrics) SetDataMetrics(key string, value Gauge) {
	e.mut.Lock()
	e.data[key] = value
	e.mut.Unlock()
}

// GetDataMetrics получение всех данных из Metrics
func (e *Metrics) GetDataMetrics() map[string]Gauge {
	e.mut.RLock()
	defer e.mut.RUnlock()
	mapCopy := make(map[string]Gauge, len(e.data))
	for key, value := range e.data {
		mapCopy[key] = value
	}
	return mapCopy
}
