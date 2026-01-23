package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"accounting/internal/pkg/logger"
)

func TestRecoveryMiddleware_NoPanic(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithWriter(&buf, "json", "info")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrapped := Recovery(log)(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "success" {
		t.Errorf("expected body 'success', got %s", w.Body.String())
	}
}

func TestRecoveryMiddleware_WithPanic(t *testing.T) {
	var buf bytes.Buffer
	log := logger.NewWithWriter(&buf, "json", "error")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrapped := Recovery(log)(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic
	wrapped.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Verify response is RFC 7807 problem detail
	var problem ProblemDetail
	if err := json.Unmarshal(w.Body.Bytes(), &problem); err != nil {
		t.Fatalf("failed to unmarshal problem detail: %v", err)
	}

	if problem.Status != http.StatusInternalServerError {
		t.Errorf("expected problem status 500, got %d", problem.Status)
	}
	if problem.Title != "Internal Server Error" {
		t.Errorf("expected title 'Internal Server Error', got %s", problem.Title)
	}

	// Verify Content-Type header
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/problem+json" {
		t.Errorf("expected Content-Type 'application/problem+json', got %s", contentType)
	}

	// Verify panic was logged
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Panic recovered") {
		t.Errorf("expected 'Panic recovered' in log, got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "test panic") {
		t.Errorf("expected panic message in log, got: %s", logOutput)
	}
}
