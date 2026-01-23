package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httptesting "accounting/internal/handler/http"
)

func TestDeleteUserHandlerSuccess(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewDeleteUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	if mockService.DeleteUserCalls != 1 {
		t.Errorf("expected 1 deleteUser call, got %d", mockService.DeleteUserCalls)
	}
}

func TestDeleteUserHandlerInvalidID(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewDeleteUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodDelete,
		"/api/v1/users/not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteUserHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewDeleteUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
