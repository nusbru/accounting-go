package account

import (
	"net/http"

	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type DeleteAccountHandler struct {
	service interfaces.AccountService
}

func NewDeleteAccountHandler(service interfaces.AccountService) *DeleteAccountHandler {
	return &DeleteAccountHandler{service: service}
}

// DeleteAccount godoc
// @Summary Delete an account
// @Description Delete an existing account by ID
// @Tags account
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID (UUID)"
// @Success 204 "Account deleted successfully"
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Account not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/accounts/{account_id} [delete]
func (h *DeleteAccountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.RequestURI))
		return
	}

	id := extractID(r.URL.Path, "/api/v1/accounts/")
	if id == "" {
		validationErrors := common.CollectErrors(
			common.ValidateRequired(id, "account_id"),
		)
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.RequestURI, validationErrors))
		return
	}

	// Validate account ID is a UUID
	if err := common.ValidateUUID(id, "account_id"); err != nil {
		validationErrors := common.CollectErrors(err)
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.RequestURI, validationErrors))
		return
	}

	if err := h.service.DeleteAccount(r.Context(), id); err != nil {
		common.WriteProblem(w, common.NewInternalErrorProblem(r.RequestURI))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
