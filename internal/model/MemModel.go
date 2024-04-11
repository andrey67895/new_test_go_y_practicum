package model

import "sync"

type Gauge struct {
	name    string
	isGauge bool
	metrics float64
	mux     sync.RWMutex
}

type Count struct {
	name    string
	isGauge bool
	metrics int64
	mux     sync.RWMutex
}

func (e *Count) UpdateCountPlusOne() {
	e.mux.Lock()
	e.metrics = e.metrics + 1
	e.mux.Unlock()
}

func (e *Count) ClearCount() {
	e.mux.Lock()
	e.metrics = 0
	e.mux.Unlock()
}

func (e *Gauge) GetMetrics() float64 {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.metrics
}

func (e *Count) GetName() string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.name
}

func (e *Count) GetMetrics() int64 {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.metrics
}

func NewGauge(name string, metrics float64) Gauge {
	return Gauge{name: name, isGauge: true, metrics: metrics}
}

func NewCount(name string, metrics int64) Count {
	return Count{name: name, metrics: metrics}
}
