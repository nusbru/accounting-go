package account

import (
	"encoding/json"
	"errors"
	"net/http"

	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/domain/interfaces"
	"accounting/internal/handler/http/common"
)

type CreateAccountHandler struct {
	service interfaces.AccountService
}

func NewCreateAccountHandler(service interfaces.AccountService) *CreateAccountHandler {
	return &CreateAccountHandler{service: service}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account for a user with the specified details
// @Tags account
// @Accept json
// @Produce json
// @Param request body CreateAccountRequest true "Account creation request"
// @Success 201 {object} AccountResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/accounts [post]
func (h *CreateAccountHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.WriteProblem(w, common.NewMethodNotAllowedProblem(r.RequestURI))
		return
	}

	var req CreateAccountRequest
	if r.Body == nil {
		common.WriteProblem(w, common.NewBadRequestProblem("invalid request body", r.RequestURI))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.WriteProblem(w, common.NewBadRequestProblem("invalid request body", r.RequestURI))
		return
	}

	// Validate request fields
	validationErrors := common.CollectErrors(
		common.ValidateUUID(req.UserID, "user_id"),
		common.ValidateRequired(req.Name, "name"),
		common.ValidateStringLength(req.Name, "name", 1, 100),
		common.ValidateEnum(string(req.Type), []string{"CHECKING", "SAVINGS", "CREDIT_CARD", "CASH", "INVESTMENT"}, "type"),
		common.ValidateCurrency(req.Currency, "currency"),
	)

	if len(validationErrors) > 0 {
		common.WriteProblem(w, common.NewValidationProblemWithErrors(r.RequestURI, validationErrors))
		return
	}

	account, err := h.service.CreateAccount(r.Context(), req.UserID, req.Name, req.Type, req.Currency)
	if err != nil {
		var dupErr *domainerrors.ErrDuplicateAccount
		if errors.As(err, &dupErr) {
			common.WriteProblem(w, common.NewValidationProblem(err.Error(), r.RequestURI))
			return
		}
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

	common.WriteJSON(w, http.StatusCreated, toAccountResponse(account))
}
