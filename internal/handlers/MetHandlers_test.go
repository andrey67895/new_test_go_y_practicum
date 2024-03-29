package handlers

import (
	"context"
	"github.com/andrey67895/new_test_go_y_practicum/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func AddChiURLParams(r *http.Request, params map[string]string) *http.Request {
	ctx := chi.NewRouteContext()
	for k, v := range params {
		ctx.URLParams.Add(k, v)
	}

	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
}

func TestMetHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/gauge/Test/10", nil), map[string]string{
					"type": "gauge", "name": "Test", "value": "10",
				}),
			},
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive test #2",
			args: args{
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Test/10", nil), map[string]string{
					"type": "counter", "name": "Test", "value": "10",
				}),
			},
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #1",
			args: args{
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/random/Test/10", nil), map[string]string{
					"type": "random", "name": "Test", "value": "10",
				}),
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative test #2",
			args: args{
				req: AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Test/none", nil), map[string]string{
					"type": "counter", "name": "Test", "value": "none",
				}),
			},
			want: want{
				code:        400,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			MetHandler(w, test.args.req)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			err := res.Body.Close()
			if err != nil {
				return
			}
		})
	}
}

func TestCountValue(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		want want
		args args
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := storage.LocalNewMemStorageCounter.SetCounter("Test", 100)
			assert.NoError(t, err)
			w := httptest.NewRecorder()
			req := AddChiURLParams(httptest.NewRequest("POST", "/update/counter/Test/100", nil), map[string]string{
				"type": "counter", "name": "Test", "value": "100",
			})

			MetHandler(w, req)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			value, err := storage.LocalNewMemStorageCounter.GetCounter("Test")
			assert.NoError(t, err)
			assert.Equal(t, 200, int(value))
			errRequest := res.Body.Close()
			if errRequest != nil {
				return
			}
		})
	}
}
