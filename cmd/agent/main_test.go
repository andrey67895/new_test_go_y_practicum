package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/model"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	handlers2 "github.com/andrey67895/new_test_go_y_practicum/internal/transport/handlers"
	"github.com/andrey67895/new_test_go_y_practicum/internal/transport/middlewares"
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

func Test_sendRequestJSON(t *testing.T) {
	type args struct {
		tJSON model.JSONMetrics
	}
	delta := int64(10)
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive test #1",
			args: args{
				tJSON: model.JSONMetrics{
					ID:    "TEST_ID",
					MType: "counter",
					Delta: &delta,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(middleware.RealIP, handlers2.WithLogging, middleware.Recoverer, middlewares.GzipHandleResponse, middlewares.WithSendsGzip, middlewares.CheckHeaderCrypto, middlewares.ResponseAddHeaderCrypto)
			r.Post("/update/", handlers2.SaveMetDataForJSON(storage.InMemStorage{}))
			server := httptest.NewServer(r)
			url := strings.ReplaceAll(server.URL, "http://", "")
			assert.NotPanics(t, func() {
				err := sendRequestJSON(url, tt.args.tJSON)
				assert.NoError(t, err)
			})
			server.CloseClientConnections()
			server.Close()
		})
	}
}

func Test_retrySendRequestJSON(t *testing.T) {
	type args struct {
		host  string
		tJSON model.JSONMetrics
	}
	delta := int64(10)
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "positive test #1",
			args: args{
				tJSON: model.JSONMetrics{
					ID:    "TEST_ID",
					MType: "counter",
					Delta: &delta,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Use(middleware.RealIP, handlers2.WithLogging, middleware.Recoverer, middlewares.GzipHandleResponse, middlewares.WithSendsGzip, middlewares.CheckHeaderCrypto, middlewares.ResponseAddHeaderCrypto)
			r.Post("/update/", handlers2.SaveMetDataForJSON(storage.InMemStorage{}))
			server := httptest.NewServer(r)
			url := strings.ReplaceAll(server.URL, "http://", "")
			assert.NotPanics(t, func() {
				err := retrySendRequestJSON(url, tt.args.tJSON)
				assert.NoError(t, err)
			})
			server.CloseClientConnections()
			server.Close()
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
