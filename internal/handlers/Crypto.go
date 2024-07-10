package handlers

import (
	"bytes"
	"crypto/hmac"
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
			h := hmac.New(sha256.New, []byte(config.HashKeyServer))
			h.Write(hBody)
			w.Header().Add("HashSHA256", fmt.Sprintf("%x", h.Sum(nil)))

			if !strings.EqualFold(r.Header.Get("HashSHA256"), fmt.Sprintf("%x", h)) {
				log.Error("Не соответсвует hash: сгенерированному и полученному")
				//TODO Ошибка в АТ
				//	w.WriteHeader(http.StatusBadRequest)
				//	return
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
			if r.Header.Get("HashSHA256") == "" {
				log.Error("Отсутсвует header ключ HashSHA256")
				//TODO Ошибка в АТ
				// http.Error(w, "Отсутсвует header ключ HashSHA256", http.StatusBadRequest)
				// return
			}
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(cryptoFn)
}
