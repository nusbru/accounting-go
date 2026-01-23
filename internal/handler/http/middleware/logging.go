package middleware

import (
	"context"
	"net/http"
	"time"

	"accounting/internal/pkg/logger"
	"github.com/google/uuid"
)

// RequestLoggingMiddleware logs HTTP requests and responses
func RequestLoggingMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate request ID
			requestID := uuid.New().String()
			ctx := context.WithValue(r.Context(), "request-id", requestID)
			r = r.WithContext(ctx)

			// Create response writer wrapper to capture status code
			wrapped := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

			// Log request
			logger.WithField("request_id", requestID).Info(
				"HTTP request started",
				"method", r.Method,
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr,
			)

			// Record start time
			start := time.Now()

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Log response
			duration := time.Since(start)
			logger.WithField("request_id", requestID).Info(
				"HTTP request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// AddRequestIDToContext adds request ID to context
func AddRequestIDToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), "request-id", requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
