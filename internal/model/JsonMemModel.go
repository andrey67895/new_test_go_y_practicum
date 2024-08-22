// Package model набор стркутур для общения с клиентом, сервисом и между собой
package model

// JSONMetrics создание структуры для JSON варианта
type JSONMetrics struct {
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
}

// SetValue обновление поля Value
func (e *JSONMetrics) SetValue(v float64) {
	e.Value = &v
}

// SetDelta обновление поля Delta
func (e *JSONMetrics) SetDelta(d int64) {
	e.Delta = &d
}

// GetValue получение поля Value
func (e *JSONMetrics) GetValue() float64 {
	return *e.Value
}

// GetDelta получение поля Delta
func (e *JSONMetrics) GetDelta() int64 {
	return *e.Delta
}
