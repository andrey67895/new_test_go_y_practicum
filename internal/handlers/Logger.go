package handlers

import (
	"net/http"
	"time"

	"github.com/andrey67895/new_test_go_y_practicum/internal/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	log := logger.Log()
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().Local()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start).Milliseconds()
		log.Infoln(
			"Uri: ", r.RequestURI,
			"Method: ", r.Method,
			"Status: ", responseData.status,
			"Duration: ", duration,
			"Size: ", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
