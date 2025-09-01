package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		fmt.Printf("%s %s took %v - %d %s\n", r.Method, r.URL.Path, duration, wrapped.statusCode, http.StatusText(wrapped.statusCode))
	})
}

// Interceptor sturct
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Intercept the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
