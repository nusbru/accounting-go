package user

import (
	"net/http"

	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type DeleteUserHandler struct {
	service interfaces.UserService
}

func NewDeleteUserHandler(service interfaces.UserService) *DeleteUserHandler {
	return &DeleteUserHandler{service: service}
}

// Handle deletes a user
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Param id path string true "User ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} common.ValidationProblem
// @Failure 500 {object} common.ProblemDetail
// @Router /api/v1/users/{id} [delete]
func (h *DeleteUserHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.URL.Path))
		return
	}

	id := extractID(r.URL.Path, "/api/v1/users/")
	if err := common.ValidateUUID(id, "id"); err != nil {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.URL.Path, []common.ValidationError{*err}))
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		common.WriteProblem(w, common.NewInternalErrorProblem(r.URL.Path))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
