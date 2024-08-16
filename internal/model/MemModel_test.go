package model

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCount_ClearCount(t *testing.T) {
	tests := []struct {
		name    string
		metrics int64
	}{
		{
			name:    "positive test #1",
			metrics: 0,
		},
		{
			name:    "positive test #2",
			metrics: 10,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mut := sync.RWMutex{}
			e := &Count{
				name:    test.name,
				isGauge: false,
				metrics: test.metrics,
				mux:     mut,
			}
			e.ClearCount()
			assert.Equal(t, 0, int(e.metrics))
		})
	}
}

func TestCount_GetMetrics(t *testing.T) {
	type fields struct {
		name    string
		isGauge bool
		metrics int64
		mux     sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive test #1",
			fields: fields{
				name:    "positive test #1",
				isGauge: false,
				metrics: 10,
				mux:     sync.RWMutex{},
			},
		},
		{
			name: "positive test #2",
			fields: fields{
				name:    "positive test #2",
				isGauge: true,
				metrics: 10,
				mux:     sync.RWMutex{},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := &Count{
				name:    test.fields.name,
				isGauge: test.fields.isGauge,
				metrics: test.fields.metrics,
				mux:     test.fields.mux,
			}
			assert.Equal(t, e.GetMetrics(), test.fields.metrics)
		})
	}
}

func TestCount_GetName(t *testing.T) {
	type fields struct {
		name    string
		isGauge bool
		metrics int64
		mux     sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "positive test #1",
			fields: fields{
				name:    "positive test #1",
				isGauge: false,
				metrics: 10,
				mux:     sync.RWMutex{},
			},
		},
		{
			name: "positive test #2",
			fields: fields{
				name:    "positive test #2",
				isGauge: true,
				metrics: 10,
				mux:     sync.RWMutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Count{
				name:    tt.fields.name,
				isGauge: tt.fields.isGauge,
				metrics: tt.fields.metrics,
				mux:     tt.fields.mux,
			}
			assert.Equal(t, e.GetName(), tt.fields.name)
		})
	}
}

func TestCount_UpdateCountPlusOne(t *testing.T) {
	type fields struct {
		name    string
		isGauge bool
		metrics int64
		mux     sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "positive test #1",
			fields: fields{
				name:    "positive test #1",
				isGauge: false,
				metrics: 10,
				mux:     sync.RWMutex{},
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Count{
				name:    tt.fields.name,
				isGauge: tt.fields.isGauge,
				metrics: tt.fields.metrics,
				mux:     tt.fields.mux,
			}
			e.UpdateCountPlusOne()
			assert.Equal(t, tt.want, int(e.GetMetrics()))
		})
	}
}

func TestGauge_GetMetrics(t *testing.T) {
	type fields struct {
		name    string
		isGauge bool
		metrics float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "positive test #1",
			fields: fields{
				name:    "positive test #1",
				isGauge: true,
				metrics: 10.0,
			},
			want: 10.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Gauge{
				name:    tt.fields.name,
				isGauge: tt.fields.isGauge,
				metrics: tt.fields.metrics,
			}
			assert.Equal(t, tt.want, e.GetMetrics())
		})
	}
}

func TestNewGauge(t *testing.T) {
	type args struct {
		name    string
		metrics float64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				name:    "positive test #1",
				metrics: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gauge := NewGauge(tt.args.name, tt.args.metrics)
			assert.Equal(t, tt.args.name, gauge.name)
			assert.Equal(t, tt.args.metrics, gauge.metrics)
		})
	}
}
