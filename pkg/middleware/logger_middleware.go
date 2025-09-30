package middleware

import (
	"flight-api/pkg/logger"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func HTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := logger.NewLogger(logger.INFO_DEBUG_LEVEL)

		// Wrap responseWriter to capture status code
		ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)

		logger.WithFields(
			logrus.Fields{
				"method":      r.Method,
				"path":        r.URL,
				"status":      ww.status,
				"duration":    time.Since(start),
				"remote_addr": r.RemoteAddr,
			},
		).Info("HTTP Request")
	})
}
