package model

import "sync"

type Metrics struct {
	mut  sync.RWMutex
	data map[string]Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		data: make(map[string]Gauge),
	}
}

func (e *Metrics) SetDataMetrics(key string, value Gauge) error {
	e.mut.Lock()
	e.data[key] = value
	e.mut.Unlock()
	return nil
}

func (e *Metrics) GetDataMetrics() map[string]Gauge {
	e.mut.RLock()
	defer e.mut.RUnlock()
	mapCopy := make(map[string]Gauge, len(e.data))
	for key, value := range e.data {
		mapCopy[key] = value
	}
	return mapCopy
}
