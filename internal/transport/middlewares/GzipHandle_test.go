package middlewares

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestWithSendsGzip(t *testing.T) {
	type TestJson struct {
		Test string `json:"test,omitempty"`
	}
	tests := []struct {
		name     string
		compress bool
	}{
		{
			name:     "positive test #1",
			compress: true,
		},
		{
			name:     "positive test #1",
			compress: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(TestJson{Test: "test"})

			var req *http.Request

			if tt.compress {
				req, _ = http.NewRequest("POST", "/", bytes.NewReader(compress(b)))
				req.Header.Add("Content-Encoding", "gzip")
			} else {
				req, _ = http.NewRequest("POST", "/", bytes.NewReader(b))
			}

			w := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Use(WithSendsGzip)

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				all, err := io.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
				}
				defer r.Body.Close()
				w.WriteHeader(http.StatusOK)
				w.Write(all)
			})
			r.ServeHTTP(w, req)

			if w.Code != 200 {
				t.Fatal("Response Code should be 200")
			}
			assert.Equal(t, string(b), w.Body.String())
		})
	}
}

func compress(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	if _, err := w.Write(data); err != nil {
		return nil
	}
	if err := w.Close(); err != nil {
		return nil
	}
	return b.Bytes()
}
