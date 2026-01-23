package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"accounting/internal/pkg/logger"
)

// ProblemDetail represents RFC 7807 Problem Details for HTTP APIs.
type ProblemDetail struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// Recovery returns a middleware that recovers from panics and returns a proper error response.
func Recovery(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Get stack trace
					stack := debug.Stack()

					// Extract request ID from context
					requestID := GetRequestID(r.Context())

					// Log the panic
					log.Error("Panic recovered",
						"request_id", requestID,
						"error", err,
						"stack", string(stack),
						"method", r.Method,
						"path", r.URL.Path,
					)

					// Return RFC 7807 problem detail response
					problem := ProblemDetail{
						Type:     "https://api.accounting.app/problems/internal-error",
						Title:    "Internal Server Error",
						Status:   http.StatusInternalServerError,
						Detail:   "An unexpected error occurred. Please try again later.",
						Instance: r.URL.Path,
					}

					w.Header().Set("Content-Type", "application/problem+json")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(problem)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
