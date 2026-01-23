package user

import (
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type GetUserByEmailHandler struct {
	service interfaces.UserService
}

func NewGetUserByEmailHandler(service interfaces.UserService) *GetUserByEmailHandler {
	return &GetUserByEmailHandler{service: service}
}

// Handle retrieves a user by email
// @Summary Get user by email
// @Description Get a user's details by their email address
// @Tags users
// @Produce json
// @Param email query string true "User email address"
// @Success 200 {object} UserResponse
// @Failure 400 {object} common.ValidationProblem
// @Failure 404 {object} common.ProblemDetail
// @Failure 500 {object} common.ProblemDetail
// @Router /api/v1/users/search [get]
func (h *GetUserByEmailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.URL.Path))
		return
	}

	email := r.URL.Query().Get("email")
	if err := common.ValidateEmail(email, "email"); err != nil {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.URL.Path, []common.ValidationError{*err}))
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), email)
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
