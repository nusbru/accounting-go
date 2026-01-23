package account

import (
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type GetAccountHandler struct {
	service interfaces.AccountService
}

func NewGetAccountHandler(service interfaces.AccountService) *GetAccountHandler {
	return &GetAccountHandler{service: service}
}

// GetAccount godoc
// @Summary Get an account by ID
// @Description Retrieve a specific account by its ID
// @Tags account
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID (UUID)"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Account not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/accounts/{account_id} [get]
func (h *GetAccountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	account, err := h.service.GetAccount(r.Context(), id)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			common.WriteProblem(w, common.NewNotFoundProblem(err.Error(), r.RequestURI))
			return
		}
		common.WriteProblem(w, common.NewInternalErrorProblem(r.RequestURI))
		return
	}
	if account == nil {
		common.WriteProblem(w, common.NewNotFoundProblem("account not found", r.RequestURI))
		return
	}

	common.WriteJSON(w, http.StatusOK, toAccountResponse(account))
}
