package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONMetrics_GetDelta(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	delta := int64(10)
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "positive test #1",
			fields: fields{
				ID:    "positive test #1",
				MType: "MType",
				Delta: &delta,
			},
			want: delta,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &JSONMetrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			assert.Equal(t, tt.want, e.GetDelta())
		})
	}
}

func TestJSONMetrics_GetValue(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	value := 10.0
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "positive test #1",
			fields: fields{
				ID:    "positive test #1",
				MType: "MType",
				Value: &value,
			},
			want: value,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &JSONMetrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			assert.Equal(t, tt.want, e.GetValue())
		})
	}
}

func TestJSONMetrics_SetDelta(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	type args struct {
		d int64
	}
	delta := int64(10)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			fields: fields{
				ID:    "positive test #1",
				MType: "MType",
				Delta: &delta,
			},
			args: args{
				d: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &JSONMetrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			assert.Equal(t, *tt.fields.Delta, e.GetDelta())
			e.SetDelta(tt.args.d)
			assert.Equal(t, tt.args.d, e.GetDelta())
		})
	}
}

func TestJSONMetrics_SetValue(t *testing.T) {
	type fields struct {
		ID    string
		MType string
		Delta *int64
		Value *float64
	}
	type args struct {
		v float64
	}
	value := 10.0
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "positive test #1",
			fields: fields{
				ID:    "positive test #1",
				MType: "MType",
				Value: &value,
			},
			args: args{
				v: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &JSONMetrics{
				ID:    tt.fields.ID,
				MType: tt.fields.MType,
				Delta: tt.fields.Delta,
				Value: tt.fields.Value,
			}
			assert.Equal(t, *tt.fields.Value, e.GetValue())
			e.SetValue(tt.args.v)
			assert.Equal(t, tt.args.v, e.GetValue())
		})
	}
}
