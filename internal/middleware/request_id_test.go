package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is in context
		id := GetRequestID(r.Context())
		if id == "" {
			t.Error("expected request ID in context, got empty string")
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := RequestID(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	// Verify response header is set
	responseID := w.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("expected X-Request-ID header, got empty string")
	}

	// Verify it looks like a UUID (36 characters with dashes)
	if len(responseID) != 36 {
		t.Errorf("expected UUID format (36 chars), got %d chars", len(responseID))
	}
}

func TestRequestIDMiddleware_PreservesExistingID(t *testing.T) {
	existingID := "existing-request-id-12345"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		if id != existingID {
			t.Errorf("expected preserved ID %s, got %s", existingID, id)
		}
		w.WriteHeader(http.StatusOK)
	})

	wrapped := RequestID(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	responseID := w.Header().Get("X-Request-ID")
	if responseID != existingID {
		t.Errorf("expected response header %s, got %s", existingID, responseID)
	}
}

func TestGetRequestID_NotInContext(t *testing.T) {
	ctx := context.Background()
	id := GetRequestID(ctx)
	if id != "" {
		t.Errorf("expected empty string for missing ID, got %s", id)
	}
}
