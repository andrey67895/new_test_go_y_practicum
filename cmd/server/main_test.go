package main

import (
	"testing"

	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRestoringDataFromFile(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				fname: "tmp/metrics-db.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				RestoringDataFromFile(tt.args.fname)
			})
		})
	}
}

func TestSaveDataInFile(t *testing.T) {
	type args struct {
		fname string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				fname: "tmp/metrics-db.json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				SaveDataInFile(tt.args.fname)
			})
		})
	}
}

func TestSaveData(t *testing.T) {
	type args struct {
		tModel model.JSONMetrics
	}
	value := 10.0
	delta := int64(10)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				model.JSONMetrics{
					ID:    "ID1",
					MType: "gauge",
					Value: &value,
				},
			},
		},
		{
			name: "positive test #2",
			args: args{
				model.JSONMetrics{
					ID:    "ID2",
					MType: "counter",
					Delta: &delta,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t,
				func() {
					SaveData(tt.args.tModel)
				})
		})
	}
}
