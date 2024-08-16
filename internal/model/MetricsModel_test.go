package model

import (
	"github.com/stretchr/testify/assert"
	"slices"
	"sync"
	"testing"
)

func TestMetrics_GetDataMetrics(t *testing.T) {
	type fields struct {
		mut  sync.RWMutex
		data map[string]Gauge
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test #1",
			fields: fields{
				mut:  sync.RWMutex{},
				data: map[string]Gauge{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Metrics{
				mut:  tt.fields.mut,
				data: tt.fields.data,
			}
			key := "1"
			gauge := Gauge{name: "TEST", metrics: 10, isGauge: true}
			e.SetDataMetrics(key, gauge)
			for tKey, value := range e.GetDataMetrics() {
				assert.Equal(t, tKey, key)
				assert.Equal(t, gauge.name, value.name)
				assert.Equal(t, gauge.metrics, value.metrics)
				assert.Equal(t, gauge.isGauge, value.isGauge)
			}
		})
	}
}

func TestMetrics_SetDataMetricsForMap(t *testing.T) {
	type fields struct {
		mut  sync.RWMutex
		data map[string]Gauge
	}
	metricsName := []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
		"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
		"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue",
	}
	type args struct {
		metricsName []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			fields: fields{
				mut:  sync.RWMutex{},
				data: map[string]Gauge{},
			},
			args: args{
				metricsName: metricsName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Metrics{
				mut:  tt.fields.mut,
				data: tt.fields.data,
			}
			e.SetDataMetricsForMap(tt.args.metricsName)
			assert.True(t, len(e.GetDataMetrics()) > 0)
			for tKey := range e.GetDataMetrics() {
				found := slices.Contains(tt.args.metricsName, tKey)
				assert.True(t, found)
			}
		})
	}
}
