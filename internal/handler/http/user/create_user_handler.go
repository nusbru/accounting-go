package user

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type CreateUserHandler struct {
	service interfaces.UserService
}

func NewCreateUserHandler(service interfaces.UserService) *CreateUserHandler {
	return &CreateUserHandler{service: service}
}

// Handle creates a new user
// @Summary Create a new user
// @Description Create a new user with name and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User creation request"
// @Success 201 {object} UserResponse
// @Failure 400 {object} common.ValidationProblem
// @Failure 500 {object} common.ProblemDetail
// @Router /api/v1/users [post]
func (h *CreateUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.URL.Path))
		return
	}

	var req CreateUserRequest
	if r.Body == nil {
		common.WriteProblem(w, common.NewBadRequestProblem("Invalid JSON format", r.URL.Path))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteProblem(w, common.NewBadRequestProblem("Invalid JSON format", r.URL.Path))
		return
	}

	// Validate input
	validationErrors := common.CollectErrors(
		common.ValidateRequired(req.Name, "name"),
		common.ValidateStringLength(req.Name, "name", 1, 100),
		common.ValidateEmail(req.Email, "email"),
	)

	if len(validationErrors) > 0 {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.URL.Path, validationErrors))
		return
	}

	user, err := h.service.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		var dupErr *domainerrors.ErrDuplicateEmail
		if errors.As(err, &dupErr) {
			common.WriteProblem(w, common.NewValidationProblem(err.Error(), r.URL.Path))
			return
		}
		var invalidErr *domainerrors.ErrInvalidInput
		if errors.As(err, &invalidErr) {
			common.WriteProblem(w, common.NewValidationProblem(err.Error(), r.URL.Path))
			return
		}
		common.WriteProblem(w, common.NewInternalErrorProblem(r.URL.Path))
		return
	}

	common.WriteJSON(w, http.StatusCreated, toUserResponse(user))
}
