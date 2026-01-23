package transaction

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"accounting/internal/domain/interfaces"
	domainerrors "accounting/internal/domain/errors"
	"accounting/internal/handler/http/common"
)

type CreateTransactionHandler struct {
	service interfaces.TransactionService
}

func NewCreateTransactionHandler(service interfaces.TransactionService) *CreateTransactionHandler {
	return &CreateTransactionHandler{service: service}
}

// @Summary Create a new transaction
// @Description Create a new transaction for an account
// @Tags transactions
// @Accept json
// @Produce json
// @Param body body CreateTransactionRequest true "Transaction request"
// @Success 201 {object} TransactionResponse
// @Failure 400 {object} common.ValidationProblem "Validation error"
// @Failure 500 {object} common.ProblemDetail "Internal server error"
// @Router /api/v1/transactions [post]
func (h *CreateTransactionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		problem := common.NewMethodNotAllowedProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	var req CreateTransactionRequest
	if r.Body == nil {
		problem := common.NewBadRequestProblem("invalid request body", r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		problem := common.NewBadRequestProblem("invalid request body", r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	// Validate input
	validationErrors := common.CollectErrors(
		common.ValidateUUID(req.AccountID, "account_id"),
		common.ValidatePositive(req.Amount, "amount"),
		common.ValidateCurrency(req.Currency, "currency"),
		common.ValidateEnum(string(req.Type), []string{"INCOME", "EXPENSE", "TRANSFER"}, "type"),
	)
	if len(validationErrors) > 0 {
		problem := common.NewValidationProblemWithErrors(r.RequestURI, validationErrors)
		common.WriteProblem(w, problem)
		return
	}

	var date time.Time
	if req.Date != nil {
		date = *req.Date
	}

	transaction, err := h.service.CreateTransaction(
		r.Context(),
		req.AccountID,
		req.Amount,
		req.Currency,
		req.Description,
		req.Category,
		req.Type,
		date,
	)
	if err != nil {
		var notFoundErr *domainerrors.ErrNotFound
		if errors.As(err, &notFoundErr) {
			problem := common.NewNotFoundProblem(err.Error(), r.RequestURI)
			common.WriteProblem(w, problem)
			return
		}
		var invalidErr *domainerrors.ErrInvalidInput
		if errors.As(err, &invalidErr) {
			problem := common.NewValidationProblem(err.Error(), r.RequestURI)
			common.WriteProblem(w, problem)
			return
		}
		problem := common.NewInternalErrorProblem(r.RequestURI)
		common.WriteProblem(w, problem)
		return
	}

	common.WriteJSON(w, http.StatusCreated, toTransactionResponse(transaction))
}
