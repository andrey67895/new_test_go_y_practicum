package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
)

func Test_getMemByGopsutil(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				getMemByGopsutil()
			})
		})
	}
}

func Test_sendHashKey(t *testing.T) {
	type args struct {
		data []byte
		key  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				data: []byte("TEST"),
				key:  "KEY",
			},
		},
		{
			name: "positive test #2",
			args: args{
				data: []byte("TEST"),
				key:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.HashKeyAgent = tt.args.key
			assert.NotPanics(t, func() {
				sendHashKey(httptest.NewRequest(http.MethodPost, "http://localhost:8080/", nil), tt.args.data)
			})

		})
	}
}
