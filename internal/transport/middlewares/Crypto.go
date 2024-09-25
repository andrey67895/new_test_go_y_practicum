// Package handlers работа с логикой handlers
package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/andrey67895/new_test_go_y_practicum/internal/config"
	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
)

var log = logger.Log()

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

func CheckRSAAndDecrypt(h http.Handler) http.Handler {
	cryptoFn := func(w http.ResponseWriter, r *http.Request) {
		if config.CryptoKeyServer != "" {
			all, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка шифрования", http.StatusInternalServerError)
				return
			}
			private := importPrivateKey()
			v15, err := rsa.DecryptPKCS1v15(rand.Reader, private, all)
			if err != nil {
				http.Error(w, "Ошибка шифрования", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(v15))
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(cryptoFn)
}

func importPrivateKey() *rsa.PrivateKey {
	file, err := os.ReadFile(config.CryptoKeyServer)
	if err != nil {
		return nil
	}
	block, _ := pem.Decode(file)
	if block == nil {
		return nil
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Error("Error: ", err.Error())
		return nil
	}
	return private
}
