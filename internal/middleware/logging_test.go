package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"accounting/internal/pkg/logger"
)

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithWriter(&buf, "json", "info")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Chain with RequestID first
	wrapped := RequestID(Logging(log)(handler))

	req := httptest.NewRequest(http.MethodGet, "/test/path", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	// Verify log output contains expected fields
	logOutput := buf.String()
	if !strings.Contains(logOutput, "HTTP request") {
		t.Errorf("expected log message 'HTTP request', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "/test/path") {
		t.Errorf("expected path in log, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "GET") {
		t.Errorf("expected method in log, got: %s", logOutput)
	}
}

func TestLoggingMiddleware_CapturesStatusCode(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithWriter(&buf, "json", "info")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	wrapped := Logging(log)(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "404") {
		t.Errorf("expected status 404 in log, got: %s", logOutput)
	}
}
