package middleware

import (
	"net/http"
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
func (m *Middleware) PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, route := range m.ExcludeRoutes {
			if r.URL.Path == route {
				next.ServeHTTP(w, r)
				return
			}
		}
		start := time.Now()

		// Wrap the response writer to capture status code
		wrappedWriter := &responseWriter{ResponseWriter: w}

		// Call the next handler in the chain
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)

		// Record metrics
		m.RequestCounter.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrappedWriter.statusCode)).Inc()
		m.RequestDurationHistogram.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(wrappedWriter.statusCode)).Observe(duration.Seconds())
	})
}
