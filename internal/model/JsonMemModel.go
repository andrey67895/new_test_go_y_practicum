package model

type JSONMetrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (e *JSONMetrics) SetValue(v float64) {
	e.value = &v
}

func (e *JSONMetrics) SetDelta(d int64) {
	e.delta = &d
}

func (e *JSONMetrics) GetValue() float64 {
	return *e.value
}

func (e *JSONMetrics) GetDelta() int64 {
	return *e.delta
}
