package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"accounting/internal/domain/entity"
	"accounting/internal/domain/errors"
	httptesting "accounting/internal/handler/http"
)

func TestGetUserByEmailHandlerSuccess(t *testing.T) {
	testUser := &entity.User{
		ID:    "123e4567-e89b-12d3-a456-426614174000",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockService := &httptesting.MockUserService{
		UserToReturn: testUser,
	}
	handler := NewGetUserByEmailHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/search?email=john@example.com",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response UserResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Email != "john@example.com" {
		t.Errorf("expected email %q, got %q", "john@example.com", response.Email)
	}

	if mockService.GetUserByEmailCalls != 1 {
		t.Errorf("expected 1 getUserByEmail call, got %d", mockService.GetUserByEmailCalls)
	}
}

func TestGetUserByEmailHandlerNotFound(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserByEmailHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/search?email=nonexistent@example.com",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetUserByEmailHandlerMissingEmail(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserByEmailHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/search",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetUserByEmailHandlerInvalidEmail(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserByEmailHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/search?email=invalid-email",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetUserByEmailHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserByEmailHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/users/search?email=john@example.com",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetUserByEmailHandlerServiceError(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("user", "john@example.com")
	mockService := &httptesting.MockUserService{
		LastGetUserByEmailErr: notFoundErr,
	}
	handler := NewGetUserByEmailHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/search?email=john@example.com",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
