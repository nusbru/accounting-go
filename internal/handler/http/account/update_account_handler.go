package account

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type UpdateAccountHandler struct {
	service interfaces.AccountService
}

func NewUpdateAccountHandler(service interfaces.AccountService) *UpdateAccountHandler {
	return &UpdateAccountHandler{service: service}
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update an existing account with the specified details
// @Tags account
// @Accept json
// @Produce json
// @Param account_id path string true "Account ID (UUID)"
// @Param request body UpdateAccountRequest true "Account update request"
// @Success 200 {object} AccountResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 404 {object} common.ProblemDetail "Account not found"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/accounts/{account_id} [put]
func (h *UpdateAccountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	var req UpdateAccountRequest
	if r.Body == nil {
		common.WriteProblem(w, common.NewBadRequestProblem("invalid request body", r.RequestURI))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteProblem(w, common.NewBadRequestProblem("invalid request body", r.RequestURI))
		return
	}

	// Validate optional request fields if provided
	validationErrors := common.CollectErrors()

	if req.Name != "" {
		validationErrors = append(validationErrors, common.CollectErrors(
			common.ValidateStringLength(req.Name, "name", 1, 100),
		)...)
	}

	if req.Type != "" {
		validationErrors = append(validationErrors, common.CollectErrors(
			common.ValidateEnum(string(req.Type), []string{"CHECKING", "SAVINGS", "CREDIT_CARD", "CASH", "INVESTMENT"}, "type"),
		)...)
	}

	if req.Currency != "" {
		validationErrors = append(validationErrors, common.CollectErrors(
			common.ValidateCurrency(req.Currency, "currency"),
		)...)
	}

	if len(validationErrors) > 0 {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.RequestURI, validationErrors))
		return
	}

	account, err := h.service.UpdateAccount(r.Context(), id, req.Name, req.Type, req.Currency)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			common.WriteProblem(w, common.NewNotFoundProblem(err.Error(), r.RequestURI))
			return
		}
		var invalidErr *domainerrors.ErrInvalidInput
		if errors.As(err, &invalidErr) {
			common.WriteProblem(w, common.NewValidationProblem(err.Error(), r.RequestURI))
			return
		}
		common.WriteProblem(w, common.NewInternalErrorProblem(r.RequestURI))
		return
	}

	common.WriteJSON(w, http.StatusOK, toAccountResponse(account))
}
