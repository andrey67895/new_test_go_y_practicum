package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
)

func TestGetPing(t *testing.T) {
	type args struct {
		iStorage storage.IStorageData
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		{
			name: "positive test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			handler := GetPing(storage.InMemStorage{})
			handler(w, req)
			res := w.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)
			res.Body.Close()
		})
	}
}
