package user

import (
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type GetUserHandler struct {
	service interfaces.UserService
}

func NewGetUserHandler(service interfaces.UserService) *GetUserHandler {
	return &GetUserHandler{service: service}
}

// Handle retrieves a user by ID
// @Summary Get user by ID
// @Description Get a user's details by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} UserResponse
// @Failure 400 {object} common.ValidationProblem
// @Failure 404 {object} common.ProblemDetail
// @Failure 500 {object} common.ProblemDetail
// @Router /api/v1/users/{id} [get]
func (h *GetUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.URL.Path))
		return
	}

	id := extractID(r.URL.Path, "/api/v1/users/")
	if err := common.ValidateUUID(id, "id"); err != nil {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.URL.Path, []common.ValidationError{*err}))
		return
	}

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			common.WriteProblem(w, common.NewNotFoundProblem(err.Error(), r.URL.Path))
			return
		}
		common.WriteProblem(w, common.NewInternalErrorProblem(r.URL.Path))
		return
	}
	if user == nil {
		common.WriteProblem(w, common.NewNotFoundProblem("User not found", r.URL.Path))
		return
	}

	common.WriteJSON(w, http.StatusOK, toUserResponse(user))
}
