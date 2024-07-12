package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
)

func ResponseAddHeaderCrypto(h http.Handler) http.Handler {
	cryptoFn := func(w http.ResponseWriter, r *http.Request) {
		if config.HashKeyServer != "" {
			body, _ := io.ReadAll(r.Body)
			hBody := bytes.Clone(body)
			h := hmac.New(sha256.New, []byte(config.HashKeyServer))
			h.Write(hBody)
			w.Header().Add("HashSHA256", fmt.Sprintf("%x", h.Sum(nil)))
			hash := r.Header.Get("HashSHA256")
			if hash != "" {
				if hash != fmt.Sprintf("%x", h.Sum(nil)) {
					log.Error("Не соответсвует hash: сгенерированному и полученному")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			r.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(cryptoFn)
}

func CheckHeaderCrypto(h http.Handler) http.Handler {
	cryptoFn := func(w http.ResponseWriter, r *http.Request) {
		if config.HashKeyServer != "" {
			hash := r.Header.Get("HashSHA256")
			if hash != "" {
				matched, _ := regexp.Match(`^[a-fA-F0-9]{64}$`, []byte(hash))
				if !matched {
					log.Error("HashSHA256 не валиден:", hash)
					http.Error(w, "HashSHA256 не валиден", http.StatusBadRequest)
					return
				}
			}
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(cryptoFn)
}
