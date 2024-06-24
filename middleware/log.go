package middleware

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	Logger        *log.Logger
	ExcludeRoutes []string
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		startTime := time.Now()

		// Wrap the response writer
		lrw := NewLogResponseWriter(w)

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log the request
		m.Logger.Printf("%s %s %d %s", r.Method, r.URL.EscapedPath(), lrw.statusCode, time.Since(startTime))
	})
}

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	// Write the default status code
	lrw := &LogResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	return lrw
}

func (lrw *LogResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
