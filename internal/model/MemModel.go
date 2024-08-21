package model

import "sync"

// Gauge создание структуры
type Gauge struct {
	name    string
	isGauge bool
	metrics float64
}

// Count создание структуры
type Count struct {
	name    string
	isGauge bool
	metrics int64
	mux     sync.RWMutex
}

// UpdateCountPlusOne обновление метрики Count + 1
func (e *Count) UpdateCountPlusOne() {
	e.mux.Lock()
	e.metrics = e.metrics + 1
	e.mux.Unlock()
}

// ClearCount обнуление метрики Count
func (e *Count) ClearCount() {
	e.mux.Lock()
	e.metrics = 0
	e.mux.Unlock()
}

// GetMetrics получение данным метрики для Gauge
func (e *Gauge) GetMetrics() float64 {
	return e.metrics
}

// GetName получение названия для Count
func (e *Count) GetName() string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.name
}

// GetMetrics получение метрики для Count
func (e *Count) GetMetrics() int64 {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.metrics
}

// NewGauge инициализация объекта Gauge
func NewGauge(name string, metrics float64) Gauge {
	return Gauge{name: name, isGauge: true, metrics: metrics}
}

// NewCount инициализация объекта Count
func NewCount(name string, metrics int64) Count {
	return Count{name: name, metrics: metrics}
}
