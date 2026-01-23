package user

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type UpdateUserHandler struct {
	service interfaces.UserService
}

func NewUpdateUserHandler(service interfaces.UserService) *UpdateUserHandler {
	return &UpdateUserHandler{service: service}
}

// Handle updates an existing user
// @Summary Update a user
// @Description Update user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param user body UpdateUserRequest true "User update request"
// @Success 200 {object} UserResponse
// @Failure 400 {object} common.ValidationProblem
// @Failure 404 {object} common.ProblemDetail
// @Failure 500 {object} common.ProblemDetail
// @Router /api/v1/users/{id} [put]
func (h *UpdateUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.URL.Path))
		return
	}

	id := extractID(r.URL.Path, "/api/v1/users/")
	if err := common.ValidateUUID(id, "id"); err != nil {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.URL.Path, []common.ValidationError{*err}))
		return
	}

	var req UpdateUserRequest
	if r.Body == nil {
		common.WriteProblem(w, common.NewBadRequestProblem("Invalid JSON format", r.URL.Path))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteProblem(w, common.NewBadRequestProblem("Invalid JSON format", r.URL.Path))
		return
	}

	// Validate input (at least one field should be provided)
	validationErrors := []common.ValidationError{}
	if req.Name != "" {
		if err := common.ValidateStringLength(req.Name, "name", 1, 100); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}
	if req.Email != "" {
		if err := common.ValidateEmail(req.Email, "email"); err != nil {
			validationErrors = append(validationErrors, *err)
		}
	}

	if req.Name == "" && req.Email == "" {
		validationErrors = append(validationErrors, common.ValidationError{
			Field:   "request",
			Message: "At least one field (name or email) must be provided",
		})
	}

	if len(validationErrors) > 0 {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.URL.Path, validationErrors))
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, req.Name, req.Email)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			common.WriteProblem(w, common.NewNotFoundProblem(err.Error(), r.URL.Path))
			return
		}
		common.WriteProblem(w, common.NewInternalErrorProblem(r.URL.Path))
		return
	}

	common.WriteJSON(w, http.StatusOK, toUserResponse(user))
}
