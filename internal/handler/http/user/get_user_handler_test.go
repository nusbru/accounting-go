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

func TestGetUserHandlerSuccess(t *testing.T) {
	testUser := &entity.User{
		ID:    "123e4567-e89b-12d3-a456-426614174000",
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockService := &httptesting.MockUserService{
		UserToReturn: testUser,
	}
	handler := NewGetUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000",
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

	if response.ID != testUser.ID {
		t.Errorf("expected ID %q, got %q", testUser.ID, response.ID)
	}

	if response.Name != testUser.Name {
		t.Errorf("expected name %q, got %q", testUser.Name, response.Name)
	}

	if response.Email != testUser.Email {
		t.Errorf("expected email %q, got %q", testUser.Email, response.Email)
	}

	if mockService.GetUserCalls != 1 {
		t.Errorf("expected 1 getUser call, got %d", mockService.GetUserCalls)
	}
}

func TestGetUserHandlerNotFound(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetUserHandlerInvalidID(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/not-a-uuid",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetUserHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewGetUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodPost,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGetUserHandlerServiceError(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("user", "123e4567-e89b-12d3-a456-426614174000")
	mockService := &httptesting.MockUserService{
		LastGetUserErr: notFoundErr,
	}
	handler := NewGetUserHandler(mockService)

	req, _ := http.NewRequest(
		http.MethodGet,
		"/api/v1/users/123e4567-e89b-12d3-a456-426614174000",
		nil,
	)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
