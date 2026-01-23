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

func TestUpdateUserHandlerSuccess(t *testing.T) {
	testUser := &entity.User{
		ID:    "123e4567-e89b-12d3-a456-426614174000",
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
	mockService := &httptesting.MockUserService{
		UserToReturn: testUser,
	}
	handler := NewUpdateUserHandler(mockService)

	reqBody := UpdateUserRequest{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response UserResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Name != "Jane Doe" {
		t.Errorf("expected name %q, got %q", "Jane Doe", response.Name)
	}

	if mockService.UpdateUserCalls != 1 {
		t.Errorf("expected 1 updateUser call, got %d", mockService.UpdateUserCalls)
	}
}

func TestUpdateUserHandlerPartialUpdate(t *testing.T) {
	testUser := &entity.User{
		ID:    "123e4567-e89b-12d3-a456-426614174000",
		Name:  "Jane Doe",
		Email: "test@example.com",
	}
	mockService := &httptesting.MockUserService{
		UserToReturn: testUser,
	}
	handler := NewUpdateUserHandler(mockService)

	reqBody := UpdateUserRequest{
		Name: "Jane Doe",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateUserHandlerEmptyBody(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewUpdateUserHandler(mockService)

	reqBody := UpdateUserRequest{}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserHandlerInvalidEmail(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewUpdateUserHandler(mockService)

	reqBody := UpdateUserRequest{
		Email: "invalid-email",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUserHandlerNotFound(t *testing.T) {
	notFoundErr := errors.NewErrNotFound("user", "123e4567-e89b-12d3-a456-426614174000")
	mockService := &httptesting.MockUserService{
		LastUpdateUserErr: notFoundErr,
	}
	handler := NewUpdateUserHandler(mockService)

	reqBody := UpdateUserRequest{
		Name: "Jane Doe",
	}

	req := httptesting.NewTestRequest(http.MethodPut, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000", reqBody)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateUserHandlerInvalidMethod(t *testing.T) {
	mockService := &httptesting.MockUserService{}
	handler := NewUpdateUserHandler(mockService)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()

	handler.Handle(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}
