package account

import (
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type ListUserAccountsHandler struct {
	service interfaces.AccountService
}

func NewListUserAccountsHandler(service interfaces.AccountService) *ListUserAccountsHandler {
	return &ListUserAccountsHandler{service: service}
}

// ListUserAccounts godoc
// @Summary List all accounts for a user
// @Description Retrieve all accounts associated with a specific user
// @Tags account
// @Accept json
// @Produce json
// @Param user_id path string true "User ID (UUID)"
// @Success 200 {array} AccountResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/users/{user_id}/accounts [get]
func (h *ListUserAccountsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.RequestURI))
		return
	}

	userID := extractID(r.URL.Path, "/api/v1/users/")
	if userID == "" {
		validationErrors := common.CollectErrors(
			common.ValidateRequired(userID, "user_id"),
		)
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.RequestURI, validationErrors))
		return
	}

	// Validate user ID is a UUID
	if err := common.ValidateUUID(userID, "user_id"); err != nil {
		validationErrors := common.CollectErrors(err)
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.RequestURI, validationErrors))
		return
	}

	accounts, err := h.service.ListUserAccounts(r.Context(), userID)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			common.WriteProblem(w, common.NewNotFoundProblem(err.Error(), r.RequestURI))
			return
		}
		common.WriteProblem(w, common.NewInternalErrorProblem(r.RequestURI))
		return
	}

	response := make([]*AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		response = append(response, toAccountResponse(acc))
	}

	common.WriteJSON(w, http.StatusOK, response)
}
