package router

import (
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
)

func WithCrypto(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		if config.HashKeyServer != "" {
			h := sha256.New()
			buf := make([]byte, 8)
			if _, err := io.ReadFull(r.Body, buf); err != nil {
				return
			}
			h.Write(buf)
			dst := h.Sum(nil)
			w.Header().Add("HashSHA256", string(dst))
		}
		h.ServeHTTP(w, r)

	}

	return http.HandlerFunc(logFn)
}
