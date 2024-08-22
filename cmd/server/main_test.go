package main

import (
	"testing"

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

func TestSaveDataForInterval(t *testing.T) {
	type args struct {
		fname         string
		storeInterval int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				fname:         "tmp/metrics-db.json",
				storeInterval: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				SaveDataForInterval(tt.args.fname, tt.args.storeInterval)
			})
		})
	}
}
