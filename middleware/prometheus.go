package middleware

import (
	"net/http"
	"src/etheus"
	"strconv"
	"time"
)

// Wrap the response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Middleware to instrument handlers with Prometheus metrics
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w}

		// Call the next handler in the chain
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		// Record metrics
		etheus.RequestCounter.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrappedWriter.statusCode)).Inc()
		etheus.RequestDurationHistogram.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrappedWriter.statusCode)).Observe(duration.Seconds())
	})
}
