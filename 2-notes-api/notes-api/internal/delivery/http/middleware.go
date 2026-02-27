package http

import (
	"net/http"
	"time"

	"notes-api/internal/logger"
)

// LoggingMiddleware logs every incoming HTTP request.
func LoggingMiddleware(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		// Process request
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		log.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", duration.Milliseconds())
	})
}
