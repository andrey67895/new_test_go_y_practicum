package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMemByStats(t *testing.T) {
	type args struct {
		name []string
	}
	var metricsName = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc",
		"HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC",
		"Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs",
		"NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
		"StackSys", "Sys", "TotalAlloc", "RandomValue",
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "positive test #1",
			args: args{
				name: metricsName,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, value := range tt.args.name {
				assert.NotPanics(t, func() {
					GetMemByStats(value)
				})
			}
		})
	}
}
