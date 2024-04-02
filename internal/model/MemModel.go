package model

type Gauge struct {
	name    string
	isGauge bool
	metrics float64
}

type Count struct {
	name    string
	isGauge bool
	metrics int64
}

func (e *Count) UpdateCountPlusOne() {
	e.metrics = e.metrics + 1
}

func (e *Gauge) GetMetrics() float64 {
	return e.metrics
}

func (e *Count) GetMetrics() int64 {
	return e.metrics
}

func NewGauge(name string, metrics float64) Gauge {
	return Gauge{name, true, metrics}
}

func NewCount(name string, metrics int64) Count {
	return Count{name, false, metrics}
}
