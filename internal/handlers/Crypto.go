package handlers

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
)

func WithCrypto(h http.Handler) http.Handler {
	cryptoFn := func(w http.ResponseWriter, r *http.Request) {
		if config.HashKeyServer != "" {
			body, _ := io.ReadAll(r.Body)
			hBody := bytes.Clone(body)
			hBody = append(hBody, []byte(config.HashKeyServer)...)
			h := sha256.Sum256(hBody)
			w.Header().Add("HashSHA256", fmt.Sprintf("%x", h))
			if !strings.EqualFold(r.Header.Get("HashSHA256"), fmt.Sprintf("%x", h)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(cryptoFn)
}
